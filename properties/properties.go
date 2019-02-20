package properties

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// LoadPropertiesFile Loads the contents of a properties file into a configuration struct
func LoadPropertiesFile(fileName string, optional bool, object interface{}, metaDataKey string) error {
	properties, err := ReadPropertiesFile(fileName, false)
	if err != nil {
		if optional {
			return nil
		}
		return err
	}

	return LoadProperties(properties, object, metaDataKey)
}

// LoadProperties Loads the contents of a map into a configuration struct
func LoadProperties(properties map[string]string, object interface{}, metaDataKey string) error {
	var values = func(key string) (string, bool) {
		value, ok := properties[key]
		return value, ok
	}
	return commonLoad(values, object, metaDataKey)
}

// LoadEnvironment Loads a configuration struct from environment variables
func LoadEnvironment(object interface{}, metaDataKey string) error {
	var values = func(key string) (string, bool) {
		return os.LookupEnv(key)
	}
	return commonLoad(values, object, metaDataKey)
}

// commonLoad Loads values from a helper function into a configuration struct
func commonLoad(values func(string) (string, bool), object interface{}, metaDataKey string) error {
	objectType := reflect.TypeOf(object)
	if objectType.Kind() != reflect.Ptr {
		return errors.New("utility.commonLoad was called with non-pointer object")
	}

	pointeeType := objectType.Elem()
	if pointeeType.Kind() != reflect.Struct {
		return errors.New("utility.commonLoad was called with an object that wasn't a pointer to a struct")
	}
	pointeeValue := reflect.ValueOf(object).Elem()

	fieldCount := pointeeType.NumField()
	for fieldIndex := 0; fieldIndex < fieldCount; fieldIndex++ {
		field := pointeeType.Field(fieldIndex)
		value, ok := values(field.Name)
		if !ok {
			var tagValue string
			tagValue, ok = field.Tag.Lookup(metaDataKey)
			if ok {
				value, ok = values(tagValue)
			}
		}

		if fieldValue := pointeeValue.Field(fieldIndex); ok && fieldValue.CanSet() {
			switch field.Type.Kind() {
			case reflect.Bool:
				var boolValue bool
				switch strings.ToLower(value) {
				case "1":
					boolValue = true

				case "true":
					boolValue = true

				case "t":
					boolValue = true
				}
				fieldValue.SetBool(boolValue)

			case reflect.Int:
				fallthrough
			case reflect.Int8:
				fallthrough
			case reflect.Int16:
				fallthrough
			case reflect.Int32:
				fallthrough
			case reflect.Int64:
				var intValue int64
				if 0 != len(value) {
					_, err := fmt.Sscanf(value, "%d", &intValue)
					if err != nil {
						return err
					}
				}
				fieldValue.SetInt(intValue)

			case reflect.Uint:
				fallthrough
			case reflect.Uint8:
				fallthrough
			case reflect.Uint16:
				fallthrough
			case reflect.Uint32:
				fallthrough
			case reflect.Uint64:
				var uintValue uint64
				if 0 != len(value) {
					_, err := fmt.Sscanf(value, "%d", &uintValue)
					if err != nil {
						return err
					}
				}
				fieldValue.SetUint(uintValue)

			case reflect.String:
				fieldValue.SetString(value)
			}
		}
	}

	return nil
}

// ReadPropertiesFile Reads a properties file into a map[string]string
func ReadPropertiesFile(fileName string, optional bool) (map[string]string, error) {
	result := make(map[string]string)

	rdr, err := os.Open(fileName)
	if err != nil {
		if optional {
			return result, nil
		}
		return nil, err
	}
	defer rdr.Close()

	fileScanner := bufio.NewScanner(rdr)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if len(line) > 0 && line[0] != '#' {
			parts := strings.Fields(line)
			var value []string
			if len(parts) >= 1 {
				key := parts[0]

				if len(parts) > 1 {
					value = parts[1:]
				} else {
					value = []string{""}
				}

				_, ok := result[key]
				if ok {
					return nil, errors.New("The property '" + key + "' is found twice in the file '" + fileName + "'")
				}
				result[key] = strings.Join(value, "")
			}
		}
	}
	if err := fileScanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
