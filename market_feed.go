package go5paisa

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// MarketFeedRequestData represents a single MarketFeedRequestData
type MarketFeedRequestData struct {
	Exchange     string `json:"Exch"`
	ExchangeType string `json:"ExchType"`
	Symbol       string `json:"Symbol"`
	Expiry       string `json:"Expiry"`
	StrikePrice  string `json:"StrikePrice"`
	OptionType   string `json:"OptionType"`

	ScripCode int `json:"ScripCode"`
}

type MarketFeedRequest struct {
	Count           int                     `json:"Count"`
	ClientLoginType int                     `json:"ClientLoginType"`
	LastRequestTime string                  `json:"LastRequestTime"`
	RefreshRate     string                  `json:"RefreshRate"`
	MarketFeedData  []MarketFeedRequestData `json:"MarketFeedData"`
}

type marketFeedPayload struct {
	Head *payloadHead      `json:"head"`
	Body MarketFeedRequest `json:"body"`
}

// Holding represents a single holding
type MarketFeed struct {
	Symbol       string  `json:"Symbol"`
	LastRate     float32 `json:"LastRate"`
	Exchange     string  `json:"Exch"`
	ExchangeType string  `json:"ExchType"`
	Chg          float32 `json:"Chg"`
	ChgPcnt      float32 `json:"ChgPcnt"`
	High         float32 `json:"High"`
	Low          float32 `json:"Low"`
	PClose       float32 `json:"PClose"`
	TickDt       string  `json:"TickDt"`
	Time         int     `json:"Time"`
	Token        int     `json:"Token"`
	TotalQty     int     `json:"TotalQty"`
}

type MarketFeedResponseHead struct {
	ResponseCode      string `json:"responseCode"`
	Status            string `json:"status"`
	StatusDescription string `json:"statusDescription"`
}

type MarketFeedResponseBody struct {
	CacheTime int          `json:"CacheTime"`
	Status    int          `json:"Status"`
	Message   string       `json:"Message"`
	TimeStamp string       `json:"TimeStamp"`
	Data      []MarketFeed `json:"Data"`
}

// Data has all holdings for a user
type MarketFeedResponse struct {
	Head MarketFeedResponseHead `json:"head"`
	Body MarketFeedResponseBody `json:"body"`
}

func parseMarketFeedsResponse(resBody []byte, obj Holdings) {
	var body responseData
	body.Body = obj
	if err := json.Unmarshal(resBody, &body); err != nil {
		log.Fatal("Error parsing JSON response:", err)
	}
}

func (r *MarketFeedRequest) AddMarketFeedRequestData(data MarketFeedRequestData) {

	r.MarketFeedData = append(r.MarketFeedData, data)
	r.Count = len(r.MarketFeedData)
}

// GetHoldings fetches holdings of the user
func (c *Client) GetMarketFeed(request MarketFeedRequest) (MarketFeedResponse, error) {
	var feed MarketFeedResponse

	epochTime := time.Now().UnixNano() / int64(time.Millisecond)
	dateString := fmt.Sprintf("/Date(%d)/", epochTime)

	request.LastRequestTime = dateString

	head := c.buildHeader("5PMF")

	payload := marketFeedPayload{
		Head: &head,
		Body: request,
	}

	resBody, err := c.postRequest(payload, marketFeedRoute)
	if err != nil {
		return feed, err
	}

	parseResBody(resBody, &feed)

	if err := json.Unmarshal(resBody, &feed); err != nil {
		log.Fatal("Error parsing JSON response:", err)
	}

	return feed, nil
}

func (c *Client) GetMarketFeedV1(request MarketFeedRequest) (MarketFeedResponse, error) {
	var feed MarketFeedResponse

	epochTime := time.Now().UnixNano() / int64(time.Millisecond)
	dateString := fmt.Sprintf("/Date(%d)/", epochTime)

	request.LastRequestTime = dateString

	head := c.buildHeader("5PMF")

	payload := marketFeedPayload{
		Head: &head,
		Body: request,
	}

	resBody, err := c.postRequest(payload, marketFeedRouteV1)
	if err != nil {
		return feed, err
	}

	parseResBody(resBody, &feed)

	if err := json.Unmarshal(resBody, &feed); err != nil {
		log.Fatal("Error parsing JSON response:", err)
	}

	return feed, nil
}
