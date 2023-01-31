package utils

import (
	"context"
	"strconv"
)

// PaginatorParams This function calculates the pagination for list data
func PaginatorParams(page int32, limit int32) (int32, int32) {
	Limit := limit + 1
	Offset := (limit * page) - limit

	return Limit, Offset
}

// PaginatorPages Function that determines if the data has next page or not
func PaginatorPages(ctx context.Context, page int32, limit int32, dataLength int) (string, string) {
	// This returns django like next_page and previous_page on the pagination
	// This is not currently being used as the host is not required for frontend pagination
	//url := ctx.Value("url").(string)

	var nextPage string
	var previousPage string

	if page > 1 {
		previousPage = "Page=" + strconv.Itoa(int(page)-1) + "&ItemsPerPage=" + strconv.Itoa(int(limit))
	} else {
		previousPage = ""
	}
	if limit >= int32(dataLength) {
		nextPage = ""
	} else {
		nextPage = "Page=" + strconv.Itoa(int(page)+1) + "&ItemsPerPage=" + strconv.Itoa(int(limit))
	}

	return nextPage, previousPage

}
