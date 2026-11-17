package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpGetFavorites получает избранное с фильтром по названию объявления.
func HttpGetFavorites(t *testing.T, token, title string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Query("search", title).
		Expect(t).End().Response
}

// HttpPostFavorite добавляет объявление в избранное.
func HttpPostFavorite(t *testing.T, token, advertisementID string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(map[string]string{"advertisement_id": advertisementID}).
		Expect(t).End().Response
}

// HttpDeleteFavorite удаляет объявление из избранного.
func HttpDeleteFavorite(t *testing.T, token, advertisementID string) *http.Response {
	endpoint := path.FavoritesPath + "/" + advertisementID
	return apiRunner.GetRunner().Auth(token).Create().Delete(endpoint).Expect(t).End().Response
}

// HttpPostBulkFavorites отправляет список ID для добавления в избранное.
func HttpPostBulkFavorites(t *testing.T, token string, advertisementIDs []string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Post(path.BulkFavoritesPath).
		ContentType("application/json").
		JSON(map[string][]string{"advertisement_ids": advertisementIDs}).
		Expect(t).End().Response
}
