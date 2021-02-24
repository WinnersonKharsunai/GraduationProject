package protocol

// Request ...
type Request struct {
	Header Header `json:"header"`
	Body   string `json:"body"`
}

// Header ...
type Header struct {
	Version     string `json:"version"`
	RemoteAddr  string `json:"remoteAddr"`
	ContentType string `json:"contentType"`
	Method      string `json:"method"`
}

// Response ...
type Response struct {
	Error string `json:"error"`
	Body  []byte `json:"body"`
}
