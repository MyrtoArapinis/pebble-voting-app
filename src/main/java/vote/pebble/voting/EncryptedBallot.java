package vote.pebble.voting;

import vote.pebble.common.HashValue;
import vote.pebble.common.ParseException;
import vote.pebble.vdf.VDF;

import javax.crypto.Cipher;
import javax.crypto.spec.GCMParameterSpec;
import javax.crypto.spec.SecretKeySpec;
import java.nio.BufferUnderflowException;
import java.nio.ByteBuffer;
import java.security.GeneralSecurityException;
import java.util.Arrays;

public final class EncryptedBallot {
    private static final String STRUCT = "EncryptedBallot";

    public final byte[] vdfInput, payload;

    public EncryptedBallot(byte[] vdfInput, byte[] payload) {
        this.vdfInput = vdfInput;
        this.payload = payload;
    }

    public static EncryptedBallot fromBytes(byte[] bytes) throws ParseException {
        try {
            var buffer = ByteBuffer.wrap(bytes);
            int len = buffer.getShort();
            if (len < 0 || len > 4096)
                throw new ParseException(STRUCT, "Invalid VDF input size");
            var vdfInput = new byte[len];
            buffer.get(vdfInput);
            len = buffer.remaining();
            if (len < 0 || len > 4096)
                throw new ParseException(STRUCT, "Invalid payload size");
            var payload = new byte[len];
            buffer.get(payload);
            return new EncryptedBallot(vdfInput, payload);
        } catch (BufferUnderflowException e) {
            throw new ParseException(STRUCT, e);
        }
    }

    public byte[] toBytes() {
        assert vdfInput.length <= 4096 && payload.length <= 4096;
        return ByteBuffer.allocate(2 + vdfInput.length + payload.length)
                .putShort((short) vdfInput.length)
                .put(vdfInput)
                .put(payload)
                .array();
    }

    private static Cipher createCipher(int mode, VDF.Solution vdfSol) throws GeneralSecurityException {
        var md = HashValue.createMessageDigest();
        md.update(vdfSol.input);
        md.update(vdfSol.output);
        var bytes = md.digest();
        var key = new SecretKeySpec(bytes, 0, 16, "AES");
        var params = new GCMParameterSpec(128, bytes, 16, 12);
        var cipher = Cipher.getInstance("AES/GCM/NoPadding");
        cipher.init(mode, key, params);
        return cipher;
    }

    public static EncryptedBallot encrypt(Ballot ballot, VDF.Solution vdfSol) {
        try {
            var cipher = createCipher(Cipher.ENCRYPT_MODE, vdfSol);
            var payload = cipher.doFinal(ballot.content);
            return new EncryptedBallot(vdfSol.input, payload);
        } catch (GeneralSecurityException e) {
            throw new RuntimeException(e);
        }
    }

    public Ballot decrypt(VDF.Solution vdfSol) throws BallotDecryptionFailedException {
        assert Arrays.equals(vdfInput, vdfSol.input);
        try {
            var cipher = createCipher(Cipher.DECRYPT_MODE, vdfSol);
            return new Ballot(cipher.doFinal(payload));
        } catch (GeneralSecurityException e) {
            throw new BallotDecryptionFailedException(e);
        }
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        EncryptedBallot that = (EncryptedBallot) o;
        return Arrays.equals(vdfInput, that.vdfInput) && Arrays.equals(payload, that.payload);
    }

    @Override
    public int hashCode() {
        return Arrays.hashCode(vdfInput) * 37 + Arrays.hashCode(payload);
    }
}
