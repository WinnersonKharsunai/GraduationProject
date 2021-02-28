package protocol

// Request is accepted request type for IMQ
type Request struct {
	Header Header `json:"header"`
	Body   string `json:"body"`
}

// Header is accepted header type for IMQ
type Header struct {
	Version     string `json:"version"`
	RemoteAddr  string `json:"remoteAddr"`
	ContentType string `json:"contentType"`
	Method      string `json:"method"`
}

// Response is accepted response type for IMQ
type Response struct {
	Error string `json:"error"`
	Body  []byte `json:"body"`
}
