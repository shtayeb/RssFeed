package types

type Pagination struct {
	PerPage      int
	CurrentPage  int
	FirstPageUrl string
	LastPageUrl  string
	NextPageUrl  string
	PrevPageUrl  string
	Next         int
	Previous     int
	TotalPage    int
}
