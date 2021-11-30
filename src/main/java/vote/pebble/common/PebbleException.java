package vote.pebble.common;

public class PebbleException extends Exception {
    public PebbleException(String message) {
        super(message);
    }

    public PebbleException(Throwable cause) {
        super(cause);
    }

    public PebbleException(String message, Throwable cause) {
        super(message, cause);
    }
}
