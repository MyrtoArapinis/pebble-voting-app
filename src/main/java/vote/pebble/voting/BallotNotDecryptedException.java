package vote.pebble.voting;

import vote.pebble.common.PebbleException;
import vote.pebble.voting.structs.EncryptedBallot;

public class BallotNotDecryptedException extends PebbleException {
    public final EncryptedBallot encryptedBallot;

    public BallotNotDecryptedException(EncryptedBallot encryptedBallot) {
        super("A ballot has not been decrypted yet");
        this.encryptedBallot = encryptedBallot;
    }
}
