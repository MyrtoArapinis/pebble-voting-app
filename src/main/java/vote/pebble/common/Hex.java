package vote.pebble.common;

public final class Hex {
    private static final char[] HEX = "0123456789abcdef".toCharArray();

    public static char[] encodeHex(byte[] data) {
        var res = new char[data.length * 2];
        for (int i = 0; i < data.length; i++) {
            int b = data[i] & 255;
            res[i * 2] = HEX[b >>> 4];
            res[i * 2 + 1] = HEX[b & 15];
        }
        return res;
    }

    public static String encodeHexString(byte[] data) {
        return new String(encodeHex(data));
    }
}
