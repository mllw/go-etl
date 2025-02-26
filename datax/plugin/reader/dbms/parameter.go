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

package dbms

import (
	"bytes"
	"database/sql"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/pingcap/errors"
)

// TableParamConfig 表参数配置
type TableParamConfig interface {
	GetColumns() []Column              //获取列信息
	GetBaseTable() *database.BaseTable //获取表信息
}

// TableParamTable 通过表参数获取对应数据库的表
type TableParamTable interface {
	Table(*database.BaseTable) database.Table //通过表参数获取对应数据库的表
}

// TableParam 表参数
type TableParam struct {
	*database.BaseParam

	Config TableParamConfig
}

// NewTableParam 获取表参数配置config，通过表参数获取对应数据库的表table和事务选项opts获取表参数
func NewTableParam(config TableParamConfig, table TableParamTable, opts *sql.TxOptions) *TableParam {
	return &TableParam{
		BaseParam: database.NewBaseParam(table.Table(config.GetBaseTable()), opts),

		Config: config,
	}
}

// Query 获取查询语句
func (t *TableParam) Query(_ []element.Record) (string, error) {
	buf := bytes.NewBufferString("select ")
	if len(t.Config.GetColumns()) == 0 {
		return "", errors.NewNoStackError("column is empty")
	}
	for i, v := range t.Config.GetColumns() {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(v.GetName())
	}
	buf.WriteString(" from ")
	buf.WriteString(t.Table().Quoted())
	buf.WriteString(" where 1 = 2")
	return buf.String(), nil
}

// Agrs 获取查询参数
func (t *TableParam) Agrs(_ []element.Record) ([]interface{}, error) {
	return nil, nil
}

// QueryParam 查询参数
type QueryParam struct {
	*database.BaseParam

	Config Config
}

// NewQueryParam 通过关系型数据库输入配置config，对应数据库表table和事务选项opts获取查询参数
func NewQueryParam(config Config, table database.Table, opts *sql.TxOptions) *QueryParam {
	return &QueryParam{
		BaseParam: database.NewBaseParam(table, opts),

		Config: config,
	}
}

// Query 获取查询语句
func (q *QueryParam) Query(_ []element.Record) (string, error) {
	if len(q.Config.GetQuerySQL()) > 1 {
		return "", errors.NewNoStackError("too much querySQL")
	}

	if len(q.Config.GetQuerySQL()) == 1 {
		return q.Config.GetQuerySQL()[0], nil
	}

	buf := bytes.NewBufferString("select ")
	if len(q.Table().Fields()) == 0 {
		return "", errors.NewNoStackError("column is empty")
	}
	for i, v := range q.Table().Fields() {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(v.Quoted())
	}
	buf.WriteString(" from ")
	buf.WriteString(q.Table().Quoted())
	if q.Config.GetWhere() != "" {
		buf.WriteString(" where ")
		buf.WriteString(q.Config.GetWhere())
	}
	return buf.String(), nil
}

// Agrs 获取查询参数
func (q *QueryParam) Agrs(_ []element.Record) (a []interface{}, err error) {
	if len(q.Config.GetQuerySQL()) > 0 {
		return nil, nil
	}

	if q.Config.GetSplitConfig().Key != "" {
		for _, v := range q.Table().Fields() {
			if q.Config.GetSplitConfig().Key == v.Name() {
				var left, right element.Column
				if left, err = q.Config.GetSplitConfig().Range.leftColumn(v.Name()); err != nil {
					return
				}
				if right, err = q.Config.GetSplitConfig().Range.rightColumn(v.Name()); err != nil {
					return
				}
				var li, ri interface{}
				if li, err = v.Valuer(left).Value(); err != nil {
					return
				}
				if ri, err = v.Valuer(right).Value(); err != nil {
					return
				}
				a = append(a, li, ri)
				return
			}
		}
	}
	return nil, nil
}

// SplitParam 切分参数
type SplitParam struct {
	*database.BaseParam

	Config Config
}

// NewSplitParam 获取表参数配置config，通过表参数获取对应数据库的表table和事务选项opts获取切分表参数
func NewSplitParam(config Config, table TableParamTable, opts *sql.TxOptions) *SplitParam {
	return &SplitParam{
		BaseParam: database.NewBaseParam(table.Table(config.GetBaseTable()), opts),

		Config: config,
	}
}

// Query 获取查询语句
func (s *SplitParam) Query(_ []element.Record) (string, error) {
	buf := bytes.NewBufferString("select ")

	buf.WriteString(s.Config.GetSplitConfig().Key)
	buf.WriteString(" from ")
	buf.WriteString(s.Table().Quoted())
	buf.WriteString(" where 1 = 2")

	return buf.String(), nil
}

// Agrs 获取查询参数
func (s *SplitParam) Agrs(_ []element.Record) ([]interface{}, error) {
	return nil, nil
}

// MinParam 最小值参数
type MinParam struct {
	*database.BaseParam

	Config Config
}

// NewMinParam 通过关系型数据库输入配置config，对应数据库表table和事务选项opts获取最小值参数
func NewMinParam(config Config, table database.Table, opts *sql.TxOptions) *MinParam {
	return &MinParam{
		BaseParam: database.NewBaseParam(table, opts),

		Config: config,
	}
}

// Query 获取查询语句
func (m *MinParam) Query(_ []element.Record) (string, error) {
	buf := bytes.NewBufferString("select min(")
	buf.WriteString(m.Config.GetSplitConfig().Key)

	buf.WriteString(") from ")
	buf.WriteString(m.Table().Quoted())
	if m.Config.GetWhere() != "" {
		buf.WriteString(" where ")
		buf.WriteString(m.Config.GetWhere())
	}
	return buf.String(), nil
}

// Agrs 获取查询参数
func (m *MinParam) Agrs(_ []element.Record) ([]interface{}, error) {
	return nil, nil
}

// MaxParam 最大值参数
type MaxParam struct {
	*database.BaseParam

	Config Config
}

// NewMaxParam 通过关系型数据库输入配置config，对应数据库表table和事务选项opts获取查询参数
func NewMaxParam(config Config, table database.Table, opts *sql.TxOptions) *MaxParam {
	return &MaxParam{
		BaseParam: database.NewBaseParam(table, opts),

		Config: config,
	}
}

// Query 获取查询语句
func (m *MaxParam) Query(_ []element.Record) (string, error) {
	buf := bytes.NewBufferString("select max(")
	buf.WriteString(m.Config.GetSplitConfig().Key)

	buf.WriteString(") from ")
	buf.WriteString(m.Table().Quoted())
	if m.Config.GetWhere() != "" {
		buf.WriteString(" where ")
		buf.WriteString(m.Config.GetWhere())
	}
	return buf.String(), nil
}

// Agrs 获取查询参数
func (m *MaxParam) Agrs(_ []element.Record) ([]interface{}, error) {
	return nil, nil
}
