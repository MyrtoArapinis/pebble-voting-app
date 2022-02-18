package vote.pebble.voting.structs;

import java.time.Instant;

public final class ElectionParams {
    public final EligibilityList eligibilityList;
    public final Instant voteStart, tallyStart;
    public final long vdfDifficulty;
    public final String votingMethod;
    public final String[] choices;

    public ElectionParams(EligibilityList eligibilityList, Instant voteStart, Instant tallyStart, long vdfDifficulty, String votingMethod, String[] choices) {
        this.eligibilityList = eligibilityList;
        this.voteStart = voteStart;
        this.tallyStart = tallyStart;
        this.vdfDifficulty = vdfDifficulty;
        this.votingMethod = votingMethod;
        this.choices = choices;
    }

    public ElectionPhase phase() {
        var now = Instant.now();
        if (now.compareTo(voteStart) < 0)
            return ElectionPhase.CRED_GEN;
        if (now.compareTo(tallyStart) < 0)
            return ElectionPhase.VOTE;
        return ElectionPhase.TALLY;
    }
}
