package data

import (
	"context"
	"time"

	"golang.org/x/exp/constraints"
)

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

func GetDefaultTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}
