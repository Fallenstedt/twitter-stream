package httpclient

type RequestOpts struct {
	Retries uint8
	Method  string
	Url     string
	Body    string
	Headers []struct {
		Key   string
		Value string
	}
}
