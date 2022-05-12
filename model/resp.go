package model

type ResponseForm struct {
	Success    bool             `json:"success"`
	Result     interface{}      `json:"result"`
	Messages   []string         `json:"messages"`
	Errors     []*ResponseError `json:"erros"`
	ResultInfo *ResultInfo      `json:"result_info,omitempty"`
}

type ResponseError struct {
	Code    int         `json:"code"`
	Source  interface{} `json:"source,omitempty"`
	Title   string      `json:"title,omitempty"`
	Message string      `json:"message"`
}

type ResultInfo struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Count      int `json:"count"`
	TotalCount int `json:"total_count"`
}
