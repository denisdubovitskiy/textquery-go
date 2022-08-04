package textquery

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	cases := []struct {
		query string
		input string
		match bool
	}{
		{
			query: "",
			input: "some input",
			match: false,
		},
		{
			query: "a AND b",
			input: "a",
			match: false,
		},
		{
			query: "a AND b",
			input: "b",
			match: false,
		},
		{
			query: "a AND b",
			input: "a b",
			match: true,
		},
		{
			query: "a AND b AND NOT c",
			input: "a b",
			match: true,
		},
		{
			query: "a OR b",
			input: "a n",
			match: true,
		},
		{
			query: "b OR a",
			input: "a n",
			match: true,
		},
		{
			query: "a AND b AND NOT c",
			input: "a b c",
			match: false,
		},
		{
			query: "a AND b AND NOT c",
			input: "a b d",
			match: true,
		},
		{
			query: "(a AND b AND NOT c) OR (c AND g)",
			input: "a b d",
			match: true,
		},
		{
			query: "a AND b AND NOT d",
			input: "a b d",
			match: false,
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("q:%s,i:%s", tc.query, tc.input), func(t *testing.T) {
			tree := Parse(tc.query)

			assert.Equal(t, tc.match, tree.Match(tc.input))
		})
	}
}
