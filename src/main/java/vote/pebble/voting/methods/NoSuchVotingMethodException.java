package vote.pebble.voting.methods;

import vote.pebble.common.PebbleException;

public class NoSuchVotingMethodException extends PebbleException {
    public NoSuchVotingMethodException(String method) {
        super("No such Voting Method: " + method);
    }
}
