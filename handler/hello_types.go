package handler

type BaseRepo struct {
}

type HelloReq struct {
	Name string `uri:"name"` // uri,form,json
}

type HelloResp struct {
	Msg string `json:"msg"`
}
