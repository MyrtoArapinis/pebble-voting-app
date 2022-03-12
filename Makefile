all:
	make -C zkp/anoncred1/
	cp zkp/anoncred1/libanoncred1-jni.so src/main/resources/libanoncred1-jni.so
	./gradlew build
