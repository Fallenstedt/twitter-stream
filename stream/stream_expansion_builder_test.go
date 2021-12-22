package stream

import "testing"

func TestStreamQueryParamsBuilderBuildsQueryParams(t *testing.T) {
	builder := NewStreamQueryParamsBuilder()

	result := builder.
		AddExpansion("expansion1").
		AddExpansion("expansion2").
		AddBackFillMinutes(1).
		AddMediaField("mediaField1").
		AddMediaField("mediaField2").
		AddPlaceField("placeField1").
		AddPlaceField("placeField2").
		AddPollField("pollField1").
		AddPollField("pollField2").
		AddTweetField("tweetField1").
		AddTweetField("tweetField2").
		AddUserField("userField1").
		AddUserField("userField2").
		Build().Encode()
	expected := "backfill_minutes=1&expansions=expansion1%2Cexpansion2&media.fields=mediaField1%2CmediaField2&place.fields=placeField1%2CplaceField2&poll.fields=pollField1%2CpollField2&tweet.fields=tweetField1%2CtweetField2&user.fields=userField1%2CuserField2"
	if result != expected {
		t.Errorf("ahh")
	}

}