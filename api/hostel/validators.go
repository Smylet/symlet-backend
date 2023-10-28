package hostel

import (
	"mime/multipart"
	"net/http"

	"github.com/go-playground/validator"

	"github.com/Smylet/symlet-backend/utilities/common"
)

func ValidateImageExtension(fl validator.FieldLevel) bool {
	logger := common.NewLogger()
	fileHeader := fl.Field().Interface().([]*multipart.FileHeader)

	for _, file := range fileHeader {
		fileContent, err := file.Open()
		if err != nil {
			logger.Error("failed to open file: ", err.Error())
			return false
		}
		defer fileContent.Close()

		buffer := make([]byte, 512)
		n, err := fileContent.Read(buffer)
		if err != nil {
			logger.Error("failed to read file: ", err.Error())
			return false
		}

		contentType := http.DetectContentType(buffer[:n])
		valid := contentType == "image/jpeg" || contentType == "image/png" || contentType == "image/gif"

		if !valid {
			return false
		}
	}

	return true
}
