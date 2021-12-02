package vote.pebble.voting.methods;

import java.util.Objects;

public final class TallyCount implements Comparable<TallyCount> {
    public final int index;
    public long count = 0;

    public TallyCount(int index) {
        this.index = index;
    }

    public void add(int count) {
        this.count += count;
    }

    public void addOne() {
        count++;
    }

    @Override
    public int compareTo(TallyCount other) {
        if (count < other.count) return -1;
        if (count > other.count) return 1;
        return Integer.compare(index, other.index);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        TallyCount that = (TallyCount) o;
        return index == that.index && count == that.count;
    }

    @Override
    public int hashCode() {
        return Objects.hash(index, count);
    }
}
