package entities

type ResponseData struct {
    ErrCode int64 `json:"errCode"`
    ErrMsg string `json:"errMsg"`
    Data interface{} `json:"data"`
}
