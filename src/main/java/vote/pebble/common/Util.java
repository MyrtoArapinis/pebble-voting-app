package vote.pebble.common;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;

import java.math.BigInteger;

public final class Util {
    public static final Gson GSON = new GsonBuilder().disableHtmlEscaping().create();

    public static <T> int indexOf(T[] a, Object o) {
        int i = 0;
        while (!a[i].equals(o))
            i++;
        return i;
    }

    public static byte[] concat(Iterable<byte[]> arrays) {
        int length = 0;
        for (var arr : arrays)
            length += arr.length;
        var result = new byte[length];
        int i = 0;
        for (var arr : arrays) {
            System.arraycopy(arr, 0, result, i, arr.length);
            i += arr.length;
        }
        return result;
    }

    public static byte[] concat(byte[] ...arrays) {
        int length = 0;
        for (var arr : arrays)
            length += arr.length;
        var result = new byte[length];
        int i = 0;
        for (var arr : arrays) {
            System.arraycopy(arr, 0, result, i, arr.length);
            i += arr.length;
        }
        return result;
    }

    public static byte[] longToBytes(long val) {
        var result = new byte[8];
        for (int i = 8; i --> 0;) {
            result[i] = (byte) val;
            val >>>= 8;
        }
        return result;
    }

    public static BigInteger natFromBytes(byte[] input) {
        return new BigInteger(1, input);
    }

    public static BigInteger natFromBytes(byte[] input, int off, int length) {
        return new BigInteger(1, input, off, length);
    }

    public static byte[] natToBytes(BigInteger n, int length) {
        assert n.signum() >= 0;
        var bytes = n.toByteArray();
        if (bytes.length == length)
            return bytes;
        var result = new byte[length];
        if (bytes.length == length + 1 && bytes[0] == 0) {
            System.arraycopy(bytes, 1, result, 0, length);
        } else {
            assert bytes.length <= length;
            System.arraycopy(bytes, 0, result, length - bytes.length, bytes.length);
        }
        return result;
    }

    public static <T> T fromJson(String s, Class<T> c) {
        return GSON.fromJson(s, c);
    }

    public static String toJson(Object o) {
        return GSON.toJson(o);
    }
}
