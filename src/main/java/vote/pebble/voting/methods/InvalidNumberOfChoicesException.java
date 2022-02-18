package vote.pebble.voting.methods;

import vote.pebble.common.PebbleException;

public class InvalidNumberOfChoicesException extends PebbleException {
    public final String votingMethod;

    public InvalidNumberOfChoicesException(String votingMethod) {
        super("Invalid number of choices for '" + votingMethod + "' voting method");
        this.votingMethod = votingMethod;
    }
}
