# Building Pebble

The Pebble software consists of three components:

1. Client command line app (Java)
2. Anonymous Credentials library (Go & C)
3. Server (Go)

## Requirements

* Go 1.16 or higher
* Java Development Kit 11 or higher

Java dependencies are managed by Gradle.

## Instructions

Functions within the Anonymous Credentials library are called from the Java app using JNI. To compile the library ensure that the `$JAVA_HOME` environment variable is set, e.g. by running `echo $JAVA_HOME`. An example setting would be `/usr/lib/jvm/java-11-openjdk-amd64`.  If it is not set and you have identified your Java installation directory, you can set it with e.g. `export JAVA_HOME=/usr/lib/jvm/java-11-openjdk-amd64`.

The command line app and its library can be built using `make`. The resulting `jar` will be in the `build/libs` directory.

The server can be built with `go build` in the `server` directory.
