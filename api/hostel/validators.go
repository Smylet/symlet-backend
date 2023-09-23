package hostel

import (
	"mime/multipart"
	"net/http"

	"github.com/go-playground/validator"
)


func ValidateImageExtension(fl validator.FieldLevel) bool {
    fileHeader := fl.Field().Interface().([]*multipart.FileHeader)

    for _, file := range fileHeader {
        fileContent, err := file.Open()
        if err != nil {
            return false
        }
        defer fileContent.Close()

        buffer := make([]byte, 512)
        _, err = fileContent.Read(buffer)
        if err != nil {
            return false
        }

        contentType := http.DetectContentType(buffer)
        valid := contentType == "image/jpeg" || contentType == "image/png" || contentType == "image/gif"

        if !valid {
            return false
        }
    }

    return true
}
