package go5paisa

type payloadHead struct {
	AppName     string `json:"appName"`
	AppVer      string `json:"appVer"`
	Key         string `json:"key"`
	OsName      string `json:"osName"`
	RequestCode string `json:"requestCode"`
	UserID      string `json:"userId"`
	Password    string `json:"password"`
}

type loginHead struct {
	Key string `json:"Key"`
}

type loginBody struct {
	Email    string `json:"Email_ID"`
	PIN      string `json:"PIN"`
	LocalIP  string `json:"LocalIP"`
	PublicIP string `json:"PublicIP"`
	TOTP     string `json:"TOTP"`
}

type requestAccessTokenBody struct {
	UserId       string `json:"UserId"`
	RequestToken string `json:"RequestToken"`
	EncryKey     string `json:"EncryKey"`
	LocalIP      string `json:"LocalIP"`
	PublicIP     string `json:"PublicIP"`
}

type loginPayload struct {
	Head *loginHead `json:"head"`
	Body loginBody  `json:"body"`
}

type accessTokenPayload struct {
	Head *loginHead             `json:"head"`
	Body requestAccessTokenBody `json:"body"`
}

type genericPayload struct {
	Head *payloadHead       `json:"head"`
	Body genericPayloadBody `json:"body"`
}

type genericPayloadBody struct {
	ClientCode string `json:"ClientCode"`
}

type orderStatusPayloadBody struct {
	ClientCode string           `json:"ClientCode"`
	OrdList    []OrderForStatus `json:"OrdStatusReqList"`
}

type orderStatusPayload struct {
	Head *payloadHead           `json:"head"`
	Body orderStatusPayloadBody `json:"body"`
}
