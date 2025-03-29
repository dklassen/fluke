package main

import (
	"strings"
	"testing"
)

func TestGenerateMigration(t *testing.T) {
	tests := []struct {
		name     string
		oldSQL   string
		newSQL   string
		expected string
	}{
		{
			name:   "New table",
			oldSQL: "",
			newSQL: `
				CREATE TABLE users (
					id UInt32,
					name String
				) ENGINE = MergeTree
				PRIMARY KEY id
				ORDER BY id;`,
			expected: `CREATE TABLE users (id UInt32, name String) ENGINE = MergeTree PRIMARY KEY id ORDER BY id;`,
		},
		{
			name: "Modify table structure",
			oldSQL: `
				CREATE TABLE users (
					id UInt32
				) ENGINE = MergeTree
				PRIMARY KEY (id)
				ORDER BY id;`,
			newSQL: `
				CREATE TABLE users (
					id UInt32,
					name String
				) ENGINE = MergeTree
				PRIMARY KEY (id)
				ORDER BY id;`,
			expected: `ALTER TABLE users ADD COLUMN name String;`,
		},
		{
			name: "Drop table",
			oldSQL: `
				CREATE TABLE users (
					id UInt32
				) ENGINE = MergeTree
				ORDER BY (id);`,
			newSQL:   "",
			expected: `DROP TABLE IF EXISTS users;`,
		},
		{
			name: "Modify table engine structure",
			oldSQL: `
				CREATE TABLE users (
					id UInt32,
					name String
				) ENGINE = MergeTree
				PRIMARY KEY id
				ORDER BY id DESC;`,
			newSQL: `
				CREATE TABLE users (
					id UInt32,
					name String
				) ENGINE = Kafka()
				PRIMARY KEY (id, name)
				ORDER BY (id, name);`,
			expected: `DROP TABLE IF EXISTS users;
CREATE TABLE users (id UInt32, name String) ENGINE = Kafka PRIMARY KEY (id, name) ORDER BY (id, name);`,
		},
		{
			name: "No changes",
			oldSQL: `
				CREATE TABLE users (
					id UInt32
				) ENGINE = MergeTree
				ORDER BY id;`,
			newSQL: `
				CREATE TABLE users (
					id UInt32
				) ENGINE = MergeTree
				ORDER BY id;`,
			expected: ``,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			subject := NewSQLParser(test.oldSQL).Parse()
			target := NewSQLParser(test.newSQL).Parse()

			output := strings.TrimSpace(GenerateMigration(subject, target))
			expected := strings.TrimSpace(test.expected)

			if output != expected {
				t.Errorf("\n---\nGot:\n%s\n---\nWant:\n%s\n", output, expected)
			}
		})
	}
}
