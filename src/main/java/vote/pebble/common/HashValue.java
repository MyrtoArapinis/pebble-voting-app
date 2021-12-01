package vote.pebble.common;

import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

public final class HashValue extends ByteString {
    public static final int SIZE_BYTES = 32;
    public static final HashValue ZERO = new HashValue(new byte[32]);

    private static final MessageDigest md = createMessageDigest();

    public HashValue(byte[] bytes) {
        super(bytes);
        assert bytes.length == SIZE_BYTES;
    }

    public HashValue(byte[] bytes, int off) {
        super(bytes, off, SIZE_BYTES);
    }

    public HashValue(ByteString input) {
        super(input.bytes);
        assert bytes.length == SIZE_BYTES;
    }

    public static MessageDigest createMessageDigest() {
        try {
            return MessageDigest.getInstance("SHA-256");
        } catch (NoSuchAlgorithmException e) {
            throw new RuntimeException(e);
        }
    }

    public static byte[] digest(byte[] message) {
        synchronized (md) {
            md.reset();
            return md.digest(message);
        }
    }

    public static HashValue hash(byte[] message) {
        return new HashValue(digest(message));
    }
}
