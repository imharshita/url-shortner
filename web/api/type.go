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

type expandReq struct {
	ShortURL string `json:"shortURL"`
}

type expandResp struct {
	LongURL string `json:"longURL"`
}

type domainCount struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
}
