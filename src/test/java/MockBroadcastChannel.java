import vote.pebble.voting.broadcast.BroadcastChannel;
import vote.pebble.voting.structs.*;

import java.util.ArrayList;

public class MockBroadcastChannel implements BroadcastChannel {
    final ArrayList<CredentialMessage> credentials = new ArrayList<>();
    final ArrayList<SignedBallot> ballots = new ArrayList<>();
    final ArrayList<DecryptionMessage> decryptions = new ArrayList<>();
    final ElectionParams params;

    public MockBroadcastChannel(ElectionParams params) {
        this.params = params;
    }

    @Override
    public void postCredential(CredentialMessage msg) {
        if (params.phase() == ElectionPhase.CRED_GEN)
            credentials.add(msg);
    }

    @Override
    public void postSignedBallot(SignedBallot msg) {
        if (params.phase() == ElectionPhase.VOTE)
            ballots.add(msg);
    }

    @Override
    public void postBallotDecryption(DecryptionMessage msg) {
        if (params.phase() == ElectionPhase.TALLY)
            decryptions.add(msg);
    }

    @Override
    public Iterable<CredentialMessage> getCredentials() {
        return credentials;
    }

    @Override
    public Iterable<SignedBallot> getSignedBallots() {
        return ballots;
    }

    @Override
    public Iterable<DecryptionMessage> getBallotDecryptions() {
        return decryptions;
    }
}
