package vote.pebble.zkp;

public interface SecretCredential {
    byte[] toBytes();

    PublicCredential getPublicCredential();

    byte[] getSerialNumber();
}
