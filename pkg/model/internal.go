package model

// User is our enhanced user, containing all the aditional info
type User struct {
	IP         string `json:"ip"`
	Time       string `json:"time"`
	Country    string `json:"country"`
	ISOCountry string `json:"iso_country"`
	Distance   int    `json:"distance"`
	IsAWS      bool   `json:"is_aws"`
}

// Statistic response for service
type Statistic struct {
	Country  string `json:"country"`
	Distance int    `json:"distance"`
}
