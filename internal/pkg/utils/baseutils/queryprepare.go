package baseutils

import (
	"strings"
)

func PrepareQueryForSearch(query string) string {
	query = strings.Join(strings.Fields(query), "|")
	query = "%(" + query + ")%"
	return query
}
