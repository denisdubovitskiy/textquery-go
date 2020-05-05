package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckNeedsWrapping(t *testing.T) {
	assert.True(t, New().checkNeedsWrapping("word1"))
	assert.True(t, New().checkNeedsWrapping("word1 OR word2"))
	assert.True(t, New().checkNeedsWrapping("(word1 OR word2) AND (word3 OR word45)"))
	assert.False(t, New().checkNeedsWrapping("((word1 OR word2) AND (word3 OR word4))"))
}

func TestOpenParenthesis(t *testing.T) {
	cases := []struct {
		query    string
		msg      string
		position int
		char     string
		err      bool
	}{
		{
			query:    "field1[exact]:word1) OR word2 AND (word3 OR word45)",
			msg:      "open parenthesis is not found",
			position: 19,
			char:     ")",
			err:      true,
		},
		{
			query:    "(field1[exact]:word1 OR word2 AND (word3 OR word45)",
			msg:      "close parenthesis is not found",
			position: 0,
			char:     "(",
			err:      true,
		},
		{
			query: "(field1[exact]:word1) OR word2 AND (word3 OR word45)",
		},
		{
			query:    "(field1exact]:word1) OR word2 AND (word3 OR word45)",
			msg:      "open parenthesis is not found",
			position: 12,
			char:     "]",
			err:      true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.query, func(t *testing.T) {
			_, err := New().Parse(tc.query)

			if tc.err {
				require.Error(t, err)
				assert.Equal(t, tc.msg, err.Error())
				assert.Equal(t, tc.position, err.(ValidationError).Pos())
				assert.Equal(t, tc.char, err.(ValidationError).Char())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
