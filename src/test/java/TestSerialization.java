import vote.pebble.common.ByteString;
import vote.pebble.common.HashValue;
import vote.pebble.common.ParseException;
import vote.pebble.voting.EligibilityList;
import vote.pebble.voting.EncryptedBallot;
import vote.pebble.voting.SignedBallot;

import cafe.cryptography.ed25519.Ed25519PrivateKey;
import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;

import java.security.SecureRandom;

public class TestSerialization {
    static final SecureRandom secureRandom = new SecureRandom();

    static byte[] randomBytes(int size) {
        var res = new byte[size];
        secureRandom.nextBytes(res);
        return res;
    }

    static void testEligibilityList(int size) throws ParseException {
        var ell1 = new EligibilityList();
        for (int i = 0; i < 10; i++) {
            var idCom = new byte[32];
            secureRandom.nextBytes(idCom);
            var publicKey = Ed25519PrivateKey.generate(secureRandom).derivePublic();
            ell1.add(publicKey, new HashValue(idCom));
        }
        var ell2 = EligibilityList.fromBytes(ell1.toBytes());
        assertEquals(ell1.hash(), ell2.hash());
    }

    @Test
    void testEligibilityList() throws ParseException {
        testEligibilityList(11);
    }

    @Test
    void testEmptyEligibilityList() throws ParseException {
        testEligibilityList(0);
    }

    @Test
    void testEncryptedBallot() throws ParseException {
        var ballot1 = new EncryptedBallot(randomBytes(512), randomBytes(3));
        var ballot2 = EncryptedBallot.fromBytes(ballot1.toBytes());
        assertEquals(ballot1, ballot2);
    }

    @Test
    void testSignedBallot() throws ParseException {
        var ballot1 = new SignedBallot(
                new EncryptedBallot(randomBytes(512), randomBytes(3)),
                new ByteString(randomBytes(16)),
                randomBytes(256));
        var ballot2 = SignedBallot.fromBytes(ballot1.toBytes());
        assertEquals(ballot1, ballot2);
    }
}
