package vote.pebble.common;

import java.util.Arrays;

public class ByteString implements Comparable<ByteString> {
    public final byte[] bytes;

    public ByteString(byte[] bytes) {
        assert bytes != null;
        this.bytes = bytes;
    }

    public ByteString(byte[] bytes, int off, int length) {
        this.bytes = Arrays.copyOfRange(bytes, off, off + length);
    }

    @Override
    public int compareTo(ByteString other) {
        return Arrays.compare(bytes, other.bytes);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o instanceof ByteString)
            return Arrays.equals(((ByteString) o).bytes, bytes);
        return false;
    }

    @Override
    public int hashCode() {
        return Arrays.hashCode(bytes);
    }

    @Override
    public String toString() {
        return Hex.encodeHexString(bytes);
    }
}
