package vote.pebble.voting.structs;

import vote.pebble.common.HashValue;
import vote.pebble.common.ParseException;
import vote.pebble.vdf.VDF;

import java.nio.BufferUnderflowException;
import java.nio.ByteBuffer;
import java.util.Arrays;

public final class DecryptionMessage {
    private static final String STRUCT = "DecryptionMessage";

    public final HashValue inputHash;
    public final byte[] output, proof;

    public DecryptionMessage(HashValue inputHash, byte[] output, byte[] proof) {
        this.inputHash = inputHash;
        this.output = output;
        this.proof = proof;
    }

    public DecryptionMessage(VDF.Solution vdfSol) {
        inputHash = HashValue.hash(vdfSol.input);
        output = vdfSol.output;
        proof = vdfSol.proof;
    }

    public static DecryptionMessage fromBytes(byte[] bytes) throws ParseException {
        try {
            var buffer = ByteBuffer.wrap(bytes);
            var inputHash = new byte[32];
            buffer.get(inputHash);
            int outputLen = buffer.getShort();
            if (outputLen < 0 || outputLen > 1024)
                throw new ParseException(STRUCT, "Invalid output size");
            var output = new byte[outputLen];
            buffer.get(output);
            var proof = new byte[buffer.remaining()];
            buffer.get(proof);
            return new DecryptionMessage(new HashValue(inputHash), output, proof);
        } catch (BufferUnderflowException e) {
            throw new ParseException(STRUCT, e);
        }
    }

    public byte[] toBytes() {
        assert output.length <= 1024;
        return ByteBuffer.allocate(34 + output.length + proof.length)
                .put(inputHash.bytes)
                .putShort((short) output.length)
                .put(output)
                .put(proof)
                .array();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        DecryptionMessage that = (DecryptionMessage) o;
        return that.inputHash.equals(inputHash) && Arrays.equals(output, that.output) && Arrays.equals(proof, that.proof);
    }

    @Override
    public int hashCode() {
        return (inputHash.hashCode() * 37 + Arrays.hashCode(output)) * 37 + Arrays.hashCode(proof);
    }
}
