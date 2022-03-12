package vote.pebble.voting;

import vote.pebble.common.ParseException;
import vote.pebble.vdf.VDF;
import vote.pebble.voting.structs.SignedBallot;
import vote.pebble.zkp.CredentialException;
import vote.pebble.zkp.CredentialSystem;
import vote.pebble.zkp.SecretCredential;

import cafe.cryptography.ed25519.Ed25519PrivateKey;

public interface SecretsManager {
    Ed25519PrivateKey getPrivateKey();

    SecretCredential getSecretCredential(CredentialSystem system) throws CredentialException;

    SignedBallot getBallot() throws ParseException;

    void setBallot(SignedBallot ballot);

    VDF.Solution getVDFSolution();

    void setVDFSolution(VDF.Solution solution);
}
