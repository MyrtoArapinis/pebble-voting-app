package vote.pebble.voting.broadcast;

import vote.pebble.voting.structs.CredentialMessage;
import vote.pebble.voting.structs.DecryptionMessage;
import vote.pebble.voting.structs.SignedBallot;

public interface BroadcastChannel {
    void postCredential(CredentialMessage msg);

    void postSignedBallot(SignedBallot msg);

    void postBallotDecryption(DecryptionMessage msg);

    Iterable<CredentialMessage> getCredentials();

    Iterable<SignedBallot> getSignedBallots();

    Iterable<DecryptionMessage> getBallotDecryptions();
}
