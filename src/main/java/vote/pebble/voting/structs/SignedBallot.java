package vote.pebble.voting.structs;

import vote.pebble.common.ByteString;
import vote.pebble.common.ParseException;
import vote.pebble.zkp.CredentialException;
import vote.pebble.zkp.CredentialSet;
import vote.pebble.zkp.SecretCredential;

import java.nio.BufferUnderflowException;
import java.nio.ByteBuffer;
import java.util.Arrays;

public final class SignedBallot {
    private static final String STRUCT = "SignedBallot";

    public final EncryptedBallot encryptedBallot;
    public final ByteString serialNo;
    public final byte[] signature;

    public SignedBallot(EncryptedBallot encryptedBallot, ByteString serialNo, byte[] signature) {
        assert serialNo.bytes.length <= 80 && signature.length <= 30000;
        this.encryptedBallot = encryptedBallot;
        this.serialNo = serialNo;
        this.signature = signature;
    }

    public static SignedBallot sign(EncryptedBallot encryptedBallot, CredentialSet credentialSet, SecretCredential secretCredential) throws CredentialException {
        var signature = credentialSet.sign(secretCredential, encryptedBallot.toBytes());
        return new SignedBallot(encryptedBallot, new ByteString(secretCredential.getSerialNumber()), signature);
    }

    public static SignedBallot fromBytes(byte[] bytes) throws ParseException {
        try {
            var buffer = ByteBuffer.wrap(bytes);
            int len = buffer.get();
            if (len < 0 || len > 80)
                throw new ParseException(STRUCT, "Invalid serial no. size");
            var serialNo = new byte[len];
            buffer.get(serialNo);
            len = buffer.getShort();
            if (len < 0 || len > 30000)
                throw new ParseException(STRUCT, "Invalid signature size");
            var signature = new byte[len];
            buffer.get(signature);
            var encBallotBytes = new byte[buffer.remaining()];
            buffer.get(encBallotBytes);
            return new SignedBallot(
                    EncryptedBallot.fromBytes(encBallotBytes),
                    new ByteString(serialNo),
                    signature);
        } catch (BufferUnderflowException e) {
            throw new ParseException(STRUCT, e);
        }
    }

    public byte[] toBytes() {
        var encBallotBytes = encryptedBallot.toBytes();
        return ByteBuffer.allocate(3 + serialNo.bytes.length + signature.length + encBallotBytes.length)
                .put((byte) serialNo.bytes.length)
                .put(serialNo.bytes)
                .putShort((short) signature.length)
                .put(signature)
                .put(encBallotBytes)
                .array();
    }

    public boolean verify(CredentialSet credentialSet) {
        return credentialSet.verify(serialNo.bytes, signature, encryptedBallot.toBytes());
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        SignedBallot that = (SignedBallot) o;
        return that.encryptedBallot.equals(encryptedBallot) && that.serialNo.equals(serialNo) && Arrays.equals(signature, that.signature);
    }

    @Override
    public int hashCode() {
        return (encryptedBallot.hashCode() * 37 + serialNo.hashCode()) * 37 + Arrays.hashCode(signature);
    }
}
