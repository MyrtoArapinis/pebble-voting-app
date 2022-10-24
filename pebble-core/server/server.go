package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/giry-dev/pebble-voting-app/pebble-core/util"
	"github.com/giry-dev/pebble-voting-app/pebble-core/voting"
)

type SetupStatus uint8

const (
	SetupError SetupStatus = iota
	SetupInProgress
	SetupDone
)

type ServerError struct {
	StatusCode int
	Body       string
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("pebble: server error %d: %s", e.StatusCode, e.Body)
}

type SetupInfo struct {
	Status     SetupStatus
	Error      string
	BackendId  string
	Invitation string
}

type ElectionService interface {
	Create(params ElectionSetupParams) error
	Setup(adminId string) SetupInfo
	Election(backendId string) (*voting.Election, error)
}

type Server struct {
	srv          ElectionService
	create, post bool
}

func respondText(w http.ResponseWriter, statusCode int, body string) {
	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
}

func respondJson(w http.ResponseWriter, o interface{}) {
	content, err := json.Marshal(o)
	if err != nil {
		respondText(w, 500, err.Error())
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Content-Length", strconv.Itoa(len(content)))
		w.WriteHeader(200)
		w.Write(content)
	}
}

func decodeJson(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	const (
		GET  = "GET"
		POST = "POST"
	)
	ctx := context.Background()
	path := req.URL.Path
	if path == "/create" {
		if req.Method != POST {
			respondText(w, 405, "Method not allowed")
			return
		}
		if !s.create {
			respondText(w, 403, "Server does not create elections")
			return
		}
		var params ElectionSetupParams
		err := decodeJson(req.Body, &params)
		if err != nil {
			respondText(w, 400, err.Error())
			return
		}
		err = s.srv.Create(params)
		if err != nil {
			respondText(w, 500, err.Error())
			return
		}
		respondText(w, 200, "Election creation enqueued")
		return
	} else if adminId, ok := util.GetSuffix(path, "/setup/"); ok {
		if req.Method != GET {
			respondText(w, 405, "Method not allowed")
			return
		}
		info := s.srv.Setup(adminId)
		var resp struct {
			Status     string `json:"status"`
			Message    string `json:"message,omitempty"`
			BackendId  string `json:"backendId,omitempty"`
			Invitation string `json:"invitation,omitempty"`
		}
		switch info.Status {
		case SetupError:
			resp.Status = "SetupError"
			resp.Message = info.Error
		case SetupInProgress:
			resp.Status = "InProgress"
		case SetupDone:
			resp.Status = "Done"
			resp.BackendId = info.BackendId
			resp.Invitation = info.Invitation
		default:
			respondText(w, 500, "Unknown status")
			return
		}
		respondJson(w, resp)
	} else if backendId, ok := util.GetSuffix(path, "/election/"); ok {
		if req.Method != GET {
			respondText(w, 405, "Method not allowed")
			return
		}
		election, err := s.srv.Election(backendId)
		if err != nil {
			respondText(w, 500, err.Error())
			return
		}
		prog, err := election.Progress(ctx)
		if err != nil {
			respondText(w, 500, err.Error())
			return
		}
		switch prog.Phase {
		case voting.Setup:
			var resp struct {
				Status string `json:"status"`
			}
			resp.Status = "Setup"
			respondJson(w, resp)
		case voting.CredGen:
			var resp struct {
				Status string `json:"status"`
			}
			resp.Status = "CredGen"
			respondJson(w, resp)
		case voting.Cast:
			var resp struct {
				Status   string `json:"status"`
				Progress int    `json:"progress"`
				Total    int    `json:"total"`
			}
			resp.Status = "Cast"
			resp.Progress = prog.Count
			resp.Total = prog.Total
			respondJson(w, resp)
		case voting.Tally:
			var resp struct {
				Status   string         `json:"status"`
				Progress int            `json:"progress"`
				Total    int            `json:"total"`
				Counts   map[string]int `json:"counts"`
			}
			resp.Status = "Tally"
			respondJson(w, resp)
		case voting.End:
			var resp struct {
				Status string         `json:"status"`
				Valid  int            `json:"valid"`
				Total  int            `json:"total"`
				Counts map[string]int `json:"counts"`
			}
			resp.Status = "End"
			respondJson(w, resp)
		}
	} else if backendId, ok := util.GetSuffix(path, "/params/"); ok {
		if req.Method != GET {
			respondText(w, 405, "Method not allowed")
			return
		}
		election, err := s.srv.Election(backendId)
		if err != nil {
			respondText(w, 500, err.Error())
			return
		}
		body := election.Params().Bytes()
		w.Header().Add("Content-Length", strconv.Itoa(len(body)))
		w.WriteHeader(200)
		w.Write(body)
	} else if backendId, ok := util.GetSuffix(path, "/messages/"); ok {
		election, err := s.srv.Election(backendId)
		if err != nil {
			respondText(w, 500, err.Error())
			return
		}
		if req.Method == GET {
			msgs, err := election.Channel().Get(ctx)
			if err != nil {
				respondText(w, 500, err.Error())
				return
			}
			w.WriteHeader(200)
			l := []byte{0, 0}
			for _, msg := range msgs {
				p := msg.Bytes()
				if len(p) < 128 {
					l[0] = byte(len(p))
					w.Write(l[:1])
				} else {
					l[0] = byte(len(p)>>8) | 128
					l[1] = byte(len(p))
					w.Write(l)
				}
				w.Write(p)
			}
		} else if req.Method == POST {
			if !s.post {
				respondText(w, 403, "Server does not post messages")
				return
			}
			p, err := io.ReadAll(req.Body)
			if err != nil {
				respondText(w, 400, err.Error())
				return
			}
			msg, err := voting.MessageFromBytes(p)
			if err != nil {
				respondText(w, 400, err.Error())
				return
			}
			err = election.Channel().Post(ctx, msg)
			if err != nil {
				respondText(w, 500, err.Error())
			} else {
				respondText(w, 200, "Message posted")
			}
		} else {
			respondText(w, 405, "Method not allowed")
		}
	} else {
		respondText(w, 404, "Endpoint not found")
	}
}
