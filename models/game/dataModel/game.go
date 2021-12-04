package dataModel

type GameSales struct {
	Rank        int     `json:"rank"`
	Name        string  `json:"name"`
	Platform    string  `json:"platform"`
	Year        int     `json:"year"`
	Genre       string  `json:"genre"`
	Publisher   string  `json:"publisher"`
	NASales     float64 `json:"NA_Sales"`
	EUSales     float64 `json:"EU_Sales"`
	JPSales     float64 `json:"JP_Sales"`
	OtherSales  float64 `json:"other_Sales"`
	GlobalSales float64 `json:"global_Sales"`
}
