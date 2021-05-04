package baseutils

import "strings"

func PrepareQueryForSearch(query string) string {
	query = strings.ToLower(strings.Trim(query, " "))
	query = strings.Join(strings.Split(query, " "), "%")
	query = "%" + query + "%"
	return query
}
