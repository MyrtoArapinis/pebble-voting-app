The Pebble software consists of three components:

1. Client command line app (Java)
2. Anonymous Credentials library (Go & C)
3. Server (Go)

The first two can be built using `make`. The resulting `jar` will be in the `build/libs` directory.

The server can be built with `go build` in the `server` directory.
