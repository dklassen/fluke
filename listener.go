package main

import (
	"fluke/parser"
)

type ClickHouseSQLListener struct {
	*parser.BaseClickHouseParserListener
	CurrentTable Schema
	Tables       map[string]Schema
}

func NewClickHouseSQLListener() *ClickHouseSQLListener {
	return &ClickHouseSQLListener{
		Tables: make(map[string]Schema),
	}
}

func ParseOrderByClause(ctx *parser.OrderByClauseContext) []OrderByKeyInfo {
	var orderBy []OrderByKeyInfo

	if ctx == nil || ctx.OrderExprList() == nil {
		return orderBy
	}

	for _, orderExpr := range ctx.OrderExprList().AllOrderExpr() {
		colExpr := orderExpr.ColumnExpr(0)
		if colExpr == nil {
			continue
		}

		switch actual := colExpr.(type) {
		case *parser.ColumnExprIdentifierContext:
			orderBy = append(orderBy, OrderByKeyInfo{
				Value:     actual.GetText(),
				SortOrder: parseSortOrder(orderExpr.(*parser.OrderExprContext)),
			})

		case *parser.ColumnExprTupleContext:
			for _, child := range actual.ColumnExprList().GetChildren() {
				if ident, ok := child.(*parser.ColumnsExprColumnContext); ok {
					orderBy = append(orderBy, OrderByKeyInfo{
						Value:     ident.GetText(),
						SortOrder: parseSortOrder(orderExpr.(*parser.OrderExprContext)),
					})
				}
			}
		}
	}

	return orderBy
}

func ParsePrimaryKeyClause(clause *parser.PrimaryKeyClauseContext) []PrimaryKeyInfo {
	var keys []PrimaryKeyInfo

	if clause == nil || clause.ColumnExpr() == nil {
		return keys
	}

	switch expr := clause.ColumnExpr().(type) {
	case *parser.ColumnExprIdentifierContext:
		keys = append(keys, PrimaryKeyInfo{Value: expr.GetText()})
	case *parser.ColumnExprTupleContext:
		for _, child := range expr.ColumnExprList().GetChildren() {
			if ident, ok := child.(*parser.ColumnsExprColumnContext); ok {
				keys = append(keys, PrimaryKeyInfo{Value: ident.GetText()})
			}
		}
	}
	return keys
}

func ParsePartitionByClause(clause *parser.PartitionByClauseContext) []PartitionByInfo {
	var partitionBy []PartitionByInfo

	if clause == nil || clause.ColumnExpr() == nil {
		return partitionBy
	}

	for _, child := range clause.ColumnExpr().GetChildren() {
		if col, ok := child.(*parser.ColumnExprIdentifierContext); ok {
			partitionBy = append(partitionBy, PartitionByInfo{
				Value: col.GetText(),
			})
		}
	}

	return partitionBy
}

func (l *ClickHouseSQLListener) EnterEngineClause(ctx *parser.EngineClauseContext) {
	l.CurrentTable.Engine = Engine{Value: ctx.EngineExpr().IdentifierOrNull().GetText()}
}

func (l *ClickHouseSQLListener) EnterOrderByClause(ctx *parser.OrderByClauseContext) {
	l.CurrentTable.OrderBy = ParseOrderByClause(ctx)
}

func (l *ClickHouseSQLListener) EnterPrimaryKeyClause(ctx *parser.PrimaryKeyClauseContext) {
	l.CurrentTable.PrimaryKey = ParsePrimaryKeyClause(ctx)
}

func (l *ClickHouseSQLListener) EnterPartitionByClause(ctx *parser.PartitionByClauseContext) {
	l.CurrentTable.PartitionBy = ParsePartitionByClause(ctx)
}

func (l *ClickHouseSQLListener) EnterTableColumnDfnt(ctx *parser.TableColumnDfntContext) {
	colName := ctx.NestedIdentifier().GetText()
	colType := ctx.ColumnTypeExpr().GetText()

	l.CurrentTable.Columns = append(l.CurrentTable.Columns, Column{
		Value: colName,
		Type:  colType,
	})
}

// NOTE:: Parsing isn't error proof if schema is malformed and when we encounter we panic instead
// of recovering gracefully and pointing out where the error is.
func (l *ClickHouseSQLListener) EnterCreateTableStmt(ctx *parser.CreateTableStmtContext) {
	tableName := ctx.TableIdentifier().GetText()

	l.CurrentTable = Schema{
		TableName: tableName,
		Columns:   []Column{},
	}
}

func (l *ClickHouseSQLListener) ExitCreateTableStmt(ctx *parser.CreateTableStmtContext) {
	l.Tables[l.CurrentTable.TableName] = l.CurrentTable
}
