package main

import (
	"os"
	"fmt"
	"time"
	"strconv"
)

type Values interface {
	Compare() int8
}

func GreaterThan(a, b string) bool {
	if comp, ok := Compare(a, b); ok {
		return comp > 0
	}
	return false
}

func GreaterOrEqual(a, b string) bool {
	if comp, ok := Compare(a, b); ok {
		return comp >= 0
	}
	return false
}

func LessThan(a, b string) bool {
	if comp, ok := Compare(a, b); ok {
		return comp < 0
	}
	return false
}

func LessOrEqual(a, b string) bool {
	if comp, ok := Compare(a, b); ok {
		return comp <= 0
	}
	return false
}

func Compare(a, b string) (int8, bool) {
	values, err := Parse(a, b)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 0, false
	}

	return values.Compare(), true
}

func Parse(a, b string) (Values, error) {
	if intValues, err := ParseInts(a, b); err == nil {
		return intValues, nil
	}

	if floatValues, err := ParseFloats(a, b); err == nil {
		return floatValues, nil
	}

	if timeValues, err := ParseTimes(a, b); err == nil {
		return timeValues, nil
	}

	return nil, fmt.Errorf("Cannot convert both '%v' and '%v' to same type (tried int, float, and date)", a, b)
}

/*
 * Int Values
 */

type IntValues struct {
	a, b int64
}

func (i IntValues) String() string {
	return fmt.Sprintf("%T#{a=%v, b=%v}", i, i.a, i.b)
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

func (f FloatValues) String() string {
	return fmt.Sprintf("%T#{a=%v, b=%v}", f, f.a, f.b)
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

func (t TimeValues) String() string {
	return fmt.Sprintf("%T#{a=%v, b=%v}", t, t.a, t.b)
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
