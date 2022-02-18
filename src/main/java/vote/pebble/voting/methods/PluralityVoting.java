package vote.pebble.voting.methods;

import vote.pebble.voting.structs.Ballot;

import java.util.ArrayList;

public class PluralityVoting extends VotingMethod {
    protected PluralityVoting(int numChoices) throws InvalidNumberOfChoicesException {
        super(numChoices);
        if (numChoices <= 0 || numChoices > 256)
            throw new InvalidNumberOfChoicesException(toString());
    }

    @Override
    public Ballot vote(int[] choices) {
        assert choices.length == 1 && choices[0] >= 0 && choices[0] < 256;
        return new Ballot(new byte[] {(byte) choices[0]});
    }

    @Override
    public ArrayList<TallyCount> tally(Iterable<Ballot> ballots) {
        var counts = new ArrayList<TallyCount>(numChoices);
        for (int i = 0; i < numChoices; i++)
            counts.add(new TallyCount(i));
        for (var ballot : ballots) {
            if (ballot.content.length != 1)
                continue;
            int index = Byte.toUnsignedInt(ballot.content[0]);
            if (index >= numChoices)
                continue;
            counts.get(index).addOne();
        }
        return counts;
    }

    @Override
    public String toString() {
        return "Plurality";
    }
}
