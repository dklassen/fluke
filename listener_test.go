package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderByKeyStmtString(t *testing.T) {
	tests := []struct {
		name           string
		orderByKeyStmt OrderByKeyList
		expected       string
	}{
		{
			name: "Order by one column",
			orderByKeyStmt: OrderByKeyList{
				{Value: "a", SortOrder: ASC},
			},
			expected: "a",
		},
		{
			name: "Order by two columns",
			orderByKeyStmt: OrderByKeyList{
				{Value: "a", SortOrder: ASC},
				{Value: "b", SortOrder: DESC},
			},
			expected: "(a, b DESC)",
		},
		{
			name: "Order by two column with ASC",
			orderByKeyStmt: OrderByKeyList{
				{Value: "a", SortOrder: ASC},
				{Value: "b", SortOrder: ASC},
			},
			expected: "(a, b)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.orderByKeyStmt.String())
		})
	}
}
