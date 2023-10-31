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

package core

import (
	"reflect"
	"testing"

	"github.com/mllw/go-etl/config"
	"github.com/mllw/go-etl/datax/core/statistics/container"
)

func testJSONFromString(s string) *config.JSON {
	j, err := config.NewJSONFromString(s)
	if err != nil {
		panic(err)
	}
	return j
}

func TestBaseCotainer_SetConfig(t *testing.T) {
	type args struct {
		conf *config.JSON
	}
	tests := []struct {
		name string
		b    *BaseCotainer
		args args
		want *config.JSON
	}{
		{
			name: "1",
			b:    NewBaseCotainer(),
			args: args{
				conf: testJSONFromString("{}"),
			},
			want: testJSONFromString("{}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetConfig(tt.args.conf)
			if got := tt.b.Config(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config() = %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestBaseCotainer_SetGetMetrics(t *testing.T) {
	type testStruct struct {
		Path string `json:"path"`
	}

	type args struct {
		metrics *container.Metrics
	}
	m := container.NewMetrics()
	m.Set("test", testStruct{Path: "value"})
	tests := []struct {
		name string
		b    *BaseCotainer
		args args
		want string
	}{
		{
			name: "1",
			b:    NewBaseCotainer(),
			args: args{
				metrics: m,
			},
			want: `{"test":{"path":"value"}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetMetrics(tt.args.metrics)
			if got := tt.b.Metrics().JSON().String(); got != tt.want {
				t.Errorf("Metrics() = %v, want: %v", got, tt.want)
			}
		})
	}
}
