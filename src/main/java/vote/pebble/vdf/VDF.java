package vote.pebble.vdf;

import vote.pebble.common.HashValue;
import vote.pebble.common.ParseException;
import vote.pebble.voting.structs.DecryptionMessage;

import java.nio.BufferUnderflowException;
import java.nio.ByteBuffer;

public interface VDF {
    final class Solution {
        public final byte[] input, output, proof;

        public Solution(byte[] input, byte[] output, byte[] proof) {
            this.input = input;
            this.output = output;
            this.proof = proof;
        }

        public static Solution fromBytes(byte[] bytes) {
            try {
                var buffer = ByteBuffer.wrap(bytes);
                int len = buffer.getShort();
                if (len < 0 || len > 2048)
                    throw new RuntimeException("Invalid input size: " + len);
                var input = new byte[len];
                buffer.get(input);
                len = buffer.getShort();
                if (len < 0 || len > 1024)
                    throw new RuntimeException("Invalid output size: " + len);
                var output = new byte[len];
                buffer.get(output);
                var proof = new byte[buffer.remaining()];
                buffer.get(proof);
                return new Solution(input, output, proof);
            } catch (BufferUnderflowException e) {
                throw new RuntimeException(e);
            }
        }

        public byte[] toBytes() {
            assert input.length <= 2048;
            assert output.length <= 1024;
            return ByteBuffer.allocate(4 + input.length + output.length + proof.length)
                    .putShort((short) input.length)
                    .put(input)
                    .putShort((short) output.length)
                    .put(output)
                    .put(proof)
                    .array();
        }
    }

    Solution create();

    Solution solve(byte[] input);

    boolean verify(Solution solution);
}
