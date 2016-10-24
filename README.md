# go-envparse

Environment variable parsing for Go. Parses ENV files and calls `os.Setenv(...)`. We run our applications in docker and all of our configuration values are set to environment variables. When running in development mode we need a way to set environment variable while running `make run`. Though there are ways to do this in a Makefile, none of them were desireable. This package provides us with a way to provide our applications with `-env-file dev.env.list` and have the application parse that file and set the environment variables during `init()`. The rest of the application remains the same; reading and parsing environment variables. 

## Installation

```bash
go get -u github.com/blockloop/go-envparse
# or
govendor fetch github.com/blockloop/go-envparse
```

## Usage

See godoc: https://godoc.org/github.com/blockloop/go-envparse
