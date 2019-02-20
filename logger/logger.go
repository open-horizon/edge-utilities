package logger

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	golog "log"
	"log/syslog"
	"os"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/golang/glog"
)

// Parameters parameters for logger setup
type Parameters struct {
	RootPath                 string
	FileName                 string
	MaxFileSize              int
	MaxCompressedFilesNumber int
	Destinations             string
	Prefix                   string
	Level                    string
	MaintenanceInterval      int16
}

// Logger information needed for a logger (or trace)
type Logger struct {
	Tracing                  bool
	Logger                   *golog.Logger
	Level                    int
	MaxFileSize              int64
	MaxCompressedFilesNumber int
	CurrentFile              *os.File
	useLogger                bool
	glog                     bool
	prefix                   string
	Stdout                   bool
	Syslog                   io.Writer
	ticker                   *time.Ticker
	lockChannel              chan int
}

// Error is the error struct used by the logger code
type Error struct {
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

// Log levels
const (
	NONE    = 0
	STATUS  = 1
	FATAL   = 2
	ERROR   = 3
	WARNING = 4
	INFO    = 5
	DEBUG   = 6
	TRACE   = 7
)

var logLevels = map[string]int{
	"NONE": NONE, "STATUS": STATUS, "FATAL": FATAL, "ERROR": ERROR,
	"WARNING": WARNING, "INFO": INFO, "DEBUG": DEBUG, "TRACE": TRACE,
	"XTRACE": TRACE,
}

var logLevelPrefix = []string{"NONE: ", "STATUS: ", "FATAL: ", "ERROR: ", "WARNING: ", "INFO: ", "DEBUG: ", "TRACE: "}
var logLevel2glog = []int{0, 0, 0, 0, 0, 3, 5, 6}

// meaning: STATUS, FATAL, ERROR and WARNING are "gloged" when glog verbosity >= 0 (i.e., always)
//          INFO  is "gloged" when glog verbosity >= 3
//          DEBUG is "gloged" when glog verbosity >= 5
//          TRACE is "gloged" when glog verbosity >= 6

// Init Initialize Logger
func (log *Logger) Init(parameters Parameters) error {
	dests := strings.Split(parameters.Destinations, ",")
	var file, writeToStdout, writeToSyslog, glog bool
	if len(dests) == 0 {
		file = true
	} else {
		for _, dest := range dests {
			if strings.EqualFold(dest, "file") {
				file = true
			} else if strings.EqualFold(dest, "stdout") {
				writeToStdout = true
			} else if strings.EqualFold(dest, "syslog") {
				writeToSyslog = true
			} else if strings.EqualFold(dest, "glog") {
				glog = true
			}
		}
	}

	writers := make([]io.Writer, 0)
	if file {
		info, err := os.Stat(parameters.RootPath)
		if os.IsNotExist(err) {
			err = os.MkdirAll(parameters.RootPath, 0755)
			if err != nil {
				return &Error{fmt.Sprintf("Failed to open log file at %s. Error: %s\n", parameters.RootPath, err)}
			}
		} else {
			if !info.IsDir() {
				return &Error{fmt.Sprintf("Failed to open log file at %s. %s isn't a directory.\n",
					parameters.RootPath, parameters.RootPath)}
			}
		}

		f, err := os.OpenFile(parameters.RootPath+"/"+parameters.FileName+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return &Error{fmt.Sprintf("Failed to open log file at %s. Error: %s\n", parameters.RootPath, err)}
		}
		writers = append(writers, f)
		log.CurrentFile = f
	}
	if writeToStdout {
		writers = append(writers, os.Stdout)
		log.Stdout = true
	}
	if writeToSyslog {
		slWriter, err := syslog.New(syslog.LOG_NOTICE, parameters.FileName)
		if err != nil {
			return &Error{fmt.Sprintf("Failed to create syslog writer. Error: %s\n", err)}
		}
		writers = append(writers, slWriter)
		log.Syslog = slWriter
	}
	if len(writers) == 0 && !glog {
		return &Error{fmt.Sprintf("Invalid log/trace destinations list: %s\n", parameters.Destinations)}
	}

	if len(writers) > 0 {
		mw := io.MultiWriter(writers...)
		if log.Tracing {
			parameters.Prefix = "* " + parameters.Prefix
		}
		log.Logger = golog.New(mw, parameters.Prefix, golog.LstdFlags)
		log.useLogger = true
		log.Level = logLevel(parameters.Level)
		log.MaxFileSize = int64(parameters.MaxFileSize) * 1024
		log.MaxCompressedFilesNumber = parameters.MaxCompressedFilesNumber

		log.lockChannel = make(chan int, 1)
		log.lockChannel <- 1

		log.ticker = time.NewTicker(time.Second * time.Duration(parameters.MaintenanceInterval))
		go func() {
			for {
				select {
				case <-log.ticker.C:
					log.checkFiles()
				}
			}
		}()
	}

	if glog {
		log.Level = logLevel(parameters.Level)
		log.prefix = parameters.Prefix
		log.glog = true
	}
	return nil
}

func (log *Logger) getOldestZipFileNumber() int {
	if log.CurrentFile == nil {
		return 0
	}
	for i := 1; ; i++ {
		fileName := fmt.Sprintf("%s.%d.gz", log.CurrentFile.Name(), i)
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			return i - 1
		}
	}
}

func (log *Logger) checkFiles() {
	if log.CurrentFile == nil {
		return
	}
	fi, err := log.CurrentFile.Stat()
	if err != nil {
		fmt.Printf("Failed to get log file information. Error: %s\n", err)
		return
	}

	if fi.Size() > log.MaxFileSize {
		compressedFiles := log.getOldestZipFileNumber()

		if compressedFiles >= log.MaxCompressedFilesNumber {
			for i := compressedFiles; i > log.MaxCompressedFilesNumber-1; i-- {
				fileName := fmt.Sprintf("%s.%d.gz", log.CurrentFile.Name(), i)
				if err := os.Remove(fileName); err != nil {
					fmt.Printf("Failed to remove compressed log file. Error: %s\n", err)
				}
				compressedFiles--
			}
		}
		for i := compressedFiles; i > 0; i-- {
			fileName := fmt.Sprintf("%s.%d.gz", log.CurrentFile.Name(), i)
			newFileName := fmt.Sprintf("%s.%d.gz", log.CurrentFile.Name(), i+1)
			if os.Rename(fileName, newFileName); err != nil {
				fmt.Printf("Failed to rename compressed log file. Error: %s\n", err)
			}
		}

		savFile := log.CurrentFile
		curFileName := log.CurrentFile.Name()
		savFileName := log.CurrentFile.Name() + ".1"
		zipFileName := log.CurrentFile.Name() + ".1.gz"

		log.lock()
		if err := savFile.Close(); err != nil {
			fmt.Printf("Failed to close the log file. Error: %s\n", err)
			return
		}
		if err = os.Rename(curFileName, savFileName); err != nil {
			fmt.Printf("Failed to rename the log file. Error: %s\n", err)
		}

		log.CurrentFile, err = os.OpenFile(curFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("Failed to open log file %s. Error: %s\n", curFileName, err)
			log.unLock()
			return
		}
		writers := make([]io.Writer, 0)
		writers = append(writers, log.CurrentFile)
		if log.Stdout {
			writers = append(writers, os.Stdout)
		}
		if log.Syslog != nil {
			writers = append(writers, log.Syslog)
		}
		log.Logger.SetOutput(io.MultiWriter(writers...))
		log.unLock()

		savFile, err = os.Open(savFileName)
		if err != nil {
			fmt.Printf("Failed to open log file %s. Error: %s\n", curFileName, err)
			return
		}

		var zipFile *os.File
		zipFile, err = os.Create(zipFileName)
		if err != nil {
			fmt.Printf("Failed to open file to compress log. Error: %s\n", err)
			return
		}
		defer zipFile.Close()

		w := gzip.NewWriter(zipFile)
		if _, err = io.Copy(w, savFile); err != nil {
			fmt.Printf("Failed to copy log to gzip. Error: %s\n", err)
			return
		}
		if err = w.Close(); err != nil {
			fmt.Printf("Failed to close file the compressed log. Error: %s\n", err)
			return
		}
		if err = savFile.Close(); err != nil {
			fmt.Printf("Failed to close the log file. Error: %s\n", err)
			return
		}
		if err = os.Remove(savFileName); err != nil {
			fmt.Printf("Failed to remove the log file. Error: %s\n", err)
			return
		}
	}
}

// Stop Logger
func (log *Logger) Stop() {
	if log.useLogger {
		log.CurrentFile.Close()
		log.ticker.Stop()
	}
	if log.glog {
		glog.Flush()
	}
}

// IsLogging checks if the logging level if higher or equal to the level parameter
func (log *Logger) IsLogging(level int) bool {
	return log.Level >= level || (log.glog && bool(glog.V(glog.Level(logLevel2glog[level]))))
}

func (log *Logger) printf(level int, format string, a ...interface{}) {
	if log.useLogger && log.Level >= level {
		log.lock()
		log.Logger.Printf(logLevelPrefix[level]+format, a...)
		log.unLock()
	}
	if log.glog && bool(glog.V(glog.Level(logLevel2glog[level]))) {
		var b bytes.Buffer
		b.WriteString(log.prefix)
		b.WriteString(logLevelPrefix[level])
		fmt.Fprintf(&b, format, a...)
		line := b.String()
		switch level {
		case FATAL, ERROR:
			glog.ErrorDepth(3, line)
			glog.Flush()
		case WARNING:
			glog.WarningDepth(3, line)
		default:
			glog.InfoDepth(3, line)
		}
	}
}

func (log *Logger) printfAlways(format string, a ...interface{}) {
	if log.useLogger {
		log.lock()
		log.Logger.Printf(format, a...)
		log.unLock()
	}
	if log.glog {
		var b bytes.Buffer
		b.WriteString(log.prefix)
		fmt.Fprintf(&b, format, a...)
		line := b.String()
		glog.InfoDepth(3, line)
	}
}

// Status log
func (log *Logger) Status(format string, a ...interface{}) { log.printf(STATUS, format, a...) }

// Fatal log
func (log *Logger) Fatal(format string, a ...interface{}) { log.printf(FATAL, format, a...) }

// Error log
func (log *Logger) Error(format string, a ...interface{}) { log.printf(ERROR, format, a...) }

// Warning log
func (log *Logger) Warning(format string, a ...interface{}) { log.printf(WARNING, format, a...) }

// Info log
func (log *Logger) Info(format string, a ...interface{}) { log.printf(INFO, format, a...) }

// Debug log
func (log *Logger) Debug(format string, a ...interface{}) { log.printf(DEBUG, format, a...) }

// Trace log
func (log *Logger) Trace(format string, a ...interface{}) { log.printf(TRACE, format, a...) }

// Dump a struct to the logger
func (log *Logger) Dump(label string, a interface{}) {
	objectType := reflect.TypeOf(a)
	if objectType.Kind() != reflect.Struct {
		log.printfAlways("Dump was called with an object that wasn't a struct\n")
		return
	}

	var b strings.Builder
	fmt.Fprintln(&b, label)

	dumpHelper(&b, 2, objectType, a)

	log.printfAlways("%s", b.String())
}

func dumpHelper(writer io.Writer, indent int, objectType reflect.Type, a interface{}) {
	var padBuilder strings.Builder
	padBuilder.Grow(indent)
	for i := 0; i < indent; i++ {
		padBuilder.WriteByte(' ')
	}
	padding := padBuilder.String()

	objectValue := reflect.ValueOf(a)

	fieldCount := objectType.NumField()
	for fieldIndex := 0; fieldIndex < fieldCount; fieldIndex++ {
		field := objectType.Field(fieldIndex)
		value := objectValue.Field(fieldIndex)
		if field.Type.Kind() == reflect.Struct {
			fmt.Fprintf(writer, "%s%s:\n", padding, field.Name)
			dumpHelper(writer, indent+2, field.Type, value.Interface())
		} else {
			fmt.Fprintf(writer, "%s%s  %v\n", padding, field.Name, value)
		}
	}
}

// StackTrace will log the current stack trace
func (log *Logger) StackTrace() {
	var b strings.Builder
	pc := make([]uintptr, 128)
	n := runtime.Callers(3, pc)
	if n == 0 {
		return
	}
	pc = pc[:n]
	frames := runtime.CallersFrames(pc)
	b.WriteString("STACK_TRACE:\n")
	for {
		frame, more := frames.Next()
		fmt.Fprintf(&b, "  %s\n      at %s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	log.printfAlways("%s", b.String())
}

func logLevel(stringLevel string) int {
	level, ok := logLevels[strings.ToUpper(stringLevel)]
	if !ok {
		fmt.Printf("Invalid log level %s specified. Using INFO instead\n", stringLevel)
		level = INFO
	}
	return level
}

func (log *Logger) lock() {
	<-log.lockChannel
}

func (log *Logger) unLock() {
	log.lockChannel <- 1
}

// AdjustMaxLogfileSize insures that the max log file size, when the deafult value is chosen
// is less than 1% of the file system containing the log file
func AdjustMaxLogfileSize(size int, defaultSize int, path string) (int, error) {
	if size == defaultSize {
		var info syscall.Statfs_t
		if err := syscall.Statfs(path, &info); err != nil {
			return size, err
		}
		storageSize := int(info.Blocks * uint64(info.Bsize) / 1024)
		if size > storageSize/100 {
			if storageSize > 2000 {
				return storageSize / 100, nil
			}
			return 20, nil
		}
	}
	return size, nil
}
