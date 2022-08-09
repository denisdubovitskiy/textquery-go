package textquery

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func makeTokens(t []string) []*Token {
	tokens := make([]*Token, len(t))
	for i, s := range t {
		tokens[i] = &Token{Data: s}
	}
	return tokens
}

func TestTokenize(t *testing.T) {
	cases := []struct {
		q    string
		want []*Token
	}{
		{
			q:    "a AND b AND c",
			want: makeTokens([]string{"a", "AND", "b", "AND", "c"}),
		},
		{
			q:    `"a b" AND "c d"`,
			want: makeTokens([]string{`a b`, "AND", `c d`}),
		},
		{
			q:    `"a" AND "bc" AND c AND NOT "d e"`,
			want: makeTokens([]string{`a`, `AND`, `bc`, `AND`, `c`, `AND`, `NOT`, `d e`}),
		},
		{
			q: `field{lower}[equal]:data AND d OR e[f]:g`,
			want: []*Token{
				{
					Data:     "data",
					Field:    "field",
					Operator: "equal",
					Modifier: "lower",
				},
				{Data: "AND"},
				{Data: "d"},
				{Data: "OR"},
				{
					Data:     "g",
					Field:    "e",
					Operator: "f",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.q, func(t *testing.T) {
			// act
			got := tokenize(tc.q)

			// assert
			assert.Equal(t, tc.want, got)
		})
	}
}
