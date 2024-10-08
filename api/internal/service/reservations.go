package service

import (
	"context"
	"time"

	"github.com/la-viajera/reservation-service/internal"
	"github.com/la-viajera/reservation-service/internal/appcontext"
)

func (d *DomainService) GetReservations(ctx context.Context) ([]internal.Reservation, error) {
	return d.datastore.GetReservations(ctx)
}

func (d *DomainService) CreateReservation(ctx context.Context, obj *internal.Reservation) (*internal.Reservation, error) {
	metdata, err := appcontext.Value[*internal.Metadata](ctx, appcontext.MetadataCTXKey)
	if err != nil {
		return nil, err
	}
	obj.Metadata = *metdata
	return d.datastore.CreateReservation(
		ctx,
		obj,
	)
}

func (d *DomainService) SearchReservations(ctx context.Context, s *internal.ReservationSearch) (*internal.PagedResponse[internal.Reservation], error) {
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
	records, err := d.datastore.SearchReservations(ctx, s)
	if err != nil {
		return nil, err
	}
	return &internal.PagedResponse[internal.Reservation]{
		Offset:  s.Offset,
		Results: records,
	}, nil
}

func (d *DomainService) FindOneReservation(ctx context.Context) (*internal.Reservation, error) {
	id, err := appcontext.Value[string](ctx, appcontext.ReservationIDCTXKey)
	if err != nil {
		return nil, err
	}
	return d.datastore.FindReservation(ctx, id)
}

func (d *DomainService) UpdateReservation(ctx context.Context, r *internal.Reservation) (*internal.Reservation, error) {
	id, err := appcontext.Value[string](ctx, appcontext.ReservationIDCTXKey)
	if err != nil {
		return nil, err
	}

	r.ID = &id

	metadata, err := appcontext.Value[*internal.Metadata](ctx, appcontext.MetadataCTXKey)
	if err != nil {
		return nil, err
	}
	r.Metadata = *metadata
	return d.datastore.UpdateReservation(ctx, r)
}
