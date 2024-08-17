package models

type QueryModel struct {
	All       bool   `query:"all"`
	Page      int    `query:"page"`
	PerPage   int    `query:"per_page"`
	Search    string `query:"search"`
	TenantID  string `query:"tenant_id"`
	SortField string `query:"sort_field"`
	SortOrder string `query:"sort_order"`
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
	Name      string `query:"name"`
}

type QueryModelResponse struct {
	Total    int         `json:"total"`
	PerPage  int         `json:"per_page"`
	CurPage  int         `json:"current_page"`
	LastPage int         `json:"last_page"`
	From     int         `json:"from"`
	To       int         `json:"to"`
	Length   int         `json:"length"`
	Data     interface{} `json:"data"`
}
