/*
SPDX-License-Identifier: AGPL-3.0-only

Copyright (C) 2021, Tim Möhlmann

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package crud

import (
	"testing"

	"github.com/muhlemmer/pbpgx/internal/support"
	"github.com/muhlemmer/pbpgx/internal/testlib"
	"google.golang.org/protobuf/proto"
)

func TestReadOne(t *testing.T) {
	want := &support.Simple{
		Id:    5,
		Title: "five",
		Data:  "five is a four letter word",
	}

	tab := NewTable("public", "simple_ro", nil)

	got, err := ReadOne[*support.Simple](testlib.CTX, testlib.ConnPool, tab, 5, ColumnWildcard)
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(got, want) {
		t.Errorf("ReadOne() = %v, want %v", got, want)
	}
}

func TestReadAll(t *testing.T) {
	cols := []support.SimpleColumns{
		support.SimpleColumns_title,
		support.SimpleColumns_data,
	}

	want := []*support.Simple{
		{
			Title: "one",
			Data:  "foo bar",
		},
		{
			Title: "two",
		},
		{
			Data: "golden triangle",
		},
		{
			Title: "four",
			Data:  "hello world",
		},
	}

	tab := NewTable("public", "simple_ro", nil)

	got, err := ReadAll[*support.Simple](testlib.CTX, testlib.ConnPool, tab, 4, cols)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != len(want) {
		t.Fatalf("ReadAll() =\n%v\nwant\n%v", got, want)
	}

	for i, w := range want {
		if !proto.Equal(got[i], w) {
			t.Fatalf("ReadAll() =\n%v\nwant\n%v", got[i], w)
		}
	}
}

func TestReadList(t *testing.T) {
	ids := []int{1, 4, 5}

	cols := []support.SimpleColumns{
		support.SimpleColumns_title,
		support.SimpleColumns_data,
	}

	want := []*support.Simple{
		{
			Title: "one",
			Data:  "foo bar",
		},
		{
			Title: "four",
			Data:  "hello world",
		},
		{
			Title: "five",
			Data:  "five is a four letter word",
		},
	}

	tab := NewTable("public", "simple_ro", nil)

	got, err := ReadList[*support.Simple](testlib.CTX, testlib.ConnPool, tab, ids, cols...)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != len(want) {
		t.Fatalf("ReadAll() =\n%v\nwant\n%v", got, want)
	}

	for i, w := range want {
		if !proto.Equal(got[i], w) {
			t.Fatalf("ReadAll() =\n%v\nwant\n%v", got[i], w)
		}
	}
}
