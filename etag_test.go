package httpx_test

import (
	"testing"

	"github.com/arcgolabs/httpx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseETagSet_ReturnsStableValues(t *testing.T) {
	set, err := httpx.ParseETagSet(`"abc", W/"def"`)
	require.NoError(t, err)

	values := set.Values()
	require.Len(t, values, 2)
	assert.Equal(t, httpx.ETag{Value: "abc"}, values[0])
	assert.Equal(t, httpx.ETag{Value: "def", Weak: true}, values[1])
	assert.Equal(t, `"abc", W/"def"`, set.String())

	values[0] = httpx.ETag{Value: "changed"}
	assert.True(t, set.Contains(httpx.ETag{Value: "abc"}))
}

func TestETagSet_WildcardMatchesAnyTag(t *testing.T) {
	set, err := httpx.ParseETagSet("*")
	require.NoError(t, err)

	assert.Equal(t, "*", set.String())
	assert.True(t, set.Contains(httpx.ETag{Value: "anything"}))
}

func TestETagSet_UnmarshalTextPreservesReceiverOnError(t *testing.T) {
	set, err := httpx.ParseETagSet(`"stable"`)
	require.NoError(t, err)

	require.Error(t, set.UnmarshalText([]byte("not-an-etag")))
	assert.Equal(t, `"stable"`, set.String())
}
