package blowdash

import (
	"fmt"
	"strings"
)

func Keys[T comparable, U any](m map[T]U) []T {
	var (
		ret = make([]T, len(m))
		i   int
	)
	for key := range m {
		ret[i] = key
		i++
	}
	return ret
}

func Values[T comparable, U any](m map[T]U) []U {
	var (
		ret = make([]U, len(m))
		i   int
	)
	for _, value := range m {
		ret[i] = value
		i++
	}
	return ret
}

func Map[T any, U any](elems []T, fn func(T) U) []U {
	var (
		ret = make([]U, len(elems))
		i   int
	)
	for _, elem := range elems {
		ret[i] = fn(elem)
	}
	return ret
}

func Filter[T any](elems []T, fn func(T) bool) (ret []T) {
	for _, elem := range elems {
		if fn(elem) {
			ret = append(ret, elem)
		}
	}
	return
}

// SliceStringer joins elements of a slice whose elements are stringable.
func SliceStringer[T fmt.Stringer](elems []T, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return elems[0].String()
	}
	n := len(sep) * (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i].String())
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(elems[0].String())
	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(s.String())
	}
	return b.String()
}
