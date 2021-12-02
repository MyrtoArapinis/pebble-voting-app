package vote.pebble.voting.structs;

import vote.pebble.common.PebbleException;

public class BallotDecryptionFailedException extends PebbleException {
    public BallotDecryptionFailedException(String message) {
        super(message);
    }

    public BallotDecryptionFailedException(Throwable cause) {
        super(cause);
    }

    public BallotDecryptionFailedException(String message, Throwable cause) {
        super(message, cause);
    }
}
