package protocol

// SetHeader set the header value to be sent to server
func SetHeader(version, contentType, method, addr string) Header {
	return Header{
		Version:     version,
		ContentType: contentType,
		RemoteAddr:  addr,
		Method:      method,
	}
}
