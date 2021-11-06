package vote.pebble.common;

import java.security.NoSuchAlgorithmException;
import java.util.Arrays;
import java.security.MessageDigest;

public final class HashValue {
    public static final int SIZE_BYTES = 32;
    public static final HashValue ZERO = new HashValue(new byte[32]);

    public final byte[] bytes;

    public HashValue(byte[] bytes) {
        assert bytes.length == SIZE_BYTES;
        this.bytes = bytes;
    }

    public static HashValue sha256(byte[] input) {
        MessageDigest hasher;
        try {
            hasher = MessageDigest.getInstance("SHA-256");
        } catch (NoSuchAlgorithmException e) {
            throw new RuntimeException(e);
        }
        return new HashValue(hasher.digest(input));
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
