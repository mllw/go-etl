// Copyright 2020 the go-etl Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db2

import (
	"reflect"
	"testing"

	_ "github.com/ibmdb/go_ibm_db"
	"github.com/mllw/go-etl/config"
	"github.com/mllw/go-etl/storage/database"
)

func testJSONFromString(s string) *config.JSON {
	json, err := config.NewJSONFromString(s)
	if err != nil {
		panic(err)
	}
	return json
}

func testSource(s string) *Source {
	source, err := NewSource(database.NewBaseSource(testJSONFromString(s)))
	if err != nil {
		panic(err)
	}
	return source.(*Source)
}

func TestDialect_Source(t *testing.T) {
	type args struct {
		bs *database.BaseSource
	}
	tests := []struct {
		name    string
		d       Dialect
		args    args
		want    database.Source
		wantErr bool
	}{
		{
			name: "1",
			d:    Dialect{},
			args: args{
				bs: database.NewBaseSource(testJSONFromString(`{
					"url":"HOSTNAME=192.168.0.1;PORT=50000;DATABASE=testdb",
					"username" : "user",
					"password": "passwd"
				}`)),
			},
			want: &Source{
				BaseSource: database.NewBaseSource(testJSONFromString(`{
					"url":"HOSTNAME=192.168.0.1;PORT=50000;DATABASE=testdb",
					"username" : "user",
					"password": "passwd"
				}`)),
				dsn: "DATABASE=testdb;HOSTNAME=192.168.0.1;PORT=50000;PWD=passwd;UID=user",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Source(tt.args.bs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dialect.Source() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dialect.Source() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDialect_Name(t *testing.T) {
	tests := []struct {
		name string
		d    Dialect
		want string
	}{
		{
			name: "1",
			d:    Dialect{},
			want: "db2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Name(); got != tt.want {
				t.Errorf("Dialect.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSource(t *testing.T) {
	type args struct {
		bs *database.BaseSource
	}
	tests := []struct {
		name    string
		args    args
		wantS   database.Source
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				bs: database.NewBaseSource(testJSONFromString(`{
					"url":"HOSTNAME=192.168.0.1;PORT=50000;DATABASE=testdb",
					"username" : "user",
					"password": "passwd"
				}`)),
			},
			wantS: &Source{
				BaseSource: database.NewBaseSource(testJSONFromString(`{
					"url":"HOSTNAME=192.168.0.1;PORT=50000;DATABASE=testdb",
					"username" : "user",
					"password": "passwd"
				}`)),
				dsn: "DATABASE=testdb;HOSTNAME=192.168.0.1;PORT=50000;PWD=passwd;UID=user",
			},
		},
		{
			name: "2",
			args: args{
				bs: database.NewBaseSource(testJSONFromString(`{
					"url":"HOSTNAME=192.168.0.1;PORT=50000;DATABASE=testdb",
					"username" : 1,
					"password": "passwd"
				}`)),
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				bs: database.NewBaseSource(testJSONFromString(`{
					"url":"HOSTNAME=192.168.0.1;DATABASE=testdb",
					"username" : "user",
					"password": "passwd"
				}`)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := NewSource(tt.args.bs)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("NewSource() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestSource_DriverName(t *testing.T) {
	tests := []struct {
		name string
		s    *Source
		want string
	}{
		{
			name: "1",
			s:    &Source{},
			want: "go_ibm_db",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.DriverName(); got != tt.want {
				t.Errorf("Source.DriverName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSource_ConnectName(t *testing.T) {
	tests := []struct {
		name string
		s    *Source
		want string
	}{
		{
			name: "1",
			s: testSource(`{
				"url":"HOSTNAME=192.168.0.1;PORT=50000;DATABASE=testdb",
				"username" : "user",
				"password": "passwd"
			}`),
			want: "DATABASE=testdb;HOSTNAME=192.168.0.1;PORT=50000;PWD=passwd;UID=user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.ConnectName(); got != tt.want {
				t.Errorf("Source.ConnectName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSource_Key(t *testing.T) {
	tests := []struct {
		name string
		s    *Source
		want string
	}{
		{
			name: "1",
			s: testSource(`{
				"url":"HOSTNAME=192.168.0.1;PORT=50000;DATABASE=testdb",
				"username" : "user",
				"password": "passwd"
			}`),
			want: "DATABASE=testdb;HOSTNAME=192.168.0.1;PORT=50000;PWD=passwd;UID=user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Key(); got != tt.want {
				t.Errorf("Source.Key() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSource_Table(t *testing.T) {
	type args struct {
		b *database.BaseTable
	}
	tests := []struct {
		name string
		s    *Source
		args args
		want database.Table
	}{
		{
			name: "1",
			s: testSource(`{
				"url":"HOSTNAME=192.168.0.1;PORT=50000;DATABASE=testdb",
				"username" : "user",
				"password": "passwd"
			}`),
			args: args{
				b: database.NewBaseTable("", "schema", "table"),
			},
			want: NewTable(database.NewBaseTable("", "schema", "table")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Table(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Source.Table() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuoted(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				s: "table",
			},
			want: `"table"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Quoted(tt.args.s); got != tt.want {
				t.Errorf("Quoted() = %v, want %v", got, tt.want)
			}
		})
	}
}
