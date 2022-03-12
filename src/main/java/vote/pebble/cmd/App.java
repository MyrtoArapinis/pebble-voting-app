package vote.pebble.cmd;

import vote.pebble.common.Hex;
import vote.pebble.common.ParseException;
import vote.pebble.common.PebbleException;
import vote.pebble.voting.BallotNotDecryptedException;
import vote.pebble.voting.Election;
import vote.pebble.voting.LocalFileSecretsManager;
import vote.pebble.voting.broadcast.BroadcastClient;
import vote.pebble.voting.methods.InvalidNumberOfChoicesException;
import vote.pebble.voting.methods.NoSuchVotingMethodException;
import vote.pebble.voting.structs.ElectionParams;
import vote.pebble.voting.structs.EligibilityList;

import cafe.cryptography.curve25519.InvalidEncodingException;
import cafe.cryptography.ed25519.Ed25519PublicKey;
import vote.pebble.zkp.CredentialException;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.FileWriter;
import java.io.IOException;
import java.time.Instant;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Scanner;

public final class App {
    private final File secretsFile = new File("secrets.json");
    private String baseURI;
    private ElectionParams electionParams;
    private Election election;

    private final String[] args;

    private App(String[] args) {
        this.args = args;
    }

    public static void main(String[] args) {
        var app = new App(args);
        app.run();
    }

    private void run() {
        if (args.length < 1)
            usage();
        if (args[0].equals("pubkey")) {
            showPubkey();
        } else if (args[0].equals("election")) {
            if (args.length < 2)
                usage();
            loadBaseURI();
            if (args[1].equals("id")) {
                electionID();
            } else if (args[1].equals("create")) {
                createElection();
            } else {
                try {
                    setElection(args[1].equals("info"));
                } catch (PebbleException e) {
                    e.printStackTrace();
                    System.exit(3);
                }
                if (args[1].equals("info")) {
                    showElectionInfo();
                } else if (args[1].equals("join")) {
                    try {
                        election.postCredential();
                    } catch (CredentialException e) {
                        e.printStackTrace();
                        System.exit(3);
                    }
                    System.out.println("Anonymous credential submitted");
                } else if (args[1].equals("vote")) {
                    vote();
                } else if (args[1].equals("reveal")) {
                    election.revealBallotDecryption();
                } else if (args[1].equals("open")) {
                    openBallot();
                } else if (args[1].equals("results")) {
                    showElectionResults();
                } else {
                    usage();
                }
            }
        }
    }

    private static void usage() {
        System.out.print("Pebble voting app\n"
                + "Available commands:\n\n"
                + "  pubkey                Show your public key\n"
                + "  election id ID        Sets the current election's ID\n"
                + "  election create       Create an election\n"
                + "  election info         Shows information about the current election\n"
                + "  election join         Join the current election by posting your anonymous public credential\n"
                + "  election vote CHOICE  Vote in the current election\n"
                + "  election reveal       Reveal the decryption of your ballot\n"
                + "  election open         Decrypt one ballot\n"
                + "  election results      Show the current election's results\n"
        );
        System.exit(1);
    }

    private void loadBaseURI() {
        try (var scanner = new Scanner(new File("server.txt"))) {
            if (scanner.hasNext()) {
                baseURI = scanner.next();
                return;
            }
        } catch (FileNotFoundException ignored) {}
        System.out.println("server.txt file needed");
        System.exit(2);
    }

    private void showPubkey() {
        var secretsManager = new LocalFileSecretsManager(secretsFile, null);
        var pubkey = secretsManager.getPrivateKey().derivePublic().toByteArray();
        System.out.println(Hex.encodeHexString(pubkey));
    }

    private void electionID() {
        if (args.length < 3) {
            System.out.println("Supply election ID as argument");
            System.exit(3);
            return;
        }
        try (var writer = new FileWriter("election.txt")) {
            writer.write(args[2]);
        } catch (IOException e) {
            e.printStackTrace();
            System.exit(3);
        }
    }

    private void createElection() {
        ElectionParams electionParams;
        String input;
        try (var scanner = new Scanner(System.in)) {
            scanner.useDelimiter("\n");
            System.out.print("Vote start: ");
            var voteStart = Instant.parse(scanner.nextLine());
            System.out.print("Tally start: ");
            var tallyStart = Instant.parse(scanner.nextLine());
            System.out.print("VDF difficulty: ");
            var diff = scanner.nextLong();
            scanner.nextLine();
            System.out.println("Choices:");
            var choices = new ArrayList<String>();
            while (!(input = scanner.nextLine()).isEmpty()) {
                choices.add(input);
            }
            System.out.println("Eligible voters' public keys:");
            var eligibilityList = new EligibilityList();
            while (!(input = scanner.nextLine()).isEmpty()) {
                try {
                    eligibilityList.add(Ed25519PublicKey.fromByteArray(Hex.decodeHexString(input)), null);
                } catch (IllegalArgumentException | InvalidEncodingException e) {
                    System.out.println("^^^ Ignoring malformed public key");
                }
            }
            var choicesArray = new String[choices.size()];
            electionParams = new ElectionParams(eligibilityList, voteStart, tallyStart, diff,
                    "Plurality", choices.toArray(choicesArray));
        }
        var electionID = BroadcastClient.createElection(baseURI, electionParams);
        System.out.printf("Election ID: %s\n", electionID);
    }

    private void setElection(boolean observer) throws ParseException, NoSuchVotingMethodException, InvalidNumberOfChoicesException {
        String id;
        try (var scanner = new Scanner(new File("election.txt"))) {
            id = scanner.next();
        } catch (FileNotFoundException e) {
            System.out.println("server.txt file needed");
            System.out.println("Set election ID to create it");
            System.exit(2);
            return;
        }
        var client = new BroadcastClient(baseURI, id);
        electionParams = client.getElectionParams();
        if (observer) {
            election = new Election(electionParams, client);
        } else {
            var secretsManager = new LocalFileSecretsManager(secretsFile, id);
            election = new Election(electionParams, client, secretsManager);
        }
    }

    private void showElectionInfo() {
        System.out.println("Vote start: " + electionParams.voteStart);
        System.out.println("Tally start: " + electionParams.tallyStart);
        System.out.println("VDF difficulty: " + electionParams.vdfDifficulty);
        System.out.println("Voting method: " + electionParams.votingMethod);
        System.out.println("Choices:");
        for (int i = 0; i < electionParams.choices.length; i++)
            System.out.printf("%2d. %s\n", i + 1, electionParams.choices[i]);
    }

    private void vote() {
        var scanner = new Scanner(System.in);
        if (args.length < 3) {
            System.out.println("Supply your choice number as argument");
            System.exit(1);
        }
        int choice = Integer.parseInt(args[2]) - 1;
        if (choice < 0 || choice >= electionParams.choices.length) {
            System.out.println("Choice number out of range");
            System.exit(1);
        }
        System.out.printf("Your choice:\n%2d. %s\n", choice + 1, electionParams.choices[choice]);
        if (args.length < 4 || !args[3].equals("y")) {
            System.out.print("Confirm (y/n): ");
            if (!scanner.nextLine().equalsIgnoreCase("y")) {
                System.out.println("Aborting");
                return;
            }
        }
        try {
            election.vote(electionParams.choices[choice]);
        } catch (CredentialException e) {
            e.printStackTrace();
            System.exit(3);
        }
        System.out.println("Ballot submitted");
    }

    private void openBallot() {
        try {
            election.tally();
            System.out.println("All ballots have been decrypted");
        } catch (CredentialException e) {
            e.printStackTrace();
        } catch (BallotNotDecryptedException e) {
            System.out.println("Decrypting ballot...");
            election.postBallotDecryption(e.encryptedBallot);
            System.out.println("Ballot decryption posted");
        }
    }

    private void showElectionResults() {
        try {
            var tally = election.tally();
            Collections.sort(tally);
            for (var item : tally)
                System.out.printf("[%4d] %2d. %s\n", item.count, item.index + 1, electionParams.choices[item.index]);
        } catch (CredentialException e) {
            e.printStackTrace();
        } catch (BallotNotDecryptedException e) {
            System.out.println("Some of the ballots have not been decrypted");
        }
    }
}
