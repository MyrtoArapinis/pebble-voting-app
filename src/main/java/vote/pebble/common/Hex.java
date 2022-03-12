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

    private static int hexValue(char c) {
        if (c >= '0' && c <= '9')
            return c - '0';
        if (c >= 'A' && c <= 'F')
            return c - 'A' + 10;
        if (c >= 'a' && c <= 'f')
            return c - 'a' + 10;
        throw new IllegalArgumentException("Not a hexadecimal digit");
    }

    public static byte[] decodeHexString(String s) {
        assert s.length() % 2 == 0;
        var res = new byte[s.length() / 2];
        for (int i = 0; i < res.length; i++)
            res[i] = (byte) (hexValue(s.charAt(i * 2)) * 16 + hexValue(s.charAt(i * 2 + 1)));
        return res;
    }
}
