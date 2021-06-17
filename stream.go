package twitterstream

import (
	"net/http"
)

type (
	// UnmarshalHook is a function that will unmarshal json.
	UnmarshalHook func([]byte) (interface{}, error)

	// IStream is the interface that the stream struct implements.
	IStream interface {
		StartStream(queryParams string) error
		StopStream()
		GetMessages() <-chan StreamMessage
		SetUnmarshalHook(hook UnmarshalHook)
	}

	// StreamMessage is the message that is sent from the messages channel.
	StreamMessage struct {
		Data interface{}
		Err  error
	}

	// Stream is the struct that manages a long running TCP connection with Twitter.
	// It accepts an 'unamarshalHook' which allows you to unmarshal json in a thread-safe manner.
	// It is highly encouraged to set a unmarshal hook before starting a stream. Unmarshaling json
	// in a separate goroutine is not recommended because the Go bytes.Buffer is not thread safe.
	Stream struct {
		unmarshalHook UnmarshalHook
		messages      chan StreamMessage
		httpClient    IHttpClient
		done          chan struct{}
		reader        IStreamResponseBodyReader
	}
)

func newStream(httpClient IHttpClient, reader IStreamResponseBodyReader) *Stream {
	return &Stream{
		unmarshalHook: func(bytes []byte) (interface{}, error) {
			return bytes, nil
		},
		messages:   make(chan StreamMessage),
		done:       make(chan struct{}),
		reader:     reader,
		httpClient: httpClient,
	}
}

// SetUnmarshalHook sets the function that unmarshals json. It is highly encouraged
// that you unmarshal json with this hook to promote thread-safety. Go's bytes.Buffer is not
// thread safe and can result in panics when a bytes.Buffer is shared across goroutines.
func (s *Stream) SetUnmarshalHook(hook UnmarshalHook) {
	s.unmarshalHook = hook
}

// GetMessages returns the read-only messages channel
func (s *Stream) GetMessages() <-chan StreamMessage {
	return s.messages
}

// StopStream sends a close signal to stop the stream of tweets.
func (s *Stream) StopStream() {
	close(s.done)
}

// StartStream makes an HTTP GET request to twitter and starts streaming tweets to the Messages channel.
// Accepts query params described in GET /2/tweets/search/stream to expand the payload that is returned. Query params string must begin with a ?.
// See available query params here https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/api-reference/get-tweets-search-stream.
// See an example here: https://developer.twitter.com/en/docs/twitter-api/expansions.
func (s *Stream) StartStream(optionalQueryParams string) error {
	res, err := s.httpClient.getSearchStream(optionalQueryParams)

	if err != nil {
		return err
	}

	s.reader.setStreamResponseBody(res.Body)

	go s.streamMessages(res)

	return nil
}

func (s *Stream) streamMessages(res *http.Response) {
	defer res.Body.Close()

	for !stopped(s.done) {
		b, err := s.reader.readNext()
		if err != nil {
			s.messages <- StreamMessage{
				Data: nil,
				Err:  err,
			}
			s.StopStream()
			break
		}
		if len(b) == 0 {
			// empty keep-alive
			continue
		}

		data, err := s.unmarshalHook(b)

		s.messages <- StreamMessage{
			Data: data,
			Err:  err,
		}
	}
}
