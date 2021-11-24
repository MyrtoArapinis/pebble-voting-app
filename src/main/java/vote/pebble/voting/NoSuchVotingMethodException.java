package vote.pebble.voting;

public class NoSuchVotingMethodException extends Exception {
    public NoSuchVotingMethodException(String method) {
        super("No such Voting Method: " + method);
    }
}
