package webutility

import "net/http"

type PaginationParams struct {
	URL    string
	Offset int64
	Limit  int64
	SortBy string
	Order  string
}

func (p *PaginationParams) links() PaginationLinks {
	return PaginationLinks{}
}

// PaginationLinks ...
type PaginationLinks struct {
	Count int64
	Total int64
	Base  string `json:"base"`
	Next  string `json:"next"`
	Prev  string `json:"prev"`
	Self  string `json:"self"`
}

func GetPaginationParameters(req *http.Request) (p PaginationParams) {
	p.Offset = StringToInt64(req.FormValue("offset"))
	p.Limit = StringToInt64(req.FormValue("limit"))
	p.SortBy = req.FormValue("sortBy")
	p.Order = req.FormValue("order")

	return p
}
