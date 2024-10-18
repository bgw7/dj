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

func djRoombaVoteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "trackId")
		trackID, err := strconv.Atoi(id)
		if err != nil {
			handleError(w, fmt.Errorf("invalid trackID: %s", id))
			return
		}
		ctx := context.WithValue(r.Context(), appcontext.DJRoombaVoteCTXKey, &internal.Vote{
			TrackID: trackID,
			UserID:  r.Host,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func metdataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create new context from `r` request context, and assign key `"metadataCTXKey"`
		// to value of `"internal.Metadata"`
		client := r.Host // TODO: user from a validated oauth token
		ctx := context.WithValue(r.Context(), appcontext.MetadataCTXKey, &internal.Metadata{
			CreatedBy: client,
			UpdateBy:  &client,
		})

		// call the next handler in the chain, passing the response writer and
		// the updated request object with the new context value.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
