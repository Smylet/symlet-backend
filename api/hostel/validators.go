package hostel

import (
	"mime/multipart"
	"net/http"

	"github.com/go-playground/validator/v10"
)


func validateImageExtension(fl validator.FieldLevel) bool {
    fileHeader := fl.Field().Interface().([]*multipart.FileHeader)

    for _, file := range fileHeader {
        fileContent, _ := file.Open()
        defer fileContent.Close()

        buffer := make([]byte, 512)
        _, err := fileContent.Read(buffer)
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
