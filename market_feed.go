package go5paisa

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// MarketFeedRequestData represents a single MarketFeedRequestData
type MarketFeedRequestDataV3 struct {
	Exchange     string `json:"Exch"`
	ExchangeType string `json:"ExchType"`
	ScripCode    int    `json:"ScripCode"`
}

// MarketFeedRequestData represents a single MarketFeedRequestData
type MarketFeedRequestData struct {
	Exchange     string `json:"Exch"`
	ExchangeType string `json:"ExchType"`
	ScripCode    int    `json:"ScripCode"`
	Symbol       string `json:"Symbol"`
	Expiry       string `json:"Expiry"`
	StrikePrice  string `json:"StrikePrice"`
	OptionType   string `json:"OptionType"`
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

type MarketFeedV3 struct {
	Exchange     string  `json:"Exch"`
	ExchangeType string  `json:"ExchType"`
	Token        int     `json:"Token"`
	LastRate     float64 `json:"LastRate"`
	LastQty      float32 `json:"LastQty"`
	TotalQty     int     `json:"TotalQty"`
	High         float32 `json:"High"`
	Low          float32 `json:"Low"`
	OpenRate     float32 `json:"OpenRate"`
	PClose       float32 `json:"PClose"`
	AvgRate      float32 `json:"AvgRate"`
	Time         int     `json:"Time"`
	BidQty       int     `json:"BidQty"`
	BidRate      float32 `json:"BidRate"`
	OffQty       int     `json:"OffQty"`
	OffRate      float32 `json:"OffRate"`
	TBidQ        int     `json:"TBidQ"`
	TOffQ        int     `json:"TOffQ"`
	TickDt       string  `json:"TickDt"`
	ChgPcnt      float32 `json:"ChgPcnt"`
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
		panic(err)
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
		panic(err)
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
		panic(err)
	}

	return feed, nil
}

type MarketFeedCallback func(feeds *[]MarketFeedV3)

func (c *Client) InitMarketFeedWebSocket(stocks *[]MarketFeedRequestDataV3, callback MarketFeedCallback) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Connect to the WebSocket server
	serverURL := "wss://openfeed.5paisa.com/Feeds"

	serverURL = serverURL + "/api/chat?Value1=" + c.AccessToken + "|" + c.clientCode

	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	done := make(chan struct{})
	// Send message in a separate goroutine
	go sendMessage(conn, c.clientCode, done, stocks)

	// Start a separate goroutine to listen for messages from the server

	// Wait for an interrupt signal or completion of message sending
	select {
	case <-done:
		log.Println("Message sent. Exiting.")
	case <-interrupt:
		log.Println("Interrupt signal received. Closing connection.")
		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("Error sending close message:", err)
			return
		}
		select {
		case <-done:
		}
		return
	}

	receiveMessages(conn, callback)
}

func sendMessage(conn *websocket.Conn, clientCode string, done chan<- struct{}, stocks *[]MarketFeedRequestDataV3) {

	fullJson := `{"Method":"MarketFeedV3","Operation":"Subscribe","ClientCode":"%s",
	"MarketFeedData":`

	jsonValue, _ := json.Marshal(stocks)

	fullJson += string(jsonValue)

	fullJson += "}"

	message := fmt.Sprintf(fullJson, clientCode)
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("Error sending message:", err)
	}
	done <- struct{}{}
}

func receiveMessages(conn *websocket.Conn, callback MarketFeedCallback) {
	for {
		// Read message from the server
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		var feeds []MarketFeedV3 = make([]MarketFeedV3, 0)

		//parseResBody(message, feeds)

		if err := json.Unmarshal(message, &feeds); err != nil {
			panic(err)
		}

		callback(&feeds)
	}
}
