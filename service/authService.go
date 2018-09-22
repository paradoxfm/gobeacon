package service

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/hlandau/passlib"
	"gobeacon/code"
	"gobeacon/model"
	"math/rand"
	"time"
)

const defaultPasswordLenght = 8

func RegistrationUser(r *model.RegistrationRequest) (bool, []int) {
	err := validateRegistration(r)
	if len(err) == 0 {
		_, e := getUserByEmail(r.Email)
		if e != nil {
			hash, _ := hashPassword(r.Password)
			dbErr := insertNewUser(r.Email, hash)
			if dbErr != nil {
				err = append(err, code.UserCreateUnknownError)
			}
		} else {
			err = append(err, code.UserAlreadyExist)
		}
	}
	return len(err) == 0, err
}

func validateRegistration(r *model.RegistrationRequest) []int {
	var err []int
	if len(r.Email) == 0 {
		err = append(err, code.EmailRequired)
	}
	if len(r.Password) == 0 {
		err = append(err, code.PasswordRequired)
	}
	if len(r.Password) < 6 || len(r.Password) > 14 {
		err = append(err, code.PasswordSizeInvalid)
	}
	if !valid.IsEmail(r.Email) {
		err = append(err, code.InvalidEmailFormat)
	}
	if r.Password != r.Confirm {
		err = append(err, code.PasswordsNotEqual)
	}
	return err
}

func ResetPassword(r *model.ResetPasswordRequest) (bool, []int) {
	err := validateResetPassword(r)
	if len(err) == 0 {
		usr, e := getUserByEmail(r.Email)
		if e != nil {
			err = append(err, code.UserWithEmailNotFound)//пользователь не найден
		} else {
			newPwd := randomString(defaultPasswordLenght)
			send, _ := sendNewPassword(r.Email, newPwd)
			if send {
				hash, _ := hashPassword(newPwd)
				dbErr := updateUserPassword(usr.Id, hash)

				if dbErr != nil {//ошибка обновления в базе
					err = append(err, code.UserUpdatePwdUnknownError)
				}
			} else {// ошибка отправки email
				err = append(err, code.EmailSendError)
			}
		}
	}
	return len(err) == 0, err
}

func validateResetPassword(r *model.ResetPasswordRequest) []int {
	var err []int
	if len(r.Email) == 0 {
		err = append(err, code.EmailRequired)
	}
	if !valid.IsEmail(r.Email) {
		err = append(err, code.InvalidEmailFormat)
	}
	return err
}

func LoginUser(email string, pwd string) (string, bool) {
	usr, err := getUserByEmail(email)
	if err != nil {
		return "", false
	}
	if checkHash(pwd, usr.Password) {
		return usr.Id.String(), true
	}
	return "", false
}

func ChangePassword(r *model.ChangePasswordRequest) (bool, []int) {
	var errMsg []int

	return len(errMsg) == 0, errMsg
}

/*func convertErrorsToMessages(errs interface{}, errMsg []string) ([]string) {
	if errs != nil {
		arrErr := errs.(govalidator.Errors)
		for _, element := range arrErr {
			e := element.(govalidator.Error).Err
			errMsg = AppendError(errMsg, govalidator.ToString(e))
		}
	}
	return errMsg
}*/

func checkHash(password, hash string) bool {
	_, err := passlib.Verify(password, hash)
	return err == nil
}

func hashPassword(password string) (string, error) {
	return passlib.Hash(password)
}

func randomString(strlen int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}
