package datastore

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type mockCollectableRow struct {
	fieldDescriptions func() []pgconn.FieldDescription
	scan              func(dest ...any) error
	values            func() ([]any, error)
	rawValues         func() [][]byte
}

type mockRow struct {
	mockCollectableRow
	scan func(dest ...any) error
}

func (m mockRow) Scan(dest ...any) error {
	return m.scan(dest...)
}

type mockRows struct {
	scan              func(dest ...any) error
	err               func() error
	close             func()
	next              func() bool
	fieldDescriptions func() []pgconn.FieldDescription
}

// Close implements pgx.Rows.
func (m *mockRows) Close() {
	m.close()
}

// CommandTag implements pgx.Rows.
func (m *mockRows) CommandTag() pgconn.CommandTag {
	panic("unimplemented CommandTag")
}

// Conn implements pgx.Rows.
func (m *mockRows) Conn() *pgx.Conn {
	panic("unimplemented Conn")
}

// Err implements pgx.Rows.
func (m *mockRows) Err() error {
	return m.err()
}

// FieldDescriptions implements pgx.Rows.
func (m *mockRows) FieldDescriptions() []pgconn.FieldDescription {
	return m.fieldDescriptions()
}

// Next implements pgx.Rows.
func (m *mockRows) Next() bool {
	return m.next()
}

// RawValues implements pgx.Rows.
func (m *mockRows) RawValues() [][]byte {
	panic("unimplemented RawValues")
}

// Values implements pgx.Rows.
func (m *mockRows) Values() ([]any, error) {
	panic("unimplemented Values")
}

func (m mockRows) Scan(dest ...any) error {
	return m.scan(dest...)
}

type mockConn struct {
	queryRow func(ctx context.Context, query string, args ...any) pgx.Row
	query    func(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	exec     func(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
}

func (m *mockConn) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return m.query(ctx, query, args...)
}

func (m *mockConn) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return m.exec(ctx, query, args...)
}
func (m *mockConn) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return m.queryRow(ctx, query, args...)
}

// func ptr[T any](a T) *T {
// 	return &a
// }

// func TestDB_FindOne(t *testing.T) {
// 	now := pgtype.Timestamp{}
// 	type fields struct {
// 		conn DBIface
// 	}
// 	type args struct {
// 		id string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    *internal.Reservation
// 		wantErr bool
// 	}{
// 		{
// 			name: "return no err",
// 			want: &internal.Reservation{
// 				PublicID:       ptr("25ot9"),
// 				ClientID:       ptr("client"),
// 				VenueID:        ptr("venue"),
// 				StartTimestamp: &now,
// 				EndTimestamp:   &now,
// 				Metadata: internal.Metadata{
// 					CreatedAt:   now,
// 					CreatedBy:   "by",
// 					CreatedWith: "with",
// 					UpdateAt:    &now,
// 					UpdateBy:    ptr("by"),
// 					UpdateWith:  ptr("with"),
// 				},
// 			},
// 			fields: fields{
// 				conn: &mockConn{
// 					query: func(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
// 						next := -1
// 						mockedRow := &mockRows{
// 							close: func() {},
// 							next: func() bool {
// 								next++
// 								return next == 0
// 							},
// 							err: func() error {
// 								return nil
// 							},
// 							fieldDescriptions: func() []pgconn.FieldDescription {
// 								return []pgconn.FieldDescription{
// 									{Name: "public_id"},
// 									{Name: "client_id"},
// 									{Name: "venue_id"},
// 									{Name: "start_timestamp"},
// 									{Name: "end_timestamp"},
// 									{Name: "created_at"},
// 									{Name: "created_by"},
// 									{Name: "created_with"},
// 									{Name: "updated_at"},
// 									{Name: "updated_by"},
// 									{Name: "updated_with"},
// 								}
// 							},
// 							scan: func(dest ...any) error {
// 								*(dest[0].(**string)) = ptr("25ot9")
// 								*(dest[1].(**string)) = ptr("client")
// 								*(dest[2].(**string)) = ptr("venue")
// 								*(dest[3].(**pgtype.Timestamp)) = &now
// 								*(dest[4].(**pgtype.Timestamp)) = &now
// 								*(dest[5].(*pgtype.Timestamp)) = now
// 								*(dest[6].(*string)) = "by"
// 								*(dest[7].(*string)) = "with"
// 								*(dest[8].(**pgtype.Timestamp)) = &now
// 								*(dest[9].(**string)) = ptr("by")
// 								*(dest[10].(**string)) = ptr("with")
// 								return nil
// 							},
// 						}
// 						return &mockRows{
// 							close: func() {},
// 							next: func() bool {
// 								next++
// 								return next == 0
// 							},
// 							err: func() error {
// 								return nil
// 							},
// 							scan: func(dest ...any) error {
// 								dest[0].(pgx.RowScanner).ScanRow(mockedRow)
// 								return nil
// 							},
// 						}, nil
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name:    "return ErrNoRows",
// 			wantErr: true,
// 			want:    nil,
// 			fields: fields{
// 				conn: &mockConn{
// 					query: func(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
// 						next := -1
// 						return &mockRows{
// 							close: func() {},
// 							next: func() bool {
// 								next++
// 								return next == 0
// 							},
// 							err: func() error {
// 								return nil
// 							},
// 							scan: func(dest ...any) error {
// 								return pgx.ErrNoRows
// 							},
// 						}, nil
// 					},
// 				},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			// t.Parallel()

// 			db := NewDatastore(tt.fields.conn)
// 			got, err := db.FindReservation(context.Background(), tt.args.id)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("DB.FindOne() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if diff := cmp.Diff(tt.want, got); diff != "" {
// 				t.Errorf("DB.FindOne() notEqual: -want +got \n %v", diff)
// 			}
// 		})
// 	}
// }
