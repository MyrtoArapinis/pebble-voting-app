package vote.pebble.voting.methods;

import vote.pebble.voting.structs.Ballot;

import java.util.ArrayList;

public class ApprovalVoting extends VotingMethod {
    protected ApprovalVoting(int numChoices) throws InvalidNumberOfChoicesException {
        super(numChoices);
        if (numChoices <= 0 || numChoices > 1024)
            throw new InvalidNumberOfChoicesException(toString());
    }

    @Override
    public Ballot vote(int[] choices) {
        assert choices.length <= numChoices;
        var content = new byte[numChoices];
        for (int choice : choices)
            content[choice] = 1;
        return new Ballot(content);
    }

    @Override
    public ArrayList<TallyCount> tally(Iterable<Ballot> ballots) {
        var counts = new ArrayList<TallyCount>(numChoices);
        for (int i = 0; i < numChoices; i++)
            counts.add(new TallyCount(i));
        for (var ballot : ballots) {
            if (ballot.content.length != numChoices)
                continue;
            boolean ballotInvalid = false;
            for (int i = 0; i < numChoices; i++) {
                int approval = ballot.content[i];
                if (approval < 0 || approval > 1) {
                    ballotInvalid = true;
                    break;
                }
            }
            if (ballotInvalid)
                continue;
            for (int i = 0; i < numChoices; i++)
                counts.get(i).add(ballot.content[i]);
        }
        return counts;
    }

    @Override
    public String toString() {
        return "Approval";
    }
}
