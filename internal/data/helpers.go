package data

import "golang.org/x/exp/constraints"

type Number interface {
	constraints.Integer | constraints.Float
}

func Max[T Number](a, b T) T {
	if a > b {
		return a
	} else {
		return b
	}
}
