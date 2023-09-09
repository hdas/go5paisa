package go5paisa

import (
	"bytes"
	"encoding/json"
	"io"
	// "log"
)

type TradeResponseList struct {
	TradeResList []OrderResponse `json:"TradeDetail"`
}

type TradeInfoList struct {
	TradeList []OrderForStatus `json:"TradeDetailList"`
}

// GetTradeInformation fetches order book of the user
func (c *Client) GetTradeInformation(tradeList TradeInfoList) (TradeInfoList, error) {
	var tradeReslist TradeInfoList
	// c.appConfig.head.RequestCode =
	head := c.buildHeader(tradeInfoRequestCode)
	payloadBody := orderStatusPayloadBody{
		ClientCode: c.clientCode,
		OrdList:    tradeList.TradeList,
	}
	payload := orderStatusPayload{
		Head: &head,
		Body: payloadBody,
	}
	jsonValue, _ := json.Marshal(payload)
	res, err := c.connection.Post(baseURL+tradeInfoRoute, contentType, bytes.NewBuffer(jsonValue))
	if err != nil {
		return tradeReslist, err
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return tradeReslist, err
	}
	parseResBody(resBody, &tradeReslist)
	return tradeReslist, nil
}
