package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

type ServerError struct {
	StatusCode int
	Body       string
}

type ElectionParams struct {
	EligibilityList       []byte
	VoteStart, TallyStart time.Time
	VdfDifficulty         uint64
	VotingMethod          string
	Choices               []string
}

type jsonElection struct {
	EligibilityList string   `json:"eligibilityList"`
	VoteStart       string   `json:"voteStart"`
	TallyStart      string   `json:"tallyStart"`
	VdfDifficulty   uint64   `json:"vdfDifficulty"`
	VotingMethod    string   `json:"votingMethod"`
	Choices         []string `json:"choices"`
}

func (election *ElectionParams) toJson() (result jsonElection) {
	result.EligibilityList = base64.StdEncoding.EncodeToString(election.EligibilityList)
	result.VoteStart = election.VoteStart.UTC().Format(time.RFC3339)
	result.TallyStart = election.TallyStart.UTC().Format(time.RFC3339)
	result.VdfDifficulty = election.VdfDifficulty
	result.VotingMethod = election.VotingMethod
	result.Choices = election.Choices
	return
}

func (election *ElectionParams) fromBytes(bytes []byte) error {
	var request jsonElection
	err := json.Unmarshal(bytes, &request)
	if err != nil {
		return err
	}
	election.EligibilityList, err = base64.StdEncoding.DecodeString(request.EligibilityList)
	if err != nil {
		return err
	}
	election.VoteStart, err = time.Parse(time.RFC3339, request.VoteStart)
	if err != nil {
		return err
	}
	election.TallyStart, err = time.Parse(time.RFC3339, request.TallyStart)
	if err != nil {
		return err
	}
	election.VdfDifficulty = request.VdfDifficulty
	election.VotingMethod = request.VotingMethod
	election.Choices = request.Choices
	return nil
}

type Message struct {
	Kind    string
	Content []byte
}

type Server interface {
	GetElection(id string) (*ElectionParams, *ServerError)
	CreateElection(params ElectionParams) (string, *ServerError)
	GetMessages(id string, kind string) ([]Message, *ServerError)
	PostMessage(id string, msg Message) *ServerError
}

type Handler struct {
	Server Server
}

func respondText(w http.ResponseWriter, statusCode int, body string) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
}

func respondJson(w http.ResponseWriter, o interface{}) {
	content, err := json.Marshal(o)
	if err != nil {
		respondText(w, 500, err.Error())
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.Write(content)
	}
}

func readRequestBody(req *http.Request) ([]byte, error) {
	body := make([]byte, req.ContentLength)
	if n, err := req.Body.Read(body); err != nil && n != len(body) {
		return nil, err
	}
	return body, nil
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/election":
		switch req.Method {
		case "GET":
			election, srvErr := handler.Server.GetElection(req.URL.Query().Get("id"))
			if srvErr != nil {
				respondText(w, srvErr.StatusCode, srvErr.Body)
			} else {
				respondJson(w, election.toJson())
			}
		case "POST":
			content, err := readRequestBody(req)
			if err != nil {
				respondText(w, 500, err.Error())
				return
			}
			var election ElectionParams
			err = election.fromBytes(content)
			if err != nil {
				respondText(w, 500, err.Error())
				return
			}
			id, srvErr := handler.Server.CreateElection(election)
			if srvErr != nil {
				respondText(w, srvErr.StatusCode, srvErr.Body)
			} else {
				respondText(w, 200, id)
			}
		default:
			respondText(w, 405, "Method not allowed")
		}
	case "/messages":
		id := req.URL.Query().Get("id")
		kind := req.URL.Query().Get("kind")
		switch req.Method {
		case "GET":
			messages, srvErr := handler.Server.GetMessages(id, kind)
			if srvErr != nil {
				respondText(w, srvErr.StatusCode, srvErr.Body)
			} else {
				type jsonMessage struct {
					Kind    string `json:"kind"`
					Content string `json:"content"`
				}
				var responce struct {
					Messages []jsonMessage `json:"messages"`
				}
				for _, msg := range messages {
					responce.Messages = append(responce.Messages,
						jsonMessage{msg.Kind, base64.StdEncoding.EncodeToString(msg.Content)})
				}
				respondJson(w, responce)
			}
		case "POST":
			content, err := readRequestBody(req)
			if err != nil {
				respondText(w, 500, err.Error())
				return
			}
			srvErr := handler.Server.PostMessage(id, Message{kind, content})
			if srvErr != nil {
				respondText(w, srvErr.StatusCode, srvErr.Body)
			} else {
				respondText(w, 200, "")
			}
		default:
			respondText(w, 405, "Method not allowed")
		}
	default:
		respondText(w, 404, "Endpoint not available")
	}
}

func main() {
	handler := &Handler{NewMockServer()}
	http.ListenAndServe(":8090", handler)
}
