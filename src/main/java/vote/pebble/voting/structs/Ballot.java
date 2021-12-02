package vote.pebble.voting.structs;

import java.util.Arrays;

public final class Ballot {
    public final byte[] content;

    public Ballot(byte[] content) {
        this.content = content;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        return Arrays.equals(content, ((Ballot) o).content);
    }

    @Override
    public int hashCode() {
        return Arrays.hashCode(content);
    }
}
