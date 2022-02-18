package vote.pebble.voting.methods;

import vote.pebble.voting.structs.Ballot;

import java.util.ArrayList;

public abstract class VotingMethod {
    public final int numChoices;

    protected VotingMethod(int numChoices) {
        this.numChoices = numChoices;
    }

    public static VotingMethod create(String method, int numChoices) throws NoSuchVotingMethodException, InvalidNumberOfChoicesException {
        if (method.equals("Plurality"))
            return new PluralityVoting(numChoices);
        if (method.equals("Approval"))
            return new ApprovalVoting(numChoices);
        throw new NoSuchVotingMethodException(method);
    }

    public abstract Ballot vote(int[] choices);

    public abstract ArrayList<TallyCount> tally(Iterable<Ballot> ballots);
}
