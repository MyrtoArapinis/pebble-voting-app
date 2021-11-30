package vote.pebble.common;

import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.Arrays;

public final class HashValue {
    public static final int SIZE_BYTES = 32;
    public static final HashValue ZERO = new HashValue(new byte[32]);

    private static final MessageDigest md = createMessageDigest();

    public final byte[] bytes;

    public HashValue(byte[] bytes) {
        assert bytes.length == SIZE_BYTES;
        this.bytes = bytes;
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

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o instanceof HashValue)
            return Arrays.equals(((HashValue) o).bytes, bytes);
        return false;
    }

    @Override
    public int hashCode() {
        return Arrays.hashCode(bytes);
    }
}
