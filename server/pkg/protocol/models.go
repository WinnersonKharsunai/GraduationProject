package protocol

type Request struct {
	Header Header
	Body   string
}

type Header struct {
	Version     string `json:"version"`
	RemoteAddr  string `json:"remoteAddr"`
	ContentType string `json:"contentType"`
	Method      string `json:"method"`
}

type Response struct {
	Error string `json:"error"`
	Body  string `json:"body"`
}
