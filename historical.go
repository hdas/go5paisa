package go5paisa

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HistoricalCandle represents a single historical candle
type HistoricalCandle struct {
	Timestamp time.Time `json:"-"`
	Open      float64   `json:"-"`
	High      float64   `json:"-"`
	Low       float64   `json:"-"`
	Close     float64   `json:"-"`
	Volume    int64     `json:"-"`
}

// UnmarshalJSON custom unmarshaler for HistoricalCandle since it's an array
func (hc *HistoricalCandle) UnmarshalJSON(data []byte) error {
	var arr []interface{}
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) != 6 {
		return fmt.Errorf("invalid candle data length")
	}
	ts, ok := arr[0].(string)
	if !ok {
		return fmt.Errorf("invalid timestamp")
	}

	if !strings.Contains(ts, "Z") {
		ts += "Z"
	}

	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return err
	}
	hc.Timestamp = t
	hc.Open, ok = arr[1].(float64)
	if !ok {
		return fmt.Errorf("invalid open")
	}
	hc.High, ok = arr[2].(float64)
	if !ok {
		return fmt.Errorf("invalid high")
	}
	hc.Low, ok = arr[3].(float64)
	if !ok {
		return fmt.Errorf("invalid low")
	}
	hc.Close, ok = arr[4].(float64)
	if !ok {
		return fmt.Errorf("invalid close")
	}
	vol, ok := arr[5].(float64)
	if !ok {
		return fmt.Errorf("invalid volume")
	}
	hc.Volume = int64(vol)
	return nil
}

// HistoricalData holds the candles data
type HistoricalData struct {
	Candles []HistoricalCandle `json:"candles"`
}

// HistoricalResponse is the full API response
type HistoricalResponse struct {
	Status string         `json:"status"`
	Data   HistoricalData `json:"data"`
}

// GetHistoricalCandles fetches historical candles for a given scrip
func (c *Client) GetHistoricalCandles(exchange, exchangeType string, scripCode int, interval string, from, to time.Time) ([]HistoricalCandle, error) {
	url := fmt.Sprintf("https://openapi.5paisa.com/V2/historical/%s/%s/%d/%s?from=%s&end=%s",
		exchange, exchangeType, scripCode, interval,
		from.Format(time.DateOnly), to.Format(time.DateOnly))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	res, err := c.connection.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", res.Status)
	}

	var data []byte
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response HistoricalResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	if response.Status != "success" {
		return nil, fmt.Errorf("API returned status: %s", response.Status)
	}

	return response.Data.Candles, nil
}
