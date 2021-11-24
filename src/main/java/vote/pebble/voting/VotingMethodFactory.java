package vote.pebble.voting;

import vote.pebble.voting.methods.ApprovalVoting;
import vote.pebble.voting.methods.PluralityVoting;

public final class VotingMethodFactory {
    private static final PluralityVoting pluralityVoting = new PluralityVoting();
    private static final ApprovalVoting approvalVoting = new ApprovalVoting();

    public static VotingMethod getInstance(String method) throws NoSuchVotingMethodException {
        if (method.equals(pluralityVoting.toString()))
            return pluralityVoting;
        if (method.equals(approvalVoting.toString()))
            return approvalVoting;
        throw new NoSuchVotingMethodException(method);
    }
}
