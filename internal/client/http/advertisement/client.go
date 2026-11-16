package advertisement

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"api-tests-template/internal/constants/path"
	advertisementModels "api-tests-template/internal/managers/advertisement/models"

	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpPostAdvertisement выполняет POST-запрос создания объявления (multipart/form-data).
// Библиотека apitest не умеет собирать multipart-запросы, поэтому тело формируется вручную
// через mime/multipart, а запрос отправляется напрямую через net/http.
func HttpPostAdvertisement(t *testing.T, token, title, description string, price, quantity int,
	photos []advertisementModels.Photo) *http.Response {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writeField(t, writer, "title", title)
	writeField(t, writer, "description", description)
	writeField(t, writer, "price", fmt.Sprintf("%d", price))
	writeField(t, writer, "quantity", fmt.Sprintf("%d", quantity))

	for _, photo := range photos {
		part, err := writer.CreateFormFile("photos", photo.FileName)
		if err != nil {
			t.Fatalf("не удалось создать form-file для фото: %s", err.Error())
		}
		if _, err = part.Write(photo.Content); err != nil {
			t.Fatalf("не удалось записать содержимое фото: %s", err.Error())
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("не удалось закрыть multipart writer: %s", err.Error())
	}

	// URL собирается из API_URL — доверенной переменной окружения проекта (та же, что использует
	// api-runner для всех остальных запросов), а не из пользовательского ввода, поэтому SSRF здесь неприменим.
	//nolint:gosec // G704: see comment above — API_URL задаётся в .env/.env.override, не берётся из запроса
	req, err := http.NewRequest(http.MethodPost, os.Getenv("API_URL")+path.AdvertisementPath, body)
	if err != nil {
		t.Fatalf("не удалось создать запрос на создание объявления: %s", err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", "api-tests-template")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	//nolint:gosec // G704: тот же доверенный URL, см. комментарий выше
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("не удалось выполнить запрос на создание объявления: %s", err.Error())
	}

	return resp
}

func writeField(t *testing.T, writer *multipart.Writer, name, value string) {
	if err := writer.WriteField(name, value); err != nil {
		t.Fatalf("не удалось записать поле %s: %s", name, err.Error())
	}
}

// HttpGetAdvertisementById выполняет GET-запрос получения объявления по id
func HttpGetAdvertisementById(t *testing.T, id string) *http.Response {
	return apiRunner.GetRunner().Create().Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// HttpDeleteAdvertisementById выполняет DELETE-запрос удаления объявления по id
func HttpDeleteAdvertisementById(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}
