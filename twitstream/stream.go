package twitstream

import (
	"net/http"
	"sync"
)

type (
	IStream interface {
		StartStream()
		StopStream()
		GetMessages() *chan interface{}
	}

	stream struct {
		Messages   chan interface{}
		httpClient IHttpClient
		done       chan struct{}
		group      *sync.WaitGroup
		reader     IStreamResponseBodyReader
	}
)

func newStream(httpClient IHttpClient, reader IStreamResponseBodyReader) *stream {
	return &stream{
		Messages:   make(chan interface{}),
		done:       make(chan struct{}),
		group:      new(sync.WaitGroup),
		reader:     reader,
		httpClient: httpClient,
	}
}

func (s *stream) GetMessages() *chan interface{} {
	return &s.Messages
}

// StopStream sends a close signal to stop the stream of tweets
func (s *stream) StopStream() {
	close(s.done)
}

func (s *stream) StartStream() {

	res, err := s.httpClient.newHttpRequest(&requestOpts{
		Method: "GET",
		Url:    endpoints["stream"],
	})

	if err != nil {
		panic(err)
	}

	s.reader.setStreamResponseBody(res.Body)

	s.group.Add(1)
	go s.streamMessages(res)
}

func (s *stream) streamMessages(res *http.Response) {
	defer res.Body.Close()
	defer s.group.Done()

	for !stopped(s.done) {
		data, err := s.reader.readNext()
		if err != nil {
			return
		}
		if len(data) == 0 {
			// empty keep-alive
			continue
		}

		m := string(data)
		// TODO send data or error here
		s.Messages <- m
	}
}
