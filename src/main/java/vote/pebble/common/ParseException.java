package vote.pebble.common;

public class ParseException extends Exception {
    public final String structure;

    public ParseException(String structure, String message) {
        super(message);
        this.structure = structure;
    }

    public ParseException(String structure, Throwable cause) {
        super(cause);
        this.structure = structure;
    }

    public ParseException(String structure, String message, Throwable cause) {
        super(message, cause);
        this.structure = structure;
    }

    @Override
    public String getMessage() {
        return '[' + structure + "] " + super.getMessage();
    }
}
