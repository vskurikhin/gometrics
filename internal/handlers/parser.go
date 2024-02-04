/*
 * This file was last modified at 2024-02-04 16:54 by Victor N. Skurikhin.
 * parser.go
 * $Id$
 */

package handlers

import (
	"errors"
	"github.com/vskurikhin/gometrics/api/names"
	"github.com/vskurikhin/gometrics/api/types"
	"github.com/vskurikhin/gometrics/internal/util"
	"net/http"
)

type parser struct {
	r      *http.Request
	t      types.Types
	n      names.Names
	value  *string
	status int
}

const FixedPathLength = 4

func parse(r *http.Request) (parser, error) {

	path := util.SplitPath(r)
	if len(path) < FixedPathLength || r.Method != http.MethodPost {
		return errorParser(r, http.StatusNotFound)
	}
	return parseType(r, path)
}

func parseType(r *http.Request, path []string) (parser, error) {

	switch {
	case types.COUNTER.Eq(path[1]):
		return parseName(r, types.COUNTER, path)
	case types.GAUGE.Eq(path[1]):
		return parseName(r, types.GAUGE, path)
	}
	return errorParser(r, http.StatusBadRequest)
}

func parseName(r *http.Request, t types.Types, path []string) (parser, error) {

	value := path[3]

	if _, err := t.ParseValue(value); err != nil {
		return parser{status: http.StatusBadRequest}, err
	}
	return parser{
		r: r, t: t,
		n:      names.Lookup(path[2]),
		status: http.StatusOK,
		value:  &value,
	}, nil
}

func errorParser(r *http.Request, status int) (parser, error) {
	return parser{status: status}, errors.New("can't parse request" + util.FormatRequest(r))
}
