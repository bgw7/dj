package service

import (
	"context"
	"time"

	"github.com/la-viajera/reservation-service/internal"
	"github.com/la-viajera/reservation-service/internal/appcontext"
)

func (d *DomainService) GetVenues(ctx context.Context) ([]internal.Venue, error) {
	return d.datastore.GetVenues(ctx)
}

func (d *DomainService) CreateVenue(ctx context.Context, obj *internal.Venue) (*internal.Venue, error) {
	metdata, err := appcontext.FromContext[*internal.Metadata](ctx, appcontext.MetadataCTXKey)
	if err != nil {
		return nil, err
	}
	obj.Metadata = *metdata
	return d.datastore.CreateVenue(
		ctx,
		obj,
	)
}

func (d *DomainService) SuggestVenues(ctx context.Context, query string) (*internal.PagedResponse[internal.Venue], error) {
	println(query)
	return &internal.PagedResponse[internal.Venue]{
		Offset:  0,
		Results: []internal.Venue{},
	}, nil
}

func (d *DomainService) SearchVenues(ctx context.Context, s *internal.VenueSearch) (*internal.PagedResponse[internal.Venue], error) {
	if s.Limit == 0 {
		s.Limit = 10
	}
	if s.StartTimestamp != nil {
		d, err := time.Parse(time.DateOnly, *s.StartTimestamp)
		if err != nil {
			return nil, err
		}
		*s.StartTimestamp = d.Format(time.DateOnly)
	}
	records, err := d.datastore.SearchVenues(ctx, s)
	if err != nil {
		return nil, err
	}
	return &internal.PagedResponse[internal.Venue]{
		Offset:  s.Offset,
		Results: records,
	}, nil

}
