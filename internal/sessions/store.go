package sessions

type Store interface {
	Create() (*Session, error)
	Get(id string) (*Session, error)
	Append(id string, msgs ...Message) error
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Session struct {
	ID       string    `json:"id"`
	Messages []Message `json:"messages"`
}
