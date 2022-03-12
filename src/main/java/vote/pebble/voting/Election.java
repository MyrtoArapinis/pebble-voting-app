package vote.pebble.voting;

import vote.pebble.common.ByteString;
import vote.pebble.common.HashValue;
import vote.pebble.common.Util;
import vote.pebble.vdf.PietrzakSimpleVDF;
import vote.pebble.vdf.VDF;
import vote.pebble.voting.broadcast.BroadcastChannel;
import vote.pebble.voting.methods.InvalidNumberOfChoicesException;
import vote.pebble.voting.methods.NoSuchVotingMethodException;
import vote.pebble.voting.methods.TallyCount;
import vote.pebble.voting.methods.VotingMethod;
import vote.pebble.voting.structs.*;
import vote.pebble.zkp.*;

import cafe.cryptography.ed25519.Ed25519PublicKey;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;

public final class Election {
    private static final CredentialSystem CRED_SYS = new AnonCred1();

    private final BroadcastChannel channel;
    private final SecretsManager secretsManager;
    private final VDF vdf;
    private final VotingMethod votingMethod;

    public final ElectionParams params;

    public Election(ElectionParams params, BroadcastChannel channel)
            throws NoSuchVotingMethodException, InvalidNumberOfChoicesException {
        this(params, channel, ObserverSecretsManager.getInstance());
    }

    public Election(ElectionParams params, BroadcastChannel channel, SecretsManager secretsManager)
            throws NoSuchVotingMethodException, InvalidNumberOfChoicesException {
        this.params = params;
        this.channel = channel;
        this.secretsManager = secretsManager;
        this.vdf = new PietrzakSimpleVDF(params.vdfDifficulty);
        this.votingMethod = VotingMethod.create(params.votingMethod, params.choices.length);
    }

    public void postCredential() throws CredentialException {
        assert params.phase() == ElectionPhase.CRED_GEN;
        var privateKey = secretsManager.getPrivateKey();
        var credential = secretsManager.getSecretCredential(CRED_SYS)
                .getPublicCredential().toBytes();
        var msg = CredentialMessage.sign(privateKey, credential);
        channel.postCredential(msg);
    }

    private CredentialSet getCredentialSet() throws CredentialException {
        assert params.phase().compareTo(ElectionPhase.CRED_GEN) > 0;
        var credentials = new HashMap<Ed25519PublicKey, PublicCredential>();
        for (var msg : channel.getCredentials()) {
            if (params.eligibilityList.contains(msg.publicKey) && msg.verify()) {
                try {
                    credentials.put(msg.publicKey, CRED_SYS.publicCredentialFromBytes(msg.credential));
                } catch (CredentialException ignored) {}
            }
        }
        return CRED_SYS.makeCredentialSet(credentials.values());
    }

    public void vote(String choice) throws CredentialException {
        vote(List.of(choice));
    }

    public void vote(List<String> choices) throws CredentialException {
        assert params.phase() == ElectionPhase.VOTE;
        var indexes = new int[choices.size()];
        for (int i = 0; i < indexes.length; i++)
            indexes[i] = Util.indexOf(params.choices, choices.get(i));
        var set = getCredentialSet();
        var solution = vdf.create();
        secretsManager.setVDFSolution(solution);
        var secret = secretsManager.getSecretCredential(CRED_SYS);
        var ballot = votingMethod.vote(indexes).encrypt(solution).sign(set, secret);
        secretsManager.setBallot(ballot);
        channel.postSignedBallot(ballot);
    }

    public void revealBallotDecryption() {
        postBallotDecryption(secretsManager.getVDFSolution());
    }

    public void postBallotDecryption(VDF.Solution vdfSolution) {
        assert params.phase() == ElectionPhase.TALLY;
        channel.postBallotDecryption(new DecryptionMessage(vdfSolution));
    }

    public void postBallotDecryption(EncryptedBallot encryptedBallot) {
        postBallotDecryption(vdf.solve(encryptedBallot.vdfInput));
    }

    private Ballot decryptBallot(EncryptedBallot encryptedBallot) throws BallotNotDecryptedException {
        var vdfInputHash = HashValue.hash(encryptedBallot.vdfInput);
        for (var msg : channel.getBallotDecryptions()) {
            if (msg.inputHash.equals(vdfInputHash)) {
                var vdfSol = new VDF.Solution(encryptedBallot.vdfInput, msg.output, msg.proof);
                if (vdf.verify(vdfSol)) {
                    try {
                        return encryptedBallot.decrypt(vdfSol);
                    } catch (BallotDecryptionFailedException ignored) {}
                }
            }
        }
        throw new BallotNotDecryptedException(encryptedBallot);
    }

    public ArrayList<TallyCount> tally() throws CredentialException, BallotNotDecryptedException {
        assert params.phase() == ElectionPhase.TALLY;
        var set = getCredentialSet();
        var serialNos = new HashSet<ByteString>();
        var decryptedBallots = new ArrayList<Ballot>();
        for (var signedBallot : channel.getSignedBallots()) {
            if (serialNos.contains(signedBallot.serialNo) || !signedBallot.verify(set))
                continue;
            serialNos.add(signedBallot.serialNo);
            decryptedBallots.add(decryptBallot(signedBallot.encryptedBallot));
        }
        return votingMethod.tally(decryptedBallots);
    }
}
