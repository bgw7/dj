package restapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bgw7/dj/internal"
	"github.com/bgw7/dj/internal/appcontext"
	"github.com/go-chi/chi/v5"
)

// Strongly typed context keys to avoid collisions
type ctxKey string

const (
	voteKey     ctxKey = "voteRequest"
	metadataKey ctxKey = "metadata"
)

// voteMiddleware extracts track ID from the URL and stores it in the request context.
func voteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trackIDStr := chi.URLParam(r, "trackId")
		trackID, err := strconv.Atoi(trackIDStr)
		if err != nil {
			handleError(w, fmt.Errorf("invalid trackID: %s", trackIDStr))
			return
		}

		vote := &internal.VoteRequest{
			TrackID: trackID,
			VoterID: r.Host, // Ideally, replace this with an authenticated user ID
		}

		// Store vote request in context
		ctx := appcontext.SetVoteRequest(r.Context(), vote)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// metadataMiddleware attaches metadata to the request context.
func metadataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := r.Host // TODO: Extract user from a validated auth token

		metadata := &internal.Metadata{
			CreatedBy: client,
			UpdateBy:  &client,
		}

		// Store metadata in context
		ctx := appcontext.SetMetadata(r.Context(), metadata)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func voteIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		voteIDStr := chi.URLParam(r, "voteId")
		voteID, err := strconv.Atoi(voteIDStr)
		if err != nil {
			handleError(w, fmt.Errorf("invalid voteID: %s", voteIDStr))
			return
		}

		// Get existing VoteRequest from context
		voteReq, ok := appcontext.GetVoteRequest(r.Context())
		if !ok {
			handleError(w, fmt.Errorf("vote middleware must be used before voteIDMiddleware"))
			return
		}

		// Update VoteRequest with the VoteID
		voteReq.ID = voteID
		ctx := appcontext.SetVoteRequest(r.Context(), voteReq)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
