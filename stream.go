package twitterstream

import (
	"net/http"
	"sync"
)

type (
	// IStream is the interface that the stream struct implements.
	IStream interface {
		StartStream(queryParams string) error
		StopStream()
		GetMessages() <-chan StreamMessage
	}

	// StreamMessage is the message that is sent from the messages channel.
	StreamMessage struct {
		Data []byte
		Err  error
	}


	stream struct {
		messages   chan StreamMessage
		httpClient IHttpClient
		done       chan struct{}
		group      *sync.WaitGroup
		reader     IStreamResponseBodyReader
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

// StartStream makes an HTTP GET request to twitter and starts streaming tweets to the Messages channel using Server Sent Events.
// Accepts query params described in GET /2/tweets/search/stream to expand the payload that is returned. Query params string must begin with a ?.
// See available query params here https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/api-reference/get-tweets-search-stream.
// See an example here: https://developer.twitter.com/en/docs/twitter-api/expansions.
func (s *stream) StartStream(optionalQueryParams string) error {
	res, err := s.httpClient.getSearchStream(optionalQueryParams)

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
				Err:  err,
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
			Err:  nil,
		}
	}
}
