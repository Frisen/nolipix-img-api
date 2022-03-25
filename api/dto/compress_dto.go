package dto

type CompressReq struct {
	Url        string `json:"url"`
	Cols       int    `json:"width"`
	Rows       int    `json:"height"`
	FileSuffix string `json:"suffix"`
}
