package api

type shortReq struct {
	LongURL string `json:"longURL"`
}

type errorResp struct {
	Msg string `json:"msg"`
}

type shortResp struct {
	ShortURL string `json:"shortURL"`
}
