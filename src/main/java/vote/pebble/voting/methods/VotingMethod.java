package vote.pebble.voting.methods;

import vote.pebble.voting.structs.Ballot;

import java.util.ArrayList;

public abstract class VotingMethod {
    private static final PluralityVoting pluralityVoting = new PluralityVoting();
    private static final ApprovalVoting approvalVoting = new ApprovalVoting();

    public static VotingMethod getInstance(String method) throws NoSuchVotingMethodException {
        if (method.equals(pluralityVoting.toString()))
            return pluralityVoting;
        if (method.equals(approvalVoting.toString()))
            return approvalVoting;
        throw new NoSuchVotingMethodException(method);
    }

    public abstract ArrayList<TallyCount> tally(int numChoices, Iterable<Ballot> ballots);
}
