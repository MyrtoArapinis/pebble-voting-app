package vote.pebble.voting.methods;

import vote.pebble.voting.Ballot;
import vote.pebble.voting.TallyCount;
import vote.pebble.voting.VotingMethod;

import java.util.ArrayList;

public class PluralityVoting implements VotingMethod {
    @Override
    public ArrayList<TallyCount> tally(int numChoices, Iterable<Ballot> ballots) {
        assert numChoices >= 0 && numChoices <= 256;
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
