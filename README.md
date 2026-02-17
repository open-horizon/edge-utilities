# edge-utilities

A collection of utilities used in multiple Open Horizon edge computing projects.

## Components

### Logger Package

A comprehensive logging system with support for multiple destinations and automatic log rotation.

**Features:**
- Multiple log levels: NONE, STATUS, FATAL, ERROR, WARNING, INFO, DEBUG, TRACE
- Multiple destinations: File, Stdout, Syslog, Glog
- Automatic log file rotation based on size
- Compression of rotated logs (gzip)
- Configurable retention of compressed files
- Thread-safe operations

**Usage:**
```go
import "github.com/open-horizon/edge-utilities/logger"

params := logger.Parameters{
    RootPath: "/var/log/myapp",
    FileName: "myapp",
    MaxFileSize: 10240, // KB
    MaxCompressedFilesNumber: 5,
    Destinations: "file,stdout",
    Prefix: "MYAPP",
    Level: "INFO",
    MaintenanceInterval: 60, // seconds
}

log := &logger.Logger{}
err := log.Init(params)
```

**Convenience Packages:**
- `logger/log`: Singleton logger instance for application-wide logging
- `logger/trace`: Singleton logger with tracing enabled

### Properties Package

Configuration management utilities for loading settings from multiple sources into Go structs.

**Supported Sources:**
- Properties files (key-value format)
- Environment variables
- In-memory maps

**Supported Types:**
- `bool`: Accepts "1", "true", "t" (case-insensitive)
- `int`, `int8`, `int16`, `int32`, `int64`: Signed integers
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`: Unsigned integers
- `string`: Text values
- `[]string`: String arrays (comma-separated values)

**String Array Support:**

String arrays can be loaded from comma-separated values. Whitespace is automatically trimmed, and empty values are filtered out.

**Usage:**
```go
import "github.com/open-horizon/edge-utilities/properties"

type Config struct {
    Port       int      `json:"PORT"`
    LogLevel   string   `json:"LOG_LEVEL"`
    Debug      bool     `json:"DEBUG"`
    Servers    []string `json:"SERVERS"`      // Comma-separated list
    Extensions []string `json:"EXTENSIONS"`   // Comma-separated list
}

// Load from properties file
config := &Config{}
err := properties.LoadPropertiesFile("config.properties", false, config, "json")

// Load from environment variables
err := properties.LoadEnvironment(config, "json")

// Load from map
props := map[string]string{
    "PORT": "8080",
    "LOG_LEVEL": "INFO",
    "DEBUG": "true",
    "SERVERS": "server1.example.com, server2.example.com, server3.example.com",
    "EXTENSIONS": ".pem,.crt,.cert",
}
err := properties.LoadProperties(props, config, "json")
```

**String Array Examples:**

Environment variable:
```bash
export SERVERS="server1,server2,server3"
export EXTENSIONS=".pem, .crt, .cert"  # Whitespace is trimmed
```

Properties file:
```
SERVERS server1,server2,server3
EXTENSIONS .pem,.crt,.cert
```

Result:
```go
config.Servers    // []string{"server1", "server2", "server3"}
config.Extensions // []string{".pem", ".crt", ".cert"}
```

## Installation

```bash
go get github.com/open-horizon/edge-utilities
```

## Requirements

- Go 1.25.7 or later
- Dependencies managed via Go modules

## License

See LICENSE file for details.
