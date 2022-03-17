all: src/main/resources/libanoncred1-jni.so
	./gradlew jar

check: src/main/resources/libanoncred1-jni.so
	./gradlew build

clean:
	./gradlew clean
	make -C zkp/anoncred1/ clean

src/main/resources/libanoncred1-jni.so: zkp/anoncred1/libanoncred1-jni.so
	cp $< $@

zkp/anoncred1/libanoncred1-jni.so: $(wildcard zkp/anoncred1/*.go)
	make -C zkp/anoncred1/
