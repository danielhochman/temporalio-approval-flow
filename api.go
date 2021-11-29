package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/danielhochman/temporalio-approval-flow/workflow"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.temporal.io/api/workflowservice/v1"
	temporalclient "go.temporal.io/sdk/client"
)

//go:embed frontend/dist/*
var content embed.FS

func main() {
	temporal, err := temporalclient.NewClient(temporalclient.Options{})
	if err != nil {
		panic(err)
	}

	s := &Server{
		client: temporal,
	}

	r := mux.NewRouter()

	// Frontend
	r.Handle("/", http.HandlerFunc(AssetHandler))
	r.PathPrefix("/index").Handler(http.HandlerFunc(AssetHandler))

	// Backend
	r.Handle("/api/ExecuteWorkflow", http.HandlerFunc(s.ExecuteWorkflow)).Methods(http.MethodPost)
	r.Handle("/api/ListOpenWorkflow", http.HandlerFunc(s.ListOpenWorkflow)).Methods(http.MethodGet)
	r.Handle("/api/AddReview", http.HandlerFunc(s.AddReview)).Methods(http.MethodPost)

	httpServer := &http.Server{
		Addr:         ":9000",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}

type Server struct {
	client temporalclient.Client
}

type Review struct {
	Time    string `json:"time"`
	User    string `json:"user"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type Workflow struct {
	StartTime string   `json:"startTime"`
	ID        string   `json:"id"`
	RunID     string   `json:"runId"`
	User      string   `json:"user"`
	Action    string   `json:"action"`
	Reviews   []Review `json:"reviews"`
	Status    string   `json:"status"`
}

type ListOpenWorkflowResponse struct {
	Workflows []Workflow `json:"workflows"`
}

func (s *Server) ListOpenWorkflow(w http.ResponseWriter, r *http.Request) {
	res, err := s.client.ListOpenWorkflow(r.Context(), &workflowservice.ListOpenWorkflowExecutionsRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ret := &ListOpenWorkflowResponse{
		Workflows: make([]Workflow, len(res.Executions)),
	}
	for idx, exec := range res.Executions {
		resp, err := s.client.QueryWorkflow(r.Context(), exec.Execution.WorkflowId, exec.Execution.RunId, "getState")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var result workflow.State
		if err := resp.Get(&result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

		}

		reviews := make([]Review, len(result.Reviews))
		for ridx, r := range result.Reviews {
			reviews[ridx] = Review{
				Time:    r.Timestamp.String()[:19],
				Message: r.Message,
				User:    r.Author,
				Status:  r.Status.String(),
			}
		}
		status := "open"
		if result.IsApproved() {
			status = "approve"
		} else if result.IsLocked() {
			status = "locked"
		}

		ret.Workflows[idx] = Workflow{
			StartTime: exec.StartTime.String()[:19],
			ID:        exec.Execution.WorkflowId,
			RunID:     exec.Execution.RunId,
			User:      result.User,
			Action:    result.Action,
			Reviews:   reviews,
			Status:    status,
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)
}

type ExecuteWorkflowRequest struct {
	Time   string `json:"time"`
	User   string `json:"user"`
	Action string `json:"action"`
}

func (s *Server) ExecuteWorkflow(w http.ResponseWriter, r *http.Request) {
	in := &ExecuteWorkflowRequest{}
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialize the state.
	state := &workflow.State{
		Action: in.Action,
		User:   in.User,
	}

	opts := temporalclient.StartWorkflowOptions{TaskQueue: workflow.QueueName}
	we, err := s.client.ExecuteWorkflow(r.Context(), opts, workflow.Workflow, state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&Workflow{
		ID:    we.GetID(),
		RunID: we.GetRunID(),
	})
}

type AddReviewRequest struct {
	ID      string `json:"id"`
	RunID   string `json:"runId"`
	Action  string `json:"action"`
	Message string `json:"message"`
}

func (s *Server) AddReview(w http.ResponseWriter, r *http.Request) {
	in := &AddReviewRequest{}
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var status workflow.Status
	switch in.Action {
	case "approve":
		status = workflow.Approve
	case "comment":
		status = workflow.Comment
	case "lock":
		status = workflow.Lock
	case "unlock":
		status = workflow.Unlock
	default:
		http.Error(w, fmt.Sprintf("did not recognize action '%s'", in.Action), http.StatusBadRequest)
		return
	}

	review := &workflow.Review{
		Timestamp: time.Now(),
		Author:    "Anonymous",
		Message:   in.Message,
		Status:    status,
	}

	err := s.client.SignalWorkflow(r.Context(), in.ID, in.RunID, workflow.ReviewChannel, review)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AssetHandler(w http.ResponseWriter, r *http.Request) {
	var filename string
	if r.URL.Path == "/" {
		filename = "frontend/dist/index.html"
	} else {
		w.Header().Set("Content-Type", "text/javascript")
		filename = fmt.Sprintf("frontend/dist%s", r.URL.Path)
	}

	f, err := content.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}
