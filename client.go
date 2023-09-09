package go5paisa

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"

	"golang.org/x/net/publicsuffix"
)

const (
	baseURL string = "https://Openapi.5paisa.com/VendorsAPI/Service1.svc"

	loginRoute          string = "/V2/LoginRequestMobileNewbyEmail"
	loginTotpRoute      string = "/TOTPLogin"
	accessTokenRoute    string = "/getAccessToken"
	marginRoute         string = "/V3/Margin"
	orderBookRoute      string = "/V2/OrderBook"
	holdingsRoute       string = "/V2/Holding"
	positionsRoute      string = "/V1/NetPositionNetWise"
	orderPlacementRoute string = "/V1/OrderRequest"
	orderStatusRoute    string = "/OrderStatus"
	tradeInfoRoute      string = "/TradeInformation"
	marketFeedRoute     string = "/MarketFeed"
	marketFeedRouteV1   string = "/V1/MarketFeed"

	// Request codes
	marginRequestCode         string = "5PMarginV3"
	orderBookRequestCode      string = "5POrdBkV3"
	holdingsRequestCode       string = "5PMarginV3" //"5PHoldingV2"
	positionsRequestCode      string = "5PNPNWV1"
	tradeInfoRequestCode      string = "5PTrdInfo"
	orderStatusRequestCode    string = "5POrdStatus"
	orderPlacementRequestCode string = "5POrdReq"
	loginRequestCode          string = "5PLoginV2"

	// Content Type
	contentType string = "application/json"
)

// Config is the app configuration
type Config struct {
	AppName       string
	AppSource     int64
	UserID        string
	Password      string
	UserKey       string
	EncryptionKey string

	LocalIP  string
	PublicIP string
}

// Client is the client configuration
type Client struct {
	clientCode  string
	connection  *http.Client
	AccessToken string
	config      *Config
}

// Init initializes the Client struct
func (c *Client) Init(conf *Config) error {
	c.config = conf

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return err
	}

	c.connection = &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}

	return nil
}

func (c *Client) SetAccessToken(clientCode string, accessToken string) error {
	c.AccessToken = accessToken
	c.clientCode = clientCode

	return nil
}

// Login logs in a client
func (c *Client) Login(loginId string, pin string, totp string) error {

	loginRequestBody := loginBody{
		Email:    loginId,
		PIN:      pin,
		LocalIP:  c.config.LocalIP,
		PublicIP: c.config.PublicIP,
		TOTP:     totp,
	}

	head := &loginHead{
		Key: c.config.UserKey,
	}

	loginDetails := loginPayload{
		Head: head,
		Body: loginRequestBody,
	}
	jsonValue, _ := json.Marshal(loginDetails)
	res, err := c.connection.Post(baseURL+loginTotpRoute, contentType, bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var b body
	parseResBody(resBody, &b)
	if b.ClientCode == "" || b.ClientCode == "INVALID CODE" {
		return errors.New(b.Message)
	}

	c.AccessToken, err = c.getAccessToken(b.RequestToken)
	if err != nil {
		return err
	}

	fmt.Println(c.AccessToken)

	c.clientCode = b.ClientCode

	return nil
}

func (c *Client) getAccessToken(requestToken string) (string, error) {

	accessTokenBody := requestAccessTokenBody{
		UserId:       c.config.UserID,
		RequestToken: requestToken,
		EncryKey:     c.config.EncryptionKey,
		LocalIP:      c.config.LocalIP,
		PublicIP:     c.config.PublicIP,
	}

	head := &loginHead{
		Key: c.config.UserKey,
	}

	loginDetails := accessTokenPayload{
		Head: head,
		Body: accessTokenBody,
	}
	jsonValue, _ := json.Marshal(loginDetails)
	res, err := c.connection.Post(baseURL+accessTokenRoute, contentType, bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var b AccessTokenBody
	parseResBody(resBody, &b)
	if b.AccessToken == "" || b.AccessToken == "INVALID CODE" {
		return "", errors.New(b.Message)
	}

	return b.AccessToken, nil
}

func (c *Client) buildHeader(requestCode string) payloadHead {
	head := payloadHead{
		Key:         c.config.UserKey,
		AppVer:      "1.0",
		AppName:     c.config.AppName,
		OsName:      "WEB",
		UserID:      c.config.UserID,
		Password:    c.config.Password,
		RequestCode: requestCode,
	}
	return head
}

func (c *Client) postRequest(payload any, route string) ([]byte, error) {
	jsonValue, _ := json.Marshal(payload)

	// fmt.Println(string(jsonValue))

	req, err := http.NewRequest("POST", baseURL+route, bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	res, err := c.connection.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return resBody, nil
}
