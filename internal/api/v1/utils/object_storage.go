package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

func CompressImage(file io.Reader, filename string) (*bytes.Buffer, error) {
	img, format, err := image.Decode(file)
	fmt.Println(format)
	if err != nil {
		return nil, fmt.Errorf("не удалось декодировать изображение: %v", err)
	}

	var buf bytes.Buffer

	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 20})
	case "png":
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		err = encoder.Encode(&buf, img)
	default:
		return nil, fmt.Errorf("неподдерживаемый формат: %s", format)
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка при сжатии: %v", err)
	}

	return &buf, nil
}
