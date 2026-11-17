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

func TestFavoriteWithNonexistentAdvertisement(t *testing.T) {
	base.SetupSuite()
	token := fixtures.LoginSeller(t).Token
	cases := []struct {
		name    string
		action  func(*testing.T, string, string, int) apiResponse.Error
		message string
	}{
		{"добавление несуществующего объявления", favorites.RejectAddFavorite, "Advertisement not found"},
		{"удаление несуществующего объявления", favorites.RejectRemoveFavorite, "Favorite entry not found"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.action(t, token, fixtures.MissingAdvertisementID, http.StatusNotFound)
			require.Equal(t, apiResponse.Error{Code: "not_found", Message: tc.message}, got)
		})
	}
}

func TestAddFavoriteWithInvalidAdvertisementID(t *testing.T) {
	base.SetupSuite()
	token := fixtures.LoginSeller(t).Token
	const malformedID = "invalid-id"
	cases := []struct {
		name    string
		id      string
		message string
	}{
		{"пустой ID", "", "invalid UUID length: 0"},
		{"невалидный UUID", malformedID, "invalid UUID length: 10"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := favorites.RejectAddFavorite(t, token, tc.id, http.StatusBadRequest)
			require.Equal(t, apiResponse.Error{Code: fixtures.InvalidRequestCode, Message: tc.message}, got)
		})
	}
}
