package vote.pebble.voting;

import vote.pebble.common.PebbleException;

public class BallotNotDecryptedException extends PebbleException {
    public BallotNotDecryptedException() {
        super("A ballot has not been decrypted yet");
    }
}
