package vote.pebble.voting;

import cafe.cryptography.ed25519.Ed25519PrivateKey;
import vote.pebble.vdf.VDF;
import vote.pebble.voting.structs.SignedBallot;

public final class ObserverSecretsManager implements SecretsManager {
    private static final ObserverSecretsManager INSTANCE = new ObserverSecretsManager();

    public static ObserverSecretsManager getInstance() {
        return INSTANCE;
    }

    @Override
    public Ed25519PrivateKey getPrivateKey() {
        throw new UnsupportedOperationException();
    }

    @Override
    public byte[] getSecretCredential() {
        throw new UnsupportedOperationException();
    }

    @Override
    public SignedBallot getBallot() {
        throw new UnsupportedOperationException();
    }

    @Override
    public void setBallot(SignedBallot ballot) {
        throw new UnsupportedOperationException();
    }

    @Override
    public VDF.Solution getVDFSolution() {
        throw new UnsupportedOperationException();
    }

    @Override
    public void setVDFSolution(VDF.Solution solution) {
        throw new UnsupportedOperationException();
    }
}
