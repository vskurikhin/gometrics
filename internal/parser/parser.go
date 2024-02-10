/*
 * This file was last modified at 2024-02-11 00:21 by Victor N. Skurikhin.
 * parser.go
 * $Id$
 */

package parser

import (
	"errors"
	"fmt"
	"github.com/vskurikhin/gometrics/internal/types"
	"github.com/vskurikhin/gometrics/internal/util"
	"net/http"
)

type parser struct {
	request       *http.Request
	type_         types.Types
	number        types.Name
	name          string
	originalValue string
	parsedValue   interface{}
	httpStatus    int
}

const FixedPathLength = 4

func Parse(r *http.Request) (*parser, error) {

	path := util.SplitPath(r)
	if len(path) < FixedPathLength || r.Method != http.MethodPost {
		return errorParser(r, http.StatusNotFound)
	}
	return parseType(r, path)
}

func parseType(r *http.Request, path []string) (*parser, error) {

	switch {
	case types.COUNTER.Eq(path[1]):
		return parseName(r, types.COUNTER, path)
	case types.GAUGE.Eq(path[1]):
		return parseName(r, types.GAUGE, path)
	}
	return errorParser(r, http.StatusBadRequest)
}

func parseName(r *http.Request, t types.Types, path []string) (*parser, error) {

	num := types.Lookup(path[2])
	var name string
	if num > 0 {
		name = num.String()
	} else {
		name = path[2]
	}

	value, err := t.ParseValue(path[3])
	if err != nil {
		return &parser{httpStatus: http.StatusBadRequest}, err
	}
	return &parser{
		request:       r,
		type_:         t,
		number:        num,
		name:          name,
		originalValue: path[3],
		parsedValue:   value,
		httpStatus:    http.StatusOK,
	}, nil
}

func errorParser(r *http.Request, status int) (*parser, error) {
	return &parser{httpStatus: status}, errors.New("can'type_ Parse request" + util.FormatRequest(r))
}

func (p *parser) String() string {
	return p.name
}

func (p *parser) Value() interface{} {
	return p.parsedValue
}

func (p *parser) Status() int {
	return p.httpStatus
}

func (p *parser) CalcValue(get *string) *string {

	if get == nil {
		return &p.originalValue
	}

	old, err := p.type_.ParseValue(*get)
	if err != nil {
		return nil
	}

	switch v := p.parsedValue.(type) {
	case float64:
		return &p.originalValue
	case int:
		o := typeAssertionInt(old)
		if o == nil {
			return &p.originalValue
		}
		result := fmt.Sprintf("%d", v+*o)
		return &result
	default:
		return nil
	}
}

func typeAssertionInt(i interface{}) *int {
	switch o := i.(type) {
	case int:
		return &o
	default:
		return nil
	}
}
