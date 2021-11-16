package vote.pebble.vdf;

public interface VDF {
    class Solution {
        public final byte[] input, output, proof;

        public Solution(byte[] input, byte[] output, byte[] proof) {
            this.input = input;
            this.output = output;
            this.proof = proof;
        }
    }

    Solution create();

    Solution solve(byte[] input);

    boolean verify(Solution solution);
}
