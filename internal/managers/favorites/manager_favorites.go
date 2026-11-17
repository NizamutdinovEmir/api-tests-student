package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/client/http/favorites"
	apiResponse "api-tests-template/internal/helpers/api-response"
	"api-tests-template/internal/managers/advertisements/models"
)

type Favorite struct {
	ID            string               `json:"favorite_id"`
	Advertisement models.Advertisement `json:"advertisement"`
}

type FavoriteList struct {
	Items []Favorite `json:"items"`
}

// GetFavorites возвращает избранные объявления, найденные по названию.
func GetFavorites(t *testing.T, token, title string) FavoriteList {
	t.Helper()
	return apiResponse.ReadJSON[FavoriteList](t, favorites.HttpGetFavorites(t, token, title), http.StatusOK)
}

// AddFavorite добавляет объявление и возвращает запись избранного.
func AddFavorite(t *testing.T, token, advertisementID string) Favorite {
	t.Helper()
	return apiResponse.ReadJSON[Favorite](t, favorites.HttpPostFavorite(t, token, advertisementID), http.StatusCreated)
}

// RemoveFavorite удаляет объявление из избранного и проверяет статус операции.
func RemoveFavorite(t *testing.T, token, advertisementID string, expectedStatuses ...int) {
	t.Helper()
	apiResponse.RequireStatus(t, favorites.HttpDeleteFavorite(t, token, advertisementID), expectedStatuses...)
}

// RejectGetFavorites возвращает ошибку API для отклонённого чтения избранного.
func RejectGetFavorites(t *testing.T, token string, expectedStatus int) apiResponse.Error {
	t.Helper()
	return apiResponse.ReadJSON[apiResponse.Error](t, favorites.HttpGetFavorites(t, token, ""), expectedStatus)
}

// RejectAddFavorite возвращает ошибку API для отклонённого добавления объявления.
func RejectAddFavorite(t *testing.T, token, advertisementID string, expectedStatus int) apiResponse.Error {
	t.Helper()
	return apiResponse.ReadJSON[apiResponse.Error](t, favorites.HttpPostFavorite(t, token, advertisementID), expectedStatus)
}

// RejectRemoveFavorite возвращает ошибку API для отклонённого удаления из избранного.
func RejectRemoveFavorite(t *testing.T, token, advertisementID string, expectedStatus int) apiResponse.Error {
	t.Helper()
	return apiResponse.ReadJSON[apiResponse.Error](t, favorites.HttpDeleteFavorite(t, token, advertisementID), expectedStatus)
}

// RejectBulkAddFavorites возвращает ошибку API для отклонённого массового добавления.
func RejectBulkAddFavorites(t *testing.T, token string, advertisementIDs []string, expectedStatus int) apiResponse.Error {
	t.Helper()
	return apiResponse.ReadJSON[apiResponse.Error](t, favorites.HttpPostBulkFavorites(t, token, advertisementIDs), expectedStatus)
}
