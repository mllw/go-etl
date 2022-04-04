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

import (
	"context"

	"github.com/Breeze0806/go-etl/storage/database"
)

//TableDiffer 表不同
type TableDiffer struct {
	MasterTable database.Table
	SlaveTable  database.Table
	Differ      Differ
}

//DifferStorage 差异存储
type DifferStorage interface {
	Write(ctx context.Context, differs []TableDiffer) error
	Read(ctx context.Context,
		onDiffer func(differ TableDiffer) error) error
}

//DifferStorageMaker 差异存储生成器
type DifferStorageMaker interface {
	DifferStorage() DifferStorage
}
