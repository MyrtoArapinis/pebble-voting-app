#include "vote_pebble_zkp_AnonCred1.h"

#include <stddef.h>

typedef struct { void *data; size_t len; size_t cap; } GoSlice;

int FfiGenerateCredential(GoSlice);
int FfiHashMerkleTree(GoSlice destRoot, GoSlice credentialsConcat, size_t depth);
int FfiProve(GoSlice out, GoSlice paramsBytes, GoSlice messageHash, GoSlice serialNo, GoSlice secret, size_t idx, GoSlice credentialsConcat);
int FfiVerify(GoSlice paramsBytes, GoSlice messageHash, GoSlice serialNo, GoSlice signature, GoSlice merkleRoot);

static GoSlice slice(void *data, size_t len) {
	GoSlice s;
	s.data = data;
	s.len = len;
	s.cap = len;
	return s;
}

static void releaseArray(JNIEnv *env, jarray array, jbyte *elems) {
	(*env)->ReleaseByteArrayElements(env, array, elems, 0);
}

static jsize arrayLength(JNIEnv *env, jarray array) {
	return (*env)->GetArrayLength(env, array);
}

static jbyte *arrayBytes(JNIEnv *env, jarray array) {
	return (*env)->GetByteArrayElements(env, array, NULL);
}

static void getBytes32(JNIEnv *env, jarray array, jbyte *buf) {
	(*env)->GetByteArrayRegion(env, array, 0, 32, buf);
}

JNIEXPORT jint JNICALL Java_vote_pebble_zkp_AnonCred1_jniGenerateCredential
		(JNIEnv *env, jclass cls, jbyteArray array) {
	jbyte bytes[3 * 32];
	if (arrayLength(env, array) != sizeof(bytes))
		return -2;
	int ret = FfiGenerateCredential(slice(bytes, sizeof(bytes)));
	(*env)->SetByteArrayRegion(env, array, 0, sizeof(bytes), bytes);
	return ret;
}

JNIEXPORT jint JNICALL Java_vote_pebble_zkp_AnonCred1_jniHashMerkleTree
		(JNIEnv *env, jclass cls, jbyteArray arrRoot, jbyteArray arrCredentials, jint depth) {
	jsize len = arrayLength(env, arrCredentials);
	if (len < 0 || len % 32 != 0 || arrayLength(env, arrRoot) != 32)
		return -2;
	jbyte *credentials = arrayBytes(env, arrCredentials);
	jbyte root[32];
	int ret = FfiHashMerkleTree(slice(root, 32), slice(credentials, len), depth);
	releaseArray(env, arrCredentials, credentials);
	(*env)->SetByteArrayRegion(env, arrRoot, 0, 32, root);
	return ret;
}

JNIEXPORT jint JNICALL Java_vote_pebble_zkp_AnonCred1_jniProve
		(JNIEnv *env, jclass cls, jbyteArray arrOut, jbyteArray arrParamsBytes,
		jbyteArray arrMessageHash, jbyteArray arrSerialNo, jbyteArray arrSecret,
		jint idx, jbyteArray arrCredentials) {
	if (arrayLength(env, arrMessageHash) != 32
			|| arrayLength(env, arrSerialNo) != 32
			|| arrayLength(env, arrSecret) != 32)
		return -2;
	jsize credLen = arrayLength(env, arrCredentials);
	if (credLen < 0 || credLen % 32 != 0)
		return -2;
	jbyte messageHash[32], serialNo[32], secret[32];
	getBytes32(env, arrMessageHash, messageHash);
	getBytes32(env, arrSerialNo, serialNo);
	getBytes32(env, arrSecret, secret);
	jsize outLen = arrayLength(env, arrOut);
	jsize paramsLen = arrayLength(env, arrParamsBytes);
	jbyte *paramsBytes = arrayBytes(env, arrParamsBytes);
	jbyte *credentials = arrayBytes(env, arrCredentials);
	jbyte *out = arrayBytes(env, arrOut);
	int ret = FfiProve(
		slice(out, outLen),
		slice(paramsBytes, paramsLen),
		slice(messageHash, 32),
		slice(serialNo, 32),
		slice(secret, 32),
		idx,
		slice(credentials, credLen));
	releaseArray(env, arrParamsBytes, paramsBytes);
	releaseArray(env, arrCredentials, credentials);
	releaseArray(env, arrOut, out);
	return ret;
}

JNIEXPORT jint JNICALL Java_vote_pebble_zkp_AnonCred1_jniVerify
		(JNIEnv *env, jclass cls, jbyteArray arrParamsBytes,
		jbyteArray arrMessageHash, jbyteArray arrSerialNo,
		jbyteArray arrSignature, jbyteArray arrMerkleRoot) {
	if (arrayLength(env, arrMessageHash) != 32
			|| arrayLength(env, arrSerialNo) != 32
			|| arrayLength(env, arrMerkleRoot) != 32)
		return -2;
	jbyte messageHash[32], serialNo[32], merkleRoot[32];
	getBytes32(env, arrMessageHash, messageHash);
	getBytes32(env, arrSerialNo, serialNo);
	getBytes32(env, arrMerkleRoot, merkleRoot);
	jsize paramsLen = arrayLength(env, arrParamsBytes);
	jsize sigLen = arrayLength(env, arrSignature);
	jbyte *paramsBytes = arrayBytes(env, arrParamsBytes);
	jbyte *signature = arrayBytes(env, arrSignature);
	int ret = FfiVerify(
		slice(paramsBytes, paramsLen),
		slice(messageHash, 32),
		slice(serialNo, 32),
		slice(signature, sigLen),
		slice(merkleRoot, 32));
	releaseArray(env, arrParamsBytes, paramsBytes);
	releaseArray(env, arrSignature, signature);
	return ret;
}
