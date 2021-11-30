package vote.pebble.voting;

import vote.pebble.common.HashValue;
import vote.pebble.vdf.VDF;

import javax.crypto.Cipher;
import javax.crypto.spec.GCMParameterSpec;
import javax.crypto.spec.SecretKeySpec;
import java.security.GeneralSecurityException;
import java.util.Arrays;

public class EncryptedBallot {
    public final byte[] vdfInput, payload;

    public EncryptedBallot(byte[] vdfInput, byte[] payload) {
        this.vdfInput = vdfInput;
        this.payload = payload;
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
}
