package internal

type chatChunk struct {
	SessionID string `json:"session_id,omitempty"`
	Content   string `json:"content,omitempty"`
	Error     string `json:"error,omitempty"`
	Done      bool   `json:"done,omitempty"`
}

type chatRequest struct {
	SessionID string `json:"session_id,omitempty"`
	Message   string `json:"message"`
}
