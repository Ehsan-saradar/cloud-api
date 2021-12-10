package models

type Games struct {
	Rank       int     `json:"Rank"`
	Name       string  `json:"Name"`
	Platform   string  `json:"Platform"`
	Year       int     `json:"Year"`
	Genre      string  `json:"Genre"`
	Publisher  string  `json:"Publisher"`
	NASale     float64 `json:"NA_Sales"`
	EUSale     float64 `json:"EU_Sales"`
	JPSale     float64 `json:"JP_Sales"`
	OtherSale  float64 `json:"Other_Sales"`
	GlobalSale float64 `json:"Global_Sales"`
}
