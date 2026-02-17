# AGENTS.md

This file provides guidance to agents when working with code in this repository.

## Project Overview

This is a Go utility library (`github.com/open-horizon/edge-utilities`) that provides reusable components for Open Horizon edge computing projects. It is designed to be imported as a dependency rather than run as a standalone application.

**Primary Components:**

1. **Logger Package** (`logger/`): A comprehensive logging system with advanced features
2. **Properties Package** (`properties/`): Configuration management utilities
3. **Convenience Wrappers** (`logger/log/`, `logger/trace/`): Simplified interfaces to the logger

**Technology Stack:**
- Language: Go
- Dependency Management: Vendor directory with `vendor.json`
- External Dependency: `github.com/golang/glog`

## Logger Package

The logger package provides a flexible logging system with the following capabilities:

**Log Levels:**
- NONE, STATUS, FATAL, ERROR, WARNING, INFO, DEBUG, TRACE

**Destinations:**
- File (with automatic rotation and gzip compression)
- Stdout
- Syslog
- Glog (Google's logging library)

**Key Features:**
- Multiple simultaneous log destinations
- Automatic log file rotation based on size
- Compression of rotated logs
- Configurable retention of compressed files
- Thread-safe operations with internal locking
- Stack trace logging
- Struct dumping for debugging
- Integration with glog for verbosity-based logging

**Usage Pattern:**
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
- `logger/trace`: Singleton logger with tracing enabled (prefixes messages with `*`)

## Properties Package

Provides utilities for loading configuration from multiple sources into Go structs:

**Supported Sources:**
- Properties files (key-value format)
- Environment variables
- In-memory maps

**Features:**
- Automatic type conversion (bool, int, uint, string)
- Struct tag support for custom key mapping
- Optional file loading (no error if file missing)
- Duplicate key detection

**Usage Pattern:**
```go
import "github.com/open-horizon/edge-utilities/properties"

type Config struct {
    Port     int    `json:"PORT"`
    LogLevel string `json:"LOG_LEVEL"`
    Debug    bool
}

config := &Config{}
err := properties.LoadPropertiesFile("config.properties", false, config, "json")
// Or from environment:
err := properties.LoadEnvironment(config, "json")
```

## Development Conventions

### Code Organization

**File and Module Organization Principles:**

1. **Keep Files Manageable**: Aim to keep individual source files under 1000 lines where practical
   - If a file grows beyond 1500 lines, consider splitting it into multiple files or subpackages
   - Large files are harder to navigate, review, and maintain

2. **Leverage Go's Package Structure**: 
   - Create new packages/directories when a set of related functionality becomes substantial
   - Each package should have a clear, single responsibility
   - Use internal packages for implementation details that shouldn't be exposed

3. **When to Create New Modules/Directories**:
   - When a feature or component has multiple related types and functions (>500 lines)
   - When functionality is logically distinct and could be tested independently
   - When code could potentially be reused by other parts of the system
   - When a file has grown too large and splitting within the same package isn't sufficient

4. **File Splitting Strategies**:
   - Split by functionality: `logger.go`, `logger_helpers.go`, `logger_rotation.go`
   - Split by type: `file_logger.go`, `syslog_logger.go`, `glog_logger.go`
   - Split by concern: `api.go`, `api_handlers.go`, `api_validation.go`
   - Keep related tests in corresponding `*_test.go` files

5. **Package Naming**:
   - Use short, descriptive package names (e.g., `logger`, not `loggermanagement`)
   - Avoid generic names like `util`, `common`, `helpers` for new packages
   - Package name should describe what it provides, not what it contains

### Source Maintainability Principles

**Avoid Hardcoded Duplication - Favor Centralized Configuration:**

1. **Single Source of Truth**: When the same values, constants, or structures are needed across multiple files or packages, create a centralized, configurable structure rather than duplicating hardcoded values.

2. **Configuration Over Hardcoding**: 
   - **Bad Practice**: Hardcoding `[]string{".pem", ".crt", ".cert"}` in 6 different files
   - **Good Practice**: Define `AllowedCertificateExtensions` once in configuration with defaults, reference everywhere
   - Benefits: Single point of change, consistent behavior, user customization, easier testing

3. **Reusable Structures**:
   - Create shared types, constants, and helper functions in common packages
   - Use configuration parameters for values that might need customization
   - Document the centralized structure and its usage in code comments

4. **When to Centralize**:
   - Values used in 3+ locations
   - Security-related constants (file extensions, timeouts, limits)
   - Business logic constants that might change
   - Environment-specific settings
   - Validation rules and constraints

5. **Maintainability Benefits**:
   - **Consistency**: Same behavior across all usages
   - **Flexibility**: Easy to customize per deployment
   - **Testability**: Single place to mock or override for tests
   - **Documentation**: Clear intent and purpose in one location
   - **Evolution**: Easy to extend or modify without hunting through codebase

### Code Style and Organization

**Indentation and Formatting:**
- **Use spaces, not tabs** for indentation
- **Standard indent size: 4 spaces**
- Go's `gofmt` tool enforces tabs by default, but for consistency with other languages and better cross-platform compatibility, convert tabs to spaces in non-Go files
- For Go files, follow `gofmt` conventions (which uses tabs)
- **Trim trailing whitespace** from all lines when editing files
- **End files with a single blank line** (newline at EOF)
- Ensure your editor is configured to:
  - Insert spaces when Tab key is pressed
  - Display tabs as 4 spaces for consistency
  - Convert existing tabs to spaces when editing files
  - Automatically trim trailing whitespace on save
  - Ensure newline at end of file

**Function Ordering:**
- Prefer lexical (alphabetical) ordering of functions within a file where it makes sense
- Exceptions: Group related helper functions near their primary function when it improves comprehension
- Public functions (exported) should generally appear before private functions (unexported)
- Keep `init()` functions at the top of the file after package-level variables

**Lexical Ordering for Other Constructs:**
Apply alphabetical ordering where it makes sense for:
- **Variable declarations**: Group related variables, but within groups prefer alphabetical order
- **Struct fields**: Order fields alphabetically unless logical grouping (e.g., related fields, embedded types first) improves clarity
- **Constants**: Within const blocks, prefer alphabetical ordering
- **Interface methods**: Order methods alphabetically in interface definitions
- **Configuration parameters**: In configuration structs, order fields alphabetically for easier lookup
- **Import statements**: Go automatically formats imports, but group standard library, external, and internal imports

**When NOT to use lexical ordering:**
- When fields have dependencies or initialization order matters
- When grouping by functionality significantly improves comprehension
- When following established patterns in the existing codebase
- When struct field order affects memory layout for performance reasons

**Readability Principles:**
- Prioritize code legibility and readability without sacrificing technical depth or solution quality
- Use clear, descriptive variable and function names that convey intent
- Add comments for complex logic, but prefer self-documenting code
- Break down complex functions into smaller, well-named helper functions
- Use whitespace and formatting to visually separate logical blocks
- Avoid deeply nested code; prefer early returns and guard clauses

**Code Structure:**
- Keep functions focused on a single responsibility
- Limit function length to what fits on a screen when possible (aim for < 50 lines)
- Use consistent error handling patterns throughout the codebase
- Document exported functions, types, and constants with godoc-style comments
- Place package-level constants and variables at the top of files

**Example of Good Organization:**
```go
package example

// Package-level constants
const (
    DefaultTimeout = 30
    MaxRetries     = 3
)

// Package-level variables
var (
    globalCache map[string]interface{}
)

// init function
func init() {
    globalCache = make(map[string]interface{})
}

// Exported functions (alphabetically ordered)
func CreateResource(name string) error { ... }

func DeleteResource(id string) error { ... }

func GetResource(id string) (*Resource, error) { ... }

func UpdateResource(id string, data []byte) error { ... }

// Unexported helper functions (alphabetically ordered)
func buildResourceKey(id string) string { ... }

func validateResourceData(data []byte) error { ... }
```

### Multilingual Support Requirements

**Library API Responses:**
- All public API functions that return user-facing messages MUST handle multilingual support
- Error messages and status messages must be localizable
- Use appropriate internationalization (i18n) mechanisms for API responses
- Ensure consistent message formatting across different languages

**Logging Messages:**
- Logging messages are more permissive and can default to English
- Internal debug/trace logs do not require multilingual support
- Focus multilingual efforts on user-facing API responses and error messages
- Log messages should still be clear and descriptive for debugging purposes

### Dependency Management

**Go Module Updates and Security:**

1. **Regular Dependency Checks**: 
   - Check for dependency updates during any ongoing source code work
   - Use `go list -m -u all` to check for available updates
   - Use `go mod tidy` to clean up unused dependencies

2. **Security Vulnerability Scanning**:
   - Run `go list -json -m all | nancy sleuth` or similar tools to check for CVEs
   - Use GitHub's Dependabot or similar automated scanning tools
   - Prioritize updates that address security vulnerabilities (CVEs)
   - Document CVE fixes in commit messages and pull requests

3. **Update Strategy - Prioritize Stability**:
   - **Critical Security Updates**: Apply immediately after testing
   - **Minor/Patch Updates**: Apply regularly, test thoroughly
   - **Major Version Updates**: Evaluate carefully, may require code changes
   - Always run full test suite after dependency updates: `go test ./...`
   - Always run race detector after updates: `go test -race ./...`

4. **Testing After Updates**:
   ```bash
   # Update dependencies
   go get -u ./...
   go mod tidy
   
   # Verify build
   go build ./...
   
   # Run tests
   go test ./...
   go test -race ./...
   ```

5. **Dependency Update Guidelines**:
   - Never update dependencies without running the full test suite
   - Document breaking changes in dependency updates
   - If an update breaks the build, either fix the code or pin to the previous version
   - Commit `go.mod` and `go.sum` changes together with any required code changes

6. **Handling Breaking Changes**:
   - Review changelogs and migration guides for major version updates
   - Create separate commits for dependency updates vs. code adaptations
   - Consider impact on downstream projects that depend on this library

**Go Language Version Updates:**

1. **Regular Version Checks**:
   - Check for Go language updates during ongoing development work
   - Monitor Go release notes for security patches and improvements

2. **Update Strategy - Prioritize Stability**:
   - **Security Patches**: Apply Go patch releases promptly after testing
   - **Minor Version Updates**: Evaluate and test thoroughly before upgrading
   - **Major Version Updates**: Plan carefully, may require code changes
   - **Never break the build**: Always verify builds succeed before committing version changes

3. **Testing After Go Version Updates**:
   ```bash
   # Update go.mod
   go mod edit -go=1.22  # Example version
   go mod tidy
   
   # Verify build
   go build ./...
   
   # Run tests
   go test ./...
   go test -race ./...
   ```

4. **Version Update Guidelines**:
   - Update `go.mod` file with new Go version
   - Update documentation (README.md, AGENTS.md) with new version requirement
   - Update CI/CD configurations if applicable
   - Run full test suite including race detection
   - Document any code changes required for the new version

5. **Backward Compatibility**:
   - Maintain compatibility with previous minor version when possible
   - Document minimum required Go version clearly
   - Consider impact on projects that depend on this library

### Security Considerations

1. **Path Validation**: 
   - **CRITICAL**: All file operations MUST validate paths
   - Protect against:
     - Path traversal attacks (../, ../../, etc.)
     - Null byte injection (CWE-158)
     - Symlink attacks (CWE-61)
     - Access to files outside allowed directories
   - **Usage example**:
     ```go
     // Validate file path before use
     if !isValidPath(userPath, baseDir) {
         return fmt.Errorf("invalid path: %s", userPath)
     }
     ```

2. **Input Validation**:
   - Validate all inputs from external sources
   - Use appropriate type checking and bounds checking
   - Sanitize inputs before using in file operations or system calls

3. **Error Handling**:
   - Never expose sensitive information in error messages
   - Log detailed errors internally, return generic errors to callers
   - Use custom error types for better error handling

### Testing Practices

**Test-First Development Methodology**

This project follows a test-first development approach:

1. **Tests Define Expected Behavior**: Write tests that describe the correct, desired behavior of the system
   - Tests should reflect what the code SHOULD do, not what it currently does
   - Never modify test expectations to match existing functional deficiencies
   - If a test fails due to incorrect implementation, fix the implementation, not the test

2. **Functional Improvements Over Test Adjustments**: When tests reveal issues:
   - **Correct Approach**: Improve the implementation to make the test pass
   - **Incorrect Approach**: Change the test to match the broken behavior
   - Exception: Only adjust tests if the original test expectations were genuinely incorrect

3. **Test-Driven Bug Fixes**:
   - Write a failing test that demonstrates the bug
   - Verify the test fails with current code
   - Fix the implementation to make the test pass
   - Never adjust the test to accept the buggy behavior

4. **Example Scenario**:
   ```
   BAD:  Test expects validation to reject invalid input
         → Implementation doesn't validate
         → Change test to accept invalid input ❌
   
   GOOD: Test expects validation to reject invalid input
         → Implementation doesn't validate
         → Add validation to implementation ✓
   ```

5. **When to Adjust Tests**:
   - Original test expectations were based on misunderstanding requirements
   - Requirements have legitimately changed
   - Test was testing implementation details rather than behavior
   - Never adjust tests simply because implementation is difficult to fix

**CRITICAL: Test Coverage Requirements**

All source code changes MUST be accompanied by corresponding test updates:

1. **New Features**: Add comprehensive test cases covering:
   - Happy path scenarios
   - Edge cases and boundary conditions
   - Error handling paths
   - Concurrent access patterns (with race detection)

2. **Bug Fixes**: MANDATORY test requirements:
   - Add a test case that reproduces the bug BEFORE fixing it
   - Verify the test fails with the bug present
   - Verify the test passes after the fix
   - Document the bug scenario in test comments

3. **Security Fixes**: CRITICAL test requirements:
   - Add test cases demonstrating the vulnerability (safely)
   - Test both the attack vector and the mitigation
   - Include tests for bypass attempts
   - Document CVE or security issue reference if applicable

4. **Refactoring**: Ensure existing tests still pass and add tests for:
   - New code paths introduced
   - Changed behavior or interfaces
   - Performance characteristics if relevant

**Test Isolation Principle:**

All tests MUST be designed with proper isolation to ensure reliability and maintainability:

1. **Independent Test Execution**: Each test must be able to run independently without relying on:
   - Execution order of other tests
   - State left behind by previous tests
   - Shared global state that persists between tests

2. **Clean Test Environment**: 
   - Use `t.TempDir()` for temporary file operations (automatically cleaned up)
   - Reset or mock global state in `setUp` functions
   - Clean up resources in `defer` statements or `t.Cleanup()`

3. **Avoid Test Interdependencies**:
   - Never assume test execution order
   - Each test should set up its own required state
   - Don't share test data structures between test cases
   - Use table-driven tests with independent test cases

4. **Parallel Test Safety**:
   - Mark tests as parallel-safe with `t.Parallel()` when appropriate
   - Ensure parallel tests don't share mutable state
   - Test with race detector: `go test -race`

5. **Test Data Isolation**:
   - Generate unique test data for each test case (e.g., unique IDs, timestamps)
   - Avoid hardcoded test data that could conflict across tests
   - Use test-specific prefixes or namespaces for identifiers
   - Clean up test data immediately after test completion

6. **Configuration Isolation**:
   - Create test-specific configuration instances
   - Never modify global configuration in tests
   - Use configuration mocks or test doubles when needed
   - Reset configuration state in test cleanup

7. **Time and Randomness Isolation**:
   - Mock time-dependent functions for deterministic tests
   - Use fixed seeds for random number generators in tests
   - Avoid tests that depend on wall-clock time
   - Make time-sensitive tests configurable with timeouts

8. **Example of Good Test Isolation**:
   ```go
   func TestLoggerOperation(t *testing.T) {
       t.Parallel() // Safe to run in parallel
       
       // Create isolated test environment
       tempDir := t.TempDir() // Auto-cleanup
       
       // Generate unique test data
       testID := fmt.Sprintf("test-%d", time.Now().UnixNano())
       
       // Clean up any resources
       t.Cleanup(func() {
           // Cleanup logic
       })
       
       // Test logic with isolated state
       // ...
   }
   ```

**Test Cleanup and Environmental Responsibility:**

All tests MUST clean up after themselves and avoid leaving destructive changes to the environment:

1. **Complete Cleanup**: Every test must restore the environment to its original state
   - Delete temporary files and directories created during testing
   - Close all open file handles and connections
   - Restore modified configuration or global state
   - Use `defer` statements or `t.Cleanup()` to ensure cleanup happens even on test failure

2. **No Destructive Side Effects**: Tests must not:
   - Modify files outside the test's temporary directory
   - Delete or overwrite production data or configuration
   - Leave processes running after test completion
   - Consume system resources indefinitely (memory leaks, goroutine leaks)
   - Modify shared system state that affects other tests or processes

3. **Idempotent Tests**: Tests should be repeatable without manual cleanup
   - Running the same test multiple times should produce the same results
   - Tests should not depend on previous test runs being cleaned up
   - Use unique identifiers to avoid conflicts with concurrent test runs

4. **Resource Management**:
   - Always close resources in `defer` or `t.Cleanup()` blocks
   - Use context with timeout for operations that might hang
   - Monitor and clean up goroutines to prevent leaks
   - Release locks and semaphores properly

5. **Example of Proper Cleanup**:
   ```go
   func TestWithProperCleanup(t *testing.T) {
       // Create temporary directory (auto-cleanup)
       tempDir := t.TempDir()
       
       // Create test logger
       logger := &Logger{}
       err := logger.Init(params)
       require.NoError(t, err)
       defer logger.Stop() // Ensure logger is stopped
       
       // Register cleanup
       t.Cleanup(func() {
           // Additional cleanup if needed
       })
       
       // Test logic here
       // ...
   }
   ```

**When Creating or Updating Tests:**

- **Always verify test isolation**: Run the test multiple times in different orders
- **Test in parallel**: Use `go test -parallel=10` to expose isolation issues
- **Check for race conditions**: Always run `go test -race` on concurrent code
- **Verify cleanup**: Ensure no test artifacts remain after test completion
- **Check for resource leaks**: Monitor goroutines, file handles, and memory usage
- **Document dependencies**: Clearly document any external service requirements
- **Use subtests for variations**: Group related test cases with `t.Run()` for better organization

**Test File Conventions:**
- Unit tests in `*_test.go` files alongside source
- Race detection enabled for concurrency testing: `go test -race`
- Test logger and properties loading with various configurations

**Test Naming:**
- Use descriptive test names: `TestFunctionName_Scenario_ExpectedBehavior`
- Security tests: Prefix with vulnerability type (e.g., `TestLoadFile_PathTraversal_Blocked`)
- Bug fix tests: Reference issue number if available (e.g., `TestIssue123_NullPointerFix`)

**Test Documentation Requirements:**

All test functions MUST include comprehensive documentation above the function declaration:

1. **Documentation Structure**:
   - Start with a clear one-line summary of what the test validates
   - List specific test cases and scenarios covered
   - Explain security implications and CWE references where applicable
   - Include usage notes (e.g., "run with -race detector")
   - Explain why the test is critical for the system

2. **Documentation Format**:
   ```go
   // TestFunctionName_Scenario tests [brief description]:
   // - Specific test case 1
   // - Specific test case 2
   // - Edge case or boundary condition
   //
   // Additional context about security implications, CWE references,
   // or why this test is critical for production systems.
   //
   // Usage notes: Run with go test -race for concurrency tests.
   func TestFunctionName_Scenario(t *testing.T) {
       // Test implementation
   }
   ```

3. **Required Documentation Elements**:
   - **What**: Clear description of functionality being tested
   - **How**: List of specific test cases and scenarios
   - **Why**: Security implications, CWE references, or business criticality
   - **Usage**: Special requirements (race detector, external services, etc.)

4. **Security Test Documentation**:
   - MUST include CWE reference numbers (e.g., CWE-22, CWE-158)
   - MUST explain the attack vector being tested
   - MUST explain the mitigation being validated
   - MUST note if test demonstrates vulnerability safely

5. **Documentation Benefits**:
   - Makes test purpose immediately clear to reviewers
   - Helps maintainers understand security implications
   - Provides context for why tests exist
   - Documents attack vectors and mitigations
   - Serves as inline security documentation

6. **When to Update Documentation**:
   - When adding new test cases to existing tests
   - When fixing bugs that tests should have caught
   - When security vulnerabilities are discovered
   - When test behavior or scope changes

## Integration Notes

**As a Dependency:**
This library is meant to be vendored or imported by other Open Horizon projects. To use it:

```bash
go get github.com/open-horizon/edge-utilities
```

**Common Integration Pattern:**
1. Import the appropriate package (`logger/log` for simple logging, `logger` for advanced control)
2. Initialize with `Parameters` struct
3. Use throughout application via the singleton or instance
4. Call `Stop()` on shutdown to flush logs and close files

**File System Requirements:**
- Write permissions to log directory
- Sufficient disk space for log files and compressed archives
- Log rotation requires ability to rename and delete files

**Backward Compatibility:**
- As a library, maintain strict backward compatibility
- Deprecate features before removing them
- Document breaking changes clearly in release notes
- Consider semantic versioning for releases

## Package Dependencies

- `github.com/golang/glog`: Google's logging library (used as optional destination)
- Standard library packages: `compress/gzip`, `log/syslog`, `reflect`, `syscall`, etc.

## Important Notes for AI Agents

### Efficient Tool Usage and Batching Strategies

**Principles of Efficient Tool Usage:**

When working with code changes, especially large-scale modifications, agents should optimize tool usage to minimize context usage, reduce costs, and improve performance:

1. **Batch Similar Operations**: Group similar changes together in a single tool call when possible
2. **Use the Right Tool**: Choose tools that can handle multiple operations efficiently
3. **Minimize Round Trips**: Reduce the number of tool calls by combining operations
4. **Plan Before Executing**: Analyze the scope of work before starting to identify batching opportunities

**Batching Strategies for Large Changes:**

**When to Use Batching:**

Batching is most effective for:
- Adding documentation to multiple functions in the same file
- Applying the same pattern across multiple files
- Making consistent formatting or style changes
- Adding similar test cases across multiple test files
- Updating configuration in multiple locations

**Tools That Support Batching:**

1. **apply_diff**: Can apply multiple search/replace blocks in a single call
   - Each block operates independently within the same file
   - Ideal for making multiple small, precise changes to one file
   - Example: Adding documentation to 5-10 functions in the same file
   - Each search/replace block should be separated by a blank line

2. **replace_regex**: Can apply multiple regex patterns in a single call
   - Supports multiple pattern/replacement pairs in one diff
   - Ideal for consistent pattern-based changes
   - Example: Updating import statements or renaming patterns across a file

**Batching Best Practices:**

1. **Group by File**: Batch all changes for a single file into one tool call
   - Good: One apply_diff call with 10 search/replace blocks for 10 functions
   - Bad: 10 separate apply_diff calls for the same file

2. **Limit Batch Size**: Keep batches manageable (5-15 operations per call)
   - Too small: Wastes tool calls and increases context usage
   - Too large: Harder to debug if one operation fails
   - Sweet spot: 8-12 related operations per batch

3. **Verify Before Batching**: Read the file first to ensure all targets exist
   - Use read_file to examine the file structure
   - Confirm all search strings will match exactly
   - Plan the batch based on actual file content

4. **Handle Failures Gracefully**: If a batch fails, break it into smaller batches
   - Identify which operation failed
   - Complete successful operations first
   - Retry failed operations individually or in smaller groups

**When NOT to Batch:**

Avoid batching when:
- Operations depend on each other (one must complete before the next)
- Changes span multiple files (use separate tool calls per file)
- Operations are unrelated or serve different purposes
- Debugging is needed (smaller operations are easier to troubleshoot)
- The batch would exceed 15-20 operations (too complex)

**Performance Optimization Tips:**

1. **Read Once, Write Once**: Read a file once, plan all changes, apply in one batch
2. **Use Appropriate Tools**: 
   - `apply_diff` for precise, literal replacements
   - `replace_regex` for pattern-based changes
   - `write_to_file` only for complete file rewrites
3. **Minimize File Reads**: Cache file content mentally during planning phase
4. **Parallel Planning**: While waiting for tool responses, plan the next batch
5. **Progressive Refinement**: Start with a small batch to verify approach, then scale up

### General Guidelines

1. **This is a library, not an application** - there are no build or run commands
2. **Focus on understanding the API surface** when making changes
3. **Maintain backward compatibility** as this is used by multiple Open Horizon projects
4. **The logger's file rotation logic is complex** - be careful when modifying `checkFiles()`
5. **Properties loading uses reflection** - changes must preserve type safety
6. **When adding features**, consider whether they belong in the core package or convenience wrappers
7. **Test with race detector** when modifying concurrent code
8. **NEVER commit code without tests** - especially for bug fixes and security patches
9. **Security fixes require explicit test coverage** demonstrating both the vulnerability and the fix

## AI Agent Limitations

**DO NOT attempt to provide or calculate:**

1. **Timelines and Schedules**: Do not estimate project timelines, development schedules, or completion dates
2. **Performance Metrics**: Do not calculate or estimate:
   - Latency targets or measurements
   - Throughput rates or capacity
   - Hit rates or cache efficiency
   - Response times or processing speeds
   - Resource utilization percentages
3. **Risk Evaluations**: Do not assess or quantify:
   - Security risk levels or scores
   - Business impact assessments
   - Probability of failures or incidents
   - Cost-benefit analyses
4. **Quantitative Predictions**: Avoid making numerical predictions about system behavior, user adoption, or operational metrics

**Why these limitations exist:**
- AI agents cannot accurately perform these calculations without real-world data
- Estimates and predictions require domain expertise, historical data, and context that agents lack
- Providing inaccurate numbers creates false confidence and can lead to poor decisions
- These assessments require human judgment, stakeholder input, and organizational context

**What agents CAN do:**
- Identify areas where performance testing is needed
- Suggest monitoring and measurement approaches
- Recommend best practices for performance optimization
- Point to relevant documentation or tools for proper assessment
- Implement code changes based on clear, specific requirements