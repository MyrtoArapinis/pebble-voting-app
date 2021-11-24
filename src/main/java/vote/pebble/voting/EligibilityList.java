package vote.pebble.voting;

import java.nio.ByteBuffer;
import java.util.ArrayList;
import java.util.HashMap;

import cafe.cryptography.curve25519.InvalidEncodingException;
import cafe.cryptography.ed25519.Ed25519PublicKey;
import vote.pebble.common.HashValue;
import vote.pebble.common.ParseException;

public class EligibilityList {
    private static final int MAGIC = 0x454c4c01;
    private static final String STRUCT = "EligibilityList";

    private final ArrayList<Ed25519PublicKey> publicKeys = new ArrayList<>();
    private final HashMap<Ed25519PublicKey, HashValue> idComs = new HashMap<>();

    public void add(Ed25519PublicKey publicKey, HashValue idCom) {
        if (idComs.containsKey(publicKey))
            throw new IllegalStateException("EligibilityList already contains public key");
        if (idCom == null)
            idCom = HashValue.ZERO;
        publicKeys.add(publicKey);
        idComs.put(publicKey, idCom);
    }

    public HashValue getIdentityCommitment(Ed25519PublicKey publicKey) {
        return idComs.get(publicKey);
    }

    public boolean contains(Ed25519PublicKey publicKey) {
        return idComs.containsKey(publicKey);
    }

    public byte[] toBytes() {
        var buffer = ByteBuffer.allocate(4 + 64 * publicKeys.size());
        buffer.putInt(MAGIC);
        for (var publicKey : publicKeys) {
            var idCom = idComs.get(publicKey);
            buffer.put(publicKey.toByteArray()).put(idCom.bytes);
        }
        return buffer.array();
    }

    public static EligibilityList fromBytes(byte[] bytes) throws ParseException {
        int entriesSize = bytes.length - 4;
        if (entriesSize < 0 || entriesSize % 64 != 0)
            throw new ParseException(STRUCT, "Invalid size");
        var buffer = ByteBuffer.wrap(bytes);
        if (buffer.getInt() != MAGIC)
            throw new ParseException(STRUCT, "Invalid magic");
        var bytes32 = new byte[32];
        var result = new EligibilityList();
        for (int i = 0; i < entriesSize / 64; i++) {
            Ed25519PublicKey publicKey;
            buffer.get(bytes32);
            try {
                publicKey = Ed25519PublicKey.fromByteArray(bytes32);
            } catch (InvalidEncodingException e) {
                throw new ParseException(STRUCT, e);
            }
            buffer.get(bytes32);
            var idCom = new HashValue(bytes32);
            result.add(publicKey, idCom);
        }
        return result;
    }

    public HashValue hash() {
        return HashValue.sha256(toBytes());
    }
}
