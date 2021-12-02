package vote.pebble.voting.structs;

import vote.pebble.common.ParseException;

import cafe.cryptography.curve25519.InvalidEncodingException;
import cafe.cryptography.ed25519.Ed25519PrivateKey;
import cafe.cryptography.ed25519.Ed25519PublicKey;
import cafe.cryptography.ed25519.Ed25519Signature;

import java.nio.BufferUnderflowException;
import java.nio.ByteBuffer;
import java.util.Arrays;

public final class CredentialMessage {
    private static final String STRUCT = "CredentialMessage";

    public final Ed25519PublicKey publicKey;
    public final Ed25519Signature signature;
    public final byte[] credential;

    public CredentialMessage(Ed25519PublicKey publicKey, Ed25519Signature signature, byte[] credential) {
        this.publicKey = publicKey;
        this.signature = signature;
        this.credential = credential;
    }

    public static CredentialMessage sign(Ed25519PrivateKey privateKey, byte[] credential) {
        var publicKey = privateKey.derivePublic();
        return new CredentialMessage(
                publicKey,
                privateKey.expand().sign(credential, publicKey),
                credential);
    }

    public boolean verify() {
        return publicKey.verify(credential, signature);
    }

    public static CredentialMessage fromBytes(byte[] bytes) throws ParseException {
        try {
            var buffer = ByteBuffer.wrap(bytes);
            var publicKey = new byte[32];
            buffer.get(publicKey);
            var signature = new byte[64];
            buffer.get(signature);
            var credential = new byte[buffer.remaining()];
            buffer.get(credential);
            return new CredentialMessage(
                    Ed25519PublicKey.fromByteArray(publicKey),
                    Ed25519Signature.fromByteArray(signature),
                    credential);
        } catch (BufferUnderflowException | InvalidEncodingException e) {
            throw new ParseException(STRUCT, e);
        }
    }

    public byte[] toBytes() {
        return ByteBuffer.allocate(96 + credential.length)
                .put(publicKey.toByteArray())
                .put(signature.toByteArray())
                .put(credential)
                .array();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        CredentialMessage that = (CredentialMessage) o;
        return that.publicKey.equals(publicKey) && that.signature.equals(signature) && Arrays.equals(credential, that.credential);
    }

    @Override
    public int hashCode() {
        return (publicKey.hashCode() * 37 + signature.hashCode()) * 37 + Arrays.hashCode(credential);
    }
}
