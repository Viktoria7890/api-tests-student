package utils

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

// GenerateTestImage генерирует минимальное валидное PNG-изображение 1x1 пиксель
// для загрузки в тестах — без зависимости от внешних файлов-фикстур
func GenerateTestImage() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	buf := &bytes.Buffer{}
	if err := png.Encode(buf, img); err != nil {
		panic(err)
	}

	return buf.Bytes()
}
