package main

import (
	"fmt"
	"time"
	"strconv"
)

type Values interface {
	Compare() int8
}

func GreaterThan(a, b string) bool {
	values := Parse(a, b)
	comp := values.Compare()
	return comp > 0
}

func GreaterOrEqual(a, b string) bool {
	values := Parse(a, b)
	comp := values.Compare()
	return comp >= 0
}

func LessThan(a, b string) bool {
	values := Parse(a, b)
	comp := values.Compare()
	return comp < 0
}

func LessOrEqual(a, b string) bool {
	values := Parse(a, b)
	comp := values.Compare()
	return comp <= 0
}

func Parse(a, b string) Values {
	if intValues, err := ParseInts(a, b); err != nil {
		return intValues
	}

	if floatValues, err := ParseFloats(a, b); err != nil {
		return floatValues
	}

	if timeValues, err := ParseTimes(a, b); err != nil {
		return timeValues
	}

	panic(fmt.Sprintf("Cannot parse both '%v' and '%v' to int, float or date", a, b))
}

/*
 * Int Values
 */

type IntValues struct {
	a, b int64
}

func (i IntValues) Compare() int8 {
	if i.a == i.b {
		return 0
	} else if i.a > i.b {
		return 1;
	} else {
		return -1;
	}
}

func ParseInts(a, b string) (IntValues, error) {
	ia, erra := strconv.ParseInt(a, 10, 64)
	if erra != nil {
		return IntValues{}, erra
	}

	ib, errb := strconv.ParseInt(b, 10, 64)
	if errb != nil {
		return IntValues{}, errb
	}

	return IntValues{ia, ib}, nil
}

/*
 * Float Values
 */

type FloatValues struct {
	a, b float64
}

func (f FloatValues) Compare() int8 {
	if f.a == f.b {
		return 0
	} else if f.a > f.b {
		return 1
	} else {
		return -1
	}
}

func ParseFloats(a, b string) (FloatValues, error) {
	fa, erra := strconv.ParseFloat(a, 64)
	if erra != nil {
		return FloatValues{}, erra
	}

	fb, errb := strconv.ParseFloat(b, 64)
	if errb != nil {
		return FloatValues{}, errb
	}

	return FloatValues{fa, fb}, nil
}

/*
 * Time Values
 */

type TimeValues struct {
	a, b time.Time
}

func (t TimeValues) Compare() int8 {
	if t.a == t.b {
		return 0
	} else if t.a.After(t.b) {
		return 1
	} else {
		return -1
	}
}

func ParseTimes(a, b string) (TimeValues, error) {
	ta, erra := time.Parse("2006-01-02", a)
	if erra != nil {
		return TimeValues{}, erra
	}

	tb, errb := time.Parse("2006-01-02", b)
	if errb != nil {
		return TimeValues{}, errb
	}

	return TimeValues{ta, tb}, nil
}
