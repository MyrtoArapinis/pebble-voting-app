package vote.pebble.zkp;

public interface CredentialSystem {
    SecretCredential generateSecretCredential() throws CredentialException;

    SecretCredential secretCredentialFromBytes(byte[] bytes);

    PublicCredential publicCredentialFromBytes(byte[] bytes);

    CredentialSet makeCredentialSet(Iterable<PublicCredential> publicCredentials) throws CredentialException;
}
