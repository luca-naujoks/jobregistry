# JobRegistry

A Go module for managing and registering jobs.

## Overview

`JobRegistry` provides a structured way to register, manage, and execute jobs within your Go applications.

## Getting Started

To use this module in your Go project:

```bash
go get github.com/luca-naujoks/jobRegistry
```

## Usage

Import the package into your Go code:

```text
import "github.com/luca-naujoks/jobRegistry"
```

_(Note: Depending on how your internal package structure is exposed, you may need to adjust the import path
if `jobRegistry` is not intended to be a direct import.)_

## Development

- **Language:** Go 1.26.2
- **Module:** `github.com/luca-naujoks/jobRegistry`

### Running Tests

To verify the installation and run the existing tests, execute:

```bash
go test .
```

## Project Structure

- `jobRegistry/`: Core logic containing the registry implementation and job types.
- `go.mod`: Module definition.
