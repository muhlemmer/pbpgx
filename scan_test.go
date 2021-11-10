// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

package pbpgx

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/muhlemmer/pbpgx/internal/support"
	"google.golang.org/protobuf/proto"
	pr "google.golang.org/protobuf/reflect/protoreflect"
)

func Test_convertIntValueFunc(t *testing.T) {
	f := convertIntValueFunc[uint32](pr.ValueOfInt32)
	var u uint32 = 222
	v := f(u)

	if got := v.Interface().(int32); got != 222 {
		t.Errorf("convertIntValueFunc.f() got: %v, want 222", got)
	}
}

func Test_destination_value(t *testing.T) {
	type fields struct {
		valueDecoder valueDecoder
		valueFunc    func(interface{}) pr.Value
	}
	tests := []struct {
		name   string
		fields fields
		want   pr.Value
	}{
		{
			"test",
			fields{
				valueDecoder: &pgtype.Float4{
					Float:  1.1,
					Status: pgtype.Present,
				},
				valueFunc: pr.ValueOf,
			},
			pr.ValueOfFloat32(1.1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &destination{
				valueDecoder: tt.fields.valueDecoder,
				valueFunc:    tt.fields.valueFunc,
			}
			if got := d.value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("destination.value() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testRows struct {
	names []string
	rows  [][]interface{}

	closed bool
	err    error
	pos    int
}

func newTestRows(colNames []string, rows [][]interface{}) pgx.Rows {
	r := &testRows{
		names: colNames,
		rows:  rows,
		pos:   -1,
	}

	return r
}

func (r *testRows) Close()                        { r.closed = true }
func (r *testRows) Err() error                    { return r.err }
func (r *testRows) CommandTag() pgconn.CommandTag { return nil } // no-op
func (r *testRows) RawValues() [][]byte           { return nil } // no-op

func (r *testRows) FieldDescriptions() []pgproto3.FieldDescription {
	fds := make([]pgproto3.FieldDescription, len(r.names))

	for i, n := range r.names {
		fds[i].Name = []byte(n)
	}

	return fds
}

func (r *testRows) Next() bool {
	r.pos++

	if r.pos >= len(r.rows) || r.closed {
		return false
	}

	return true
}

func (r *testRows) Values() ([]interface{}, error) {
	if r.closed {
		r.err = errors.New("Rows closed")
		return nil, r.err
	}

	values := r.rows[r.pos]
	if len(values) != len(r.names) {
		r.err = fmt.Errorf("len of values %d != len of names %d", len(values), len(r.names))
		return nil, r.err
	}

	return values, nil
}

func (r *testRows) Scan(dest ...interface{}) error {
	values, err := r.Values()
	if err != nil {
		return err
	}

	if len(values) != len(dest) {
		r.err = fmt.Errorf("len of values %d != len of dest %d", len(values), len(dest))
		return r.err
	}

	for i, dst := range dest {
		if dst == nil {
			continue
		}

		v := dst.(pgtype.Value)
		v.Set(values[i])
	}

	return nil
}

func TestScan(t *testing.T) {
	type args struct {
		names []string
		rows  [][]interface{}
	}
	tests := []struct {
		name        string
		args        args
		wantResults []*support.Supported
		wantErr     bool
	}{
		{
			"Destination error",
			args{
				[]string{"foo", "bar"},
				[][]interface{}{
					{1, 2},
				},
			},
			nil,
			true,
		},
		{
			"Scan error",
			args{
				[]string{"bl", "i32", "i64", "f", "d", "s", "bt", "u32", "u64"},
				[][]interface{}{
					{1, 2},
				},
			},
			nil,
			true,
		},
		{
			"All supported",
			args{
				[]string{"bl", "i32", "i64", "f", "d", "s", "bt", "u32", "u64"},
				[][]interface{}{
					{true, int32(1), int64(2), float32(1.1), float64(2.2), "Hello World!", []byte("Foo bar"), uint32(32), uint64(64)},
					{false, int32(3), int64(4), float32(3.1), float64(4.2), "Bye World!", []byte{}, uint32(23), uint64(46)},
				},
			},
			[]*support.Supported{
				{
					Bl:  true,
					I32: 1,
					I64: 2,
					F:   1.1,
					D:   2.2,
					S:   "Hello World!",
					Bt:  []byte("Foo bar"),
					U32: 32,
					U64: 64,
				},
				{
					Bl:  false,
					I32: 3,
					I64: 4,
					F:   3.1,
					D:   4.2,
					S:   "Bye World!",
					Bt:  []byte{},
					U32: 23,
					U64: 46,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResults, err := Scan[*support.Supported](newTestRows(tt.args.names, tt.args.rows))
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(gotResults) != len(tt.wantResults) {
				t.Errorf("Scan() =\n%s\nwant\n%s", gotResults, tt.wantResults)
			}

			for i, want := range tt.wantResults {
				if !proto.Equal(gotResults[i], want) {
					t.Errorf("Scan() =\n%s\nwant\n%s", gotResults[i], want)
				}
			}
		})
	}
}

func TestScan_unsuported2(t *testing.T) {
	type args struct {
		names []string
		rows  [][]interface{}
	}
	tests := []struct {
		name        string
		args        args
		wantResults []*support.Supported
		wantErr     bool
	}{
		{
			"Repeated fields",
			args{
				[]string{"bl", "i32"},
				[][]interface{}{
					{[]bool{true, false}, []int32{1, 2}},
					{[]bool{true, false, false}, []int32{3, 4}},
				},
			},
			nil,
			true,
		},
		{
			"Message type",
			args{
				[]string{"sup"},
				[][]interface{}{
					{&support.Supported{Bl: true}},
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResults, err := Scan[*support.Unsupported](newTestRows(tt.args.names, tt.args.rows))
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(gotResults) != len(tt.wantResults) {
				t.Errorf("Scan() =\n%s\nwant\n%s", gotResults, tt.wantResults)
			}

			for i, want := range tt.wantResults {
				if !proto.Equal(gotResults[i], want) {
					t.Errorf("Scan() =\n%s\nwant\n%s", gotResults[i], want)
				}
			}
		})
	}
}
