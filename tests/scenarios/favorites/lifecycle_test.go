package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/favorites"
	base "api-tests-template/tests"
	"api-tests-template/tests/fixtures"
)

func TestFavoriteLifecycle(t *testing.T) {
	base.SetupSuite()
	seller, buyer := fixtures.LoginSeller(t), fixtures.LoginBuyer(t)
	require.NotEqual(t, seller.User.Id, buyer.User.Id, "продавец и покупатель должны быть разными пользователями")
	ad, deleteAd := fixtures.CreateAdvertisement(t, seller.Token, fixtures.NewAdvertisementRequest(t))
	cleanupFavorite(t, buyer.Token, ad.ID)

	t.Log("Добавляем объявление в избранное и проверяем список")
	addAndCheckFavorite(t, buyer.Token, ad)

	t.Log("Удаляем объявление из избранного и проверяем список")
	favorites.RemoveFavorite(t, buyer.Token, ad.ID, http.StatusNoContent)
	require.Empty(t, favorites.GetFavorites(t, buyer.Token, ad.Title).Items, "после удаления из избранного объявление осталось в списке")

	t.Log("Добавляем объявление повторно, затем удаляем его от имени продавца")
	addAndCheckFavorite(t, buyer.Token, ad)
	deleteAd()
	require.Empty(t, favorites.GetFavorites(t, buyer.Token, ad.Title).Items, "после удаления объявления оно осталось в избранном")
	t.Log("Удалённое объявление исчезло из списка избранного")
}

// addAndCheckFavorite добавляет объявление и проверяет ID записи в списке избранного.
func addAndCheckFavorite(t *testing.T, buyerToken string, ad models.Advertisement) {
	t.Helper()
	entry := favorites.AddFavorite(t, buyerToken, ad.ID)
	require.NotEmpty(t, entry.ID, "в ответе добавления нет ID записи избранного")
	require.Equal(t, ad.ID, entry.Advertisement.ID, "ответ добавления содержит другое объявление")
	items := favorites.GetFavorites(t, buyerToken, ad.Title).Items
	require.Len(t, items, 1, "в избранном должна быть одна запись с уникальным названием %q", ad.Title)
	require.Equal(t, entry.ID, items[0].ID, "в списке избранного вернулся другой ID записи")
	require.Equal(t, ad.ID, items[0].Advertisement.ID, "в списке избранного вернулось другое объявление")
}

// cleanupFavorite регистрирует удаление тестового объявления из избранного.
func cleanupFavorite(t *testing.T, buyerToken, advertisementID string) {
	t.Helper()
	t.Cleanup(func() {
		t.Logf("Очистка: удаляем объявление %s из избранного, если оно ещё там", advertisementID)
		favorites.RemoveFavorite(t, buyerToken, advertisementID, http.StatusNoContent, http.StatusNotFound)
	})
}
