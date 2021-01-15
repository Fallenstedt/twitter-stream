package twitterstream

import (
	"net/http"
	"sync"
)

type (
	// IStream is the interface that the stream struct implements.
	IStream interface {
		StartStream() error
		StopStream()
		GetMessages() <-chan StreamMessage
	}

	stream struct {
		messages   chan StreamMessage
		httpClient IHttpClient
		done       chan struct{}
		group      *sync.WaitGroup
		reader     IStreamResponseBodyReader
	}

	// StreamMessage is the message that is sent from the messages channel.
	StreamMessage struct {
		Data []byte
		Err error
	}
)

func newStream(httpClient IHttpClient, reader IStreamResponseBodyReader) *stream {
	return &stream{
		messages:   make(chan StreamMessage),
		done:       make(chan struct{}),
		group:      new(sync.WaitGroup),
		reader:     reader,
		httpClient: httpClient,
	}
}

// GetMessages returns the read-only messages channel
func (s *stream) GetMessages() <-chan StreamMessage {
	return s.messages
}

// StopStream sends a close signal to stop the stream of tweets.
func (s *stream) StopStream() {
	close(s.done)
}

// StartStream makes an HTTP request to twitter and starts streaming tweets to the Messages channel.
func (s *stream) StartStream() error {

	res, err := s.httpClient.newHttpRequest(&requestOpts{
		Method: "GET",
		Url:    endpoints["stream"],
	})

	if err != nil {
		return err
	}

	s.reader.setStreamResponseBody(res.Body)

	s.group.Add(1)
	go s.streamMessages(res)
	return nil
}

func (s *stream) streamMessages(res *http.Response) {
	defer res.Body.Close()
	defer s.group.Done()

	for !stopped(s.done) {
		data, err := s.reader.readNext()
		if err != nil {
			s.messages <- StreamMessage{
				Data: nil,
				Err: err,
			}
			s.StopStream()
			break
		}
		if len(data) == 0 {
			// empty keep-alive
			continue
		}

		s.messages <- StreamMessage{
			Data: data,
			Err: nil,
		}
	}
}
