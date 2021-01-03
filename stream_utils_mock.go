package twitter_stream

import (
	"io"
)

type mockStreamResponseBodyReader struct {
	MockReadNext func() ([]byte, error)
	MockSetStreamResponseBody func(body io.Reader)
}

func (m mockStreamResponseBodyReader) readNext() ([]byte, error) {
	return m.MockReadNext()
}

func (m mockStreamResponseBodyReader) setStreamResponseBody(body io.Reader) {
	m.MockSetStreamResponseBody(body)
}
