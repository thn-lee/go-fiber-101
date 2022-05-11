package model

type ResponseForm struct {
	Success    bool
	Result     interface{}
	Messages   []string
	Errors     []*ResponseError
	ResultInfo *ResultInfo
}

type ResponseError struct {
	Code    int
	Source  interface{}
	Title   string
	Message string
}

type ResultInfo struct {
	Page       int
	PerPage    int
	Count      int
	TotalCount int
}
