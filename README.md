# commandments

A Go library for building type-safe CLI applications with automatic flag generation from struct tags using Cobra and Viper.

## Installation

```bash
go get github.com/pindamonhangaba/commandments
```

## Basic Usage

Create a simple command with automatic flag generation:

```go
package main

import (
    "fmt"
    "github.com/pindamonhangaba/commandments"
)

type ServerConfig struct {
    Port int    `flag:"port,Server port number"`
    Host string `flag:"host,Server host address"`
}

func main() {
    cmd := commandments.MustCMD("serve", 
        commandments.WithConfig(func(config ServerConfig) error {
            fmt.Printf("Starting server on %s:%d\n", config.Host, config.Port)
            return nil
        }),
        commandments.WithDefaultConfig(ServerConfig{
            Port: 8080,
            Host: "localhost",
        }),
    )
    
    cmd.Execute()
}
```

This automatically generates `--port` and `--host` flags with the specified descriptions and default values.

## Advanced Usage

### Multiple Commands with Shared Configuration

```go
package main

import (
    "fmt"
    "github.com/pindamonhangaba/commandments"
    "github.com/spf13/cobra"
)

type AppConfig struct {
    Debug   bool   `flag:"debug,Enable debug mode"`
    LogFile string `flag:"log-file,Path to log file"`
}

func main() {
    rootCmd := &cobra.Command{
        Use:   "myapp",
        Short: "My application",
    }

    // Serve command
    serveCmd := commandments.MustCMD("serve",
        commandments.WithConfig(func(config AppConfig) error {
            fmt.Printf("Starting server (debug: %v, log: %s)\n", config.Debug, config.LogFile)
            return nil
        }),
        commandments.WithDefaults(commandments.CMDDefaults{
            ShortDescription: "Start the server",
            Description:      "Start the HTTP server with the specified configuration",
            EnvPrefix:        "MYAPP",
        }),
    )

    // Worker command  
    workerCmd := commandments.MustCMD("worker",
        commandments.WithConfig(func(config AppConfig) error {
            fmt.Printf("Starting worker (debug: %v, log: %s)\n", config.Debug, config.LogFile)
            return nil
        }),
        commandments.WithDefaults[AppConfig](commandments.CMDDefaults{
            ShortDescription: "Start the worker",
            Description:      "Start the background worker process",
            EnvPrefix:        "MYAPP",
        }),
    )

    rootCmd.AddCommand(serveCmd, workerCmd)
    rootCmd.Execute()
}
```

## Project Structure

For larger applications, organize your commands in a structured way:

```
cmd/
├── project/
│   ├── root.go
│   └── serve.go
└── main.go
```

### cmd/project/root.go

```go
package project

import (
    "github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
    Use:   "project",
    Short: "Project management tool",
    Long:  "A CLI tool for managing projects with various commands",
}

func init() {
    RootCmd.AddCommand(ServeCmd)
}
```

### cmd/project/serve.go

```go
package project

import (
    "fmt"
    "github.com/pindamonhangaba/commandments"
    "github.com/spf13/cobra"
)

type ServeConfig struct {
    Port        int      `flag:"port,HTTP server port"`
    Host        string   `flag:"host,HTTP server host"`
    TLS         bool     `flag:"tls,Enable TLS"`
    CertFile    string   `flag:"cert-file,TLS certificate file"`
    KeyFile     string   `flag:"key-file,TLS private key file"`
    AllowedIPs  []string `flag:"allowed-ips,Allowed IP addresses"`
}

var ServeCmd = commandments.MustCMD("serve",
    commandments.WithConfig(func(config ServeConfig) error {
        fmt.Printf("Starting server on %s:%d\n", config.Host, config.Port)
        if config.TLS {
            fmt.Printf("TLS enabled (cert: %s, key: %s)\n", config.CertFile, config.KeyFile)
        }
        if len(config.AllowedIPs) > 0 {
            fmt.Printf("Allowed IPs: %v\n", config.AllowedIPs)
        }
        // Your server logic here
        return nil
    }),
    commandments.WithDefaultConfig(ServeConfig{
        Port: 8080,
        Host: "0.0.0.0",
        TLS:  false,
    }),
    commandments.WithDefaults[ServeConfig](commandments.CMDDefaults{
        ShortDescription:      "Start the HTTP server",
        Description:           "Start the HTTP server with configurable options",
        DefaultConfigFilename: "server",
        EnvPrefix:             "PROJECT",
    }),
)
```

### cmd/main.go

```go
package main

import (
    "cmd/project"
    "fmt"
    "os"
)

func main() {
    if err := project.RootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

## Configuration Options

### Flag Types

The library supports various flag types through struct tags:

```go
type Config struct {
    // String flags
    Name        string   `flag:"name,Your name"`
    Description string   `flag:"desc,Description"`
    
    // Numeric flags  
    Port    int     `flag:"port,Port number"`
    Timeout float64 `flag:"timeout,Timeout in seconds"`
    
    // Boolean flags
    Debug   bool `flag:"debug,Enable debug mode"`
    Verbose bool `flag:"verbose,Verbose output"`
    
    // Slice flags
    Tags     []string `flag:"tags,List of tags"`
    Ports    []int    `flag:"ports,List of ports"`
}
```

### Environment Variables

Flags are automatically bound to environment variables using the specified prefix:

```go
commandments.WithDefaults(commandments.CMDDefaults{
    EnvPrefix: "MYAPP", // --port becomes MYAPP_PORT
})
```

### Configuration Files

The library supports various configuration file formats (JSON, YAML, TOML, etc.):

```go
commandments.WithDefaults[ConfigType](commandments.CMDDefaults{
    DefaultConfigFilename: "myapp", // Looks for myapp.yaml, myapp.json, etc.
})
```

## Features

- **Type Safety**: Configuration is strongly typed using Go generics
- **Automatic Flag Generation**: Flags are generated from struct field tags
- **Multiple Data Types**: Support for strings, numbers, booleans, and slices
- **Environment Variables**: Automatic binding with configurable prefixes
- **Configuration Files**: Support for JSON, YAML, TOML and other formats
- **Default Values**: Easy specification of default configuration
- **Cobra Integration**: Built on top of the popular Cobra CLI framework
- **Viper Integration**: Uses Viper for configuration management

## API Reference

### Core Functions

- `NewCMD[T](name string, configs ...cmdArg[T]) (*cobra.Command, error)` - Create a new command
- `MustCMD[T](name string, configs ...cmdArg[T]) *cobra.Command` - Create a command, panic on error

### Configuration Options

- `WithConfig[T](func(T) error) cmdArg[T]` - Set the configuration handler function
- `WithDefaultConfig[T](T) cmdArg[T]` - Set default configuration values
- `WithDefaults[T](CMDDefaults) cmdArg[T]` - Set command metadata and options

### CMDDefaults Fields

- `ShortDescription string` - Short command description
- `Description string` - Long command description  
- `DefaultConfigFilename string` - Configuration file name (without extension)
- `EnvPrefix string` - Environment variable prefix