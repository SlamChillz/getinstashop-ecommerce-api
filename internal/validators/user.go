package validators

import (
	"errors"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/types"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/utils"
	"regexp"
)

func ValidateAuthPayload(input types.AuthPayload) (types.AuthPayloadErrMessage, error) {
	errMessage := types.AuthPayloadErrMessage{}
	var err error = nil
	if input.Email == "" {
		err = utils.ConcatenateErrors(err, errors.New("missing email"))
		errMessage.Email = "email is required"
	} else {
		emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		re := regexp.MustCompile(emailRegex)
		// Match email with regex
		if !re.MatchString(input.Email) {
			err = utils.ConcatenateErrors(err, errors.New("invalid email"))
			errMessage.Email = "email is invalid"
		}
	}
	if input.Password == "" {
		err = utils.ConcatenateErrors(err, errors.New("missing password"))
		errMessage.Password = "password is required"
	} else {
		if len(input.Password) < 8 {
			err = utils.ConcatenateErrors(err, errors.New("password must be at least 8 characters"))
			errMessage.Password = "password must be at least 8 characters"
		}
	}
	return errMessage, err
}
