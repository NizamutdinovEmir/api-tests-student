package advertisements

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

// HttpPostAdvertisement отправляет поля объявления и PNG-фотографию одной multipart-формой.
func HttpPostAdvertisement(t *testing.T, token string, request models.CreateRequest) *http.Response {
	t.Helper()

	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	for name, value := range request.Fields {
		require.NoError(t, form.WriteField(name, value), "запись поля %s", name)
	}
	if request.Photo != nil {
		part, err := form.CreateFormFile("photos", "test-photo.png")
		require.NoError(t, err, "добавление фотографии в форму")
		_, err = part.Write(request.Photo)
		require.NoError(t, err, "запись фотографии")
	}
	require.NoError(t, form.Close(), "завершение multipart-формы")

	return apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath).
		Header("Content-Type", form.FormDataContentType()).
		Body(body.String()).
		Expect(t).End().Response
}

// HttpGetAdvertisement получает объявление по его ID.
func HttpGetAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).End().Response
}

// HttpDeleteAdvertisement удаляет объявление по его ID.
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).End().Response
}

// HttpGetAdvertisementPhotos получает список фотографий объявления.
func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	endpoint := fmt.Sprintf("%s/%s/photos", path.AdvertisementsPath, id)
	return apiRunner.GetRunner().Auth(token).Create().Get(endpoint).Expect(t).End().Response
}

// HttpGetPhotoFile скачивает файл по публичному HTTPS-адресу из ответа TestBoard.
func HttpGetPhotoFile(t *testing.T, photoURL string) *http.Response {
	t.Helper()

	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, photoURL, nil)
	require.NoError(t, err, "не удалось составить запрос на скачивание фотографии")
	require.Equal(t, "https", request.URL.Scheme, "ожидается HTTPS-ссылка на фотографию")
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	require.NoError(t, err, "не удалось скачать фотографию")
	return response
}
