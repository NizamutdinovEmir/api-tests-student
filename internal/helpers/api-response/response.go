package api_response

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Error описывает код и сообщение ошибки TestBoard.
type Error struct {
	Code    string `json:"error"`
	Message string `json:"message"`
}

// ReadJSON проверяет статус, разбирает ответ и закрывает его тело.
func ReadJSON[T any](t *testing.T, response *http.Response, expectedStatus int) T {
	t.Helper()
	defer response.Body.Close()

	require.Equal(t, expectedStatus, response.StatusCode, "неожиданный HTTP-статус")
	var result T
	require.NoError(t, json.NewDecoder(response.Body).Decode(&result), "не удалось разобрать JSON ответа")
	return result
}

// RequireStatus проверяет статус запроса, тело ответа которого не используется.
func RequireStatus(t *testing.T, response *http.Response, expectedStatuses ...int) {
	t.Helper()
	defer response.Body.Close()

	require.Contains(t, expectedStatuses, response.StatusCode, "неожиданный HTTP-статус")
}
