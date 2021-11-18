// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

package pbpgx

import (
	"reflect"
	"testing"

	"github.com/muhlemmer/pbpgx/internal/support"
	"google.golang.org/protobuf/proto"
)

type stringer string

func (s stringer) String() string {
	return string(s)
}

func Test_colSlice(t *testing.T) {
	tests := []struct {
		name    string
		columns []stringer
		want    []string
		want1   int
	}{
		{
			"nil columns",
			nil,
			nil,
			1,
		},
		{
			"one columns",
			[]stringer{"foo"},
			[]string{"foo"},
			5,
		},
		{
			"three columns",
			[]stringer{"foo", "bar", "foo"},
			[]string{"foo", "bar", "foo"},
			19,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := colSlice(tt.columns)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("colSlice() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("colSlice() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Benchmark_colSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		colSlice([]stringer{"id", "title", "price"})
	}
}

func Test_readOneQuery(t *testing.T) {
	tests := []struct {
		name    string
		schema  string
		columns []stringer
		want    string
		want1   int
	}{
		{
			"nil columns",
			"",
			nil,
			`SELECT * FROM "public"."products" WHERE "id" = $1 LIMIT 1;`,
			58,
		},
		{
			"with columns",
			"foo",
			[]stringer{"id", "title", "price"},
			`SELECT "id", "title", "price" FROM "foo"."products" WHERE "id" = $1 LIMIT 1;`,
			76,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := readOneQuery(tt.schema, "products", tt.columns)
			if got != tt.want {
				t.Errorf("readOneQuery() got =\n%v\nwant\n%v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("readOneQuery() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Benchmark_readOneQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		readOneQuery("public", "products", []stringer{"id", "title", "price"})
	}
}

// Benchmark_colSlice-16        	 6726246	       256.3 ns/op	      96 B/op	       4 allocs/op
// Benchmark_readOneQuery-16    	 1902990	       528.8 ns/op	     176 B/op	       5 allocs/op

func TestReadOne(t *testing.T) {
	req := &support.SimpleQuery{
		Id:      2,
		Columns: []support.SimpleColumns{support.SimpleColumns_id, support.SimpleColumns_title},
	}

	want := &support.Simple{
		Id:    2,
		Title: "two",
	}

	got, err := ReadOne[int32, support.SimpleColumns, *support.Simple](testCtx, connPool, "public", "simple_ro", req)
	if err != nil {
		t.Fatal(err)
	}

	if !proto.Equal(got, want) {
		t.Errorf("ReadOne() =\n%v\nwant\n%v", got, want)
	}
}
