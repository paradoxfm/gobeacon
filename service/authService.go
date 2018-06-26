package service

import (
	"github.com/asaskevich/govalidator"
	"gobeacon/model"
)

func RegistrationUser(r *model.RegistrationRequest) (bool, []string) {
	var errMsg []string
	_, err := govalidator.ValidateStruct(r)
	errMsg = convertErrorsToMessages(err, errMsg)
	if r.Password != r.Confirm {
		errMsg = AppendError(errMsg, "Пароль не совпадает с подтверждением")
	}

	return len(errMsg) == 0, errMsg
}

func convertErrorsToMessages(errs interface{}, errMsg []string) ([]string) {
	if errs != nil {
		arrErr := errs.(govalidator.Errors)
		for _, element := range arrErr {
			e := element.(govalidator.Error).Err
			errMsg = AppendError(errMsg, govalidator.ToString(e))
		}
	}
	return errMsg
}
