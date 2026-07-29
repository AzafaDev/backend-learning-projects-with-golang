package weather

type Response struct {
	Temperature float64 `json:"temperature"`
	City        string  `json:"city"`
	Condition   string  `json:"condition"`
}
