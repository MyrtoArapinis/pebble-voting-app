package vote.pebble.voting;

import vote.pebble.vdf.VDF;
import vote.pebble.voting.structs.SignedBallot;

import cafe.cryptography.ed25519.Ed25519PrivateKey;

public interface SecretsManager {
    Ed25519PrivateKey getPrivateKey();

    byte[] getSecretCredential();

    SignedBallot getBallot();

    void setBallot(SignedBallot ballot);

    VDF.Solution getVDFSolution();

    void setVDFSolution(VDF.Solution solution);
}
