/*
 * This file was last modified at 2024-02-04 13:43 by Victor N. Skurikhin.
 * split_path.go
 * $Id$
 */

package util

import (
	"net/http"
	"strings"
)

func SplitPath(r *http.Request) []string {

	path := r.URL.Path
	path = strings.TrimSpace(path)
	//Отрезает ведущий и конечный слеш, если он существует.
	path = strings.TrimPrefix(path, "/")

	if strings.HasSuffix(path, "/") {
		cutOffLastChar := len(path) - 1
		path = path[:cutOffLastChar]
	}

	//Изолировать отдельные компоненты пути.
	components := strings.Split(path, "/")
	return components
}
