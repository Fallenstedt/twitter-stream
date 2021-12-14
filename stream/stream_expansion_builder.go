package stream


type (
	IStreamQueryParamsBuilder interface {
		AddExpansion(expansion string) *StreamQueryParamBuilder
		AddMediaField(mediaField string) *StreamQueryParamBuilder
		AddPlaceField(placeField string) *StreamQueryParamBuilder
		Build() string
	}

	StreamQueryParamBuilder struct {
		expansions []*string
		mediaFields []*string
		placeFields []*string
	}

)

func NewStreamQueryParamsBuilder() IStreamQueryParamsBuilder {
	return &StreamQueryParamBuilder{
		expansions: []*string{},
		mediaFields: []*string{},
		placeFields: []*string{},
	}
}

func (s *StreamQueryParamBuilder) Build() string {
	return "?expansions=lol"
}

// AddExpansion adds an expansion defined in https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/api-reference/get-tweets-search-stream.
// With expansions, developers can expand objects referenced in the payload. Objects available for expansion are referenced by ID.
// Add a single expansion for each invoke of `AddExpansion`.
func (s *StreamQueryParamBuilder) AddExpansion(expansion string) *StreamQueryParamBuilder {
	s.expansions = append(s.expansions, &expansion)
	return s
}

// AddMediaField adds a media field which enables you to select which specific media fields will deliver in each returned tweet.
// The Tweet will only return media fields if the Tweet contains media and if you've also included `AddExpansion("attachments.media_keys")`.
// Learn more about media fields on twitter docs https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/api-reference/get-tweets-search-stream.
// Add a single media field for each invoke of `AddMediaField`.
func (s *StreamQueryParamBuilder) AddMediaField(mediaField string) *StreamQueryParamBuilder {
	s.mediaFields = append(s.mediaFields, &mediaField)
	return s
}

// AddPlaceField adds a place field which enables you to select which specific place fields will deliver in each returned tweet.
// The Tweet will only return place fields if the Tweet contains a place and if you've also included `AddExpansion("geo.place_id")`.
// Learn more about media fields on twitter docs https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/api-reference/get-tweets-search-stream.
// Add a single place field for each invoke of `AddPlaceField`.
func (s *StreamQueryParamBuilder) AddPlaceField(placeField string) *StreamQueryParamBuilder {
	s.placeFields = append(s.placeFields, &placeField)
	return s
}