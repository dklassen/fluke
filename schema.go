package main

import (
	"fluke/parser"
	"fmt"
	"strings"
)

type SortOrder int

const (
	ASC SortOrder = iota
	DESC
)

func parseSortOrder(orderExpr *parser.OrderExprContext) SortOrder {
	if orderExpr.DESC() != nil {
		return DESC
	}
	return ASC
}

func (s SortOrder) String() string {
	switch s {
	case ASC:
		return "ASC"
	case DESC:
		return "DESC"
	default:
		panic(fmt.Sprintf("Unknown sort order: %d", s))
	}
}

type OrderByKeyInfo struct {
	Value     string
	SortOrder SortOrder
}

type PrimaryKeyInfo struct {
	Value string
}

type PartitionByInfo struct {
	Value string
}

type ColumnMap map[string]Column
type ColumnList []Column
type PrimaryKeyList []PrimaryKeyInfo
type OrderByKeyList []OrderByKeyInfo
type PartitionByList []PartitionByInfo

type Schema struct {
	TableName   string
	Columns     ColumnList
	Engine      Engine
	OrderBy     OrderByKeyList
	PrimaryKey  PrimaryKeyList
	PartitionBy PartitionByList
}

func (s Schema) ColumnMap() ColumnMap {
	m := make(ColumnMap)
	for _, col := range s.Columns {
		m[col.Value] = col
	}
	return m
}

type Engine struct {
	Value string
}

type PrimaryKeyStmt struct {
	Value string
}

type PartitionByStmt struct {
	Value string
}

type OrderByStmt struct{}

type Column struct {
	Value string
	Type  string
}

// NOTE: UGH use the Darn AST to do this
func (partitionByList PartitionByList) String() string {
	if len(partitionByList) == 1 {
		return fmt.Sprintf("%s", partitionByList[0].Value)
	}

	var parts []string
	for _, partitionByInfo := range partitionByList {
		parts = append(parts, fmt.Sprintf("%s", partitionByInfo.Value))
	}
	return fmt.Sprintf("PARTITION BY (%s)", strings.Join(parts, ", "))
}

// NOTE: UGH use the Darn AST to do this
func (primaryKeyList PrimaryKeyList) String() string {
	if len(primaryKeyList) == 1 {
		return fmt.Sprintf("%s", primaryKeyList[0].Value)
	}

	var parts []string
	for _, primaryKeyInfo := range primaryKeyList {
		parts = append(parts, fmt.Sprintf("%s", primaryKeyInfo.Value))
	}
	return fmt.Sprintf("(%s)", strings.Join(parts, ", "))
}

// NOTE: UGH use the Darn AST to do this
func (orderByKeyStmt OrderByKeyList) String() string {

	if len(orderByKeyStmt) == 1 {
		if orderByKeyStmt[0].SortOrder == ASC {
			return fmt.Sprintf("%s", orderByKeyStmt[0].Value)
		}

		return fmt.Sprintf("%s %s", orderByKeyStmt[0].Value, orderByKeyStmt[0].SortOrder)
	}

	var parts []string
	for _, orderByEntry := range orderByKeyStmt {
		if orderByEntry.SortOrder == ASC {
			parts = append(parts, fmt.Sprintf("%s", orderByEntry.Value))
		} else {
			parts = append(parts, fmt.Sprintf("%s %s", orderByEntry.Value, orderByEntry.SortOrder))
		}
	}
	return fmt.Sprintf("(%s)", strings.Join(parts, ", "))
}
