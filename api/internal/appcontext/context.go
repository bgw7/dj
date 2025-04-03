package appcontext

import (
	"context"

	"github.com/bgw7/dj/internal"
)

type ctxKey string

const (
	VoteKey     ctxKey = "voteRequest"
	MetadataKey ctxKey = "metadata"
)

// SetVoteRequest stores a VoteRequest in the context.
func SetVoteRequest(ctx context.Context, vote *internal.VoteRequest) context.Context {
	return context.WithValue(ctx, VoteKey, vote)
}

// GetVoteRequest retrieves a VoteRequest from the context.
func GetVoteRequest(ctx context.Context) (*internal.VoteRequest, bool) {
	vote, ok := ctx.Value(VoteKey).(*internal.VoteRequest)
	return vote, ok
}

// SetMetadata stores Metadata in the context.
func SetMetadata(ctx context.Context, metadata *internal.Metadata) context.Context {
	return context.WithValue(ctx, MetadataKey, metadata)
}

// GetMetadata retrieves Metadata from the context.
func GetMetadata(ctx context.Context) (*internal.Metadata, bool) {
	meta, ok := ctx.Value(MetadataKey).(*internal.Metadata)
	return meta, ok
}
