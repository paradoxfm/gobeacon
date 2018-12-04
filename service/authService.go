package service

import (
	"github.com/appleboy/gin-jwt"
	valid "github.com/asaskevich/govalidator"
	"github.com/hlandau/passlib"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"gobeacon/code"
	"gobeacon/db"
	"gobeacon/model"
	"math/rand"
	"time"
)

const defaultPasswordLength = 8

func RegistrationUser(r *model.RegistrationRequest) (bool, []int) {
	err := validateRegistration(r)
	if len(err) == 0 {
		_, e := db.LoadUserByEmail(r.Email)
		if e != nil {
			hash, _ := hashPassword(r.Password)
			dbErr := db.InsertNewUser(r.Email, hash)
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
		usr, e := db.LoadUserByEmail(r.Email)
		if e != nil {
			err = append(err, code.UserWithEmailNotFound) //пользователь не найден
		} else {
			newPwd := randomString(defaultPasswordLength)
			send, erse := sendNewPassword(r.Email, newPwd)
			if send {
				hash, _ := hashPassword(newPwd)
				logrus.WithFields(logrus.Fields{"hash": hash}).Info("New hash password")
				dbErr := db.UpdateUserPassword(usr.Id.String(), hash)

				if dbErr != nil { //ошибка обновления в базе
					err = append(err, code.UserUpdatePwdUnknownError)
				}
			} else { // ошибка отправки email
				log.Error(erse)
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

func LoginUser(email string, pwd string) (interface{}, error) {
	usr, err := db.LoadUserByEmail(email)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}
	if checkHash(pwd, usr.Password) {
		return usr, nil
	}
	return nil, jwt.ErrFailedAuthentication
}

func UserExist(id string, pwdHash string) (bool) {
	usr, err := db.LoadUserById(id)
	if err != nil {
		return false
	}
	return pwdHash == usr.Password
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
