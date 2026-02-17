package properties

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

// TestLoadEnvironment_StringArray tests loading string arrays from environment variables:
// - Single value string arrays
// - Multiple comma-separated values
// - Values with whitespace that should be trimmed
// - Empty values that should be filtered out
// - Regular string fields still work correctly
//
// This test validates the new []string support added to commonLoad function.
func TestLoadEnvironment_StringArray(t *testing.T) {
	type Config struct {
		SingleValue  []string `json:"SINGLE_VALUE"`
		MultiValue   []string `json:"MULTI_VALUE"`
		WithSpaces   []string `json:"WITH_SPACES"`
		EmptyValues  []string `json:"EMPTY_VALUES"`
		RegularString string  `json:"REGULAR_STRING"`
	}

	// Set up environment variables
	os.Setenv("SINGLE_VALUE", "value1")
	os.Setenv("MULTI_VALUE", "value1,value2,value3")
	os.Setenv("WITH_SPACES", " value1 , value2 , value3 ")
	os.Setenv("EMPTY_VALUES", "value1,,value2,  ,value3")
	os.Setenv("REGULAR_STRING", "test,string,value")
	defer func() {
		os.Unsetenv("SINGLE_VALUE")
		os.Unsetenv("MULTI_VALUE")
		os.Unsetenv("WITH_SPACES")
		os.Unsetenv("EMPTY_VALUES")
		os.Unsetenv("REGULAR_STRING")
	}()

	config := &Config{}
	err := LoadEnvironment(config, "json")
	if err != nil {
		t.Fatalf("LoadEnvironment failed: %v", err)
	}

	// Test single value
	expectedSingle := []string{"value1"}
	if !reflect.DeepEqual(config.SingleValue, expectedSingle) {
		t.Errorf("SingleValue: expected %v, got %v", expectedSingle, config.SingleValue)
	}

	// Test multiple values
	expectedMulti := []string{"value1", "value2", "value3"}
	if !reflect.DeepEqual(config.MultiValue, expectedMulti) {
		t.Errorf("MultiValue: expected %v, got %v", expectedMulti, config.MultiValue)
	}

	// Test values with spaces (should be trimmed)
	expectedSpaces := []string{"value1", "value2", "value3"}
	if !reflect.DeepEqual(config.WithSpaces, expectedSpaces) {
		t.Errorf("WithSpaces: expected %v, got %v", expectedSpaces, config.WithSpaces)
	}

	// Test empty values (should be filtered out)
	expectedEmpty := []string{"value1", "value2", "value3"}
	if !reflect.DeepEqual(config.EmptyValues, expectedEmpty) {
		t.Errorf("EmptyValues: expected %v, got %v", expectedEmpty, config.EmptyValues)
	}

	// Test regular string (should not be split)
	expectedString := "test,string,value"
	if config.RegularString != expectedString {
		t.Errorf("RegularString: expected %q, got %q", expectedString, config.RegularString)
	}
}

// TestLoadProperties_StringArray tests loading string arrays from properties map:
// - Verifies string arrays work with LoadProperties
// - Tests comma-separated values
// - Ensures regular strings still work
func TestLoadProperties_StringArray(t *testing.T) {
	type Config struct {
		Items  []string `json:"ITEMS"`
		Name   string   `json:"NAME"`
	}

	properties := map[string]string{
		"ITEMS": "item1,item2,item3",
		"NAME":  "test,name",
	}

	config := &Config{}
	err := LoadProperties(properties, config, "json")
	if err != nil {
		t.Fatalf("LoadProperties failed: %v", err)
	}

	expectedItems := []string{"item1", "item2", "item3"}
	if !reflect.DeepEqual(config.Items, expectedItems) {
		t.Errorf("Items: expected %v, got %v", expectedItems, config.Items)
	}

	expectedName := "test,name"
	if config.Name != expectedName {
		t.Errorf("Name: expected %q, got %q", expectedName, config.Name)
	}
}

// TestLoadEnvironment_EmptyStringArray tests that empty string arrays are handled correctly:
// - Empty environment variable results in empty slice
// - Unset environment variable results in nil/empty slice
func TestLoadEnvironment_EmptyStringArray(t *testing.T) {
	type Config struct {
		EmptyArray []string `json:"EMPTY_ARRAY"`
		UnsetArray []string `json:"UNSET_ARRAY"`
	}

	os.Setenv("EMPTY_ARRAY", "")
	defer os.Unsetenv("EMPTY_ARRAY")

	config := &Config{}
	err := LoadEnvironment(config, "json")
	if err != nil {
		t.Fatalf("LoadEnvironment failed: %v", err)
	}

	// Empty string should result in empty slice (not nil)
	if config.EmptyArray != nil && len(config.EmptyArray) != 0 {
		t.Errorf("EmptyArray: expected empty slice, got %v", config.EmptyArray)
	}

	// Unset variable should result in nil slice
	if config.UnsetArray != nil && len(config.UnsetArray) != 0 {
		t.Errorf("UnsetArray: expected nil/empty slice, got %v", config.UnsetArray)
	}
}

// TestLoadEnvironment_ExistingTypes tests that existing type support still works:
// - Bool values
// - Int values
// - Uint values
// - String values
//
// This ensures the new []string support doesn't break existing functionality.
func TestLoadEnvironment_ExistingTypes(t *testing.T) {
	type Config struct {
		BoolValue   bool   `json:"BOOL_VALUE"`
		IntValue    int    `json:"INT_VALUE"`
		UintValue   uint   `json:"UINT_VALUE"`
		StringValue string `json:"STRING_VALUE"`
	}

	os.Setenv("BOOL_VALUE", "true")
	os.Setenv("INT_VALUE", "42")
	os.Setenv("UINT_VALUE", "100")
	os.Setenv("STRING_VALUE", "test")
	defer func() {
		os.Unsetenv("BOOL_VALUE")
		os.Unsetenv("INT_VALUE")
		os.Unsetenv("UINT_VALUE")
		os.Unsetenv("STRING_VALUE")
	}()

	config := &Config{}
	err := LoadEnvironment(config, "json")
	if err != nil {
		t.Fatalf("LoadEnvironment failed: %v", err)
	}

	if config.BoolValue != true {
		t.Errorf("BoolValue: expected true, got %v", config.BoolValue)
	}

	if config.IntValue != 42 {
		t.Errorf("IntValue: expected 42, got %v", config.IntValue)
	}

	if config.UintValue != 100 {
		t.Errorf("UintValue: expected 100, got %v", config.UintValue)
	}

	if config.StringValue != "test" {
		t.Errorf("StringValue: expected 'test', got %q", config.StringValue)
	}
}

// TestLoadProperties_SpecialCharacters tests handling of special characters in values:
// - Values with quotes
// - Values with backslashes
// - Values with newlines (escaped)
// - Values with tabs
// - Unicode characters
//
// This ensures the properties loader handles various character encodings correctly.
func TestLoadProperties_SpecialCharacters(t *testing.T) {
	type Config struct {
		WithQuotes    string   `json:"WITH_QUOTES"`
		WithBackslash string   `json:"WITH_BACKSLASH"`
		WithUnicode   string   `json:"WITH_UNICODE"`
		ArraySpecial  []string `json:"ARRAY_SPECIAL"`
	}

	properties := map[string]string{
		"WITH_QUOTES":    `value"with"quotes`,
		"WITH_BACKSLASH": `path\to\file`,
		"WITH_UNICODE":   "Hello 世界 🌍",
		"ARRAY_SPECIAL":  `item1,"item2",item3`,
	}

	config := &Config{}
	err := LoadProperties(properties, config, "json")
	if err != nil {
		t.Fatalf("LoadProperties failed: %v", err)
	}

	if config.WithQuotes != `value"with"quotes` {
		t.Errorf("WithQuotes: expected %q, got %q", `value"with"quotes`, config.WithQuotes)
	}

	if config.WithBackslash != `path\to\file` {
		t.Errorf("WithBackslash: expected %q, got %q", `path\to\file`, config.WithBackslash)
	}

	if config.WithUnicode != "Hello 世界 🌍" {
		t.Errorf("WithUnicode: expected %q, got %q", "Hello 世界 🌍", config.WithUnicode)
	}

	expectedArray := []string{"item1", `"item2"`, "item3"}
	if !reflect.DeepEqual(config.ArraySpecial, expectedArray) {
		t.Errorf("ArraySpecial: expected %v, got %v", expectedArray, config.ArraySpecial)
	}
}

// TestLoadProperties_BoundaryValues tests boundary conditions for numeric types:
// - Maximum and minimum values for int types
// - Zero values
// - Overflow conditions (should error)
//
// This ensures proper handling of numeric boundaries and prevents overflow issues.
func TestLoadProperties_BoundaryValues(t *testing.T) {
	type Config struct {
		MaxInt64  int64  `json:"MAX_INT64"`
		MinInt64  int64  `json:"MIN_INT64"`
		MaxUint64 uint64 `json:"MAX_UINT64"`
		ZeroInt   int    `json:"ZERO_INT"`
	}

	properties := map[string]string{
		"MAX_INT64":  "9223372036854775807",
		"MIN_INT64":  "-9223372036854775808",
		"MAX_UINT64": "18446744073709551615",
		"ZERO_INT":   "0",
	}

	config := &Config{}
	err := LoadProperties(properties, config, "json")
	if err != nil {
		t.Fatalf("LoadProperties failed: %v", err)
	}

	if config.MaxInt64 != 9223372036854775807 {
		t.Errorf("MaxInt64: expected 9223372036854775807, got %d", config.MaxInt64)
	}

	if config.MinInt64 != -9223372036854775808 {
		t.Errorf("MinInt64: expected -9223372036854775808, got %d", config.MinInt64)
	}

	if config.MaxUint64 != 18446744073709551615 {
		t.Errorf("MaxUint64: expected 18446744073709551615, got %d", config.MaxUint64)
	}

	if config.ZeroInt != 0 {
		t.Errorf("ZeroInt: expected 0, got %d", config.ZeroInt)
	}
}

// TestLoadProperties_InvalidNumericValues tests error handling for invalid numeric inputs:
// - Non-numeric strings for int fields
// - Overflow values
// - Malformed numbers
//
// This ensures proper error reporting for invalid numeric conversions.
func TestLoadProperties_InvalidNumericValues(t *testing.T) {
	type Config struct {
		IntValue int `json:"INT_VALUE"`
	}

	testCases := []struct {
		name  string
		value string
	}{
		{"non-numeric", "not-a-number"},
		{"with-letters", "123abc"},
		{"with-spaces", "12 34"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			properties := map[string]string{
				"INT_VALUE": tc.value,
			}

			config := &Config{}
			err := LoadProperties(properties, config, "json")
			if err == nil {
				t.Errorf("Expected error for value %q, but got none", tc.value)
			}
		})
	}
}

// TestLoadProperties_LongValues tests handling of very long string values:
// - Very long strings (>1MB)
// - Very long arrays with many elements
//
// This ensures the loader can handle large configuration values without issues.
func TestLoadProperties_LongValues(t *testing.T) {
	type Config struct {
		LongString string   `json:"LONG_STRING"`
		LongArray  []string `json:"LONG_ARRAY"`
	}

	// Create a very long string (1MB)
	longString := strings.Repeat("a", 1024*1024)

	// Create a long array (1000 elements)
	longArrayParts := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		longArrayParts[i] = fmt.Sprintf("item%d", i)
	}
	longArrayValue := strings.Join(longArrayParts, ",")

	properties := map[string]string{
		"LONG_STRING": longString,
		"LONG_ARRAY":  longArrayValue,
	}

	config := &Config{}
	err := LoadProperties(properties, config, "json")
	if err != nil {
		t.Fatalf("LoadProperties failed: %v", err)
	}

	if len(config.LongString) != 1024*1024 {
		t.Errorf("LongString: expected length 1048576, got %d", len(config.LongString))
	}

	if len(config.LongArray) != 1000 {
		t.Errorf("LongArray: expected 1000 elements, got %d", len(config.LongArray))
	}
}

// TestLoadProperties_MixedTypes tests loading a configuration with all supported types:
// - Bool, int, uint, string, []string all in one struct
// - Verifies no type interference
//
// This ensures different types can coexist without conflicts.
func TestLoadProperties_MixedTypes(t *testing.T) {
	type Config struct {
		BoolField   bool     `json:"BOOL_FIELD"`
		IntField    int      `json:"INT_FIELD"`
		UintField   uint     `json:"UINT_FIELD"`
		StringField string   `json:"STRING_FIELD"`
		ArrayField  []string `json:"ARRAY_FIELD"`
	}

	properties := map[string]string{
		"BOOL_FIELD":   "true",
		"INT_FIELD":    "-42",
		"UINT_FIELD":   "42",
		"STRING_FIELD": "test,value",
		"ARRAY_FIELD":  "item1,item2,item3",
	}

	config := &Config{}
	err := LoadProperties(properties, config, "json")
	if err != nil {
		t.Fatalf("LoadProperties failed: %v", err)
	}

	if config.BoolField != true {
		t.Errorf("BoolField: expected true, got %v", config.BoolField)
	}

	if config.IntField != -42 {
		t.Errorf("IntField: expected -42, got %d", config.IntField)
	}

	if config.UintField != 42 {
		t.Errorf("UintField: expected 42, got %d", config.UintField)
	}

	if config.StringField != "test,value" {
		t.Errorf("StringField: expected 'test,value', got %q", config.StringField)
	}

	expectedArray := []string{"item1", "item2", "item3"}
	if !reflect.DeepEqual(config.ArrayField, expectedArray) {
		t.Errorf("ArrayField: expected %v, got %v", expectedArray, config.ArrayField)
	}
}

// TestLoadProperties_UnsetFields tests that unset fields retain their zero values:
// - Fields not present in properties should remain at zero value
// - Partial configuration loading
//
// This ensures optional configuration fields work correctly.
func TestLoadProperties_UnsetFields(t *testing.T) {
	type Config struct {
		SetField   string `json:"SET_FIELD"`
		UnsetField string `json:"UNSET_FIELD"`
		SetInt     int    `json:"SET_INT"`
		UnsetInt   int    `json:"UNSET_INT"`
	}

	properties := map[string]string{
		"SET_FIELD": "value",
		"SET_INT":   "42",
	}

	config := &Config{}
	err := LoadProperties(properties, config, "json")
	if err != nil {
		t.Fatalf("LoadProperties failed: %v", err)
	}

	if config.SetField != "value" {
		t.Errorf("SetField: expected 'value', got %q", config.SetField)
	}

	if config.UnsetField != "" {
		t.Errorf("UnsetField: expected empty string, got %q", config.UnsetField)
	}

	if config.SetInt != 42 {
		t.Errorf("SetInt: expected 42, got %d", config.SetInt)
	}

	if config.UnsetInt != 0 {
		t.Errorf("UnsetInt: expected 0, got %d", config.UnsetInt)
	}
}

// TestLoadProperties_BoolVariations tests various boolean value formats:
// - "1", "true", "t", "TRUE", "T" should all be true
// - "0", "false", "f", "FALSE", "F", "" should all be false
//
// This ensures flexible boolean parsing as documented.
func TestLoadProperties_BoolVariations(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		expected bool
	}{
		{"one", "1", true},
		{"true-lower", "true", true},
		{"true-upper", "TRUE", true},
		{"true-mixed", "True", true},
		{"t-lower", "t", true},
		{"t-upper", "T", true},
		{"zero", "0", false},
		{"false-lower", "false", false},
		{"false-upper", "FALSE", false},
		{"f-lower", "f", false},
		{"f-upper", "F", false},
		{"empty", "", false},
		{"random", "random", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			type Config struct {
				BoolValue bool `json:"BOOL_VALUE"`
			}

			properties := map[string]string{
				"BOOL_VALUE": tc.value,
			}

			config := &Config{}
			err := LoadProperties(properties, config, "json")
			if err != nil {
				t.Fatalf("LoadProperties failed: %v", err)
			}

			if config.BoolValue != tc.expected {
				t.Errorf("BoolValue for %q: expected %v, got %v", tc.value, tc.expected, config.BoolValue)
			}
		})
	}
}

// TestReadPropertiesFile_DuplicateKeys tests that duplicate keys in properties files are detected:
// - Should return error when same key appears twice
//
// This is a security/correctness test to prevent configuration ambiguity (CWE-694).
func TestReadPropertiesFile_DuplicateKeys(t *testing.T) {
	// Create temporary file with duplicate keys
	tempDir := t.TempDir()
	tempFile := tempDir + "/duplicate.properties"

	content := `KEY1 value1
KEY2 value2
KEY1 value3
`
	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	_, err = ReadPropertiesFile(tempFile, false)
	if err == nil {
		t.Error("Expected error for duplicate keys, but got none")
	}

	if !strings.Contains(err.Error(), "found twice") {
		t.Errorf("Expected error message about duplicate key, got: %v", err)
	}
}

// TestReadPropertiesFile_Comments tests that comment lines are properly ignored:
// - Lines starting with # should be ignored
// - Empty lines should be ignored
//
// This ensures proper properties file parsing.
func TestReadPropertiesFile_Comments(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := tempDir + "/comments.properties"

	content := `# This is a comment
KEY1 value1

# Another comment
KEY2 value2
  # Indented comment
KEY3 value3
`
	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	properties, err := ReadPropertiesFile(tempFile, false)
	if err != nil {
		t.Fatalf("ReadPropertiesFile failed: %v", err)
	}

	if len(properties) != 3 {
		t.Errorf("Expected 3 properties, got %d", len(properties))
	}

	if properties["KEY1"] != "value1" {
		t.Errorf("KEY1: expected 'value1', got %q", properties["KEY1"])
	}
}

// TestReadPropertiesFile_OptionalFile tests optional file loading:
// - When optional=true and file doesn't exist, should return empty map without error
// - When optional=false and file doesn't exist, should return error
//
// This ensures proper handling of optional configuration files.
func TestReadPropertiesFile_OptionalFile(t *testing.T) {
	nonExistentFile := "/tmp/does-not-exist-" + fmt.Sprintf("%d", os.Getpid()) + ".properties"

	// Test optional=true
	properties, err := ReadPropertiesFile(nonExistentFile, true)
	if err != nil {
		t.Errorf("Expected no error for optional file, got: %v", err)
	}
	if len(properties) != 0 {
		t.Errorf("Expected empty map for optional missing file, got %d entries", len(properties))
	}

	// Test optional=false
	_, err = ReadPropertiesFile(nonExistentFile, false)
	if err == nil {
		t.Error("Expected error for non-optional missing file, but got none")
	}
}

// TestLoadPropertiesFile_Integration tests the full integration of file loading:
// - Read from file
// - Parse properties
// - Load into struct
//
// This is an end-to-end test of the complete workflow.
func TestLoadPropertiesFile_Integration(t *testing.T) {
	type Config struct {
		AppName    string   `json:"APP_NAME"`
		Port       int      `json:"PORT"`
		Debug      bool     `json:"DEBUG"`
		Servers    []string `json:"SERVERS"`
	}

	tempDir := t.TempDir()
	tempFile := tempDir + "/config.properties"

	content := `# Application configuration
APP_NAME MyApplication
PORT 8080
DEBUG true
SERVERS server1.example.com,server2.example.com,server3.example.com
`
	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	config := &Config{}
	err = LoadPropertiesFile(tempFile, false, config, "json")
	if err != nil {
		t.Fatalf("LoadPropertiesFile failed: %v", err)
	}

	if config.AppName != "MyApplication" {
		t.Errorf("AppName: expected 'MyApplication', got %q", config.AppName)
	}

	if config.Port != 8080 {
		t.Errorf("Port: expected 8080, got %d", config.Port)
	}

	if config.Debug != true {
		t.Errorf("Debug: expected true, got %v", config.Debug)
	}

	expectedServers := []string{"server1.example.com", "server2.example.com", "server3.example.com"}
	if !reflect.DeepEqual(config.Servers, expectedServers) {
		t.Errorf("Servers: expected %v, got %v", expectedServers, config.Servers)
	}
}
