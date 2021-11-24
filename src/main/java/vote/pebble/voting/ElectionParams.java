package vote.pebble.voting;

import java.time.Instant;

public class ElectionParams {
    public final EligibilityList eligibilityList;
    public final Instant voteStart, tallyStart;
    public final String votingMethod;
    public final String[] choices;

    public ElectionParams(EligibilityList eligibilityList, Instant voteStart, Instant tallyStart, String votingMethod, String[] choices) {
        this.eligibilityList = eligibilityList;
        this.voteStart = voteStart;
        this.tallyStart = tallyStart;
        this.votingMethod = votingMethod;
        this.choices = choices;
    }
}
