package pricing

type CalculatePricingResp struct {
	ErrCode    int    `json:"errCode"` //errCode == 0 -> success; errCode != 0 error
	ErrMessage string `json:"errMessage"`
	Price      int    `json:"price"`
}
