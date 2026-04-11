package go5paisa

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// Expected headers for ScripMaster CSV
var expectedHeaders = []string{
	"Exch",
	"ExchType",
	"ScripCode",
	"Name",
	"Expiry",
	"ScripType",
	"StrikeRate",
	"FullName",
	"TickSize",
	"LotSize",
	"QtyLimit",
	"Multiplier",
	"SymbolRoot",
	"BOCOAllowed",
	"ISIN",
	"ScripData",
	"Series",
}

// ScripMaster represents a single scrip master entry
type ScripMaster struct {
	Exch        string  `json:"Exch"`
	ExchType    string  `json:"ExchType"`
	ScripCode   int     `json:"ScripCode"`
	Name        string  `json:"Name"`
	Expiry      string  `json:"Expiry"`
	ScripType   string  `json:"ScripType"`
	StrikeRate  float64 `json:"StrikeRate"`
	FullName    string  `json:"FullName"`
	TickSize    float64 `json:"TickSize"`
	LotSize     int     `json:"LotSize"`
	QtyLimit    int     `json:"QtyLimit"`
	Multiplier  int     `json:"Multiplier"`
	SymbolRoot  string  `json:"SymbolRoot"`
	BOCOAllowed string  `json:"BOCOAllowed"`
	ISIN        string  `json:"ISIN"`
	ScripData   string  `json:"ScripData"`
	Series      string  `json:"Series"`
}

// GetScripMaster fetches scrip master data for a given segment
func (c *Client) GetScripMaster(segment string) ([]ScripMaster, error) {
	url := fmt.Sprintf("%s/ScripMaster/segment/%s", baseURL, segment)

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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse CSV
	reader := csv.NewReader(strings.NewReader(string(body)))
	reader.Comma = ','
	reader.FieldsPerRecord = -1 // Allow variable fields

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("invalid response: no data")
	}

	// Validate header
	if !reflect.DeepEqual(records[0], expectedHeaders) {
		return nil, fmt.Errorf("header mismatch: expected %v, got %v", expectedHeaders, records[0])
	}

	// Skip header
	records = records[1:]

	var scrips []ScripMaster
	for _, record := range records {
		if len(record) != 17 {
			return nil, fmt.Errorf("invalid record length: expected 17, got %d", len(record))
		}

		scrip := ScripMaster{}
		scrip.Exch = record[0]
		scrip.ExchType = record[1]

		scripCode, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, fmt.Errorf("invalid ScripCode: %s", record[2])
		}
		scrip.ScripCode = scripCode

		scrip.Name = record[3]
		scrip.Expiry = record[4]
		scrip.ScripType = record[5]

		strikeRate, err := strconv.ParseFloat(record[6], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid StrikeRate: %s", record[6])
		}
		scrip.StrikeRate = strikeRate

		scrip.FullName = record[7]

		tickSize, err := strconv.ParseFloat(record[8], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid TickSize: %s", record[8])
		}
		scrip.TickSize = tickSize

		lotSize, err := strconv.Atoi(record[9])
		if err != nil {
			return nil, fmt.Errorf("invalid LotSize: %s", record[9])
		}
		scrip.LotSize = lotSize

		qtyLimit, err := strconv.Atoi(record[10])
		if err != nil {
			return nil, fmt.Errorf("invalid QtyLimit: %s", record[10])
		}
		scrip.QtyLimit = qtyLimit

		multiplier, err := strconv.Atoi(record[11])
		if err != nil {
			return nil, fmt.Errorf("invalid Multiplier: %s", record[11])
		}
		scrip.Multiplier = multiplier

		scrip.SymbolRoot = record[12]
		scrip.BOCOAllowed = record[13]
		scrip.ISIN = record[14]
		scrip.ScripData = record[15]
		scrip.Series = record[16]

		scrips = append(scrips, scrip)
	}

	return scrips, nil
}
