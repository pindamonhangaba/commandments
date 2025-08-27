# Commandments

Commandments is a Go library that simplifies creating CLI applications using Cobra and Viper with type-safe configuration binding. It provides a generic, declarative way to bind struct fields to command line flags, configuration files, and environment variables.

**Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.**

## Working Effectively

### Initial Setup and Dependencies
- Ensure Go 1.18+ is installed (library uses generics): `go version`
- Install all dependencies: `go mod download` -- takes ~16 seconds. NEVER CANCEL. Set timeout to 30+ seconds.
- Verify setup: `go mod tidy && go mod verify`

### Building the Library
- Build all packages: `go build ./...` -- takes ~24 seconds. NEVER CANCEL. Set timeout to 60+ seconds.
- The library builds into packages, not a standalone executable (it's a library, not an application)
- All builds should succeed without errors or warnings

### Testing and Validation
- Run all tests: `go test -v ./...` -- takes ~4 seconds. NEVER CANCEL. Set timeout to 30+ seconds.
- Run tests with race detection: `go test -race ./...` -- takes ~1 second. NEVER CANCEL. Set timeout to 30+ seconds.
- Run benchmarks: `go test -bench=. ./...` -- takes ~1 second. NEVER CANCEL. Set timeout to 30+ seconds.
- All tests must pass before committing changes

### Code Quality and Formatting
- **ALWAYS** run `go fmt ./...` before committing -- this formats all Go files
- **ALWAYS** run `go vet ./...` before committing -- this checks for common errors
- Both commands should complete without output (indicating no issues)
- These steps are **REQUIRED** before any commit or the code will not meet project standards

### Manual Validation Scenarios
After making changes to the library, **ALWAYS** test the following scenarios:

1. **Basic Command Creation**: Create a simple CLI command with struct-based configuration
2. **Flag Generation**: Verify that struct fields with `flag` tags generate appropriate command line flags
3. **Type Safety**: Test that the generic type system works correctly with different struct types
4. **Configuration Integration**: Test that Viper integration works for config files and environment variables

Example validation test:
```bash
# Create a test file in /tmp to avoid committing it
cat > /tmp/test_commandments.go << 'EOF'
package main

import (
	"fmt"
	"github.com/pindamonhangaba/commandments"
)

type TestConfig struct {
	Port     int    `flag:"port,Port to serve on"`
	Host     string `flag:"host,Host to bind to"`
	Debug    bool   `flag:"debug,Enable debug mode"`
}

func main() {
	cmd := commandments.MustCMD("test", 
		commandments.WithDefaults[TestConfig](commandments.CMDDefaults{
			ShortDescription: "Test CLI",
			Description:      "Test CLI for validation",
		}),
		commandments.WithDefaultConfig(TestConfig{
			Port: 8080,
			Host: "localhost",
			Debug: false,
		}),
		commandments.WithConfig(func(config TestConfig) error {
			fmt.Printf("Config: Port=%d, Host=%s, Debug=%t\n", 
				config.Port, config.Host, config.Debug)
			return nil
		}),
	)
	
	// REQUIRED: Add config flag when using WithConfig()
	cmd.PersistentFlags().String("config", "", "config file path")
	
	cmd.Execute()
}
EOF

# Test the help output
go run /tmp/test_commandments.go --help

# Test with custom flags  
go run /tmp/test_commandments.go --port 9090 --host example.com --debug

# Clean up
rm /tmp/test_commandments.go
```

## Validation Requirements
- **Build Validation**: Code must build successfully with `go build ./...`
- **Test Validation**: All tests must pass with `go test ./...`
- **Format Validation**: Code must be properly formatted (`go fmt ./...` should show no changes)
- **Vet Validation**: Code must pass static analysis (`go vet ./...` should show no issues)
- **Race Validation**: No race conditions (`go test -race ./...` should pass)
- **Functional Validation**: Manual testing with example CLI creation must work

## Repository Structure

### Key Files
- `custom.go` - Main library implementation with generic command creation functions
- `config.go` - Viper integration for configuration file and environment variable support  
- `custom_test.go` - Unit tests for the library functionality
- `go.mod` - Go module definition (requires Go 1.18+)
- `go.sum` - Dependency checksums
- `README.md` - Basic project documentation (minimal)
- `LICENSE` - GPL-3.0 license

### Core Components
- **MustCMD[T]()** - Creates a Cobra command with panic on error
- **NewCMD[T]()** - Creates a Cobra command with error return
- **WithConfig[T]()** - Adds configuration handler function
- **WithDefaultConfig[T]()** - Sets default configuration values
- **WithDefaults[T]()** - Sets command metadata (description, env prefix, etc.)
- **structToFlags[T]()** - Internal function that converts struct tags to CLI flags

### Directory Structure
```
/home/runner/work/commandments/commandments/
├── .git/                 # Git repository data
├── .github/              # GitHub configuration (created for Copilot)
│   └── copilot-instructions.md
├── LICENSE               # GPL-3.0 license
├── README.md             # Basic project documentation  
├── config.go             # Viper configuration integration
├── custom.go             # Main library implementation
├── custom_test.go        # Unit tests
├── go.mod                # Go module definition
└── go.sum                # Dependency checksums
```

## Common Development Tasks

### Adding New Features
1. **ALWAYS** write tests first in `custom_test.go`
2. Implement the feature in `custom.go` or `config.go`
3. Run the complete validation workflow:
   ```bash
   go fmt ./...
   go vet ./...
   go build ./...
   go test -v ./...
   go test -race ./...
   ```
4. Test manually with example CLI creation
5. Update documentation if needed

### Modifying Existing Functions
1. Check existing tests in `custom_test.go` to understand expected behavior
2. Make minimal changes to preserve existing functionality
3. Add new tests for any new behavior
4. Run full validation workflow
5. Test backwards compatibility with existing usage patterns

### Working with Generics
- This library extensively uses Go generics (`[T any]`)
- Always specify type parameters when calling generic functions: `WithDefaults[MyType](...)`
- The type `T` represents the configuration struct type
- Generic functions are defined in `custom.go`

### Understanding the Flag System
- Struct fields are converted to CLI flags using `flag` tags
- Tag format: `flag:"flag-name,Description of the flag"`
- Supported types: `string`, `int`, `bool`, `[]string`, `[]int`
- Default values come from struct field values or `WithDefaultConfig()`

### Debugging Common Issues
- **"cannot infer T" errors**: Add explicit type parameters `[YourType]`
- **"flag accessed but not defined: config" errors**: When using `WithConfig()`, manually add the config flag: `cmd.PersistentFlags().String("config", "", "config file path")`
- **Build failures**: Ensure Go 1.18+ and run `go mod tidy`
- **Test failures**: Check that struct tags match expected format

### Critical Configuration Note
**IMPORTANT**: When using `WithConfig()`, you must manually add the config flag to your command:
```go
cmd.PersistentFlags().String("config", "", "config file path")
```
This is required because the library expects this flag but does not automatically add it.

## Time Expectations and Timeouts
- **Module download**: ~16 seconds (first time), <1 second (cached) (use 30+ second timeout)
- **Build process**: ~24 seconds (first time), <1 second (cached) (use 60+ second timeout)  
- **Test execution**: ~4 seconds (first time), <1 second (cached) (use 30+ second timeout)
- **Race detection**: ~1 second (use 30+ second timeout)
- **Benchmarks**: <1 second (use 30+ second timeout)
- **Formatting**: <1 second (use 10+ second timeout)
- **Vetting**: <1 second (use 10+ second timeout)

**CRITICAL**: NEVER CANCEL any build or test command. Builds may take up to 60 seconds on slower systems. Always wait for completion.

## Integration Notes
- **Cobra Integration**: Commands are fully compatible with Cobra's command system
- **Viper Integration**: Automatic binding of flags to config files and environment variables
- **Config Files**: Supports any format Viper supports (JSON, YAML, TOML, etc.)
- **Environment Variables**: Auto-prefixed based on `EnvPrefix` setting
- **Type Safety**: Full compile-time type checking for configuration structs

## Example Usage Patterns
```go
// Basic usage
type Config struct {
    Port int `flag:"port,Server port"`
    Host string `flag:"host,Server host"`
}

cmd := commandments.MustCMD("myapp",
    commandments.WithDefaults[Config](commandments.CMDDefaults{
        ShortDescription: "My application",
        EnvPrefix: "MYAPP",
    }),
    commandments.WithDefaultConfig(Config{Port: 8080, Host: "localhost"}),
    commandments.WithConfig(func(cfg Config) error {
        // Your application logic here
        return nil
    }),
)

// CRITICAL: Add config flag when using WithConfig()
cmd.PersistentFlags().String("config", "", "config file path")
```

**Remember**: Always validate your changes work by creating and testing a real CLI example before committing.