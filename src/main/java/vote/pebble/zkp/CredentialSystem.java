package vote.pebble.zkp;

public interface CredentialSystem {
    SecretCredential generateSecretCredential() throws CredentialException;

    SecretCredential secretCredentialFromBytes(byte[] bytes) throws CredentialException;

    PublicCredential publicCredentialFromBytes(byte[] bytes) throws CredentialException;

    CredentialSet makeCredentialSet(Iterable<PublicCredential> publicCredentials) throws CredentialException;
}
