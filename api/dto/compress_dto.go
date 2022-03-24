package dto

type CompressReq struct {
	Url  string `json:"url"`
	Cols int    `json:"width"`
	Rows int    `json:"height"`
}
