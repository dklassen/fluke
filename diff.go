package main

// We are using a semantic diffing approach to capture the changes.
// This is a bit more complex than a simple set based diff, but it allows
// us to capture the changes in a way that is more human readable and
// we can pass this struct around to the various detection functions.
type SemanticSchemaDiff struct {
	AddedTables    []Schema
	RemovedTables  []Schema
	ModifiedTables []TableChange
}

// Represents a collection of changes to a existing table
type TableChange struct {
	Table           Schema
	EngineChange    EngineChange
	ColumnsRemoved  []Column
	ColumnsAdded    []Column
	ColumnsModified []Column
}

type PrimaryKeyChange struct {
	Old PrimaryKeyStmt
	New PrimaryKeyStmt
}

type PartitionByChange struct {
	Old PartitionByStmt
	New PartitionByStmt
}

type OrderByChange struct {
	Old OrderByStmt
	New OrderByStmt
}

type EngineChange struct {
	Old Engine
	New Engine
}

type Differ func(old, new SchemaTables) SemanticSchemaDiff

// NOTE:: THis is a hack, we need to preserve the order of the columns
// maps are not order preserving, so we need to use a list instead?
// TODO:: Figure out a better way to do this.
func findColumnChanges(oldCols, newCols ColumnMap) (added, removed, modified []Column) {
	seen := make(map[string]bool, len(oldCols))

	for name, newCol := range newCols {
		seen[name] = true
		if oldCol, exists := oldCols[name]; !exists {
			added = append(added, newCol)
		} else if oldCol.Type != newCol.Type {
			modified = append(modified, newCol)
		}
	}

	for name, oldCol := range oldCols {
		if !seen[name] {
			removed = append(removed, oldCol)
		}
	}

	return
}

// Logic for determining the changes between two schemas.
// This is a bit more complex than a simple set based diff, but it allows
// us to capture the changes in a way that is more human readable and
// we can pass this struct around to the various detection functions.
func SchemaDiffer(sourceSchema, targetSchema SchemaTables) SemanticSchemaDiff {
	var diff SemanticSchemaDiff

	for name := range targetSchema {
		if _, exists := sourceSchema[name]; !exists {
			diff.AddedTables = append(diff.AddedTables, targetSchema[name])
		}
	}

	for name := range sourceSchema {
		if _, exists := targetSchema[name]; !exists {
			diff.RemovedTables = append(diff.RemovedTables, sourceSchema[name])
		}
	}

	for name := range targetSchema {
		sourceTable, exists := sourceSchema[name]
		if !exists {
			break
		}

		added, removed, modified := findColumnChanges(sourceTable.ColumnMap(), targetSchema[name].ColumnMap())

		var engineChange EngineChange
		if sourceTable.Engine.Value != targetSchema[name].Engine.Value {
			engineChange = EngineChange{
				Old: sourceTable.Engine,
				New: targetSchema[name].Engine,
			}
		}

		// NOTE:: Yay no don't do this, so embarrasing
		if len(added)+len(removed)+len(modified) == 0 && engineChange == (EngineChange{}) {
			continue
		}

		tableChange := TableChange{
			Table:           targetSchema[name],
			EngineChange:    engineChange,
			ColumnsAdded:    added,
			ColumnsRemoved:  removed,
			ColumnsModified: modified,
		}

		diff.ModifiedTables = append(diff.ModifiedTables, tableChange)
	}

	return diff
}
