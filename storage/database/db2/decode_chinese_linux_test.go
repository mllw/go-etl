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

	"github.com/mllw/go-etl/element"
	"github.com/mllw/go-etl/storage/database"
)

func TestScanner_Scan_Chinese(t *testing.T) {
	type args struct {
		src interface{}
	}
	tests := []struct {
		name    string
		s       *Scanner
		args    args
		wantErr bool
		want    element.Column
	}{
		//"CHAR"  "VARCHAR"
		{
			name: "CHAR",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("CHAR")))),
			args: args{
				src: []byte("中文abc"),
			},
			want: element.NewDefaultColumn(element.NewStringColumnValue("中文abc"), "test", element.ByteSize([]byte("中文abc"))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Scan(tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("Scanner.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.s.Column(), tt.want) {
				t.Errorf("Scanner.Column() = %v, want %v", tt.s.Column(), tt.want)
			}
		})
	}
}
