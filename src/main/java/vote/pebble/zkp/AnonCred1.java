package vote.pebble.zkp;

import vote.pebble.common.HashValue;
import vote.pebble.common.Util;

import java.io.IOException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;

public class AnonCred1 implements CredentialSystem {
    private static final byte[] ZKP_PARAMS_BYTES;
    private static final int ZKP_DEPTH;

    static {
        System.load(AnonCred1.class.getClassLoader().getResource("libanoncred1-jni.so").getPath());
        try (var stream = AnonCred1.class.getClassLoader().getResourceAsStream("anoncred1-params.bin")) {
            ZKP_PARAMS_BYTES = stream.readAllBytes();
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        ZKP_DEPTH = ZKP_PARAMS_BYTES[0];
    }

    private static native int jniGenerateCredential(byte[] buffer);

    private static native int jniHashMerkleTree(byte[] root, byte[] credentials, int depth);

    private static native int jniProve(byte[] out, byte[] paramsBytes, byte[] messageHash,
                                       byte[] serialNo, byte[] secret, int idx, byte[] credentials);

    private static native int jniVerify(byte[] paramsBytes, byte[] messageHash, byte[] serialNo,
                                        byte[] signature, byte[] merkleRoot);

    private static final class PublicCredentialImpl implements PublicCredential, Comparable<PublicCredentialImpl> {
        final byte[] bytes;

        PublicCredentialImpl(byte[] bytes) {
            this.bytes = bytes;
        }

        @Override
        public byte[] toBytes() {
            return bytes;
        }

        @Override
        public int compareTo(PublicCredentialImpl other) {
            return Arrays.compare(bytes, other.bytes);
        }

        @Override
        public boolean equals(Object o) {
            if (this == o) return true;
            if (o == null || getClass() != o.getClass()) return false;
            return Arrays.equals(bytes, ((PublicCredentialImpl) o).bytes);
        }

        @Override
        public int hashCode() {
            return Arrays.hashCode(bytes);
        }
    }

    private static final class SecretCredentialImpl implements SecretCredential {
        final PublicCredentialImpl publicCredential;
        final byte[] serialNo, secret;

        SecretCredentialImpl(byte[] publicCredential, byte[] serialNo, byte[] secret) {
            this.publicCredential = new PublicCredentialImpl(publicCredential);
            this.serialNo = serialNo;
            this.secret = secret;
        }

        @Override
        public byte[] toBytes() {
            return Util.concat(publicCredential.bytes, serialNo, secret);
        }

        @Override
        public PublicCredential getPublicCredential() {
            return publicCredential;
        }

        @Override
        public byte[] getSerialNumber() {
            return serialNo;
        }
    }

    private static final class CredentialSetImpl implements CredentialSet {
        final byte[] merkleRoot;
        final byte[][] credentials;
        final byte[] credentialsConcat;

        CredentialSetImpl(Iterable<PublicCredential> publicCredentials) throws CredentialException {
            var list = new ArrayList<PublicCredentialImpl>();
            for (var item : publicCredentials)
                list.add((PublicCredentialImpl) item);
            Collections.sort(list);
            for (int i = list.size() - 1; i > 0; i--) {
                if (list.get(i).equals(list.get(i - 1)))
                    list.remove(i);
            }
            credentials = new byte[list.size()][];
            for (int i = 0; i < credentials.length; i++)
                credentials[i] = list.get(i).bytes;
            credentialsConcat = Util.concat(credentials);
            merkleRoot = new byte[32];
            int ret = jniHashMerkleTree(merkleRoot, credentialsConcat, ZKP_DEPTH);
            if (ret != 0)
                throw new CredentialException("Error while hashing Merkle tree");
        }

        @Override
        public byte[] sign(SecretCredential secretCredential, byte[] message) throws CredentialException {
            var credential = ((SecretCredentialImpl) secretCredential);
            var pubBytes = credential.publicCredential.bytes;
            int idx = 0;
            while (!Arrays.equals(pubBytes, credentials[idx]))
                idx++;
            var hash = HashValue.hash(message).bytes;
            var outBuffer = new byte[1024];
            int ret = jniProve(outBuffer, ZKP_PARAMS_BYTES, hash, credential.serialNo, credential.secret, idx, credentialsConcat);
            if (ret < 0)
                throw new CredentialException("Error while proving");
            return Arrays.copyOf(outBuffer, ret);
        }

        @Override
        public boolean verify(byte[] serialNo, byte[] signature, byte[] message) {
            var hash = HashValue.hash(message).bytes;
            int ret = jniVerify(ZKP_PARAMS_BYTES, hash, serialNo, signature, merkleRoot);
            return ret == 0;
        }
    }

    @Override
    public SecretCredential generateSecretCredential() throws CredentialException {
        var buffer = new byte[32 * 3];
        int ret = jniGenerateCredential(buffer);
        if (ret != 0)
            throw new CredentialException("Error while generating credential");
        return secretCredentialFromBytes(buffer);
    }

    @Override
    public SecretCredential secretCredentialFromBytes(byte[] bytes) throws CredentialException {
        if (bytes.length != 96)
            throw new CredentialException("Private credential must be 96 bytes");
        return new SecretCredentialImpl(
                Arrays.copyOf(bytes,32),
                Arrays.copyOfRange(bytes,32, 64),
                Arrays.copyOfRange(bytes, 64, 96));
    }

    @Override
    public PublicCredential publicCredentialFromBytes(byte[] bytes) throws CredentialException {
        if (bytes.length != 32)
            throw new CredentialException("Public credential must be 32 bytes");
        return new PublicCredentialImpl(bytes);
    }

    @Override
    public CredentialSet makeCredentialSet(Iterable<PublicCredential> publicCredentials) throws CredentialException {
        return new CredentialSetImpl(publicCredentials);
    }
}
