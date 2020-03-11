package controllers

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/volatiletech/sqlboiler/queries/qm"
)

// GetQueryModFromQuery derives db lookups from URI query parameters
func GetQueryModFromQuery(query string) []qm.QueryMod {
	var mods []qm.QueryMod
	m, _ := url.ParseQuery(query)
	for k, v := range m {
		for _, value := range v {
			if k == "limit" {
				limit, err := strconv.Atoi(value)
				if err != nil {
					continue
				}
				mods = append(mods, qm.Limit(limit))
			} else if k == "from" {
				from, err := strconv.Atoi(value)
				if err != nil {
					continue
				}
				// TODO: support order by and ASC/DESC
				mods = append(mods, qm.Where("id >= ?", from))
			} else {
				clause := fmt.Sprintf("%s=?", k)
				mods = append(mods, qm.Where(clause, v))
			}
		}
	}
	return mods
}
