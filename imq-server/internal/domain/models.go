package domain

<<<<<<< HEAD
// Message is use to hold data sent by client
=======
type Topic struct {
	name string
	msg  string
}

>>>>>>> 9fe39465475b121a78fe3f5e4b7a5638b6c0a469
type Message struct {
	MessageID string
	Data      string
	CretedAt  string
	ExpiresAt string
}
