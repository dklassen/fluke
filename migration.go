package main

import (
	"fmt"
	"sort"
	"strings"
)

type Statement interface {
	GenerateSQL() string
}

type SchemaTables map[string]Schema

type MigrationStream struct {
	Statements []Statement
}

type CreateTableStmt struct {
	Table Schema
}

type DropTableStmt struct {
	Table Schema
}

type AlterTableStmt struct {
	Change TableChange
}

func (s AlterTableStmt) GenerateSQL() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("ALTER TABLE %s", s.Change.Table.TableName))

	for _, col := range s.Change.ColumnsAdded {
		b.WriteString(fmt.Sprintf(" ADD COLUMN %s %s", col.Value, col.Type))
	}
	for _, col := range s.Change.ColumnsRemoved {
		b.WriteString(fmt.Sprintf(" DROP COLUMN %s", col.Value))
	}
	for _, col := range s.Change.ColumnsModified {
		b.WriteString(fmt.Sprintf(" MODIFY COLUMN %s %s", col.Value, col.Type))
	}

	b.WriteString(";\n")
	return b.String()
}

func (s CreateTableStmt) GenerateSQL() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("CREATE TABLE %s (", s.Table.TableName))

	var cols []string
	for _, col := range s.Table.Columns {
		cols = append(cols, fmt.Sprintf("%s %s", col.Value, col.Type))
	}
	b.WriteString(strings.Join(cols, ", "))

	b.WriteString(fmt.Sprintf(") ENGINE = %s", s.Table.Engine.Value))

	writeClause := func(keyword string, clause fmt.Stringer) {
		if str := clause.String(); str != "" {
			b.WriteString(fmt.Sprintf(" %s %s", keyword, str))
		}
	}

	if len(s.Table.PartitionBy) > 0 {
		writeClause("PARTITION BY", s.Table.PartitionBy)
	}

	if len(s.Table.PrimaryKey) > 0 {
		writeClause("PRIMARY KEY", s.Table.PrimaryKey)
	}
	if len(s.Table.OrderBy) > 0 {
		writeClause("ORDER BY", s.Table.OrderBy)
	}

	b.WriteString(";\n")
	return b.String()
}

func (s DropTableStmt) GenerateSQL() string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s;\n", s.Table.TableName)
}

// TODO:: Create schema diffing function which holds changes between two schemas
// to inform the detection logic and avoid unnecessary re-processing.
func GenerateMigrationStream(
	oldSchema, newSchema SchemaTables,
	differ Differ,
) MigrationStream {
	var stmts []Statement
	diff := differ(oldSchema, newSchema)

	for _, table := range diff.AddedTables {
		stmts = append(stmts, CreateTableStmt{Table: table})
	}

	for _, table := range diff.RemovedTables {
		stmts = append(stmts, DropTableStmt{Table: table})
	}

	for _, change := range diff.ModifiedTables {
		if change.EngineChange != (EngineChange{}) {
			stmts = append(stmts, DropTableStmt{Table: change.Table})
			stmts = append(stmts, CreateTableStmt{Table: change.Table})
		} else {
			stmts = append(stmts, AlterTableStmt{Change: change})
		}
	}

	return MigrationStream{Statements: stmts}
}

func OrderBySliceEqual(a, b []OrderByKeyInfo) bool {
	if len(a) != len(b) {
		return false
	}

	sort.Slice(a, func(i, j int) bool { return a[i].Value < a[j].Value })
	sort.Slice(b, func(i, j int) bool { return b[i].Value < b[j].Value })

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// GenerateMigration compares the old and new schemas and generates migration SQL.
func GenerateMigration(oldSchema, newSchema SchemaTables) string {
	migrationStream := GenerateMigrationStream(
		oldSchema,
		newSchema,
		SchemaDiffer,
	)

	sql := strings.Builder{}
	for _, statement := range migrationStream.Statements {
		sql.WriteString(statement.GenerateSQL())
	}

	return sql.String()
}
