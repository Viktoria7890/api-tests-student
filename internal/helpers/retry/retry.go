package retry

import (
	"net/http"
	"time"
)

const (
	// attempts сколько раз повторить запрос при 503 Service Unavailable
	attempts = 3
	// delay пауза перед повторной попыткой
	delay = 300 * time.Millisecond
)

// On503 выполняет requestFn и повторяет вызов, если сервер ответил 503 Service Unavailable.
// Тренировочный стенд при частых запросах подряд иногда временно троттлит клиента —
// это защищает тесты от таких кратковременных сбоев, не являющихся багом приложения.
func On503(requestFn func() *http.Response) *http.Response {
	var resp *http.Response

	for attempt := 0; attempt < attempts; attempt++ {
		resp = requestFn()
		if resp.StatusCode != http.StatusServiceUnavailable {
			return resp
		}

		if attempt < attempts-1 {
			_ = resp.Body.Close()
			time.Sleep(delay)
		}
	}

	return resp
}
