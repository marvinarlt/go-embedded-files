# Example: Embedded files in GO binary

This is an basic example of how to embed static files in a GO binary. The actual magic happens via the native `//go:embed directory/*` directive.

## Getting started

1. Build the binary
    ```
    go build -o bin/main main.go 
    ``` 

2. Execute the binary
    ```
    bin/main
    ```

3. Visit http://localhost:1337
