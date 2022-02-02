import vote.pebble.zkp.*;

import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.assertTrue;

import java.util.ArrayList;

public class TestAnonCred1 {
    static CredentialSystem credentialSystem;

    @BeforeAll
    static void init() {
        credentialSystem = new AnonCred1();
    }

    @Test
    void testGenerateSecretCredential() throws CredentialException {
        var secretCredential = credentialSystem.generateSecretCredential();
    }

    @Test
    void testHashMerkleTree() throws CredentialException {
        var list = new ArrayList<PublicCredential>();
        for (int i = 0; i < 562; i++)
            list.add(credentialSystem.generateSecretCredential().getPublicCredential());
        var set = credentialSystem.makeCredentialSet(list);
    }

    @Test
    void testSignVerify() throws CredentialException {
        var list = new ArrayList<PublicCredential>();
        var secretCredential = credentialSystem.generateSecretCredential();
        list.add(secretCredential.getPublicCredential());
        for (int i = 0; i < 53; i++)
            list.add(credentialSystem.generateSecretCredential().getPublicCredential());
        var set = credentialSystem.makeCredentialSet(list);
        var msg = new byte[] { 1, 2, 3, 4, 5 };
        var sig = set.sign(secretCredential, msg);
        assertTrue(set.verify(secretCredential.getSerialNumber(), sig, msg));
    }
}
