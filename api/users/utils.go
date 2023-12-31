package users

import (
	"context"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Smylet/symlet-backend/utilities/token"
)

type Args struct {
	User User
}

func buildQueryFromUser(ctx context.Context, database *gorm.DB, arg Args) *gorm.DB {
	query := database.WithContext(ctx).Model(&User{})

	v := reflect.ValueOf(arg.User)
	t := reflect.TypeOf(arg.User)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Only process non-zero fields.
		if !field.IsZero() {
			// Assumes that the column name in the DB is same as struct field name.
			// If they're different, you'd use gorm struct tags to get DB column names.
			query = query.Where(fmt.Sprintf("%s = ?", fieldType.Name), field.Interface())
		}
	}
	return query
}

func GetAuthPayloadFromCtx(ctx *gin.Context) (*token.Payload, error) {
	payload, exist := ctx.Get(AuthorizationPayloadKey)
	if !exist {
		return nil, fmt.Errorf("authorization payload does not exist")
	}
	authPayload, ok := payload.(*token.Payload)
	if !ok {
		return nil, fmt.Errorf("authorization payload is not of type *token.Payload")
	}
	return authPayload, nil
}
