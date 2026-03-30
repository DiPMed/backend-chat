package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dipmed/backend-chat/internal/llm"
	"github.com/dipmed/backend-chat/internal/sessions"
)

type Server struct {
	provider llm.Provider
	store    sessions.Store
	mux      *http.ServeMux
}

func NewServer(provider llm.Provider, store sessions.Store) *Server {
	s := &Server{provider: provider, store: store, mux: http.NewServeMux()}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("POST /chat", s.handleChat)
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		http.Error(w, "message required", http.StatusBadRequest)
		return
	}

	session, err := s.resolveSession(req.SessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	userMsg := sessions.Message{Role: "user", Content: req.Message}
	if err := s.store.Append(session.ID, userMsg); err != nil {
		http.Error(w, "failed to store message", http.StatusInternalServerError)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// First chunk includes the session ID so the client can track it.
	meta, _ := json.Marshal(chatChunk{SessionID: session.ID})
	fmt.Fprintf(w, "%s\n", meta)
	flusher.Flush()

	chatReq := &llm.ChatRequest{Messages: session.Messages}

	var fullResponse string
	err = s.provider.ChatStream(r.Context(), chatReq, func(chunk string) error {
		fullResponse += chunk
		data, marshalErr := json.Marshal(chatChunk{Content: chunk})
		if marshalErr != nil {
			return marshalErr
		}
		_, writeErr := fmt.Fprintf(w, "%s\n", data)
		if writeErr != nil {
			return writeErr
		}
		flusher.Flush()
		return nil
	})

	if err != nil {
		log.Printf("chat stream error: %v", err)
		errData, _ := json.Marshal(chatChunk{Error: "internal error"})
		fmt.Fprintf(w, "%s\n", errData)
		flusher.Flush()
		return
	}

	s.store.Append(session.ID, sessions.Message{Role: "assistant", Content: fullResponse})

	doneData, _ := json.Marshal(chatChunk{Done: true})
	fmt.Fprintf(w, "%s\n", doneData)
	flusher.Flush()
}

func (s *Server) resolveSession(id string) (*sessions.Session, error) {
	if id == "" {
		return s.store.Create()
	}
	return s.store.Get(id)
}
