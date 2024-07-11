package httpserver

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func extractSuccessData[T any](t testing.TB, r io.Reader) T {
	t.Helper()

	var result T
	s := Success{
		Data: &result, // embed pointer to parse data into result
	}

	err := json.NewDecoder(r).Decode(&s)
	require.NoError(t, err)

	return result
}

func extractErrorData(t testing.TB, r io.Reader) Errs {
	t.Helper()

	var result Errs
	err := json.NewDecoder(r).Decode(&result)
	require.NoError(t, err)
	return result
}
