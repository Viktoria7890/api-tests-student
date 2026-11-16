package utils

import (
	"crypto/rand"
	"fmt"
)

// RandomUUID генерирует случайный UUID v4 — используется для проверок с заведомо
// несуществующим, но корректным по формату идентификатором
func RandomUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("не удалось сгенерировать случайный UUID: %s", err.Error()))
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
