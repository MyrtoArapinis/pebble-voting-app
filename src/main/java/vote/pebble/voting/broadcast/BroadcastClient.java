package vote.pebble.voting.broadcast;

import vote.pebble.common.ParseException;
import vote.pebble.common.Util;
import vote.pebble.voting.structs.*;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.util.ArrayList;
import java.util.Base64;

public class BroadcastClient implements BroadcastChannel {
    private static final String MIME_BINARY = "application/octet-stream";
    private static final String MIME_JSON = "application/json";

    private transient static final Base64.Decoder DECODER = Base64.getDecoder();
    private transient static final Base64.Encoder ENCODER = Base64.getEncoder();

    private final URI electionURI, credentialsURI, ballotsURI, decryptionsURI;

    private static final class ReceivedMessage {

        private String kind;
        private String content;

        public byte[] getContent() {
            return DECODER.decode(content);
        }
    }

    private static final class ReceivedMessages {
        public ReceivedMessage[] messages;

        public static ReceivedMessages fromString(String s) {
            return Util.fromJson(s, ReceivedMessages.class);
        }
    }

    private static final class ElectionJson {
        private String eligibilityList;
        private String voteStart, tallyStart;
        private long vdfDifficulty;
        private String votingMethod;
        private String[] choices;

        public byte[] getEligibilityList() {
            return DECODER.decode(eligibilityList);
        }

        public void setEligibilityList(byte[] bytes) {
            eligibilityList = ENCODER.encodeToString(bytes);
        }

        public Instant getVoteStart() {
            return Instant.parse(voteStart);
        }

        public void setVoteStart(Instant time) {
            voteStart = time.toString();
        }

        public Instant getTallyStart() {
            return Instant.parse(tallyStart);
        }

        public void setTallyStart(Instant time) {
            tallyStart = time.toString();
        }
    }

    public BroadcastClient(String baseURI, String id) {
        if (!baseURI.endsWith("/"))
            baseURI += '/';
        electionURI = URI.create(baseURI + "election?id=" + id);
        final var messagesURI = baseURI + "messages?id=" + id;
        credentialsURI = URI.create(messagesURI + "&kind=credential");
        ballotsURI = URI.create(messagesURI + "&kind=ballot");
        decryptionsURI = URI.create(messagesURI + "&kind=decryption");
    }

    private static String httpGet(URI uri) {
        try {
            var request = HttpRequest.newBuilder(uri).GET().build();
            var client = HttpClient.newHttpClient();
            var response = client.send(request, HttpResponse.BodyHandlers.ofString());
            if (response.statusCode() != 200)
                throw new RuntimeException("HTTP status " + response.statusCode() + "\n" + response.body());
            return response.body();
        } catch (IOException | InterruptedException e) {
            throw new RuntimeException(e);
        }
    }

    private static String httpPost(URI uri, String contentType, byte[] body) {
        try {
            var request = HttpRequest.newBuilder(uri)
                    .header("Content-Type", contentType)
                    .POST(HttpRequest.BodyPublishers.ofByteArray(body))
                    .build();
            var client = HttpClient.newHttpClient();
            var response = client.send(request, HttpResponse.BodyHandlers.ofString());
            if (response.statusCode() != 200)
                throw new RuntimeException("HTTP status " + response.statusCode() + "\n" + response.body());
            return response.body();
        } catch (IOException | InterruptedException e) {
            throw new RuntimeException(e);
        }
    }

    public static String createElection(String baseURI, ElectionParams params) {
        var jsonParams = new ElectionJson();
        jsonParams.setEligibilityList(params.eligibilityList.toBytes());
        jsonParams.setVoteStart(params.voteStart);
        jsonParams.setTallyStart(params.tallyStart);
        jsonParams.vdfDifficulty = params.vdfDifficulty;
        jsonParams.votingMethod = params.votingMethod;
        jsonParams.choices = params.choices;
        if (!baseURI.endsWith("/"))
            baseURI += '/';
        return httpPost(URI.create(baseURI + "election"), MIME_JSON, Util.toJson(jsonParams).getBytes(StandardCharsets.UTF_8));
    }

    public ElectionParams getElectionParams() throws ParseException {
        var response = httpGet(electionURI);
        var params = Util.fromJson(response, ElectionJson.class);
        return new ElectionParams(
                EligibilityList.fromBytes(params.getEligibilityList()),
                params.getVoteStart(), params.getTallyStart(),
                params.vdfDifficulty, params.votingMethod, params.choices);
    }

    @Override
    public void postCredential(CredentialMessage msg) {
        httpPost(credentialsURI, MIME_BINARY, msg.toBytes());
    }

    @Override
    public void postSignedBallot(SignedBallot msg) {
        httpPost(ballotsURI, MIME_BINARY, msg.toBytes());
    }

    @Override
    public void postBallotDecryption(DecryptionMessage msg) {
        httpPost(decryptionsURI, MIME_BINARY, msg.toBytes());
    }

    private interface MessageDecoder<T> {
        T decode(byte[] bytes) throws ParseException;
    }

    private static <T> Iterable<T> getMessages(URI uri, MessageDecoder<T> decoder) {
        var recv = ReceivedMessages.fromString(httpGet(uri));
        if (recv == null || recv.messages == null)
            return new ArrayList<T>();
        var result = new ArrayList<T>(recv.messages.length);
        for (var msg : recv.messages) {
            try {
                result.add(decoder.decode(msg.getContent()));
            } catch (ParseException ignored) {}
        }
        return result;
    }

    @Override
    public Iterable<CredentialMessage> getCredentials() {
        return getMessages(credentialsURI, CredentialMessage::fromBytes);
    }

    @Override
    public Iterable<SignedBallot> getSignedBallots() {
        return getMessages(ballotsURI, SignedBallot::fromBytes);
    }

    @Override
    public Iterable<DecryptionMessage> getBallotDecryptions() {
        return getMessages(decryptionsURI, DecryptionMessage::fromBytes);
    }
}
