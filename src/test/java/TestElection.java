import vote.pebble.common.PebbleException;
import vote.pebble.voting.BallotNotDecryptedException;
import vote.pebble.voting.Election;
import vote.pebble.voting.structs.ElectionParams;
import vote.pebble.voting.structs.EligibilityList;
import vote.pebble.zkp.AnonCred1;
import vote.pebble.zkp.CredentialException;
import vote.pebble.zkp.CredentialSystem;
import vote.pebble.zkp.SecretCredential;

import cafe.cryptography.ed25519.Ed25519PrivateKey;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.BeforeAll;

import java.security.SecureRandom;
import java.time.Instant;
import java.util.ArrayList;
import java.util.Random;

import static org.junit.jupiter.api.Assertions.assertThrows;

public class TestElection {
    static final String[] CANDIDATES = { "Toby Wilkinson", "Ava McLean", "Oliver Rogers" };
    static final int VOTERS_COUNT = 10;

    static final Random random = new Random();
    static final SecureRandom secureRandom = new SecureRandom();

    static CredentialSystem credentialSystem;

    @BeforeAll
    static void init() {
        credentialSystem = new AnonCred1();
    }

    ArrayList<Ed25519PrivateKey> generatePrivateKeys(int count) {
        var result = new ArrayList<Ed25519PrivateKey>(count);
        while (count --> 0)
            result.add(Ed25519PrivateKey.generate(secureRandom));
        return result;
    }

    EligibilityList generateEligibilityList(Iterable<Ed25519PrivateKey> privateKeys) {
        var result = new EligibilityList();
        for (var key : privateKeys)
            result.add(key.derivePublic(), null);
        return result;
    }

    ElectionParams generateElectionParams(EligibilityList eligibilityList) {
        var now = Instant.now();
        return new ElectionParams(eligibilityList, now.plusSeconds(10), now.plusSeconds(20),
                10000, "Plurality", CANDIDATES);
    }

    ArrayList<SecretCredential> generateSecretCredentials(int count) throws CredentialException {
        var result = new ArrayList<SecretCredential>(count);
        while (count --> 0)
            result.add(credentialSystem.generateSecretCredential());
        return result;
    }

    @Test
    void testElection() throws PebbleException, InterruptedException {
        var privateKeys = generatePrivateKeys(VOTERS_COUNT);
        var eligibilityList = generateEligibilityList(privateKeys);
        var secretCredentials = generateSecretCredentials(VOTERS_COUNT);
        var electionParams = generateElectionParams(eligibilityList);
        var secretsManager = new MockSecretsManager();
        var election = new Election(electionParams, new MockBroadcastChannel(electionParams), secretsManager);
        for (int i = 0; i < VOTERS_COUNT; i++) {
            secretsManager.privateKey = privateKeys.get(i);
            secretsManager.secretCredential = secretCredentials.get(i);
            election.postCredential();
        }
        Thread.sleep(11000);
        int voterIdx = random.nextInt(VOTERS_COUNT);
        secretsManager.secretCredential = secretCredentials.get(voterIdx);
        election.vote(CANDIDATES[random.nextInt(CANDIDATES.length)]);
        Thread.sleep(11000);
        assertThrows(BallotNotDecryptedException.class, election::tally);
        election.revealBallotDecryption();
        election.tally();
    }
}
