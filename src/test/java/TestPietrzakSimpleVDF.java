import org.junit.jupiter.api.RepeatedTest;
import org.junit.jupiter.api.Test;
import vote.pebble.vdf.PietrzakSimpleVDF;
import vote.pebble.vdf.VDF;

import java.util.Random;

public class TestPietrzakSimpleVDF {
    static final Random random = new Random();

    static VDF createVDF() {
        int t = random.nextInt(100000) * 2 + 100000;
        return new PietrzakSimpleVDF(t);
    }

    @RepeatedTest(8)
    void testCreate() {
        var vdf = createVDF();
        var sol = vdf.create();
        assert vdf.verify(sol);
    }

    @Test
    void testSolve() {
        var vdf = createVDF();
        var puz = vdf.create();
        var sol = vdf.solve(puz.input);
        assert vdf.verify(sol);
    }
}
