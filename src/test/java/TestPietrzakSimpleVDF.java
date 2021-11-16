import org.junit.jupiter.api.Test;
import vote.pebble.vdf.PietrzakSimpleVDF;

import java.util.Random;

public class TestPietrzakSimpleVDF {
    static final Random random = new Random();

    @Test
    void testCreate() {
        int t = random.nextInt(100000) * 2 + 100000;
        var vdf = new PietrzakSimpleVDF(t);
        var sol = vdf.create();
        assert vdf.verify(sol);
    }

    @Test
    void testSolve() {
        int t = random.nextInt(100000) * 2 + 100000;
        var vdf = new PietrzakSimpleVDF(t);
        var puz = vdf.create();
        var sol = vdf.solve(puz.input);
        assert vdf.verify(sol);
    }
}
