package advertisementCreation

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	apiResponse "api-tests-template/internal/helpers/api-response"
	"api-tests-template/internal/managers/advertisements"
	base "api-tests-template/tests"
	"api-tests-template/tests/fixtures"
)

func TestRejectAdvertisementWithoutAuthorization(t *testing.T) {
	base.SetupSuite()
	cases := []struct {
		name    string
		token   string
		message string
	}{
		{"без токена", "", fixtures.MissingAuthMessage},
		{"с неверным токеном", fixtures.InvalidToken, fixtures.InvalidAuthMessage},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Отправляем валидную форму без действующей авторизации")
			request := fixtures.NewAdvertisementRequest(t)
			got := advertisements.RejectCreateAdvertisement(t, tc.token, request, http.StatusUnauthorized)
			require.Equal(t, apiResponse.Error{Code: fixtures.UnauthorizedCode, Message: tc.message}, got)
		})
	}
}
