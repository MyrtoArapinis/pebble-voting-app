package vote.pebble.zkp;

import vote.pebble.common.PebbleException;

public class CredentialException extends PebbleException {
    public CredentialException(String message) {
        super(message);
    }

    public CredentialException(Throwable cause) {
        super(cause);
    }

    public CredentialException(String message, Throwable cause) {
        super(message, cause);
    }
}
