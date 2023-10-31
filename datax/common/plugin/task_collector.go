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

package plugin

import "github.com/mllw/go-etl/element"

// TaskCollector 任务收集器
// todo 当前未使用
type TaskCollector interface {
	CollectDirtyRecordWithError(record element.Record, err error)
	CollectDirtyRecordWithMsg(record element.Record, msgErr string)
	CollectDirtyRecord(record element.Record, err error, msgErr string)
	CollectMessage(key string, value string)
}
