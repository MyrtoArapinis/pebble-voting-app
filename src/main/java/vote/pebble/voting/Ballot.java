package vote.pebble.voting;

import java.util.Arrays;

public final class Ballot {
    public final byte[] content;

    public Ballot(byte[] content) {
        this.content = content;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o instanceof Ballot)
            return Arrays.equals(content, ((Ballot) o).content);
        return false;
    }

    @Override
    public int hashCode() {
        return Arrays.hashCode(content);
    }
}
