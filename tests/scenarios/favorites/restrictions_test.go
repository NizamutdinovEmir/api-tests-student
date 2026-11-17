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

func TestCannotFavoriteOwnAdvertisement(t *testing.T) {
	base.SetupSuite()
	seller := fixtures.LoginSeller(t)
	ad, _ := fixtures.CreateAdvertisement(t, seller.Token, fixtures.NewAdvertisementRequest(t))

	t.Log("Пытаемся добавить своё объявление в избранное")
	got := favorites.RejectAddFavorite(t, seller.Token, ad.ID, http.StatusForbidden)
	require.Equal(t, apiResponse.Error{Code: "forbidden", Message: "Cannot add own advertisement to favorites"}, got)
}

func TestCannotFavoriteSameAdvertisementTwice(t *testing.T) {
	base.SetupSuite()
	seller, buyer := fixtures.LoginSeller(t), fixtures.LoginBuyer(t)
	require.NotEqual(t, seller.User.Id, buyer.User.Id, "продавец и покупатель должны быть разными пользователями")
	ad, _ := fixtures.CreateAdvertisement(t, seller.Token, fixtures.NewAdvertisementRequest(t))
	cleanupFavorite(t, buyer.Token, ad.ID)
	t.Log("Добавляем объявление в избранное покупателя")
	favorites.AddFavorite(t, buyer.Token, ad.ID)

	t.Log("Пытаемся повторно добавить объявление в избранное")
	got := favorites.RejectAddFavorite(t, buyer.Token, ad.ID, http.StatusConflict)
	require.Equal(t, apiResponse.Error{Code: "already_exists", Message: "Advertisement is already in favorites"}, got)
}
