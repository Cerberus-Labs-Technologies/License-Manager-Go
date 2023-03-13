package license

type License struct {
	ProductId int       `json:"productId"`
	UserId    int       `json:"userId"`
	Username  string    `json:"username"`
	License   string    `json:"license"`
	BoundIps  []string  `json:"boundIps"`
	Sponsored bool      `json:"sponsored"`
	Active    bool      `json:"active"`
	Permanent bool      `json:"permanent"`
	Expires   TimeStamp `json:"expires"`
	ExpiresAt TimeStamp `json:"expiresAt"`
}

type Config struct {
	LicenseKey string `json:"licenseKey"`
	ProductId  int    `json:"productId"`
	UserId     int    `json:"userId"`
}
