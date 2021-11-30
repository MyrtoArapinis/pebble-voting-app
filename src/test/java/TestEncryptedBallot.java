import org.junit.jupiter.api.RepeatedTest;
import static org.junit.jupiter.api.Assertions.*;

import vote.pebble.vdf.VDF;
import vote.pebble.voting.Ballot;
import vote.pebble.voting.BallotDecryptionFailedException;
import vote.pebble.voting.EncryptedBallot;

import java.util.Random;

public class TestEncryptedBallot {
    static final Random random = new Random();

    static Ballot createBallot() {
        var content = new byte[16];
        random.nextBytes(content);
        return new Ballot(content);
    }

    static VDF.Solution createSolution() {
        var input = new byte[512];
        var output = new byte[256];
        random.nextBytes(input);
        random.nextBytes(output);
        return new VDF.Solution(input, output, null);
    }

    static void decrypt(Ballot ballot, EncryptedBallot encBallot, VDF.Solution sol) throws BallotDecryptionFailedException {
        var decBallot = encBallot.decrypt(sol);
        assertEquals(ballot, decBallot);
    }

    @RepeatedTest(8)
    void testEncryption() throws BallotDecryptionFailedException {
        var ballot = createBallot();
        var sol1 = createSolution();
        var sol2 = new VDF.Solution(sol1.input, createSolution().output, sol1.proof);
        var encBallot = EncryptedBallot.encrypt(ballot, sol1);
        decrypt(ballot, encBallot, sol1);
        assertThrows(BallotDecryptionFailedException.class, () -> decrypt(ballot, encBallot, sol2));
    }
}
