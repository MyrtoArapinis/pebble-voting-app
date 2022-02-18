import cafe.cryptography.ed25519.Ed25519PrivateKey;
import vote.pebble.vdf.VDF;
import vote.pebble.voting.SecretsManager;
import vote.pebble.voting.structs.SignedBallot;

public class MockSecretsManager implements SecretsManager {
    Ed25519PrivateKey privateKey;
    byte[] secretCredential;
    SignedBallot ballot;
    VDF.Solution solution;

    @Override
    public Ed25519PrivateKey getPrivateKey() {
        return privateKey;
    }

    @Override
    public byte[] getSecretCredential() {
        return secretCredential;
    }

    @Override
    public SignedBallot getBallot() {
        return ballot;
    }

    @Override
    public void setBallot(SignedBallot ballot) {
        this.ballot = ballot;
    }

    @Override
    public VDF.Solution getVDFSolution() {
        return solution;
    }

    @Override
    public void setVDFSolution(VDF.Solution solution) {
        this.solution = solution;
    }
}
