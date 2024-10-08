package appcontext

import (
	"context"
	"errors"
	"fmt"
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "app context value " + k.name
}

var (
	ReservationIDCTXKey = &contextKey{"ReservationIDContext"}
	MetadataCTXKey      = &contextKey{"MetadataContext"}
	DJRoombaVoteCTXKey  = &contextKey{"DJRoombaVoteContext"}
)

var ErrCtxKeyNotFound = errors.New("key not found in context")

func Value[T any](ctx context.Context, key *contextKey) (T, error) {
	val := ctx.Value(key)
	if val == nil {
		var val T
		return val, fmt.Errorf("%s not found: %w", key, ErrCtxKeyNotFound)
	}
	return val.(T), nil
}
