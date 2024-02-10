/*
 * This file was last modified at 2024-02-10 15:27 by Victor N. Skurikhin.
 * parser.go
 * $Id$
 */

package parser

import (
	"errors"
	"fmt"
	"github.com/vskurikhin/gometrics/api/types"
	"github.com/vskurikhin/gometrics/internal/util"
	"net/http"
)

type Parser struct {
	r        *http.Request
	t        types.Types
	n        types.Name
	name     string
	original string
	value    interface{}
	status   int
}

const FixedPathLength = 4

func Parse(r *http.Request) (*Parser, error) {

	path := util.SplitPath(r)
	if len(path) < FixedPathLength || r.Method != http.MethodPost {
		return errorParser(r, http.StatusNotFound)
	}
	return parseType(r, path)
}

func parseType(r *http.Request, path []string) (*Parser, error) {

	switch {
	case types.COUNTER.Eq(path[1]):
		return parseName(r, types.COUNTER, path)
	case types.GAUGE.Eq(path[1]):
		return parseName(r, types.GAUGE, path)
	}
	return errorParser(r, http.StatusBadRequest)
}

func parseName(r *http.Request, t types.Types, path []string) (*Parser, error) {

	num := types.Lookup(path[2])
	var name string
	if num > 0 {
		name = num.String()
	} else {
		name = path[2]
	}

	value, err := t.ParseValue(path[3])
	if err != nil {
		return &Parser{status: http.StatusBadRequest}, err
	}
	return &Parser{
		r:        r,
		t:        t,
		n:        num,
		name:     name,
		original: path[3],
		value:    value,
		status:   http.StatusOK,
	}, nil
}

func errorParser(r *http.Request, status int) (*Parser, error) {
	return &Parser{status: status}, errors.New("can't Parse request" + util.FormatRequest(r))
}

func (p *Parser) String() string {
	return p.name
}

func (p *Parser) Value() interface{} {
	return p.value
}

func (p *Parser) Status() int {
	return p.status
}

func (p *Parser) CalcValue(get *string) *string {

	if get == nil {
		return &p.original
	}

	old, err := p.t.ParseValue(*get)
	if err != nil {
		return nil
	}

	switch v := p.value.(type) {
	case float64:
		return &p.original
	case int:
		o := switchCase(old)
		if o == nil {
			return &p.original
		}
		result := fmt.Sprintf("%d", v+*o)
		return &result
	default:
		return nil
	}
}

func switchCase(i interface{}) *int {
	switch o := i.(type) {
	case int:
		return &o
	default:
		return nil
	}
}
