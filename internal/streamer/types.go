package streamer

type MarkPriceData struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`

	MarkPrice            string `json:"p"`
	IndexPrice           string `json:"i"`
	EstimatedSettlePrice string `json:"P"`

	FundingRate string `json:"r"`
	FundingTime int64  `json:"T"` // 펀딩피 지불 시간
}
