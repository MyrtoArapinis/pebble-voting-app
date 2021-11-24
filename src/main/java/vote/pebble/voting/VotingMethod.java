package vote.pebble.voting;

import java.util.ArrayList;

public interface VotingMethod {
    ArrayList<TallyCount> tally(int numChoices, Iterable<Ballot> ballots);
}
