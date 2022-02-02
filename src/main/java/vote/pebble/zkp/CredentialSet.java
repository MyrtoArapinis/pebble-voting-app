package vote.pebble.zkp;

public interface CredentialSet {
    byte[] sign(SecretCredential secretCredential, byte[] message) throws CredentialException;

    boolean verify(byte[] serialNo, byte[] signature, byte[] message);
}
