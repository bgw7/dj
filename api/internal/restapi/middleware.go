package restapi

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bgw7/dj/internal"
	"github.com/bgw7/dj/internal/appcontext"
	"github.com/go-chi/chi/v5"
)

func voteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "trackId")
		trackId, err := strconv.Atoi(id)
		if err != nil {
			handleError(w, fmt.Errorf("invalid trackID: %s", id))
			return
		}
		ctx := context.WithValue(r.Context(), appcontext.DJRoombaVoteCTXKey, &internal.Vote{
			TrackID: trackId,
			VoterID: r.Host,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func metdataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := r.Host // TODO: user from a validated authorization token
		ctx := context.WithValue(r.Context(), appcontext.MetadataCTXKey, &internal.Metadata{
			CreatedBy: client,
			UpdateBy:  &client,
		})

		// call the next handler in the chain, passing the response writer and
		// the updated request object with the new context value.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
