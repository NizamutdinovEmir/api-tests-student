package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	apiResponse "api-tests-template/internal/helpers/api-response"
	"api-tests-template/internal/managers/favorites"
	base "api-tests-template/tests"
	"api-tests-template/tests/fixtures"
)

func TestRejectFavoritesWithoutAuthorization(t *testing.T) {
	base.SetupSuite()
	tokens := []struct {
		name    string
		token   string
		message string
	}{
		{"без токена", "", fixtures.MissingAuthMessage},
		{"с неверным токеном", fixtures.InvalidToken, fixtures.InvalidAuthMessage},
	}
	operations := []struct {
		name string
		call func(*testing.T, string) apiResponse.Error
	}{
		{"GET /favorites", func(t *testing.T, token string) apiResponse.Error {
			return favorites.RejectGetFavorites(t, token, http.StatusUnauthorized)
		}},
		{"POST /favorites", func(t *testing.T, token string) apiResponse.Error {
			return favorites.RejectAddFavorite(t, token, fixtures.MissingAdvertisementID, http.StatusUnauthorized)
		}},
		{"DELETE /favorites/{ad_id}", func(t *testing.T, token string) apiResponse.Error {
			return favorites.RejectRemoveFavorite(t, token, fixtures.MissingAdvertisementID, http.StatusUnauthorized)
		}},
		{"POST /favorites/bulk", func(t *testing.T, token string) apiResponse.Error {
			return favorites.RejectBulkAddFavorites(t, token, []string{fixtures.MissingAdvertisementID}, http.StatusUnauthorized)
		}},
	}
	for _, operation := range operations {
		for _, tc := range tokens {
			t.Run(operation.name+" "+tc.name, func(t *testing.T) {
				got := operation.call(t, tc.token)
				require.Equal(t, apiResponse.Error{Code: fixtures.UnauthorizedCode, Message: tc.message}, got)
			})
		}
	}
}
