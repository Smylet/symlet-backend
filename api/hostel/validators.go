package hostel

import (
	"mime/multipart"
	"path/filepath"

	"github.com/go-playground/validator/v10"
)


func validateImageExtension(fl validator.FieldLevel) bool {
    allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
    fileHeader := fl.Field().Interface().([]*multipart.FileHeader)

    for _, file := range fileHeader {
        ext := filepath.Ext(file.Filename)
        valid := false

        for _, allowedExt := range allowedExtensions {
            if ext == allowedExt {
                valid = true
                break
            }
        }

        if !valid {
            return false
        }
    }

    return true
}
