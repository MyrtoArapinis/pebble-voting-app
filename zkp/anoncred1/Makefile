LIB_SOURCES := ffi.go main.go zkp.go
C_SOURCES := $(wildcard jni/*.c)
INCLUDES := -I${JAVA_HOME}/include -I${JAVA_HOME}/include/linux

all: anoncred1 libanoncred1-jni.so

clean:
	rm -f anoncred1 *.a *.so jni/*.o

setup: anoncred1
	./anoncred1 setup 16 anoncred1-params.bin

anoncred1: $(wildcard *.go)
	go build

libanoncred1.a: $(LIB_SOURCES)
	go build -buildmode=c-archive -o $@ $^

libanoncred1-jni.so: $(C_SOURCES) libanoncred1.a
	cc $^ -o $@ -O2 -shared -fPIC $(INCLUDES)

jni/vote_pebble_zkp_AnonCred1.c: jni/vote_pebble_zkp_AnonCred1.h
