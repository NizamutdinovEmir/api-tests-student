package advertisements

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	apiResponse "api-tests-template/internal/helpers/api-response"
	"api-tests-template/internal/managers/advertisements/models"
)

func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	t.Helper()

	return apiResponse.ReadJSON[models.AdvertisementsList](t, advertisements.HttpGetAdvertisements(t, token, text), http.StatusOK)
}

// CreateAdvertisement создаёт объявление и возвращает данные созданной записи.
func CreateAdvertisement(t *testing.T, token string, request models.CreateRequest) models.Advertisement {
	t.Helper()
	return apiResponse.ReadJSON[models.Advertisement](t, advertisements.HttpPostAdvertisement(t, token, request), http.StatusCreated)
}

// RejectCreateAdvertisement проверяет статус отклонённого запроса и возвращает ошибку API.
func RejectCreateAdvertisement(t *testing.T, token string, request models.CreateRequest, expectedStatus int) apiResponse.Error {
	t.Helper()
	response := advertisements.HttpPostAdvertisement(t, token, request)
	if response.StatusCode == http.StatusCreated && expectedStatus == http.StatusBadRequest {
		created := apiResponse.ReadJSON[models.Advertisement](t, response, http.StatusCreated)
		if created.ID != "" {
			t.Cleanup(func() { DeleteAdvertisement(t, token, created.ID, http.StatusNoContent, http.StatusNotFound) })
		}
		require.FailNow(t, "сервер создал объявление с невалидными параметрами")
	}
	return apiResponse.ReadJSON[apiResponse.Error](t, response, expectedStatus)
}

// GetAdvertisement получает данные одного объявления.
func GetAdvertisement(t *testing.T, token, id string) models.Advertisement {
	t.Helper()
	return apiResponse.ReadJSON[models.Advertisement](t, advertisements.HttpGetAdvertisement(t, token, id), http.StatusOK)
}

// DeleteAdvertisement удаляет объявление и проверяет ожидаемый статус операции.
func DeleteAdvertisement(t *testing.T, token, id string, expectedStatuses ...int) {
	t.Helper()
	apiResponse.RequireStatus(t, advertisements.HttpDeleteAdvertisement(t, token, id), expectedStatuses...)
}

// GetAdvertisementPhotos возвращает данные загруженных фотографий.
func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	t.Helper()
	return apiResponse.ReadJSON[[]models.Photo](t, advertisements.HttpGetAdvertisementPhotos(t, token, id), http.StatusOK)
}

// DownloadAdvertisementPhoto возвращает содержимое загруженной фотографии.
func DownloadAdvertisementPhoto(t *testing.T, photoURL string) []byte {
	t.Helper()
	const maxPhotoSize = 5 * 1024 * 1024

	response := advertisements.HttpGetPhotoFile(t, photoURL)
	defer response.Body.Close()
	require.Equal(t, http.StatusOK, response.StatusCode, "фотография недоступна в хранилище")
	content, err := io.ReadAll(io.LimitReader(response.Body, maxPhotoSize+1))
	require.NoError(t, err, "не удалось прочитать файл фотографии")
	require.LessOrEqual(t, len(content), maxPhotoSize, "размер фотографии превышает 5 МиБ")
	return content
}
