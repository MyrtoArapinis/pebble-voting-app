package vote.pebble.voting;

import vote.pebble.common.ParseException;
import vote.pebble.common.Util;
import vote.pebble.vdf.VDF;
import vote.pebble.voting.structs.SignedBallot;
import vote.pebble.zkp.AnonCred1;
import vote.pebble.zkp.CredentialException;
import vote.pebble.zkp.CredentialSystem;
import vote.pebble.zkp.SecretCredential;

import cafe.cryptography.ed25519.Ed25519PrivateKey;

import java.io.*;
import java.security.SecureRandom;
import java.util.Base64;
import java.util.HashMap;

public class LocalFileSecretsManager implements SecretsManager {
    private final File file;
    private final String electionID;

    private static final class Secrets {
        private transient static final Base64.Decoder DECODER = Base64.getDecoder();
        private transient static final Base64.Encoder ENCODER = Base64.getEncoder();

        private String privateKey, secretCredential;
        private HashMap<String, String> ballots = new HashMap<>(), solutions = new HashMap<>();

        static Secrets load(File file) {
            try (var reader = new FileReader(file)) {
                return Util.GSON.fromJson(reader, Secrets.class);
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        }

        void save(File file) {
            try (var writer = new FileWriter(file)) {
                Util.GSON.toJson(this, writer);
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        }

        byte[] getPrivateKey() {
            return DECODER.decode(privateKey);
        }

        void setPrivateKey(byte[] pk) {
            privateKey = ENCODER.encodeToString(pk);
        }

        byte[] getSecretCredential() {
            return DECODER.decode(secretCredential);
        }

        void setSecretCredential(byte[] sc) {
            secretCredential = ENCODER.encodeToString(sc);
        }

        byte[] getBallot(String id) {
            return DECODER.decode(ballots.get(id));
        }

        void setBallot(String id, byte[] ballot) {
            ballots.put(id, ENCODER.encodeToString(ballot));
        }

        byte[] getSolution(String id) {
            return DECODER.decode(solutions.get(id));
        }

        void setSolution(String id, byte[] solution) {
            solutions.put(id, ENCODER.encodeToString(solution));
        }
    }

    public LocalFileSecretsManager(File file, String electionID) {
        this.file = file;
        this.electionID = electionID;
    }

    @Override
    public Ed25519PrivateKey getPrivateKey() {
        if (!file.exists()) {
            var random = new SecureRandom();
            var secrets = new Secrets();
            var credentialSystem = new AnonCred1();
            secrets.setPrivateKey(Ed25519PrivateKey.generate(random).toByteArray());
            try {
                secrets.setSecretCredential(credentialSystem.generateSecretCredential().toBytes());
            } catch (CredentialException e) {
                throw new RuntimeException(e);
            }
            secrets.save(file);
        }
        return Ed25519PrivateKey.fromByteArray(Secrets.load(file).getPrivateKey());
    }

    @Override
    public SecretCredential getSecretCredential(CredentialSystem system) throws CredentialException {
        return system.secretCredentialFromBytes(Secrets.load(file).getSecretCredential());
    }

    @Override
    public SignedBallot getBallot() throws ParseException {
        return SignedBallot.fromBytes(Secrets.load(file).getBallot(electionID));
    }

    @Override
    public void setBallot(SignedBallot ballot) {
        var secrets = Secrets.load(file);
        secrets.setBallot(electionID, ballot.toBytes());
        secrets.save(file);
    }

    @Override
    public VDF.Solution getVDFSolution() {
        return VDF.Solution.fromBytes(Secrets.load(file).getSolution(electionID));
    }

    @Override
    public void setVDFSolution(VDF.Solution solution) {
        var secrets = Secrets.load(file);
        secrets.setSolution(electionID, solution.toBytes());
        secrets.save(file);
    }
}
