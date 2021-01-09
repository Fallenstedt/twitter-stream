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
		GetMessages() *chan interface{}
	}

	stream struct {
		messages   chan interface{}
		httpClient IHttpClient
		done       chan struct{}
		group      *sync.WaitGroup
		reader     IStreamResponseBodyReader
	}
)

func newStream(httpClient IHttpClient, reader IStreamResponseBodyReader) *stream {
	return &stream{
		messages:   make(chan interface{}),
		done:       make(chan struct{}),
		group:      new(sync.WaitGroup),
		reader:     reader,
		httpClient: httpClient,
	}
}

// GetMessages returns the messages channel.
func (s *stream) GetMessages() *chan interface{} {
	return &s.messages
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
			s.messages <- err
			s.StopStream()
			break
		}
		if len(data) == 0 {
			// empty keep-alive
			continue
		}

		m := string(data)
		s.messages <- m
	}
}
