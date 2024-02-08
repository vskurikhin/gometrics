/*
 * This file was last modified at 2024-02-08 08:55 by Victor N. Skurikhin.
 * Parser.go
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

	value := path[3]

	v, err := t.ParseValue(value)
	if err != nil {
		return &Parser{status: http.StatusBadRequest}, err
	}
	return &Parser{
		r: r, t: t,
		n:        types.Lookup(path[2]),
		status:   http.StatusOK,
		original: value,
		value:    v,
	}, nil
}

func errorParser(r *http.Request, status int) (*Parser, error) {
	return &Parser{status: status}, errors.New("can't Parse request" + util.FormatRequest(r))
}

func (p *Parser) String() string {
	return p.n.String()
}

func (p *Parser) Name() int {
	return int(p.n)
}

func (p *Parser) MetricType() types.Types {
	return p.t
}

func (p *Parser) Original() string {
	return p.original
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

	old, err := p.MetricType().ParseValue(*get)
	if err != nil {
		return nil
	}

	switch v := p.Value().(type) {
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
