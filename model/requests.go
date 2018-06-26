package model

type RegistrationRequest struct {
	Email    string `json:"email,required" description:"Email пользователя" valid:"email~Неправильный формат email,required~Поле email обязательно для заполнения"`
	Password string `json:"password,required," description:"Пароль пользователя" valid:"required~Не заполнен пароль,length(6|14)~Длина пароля от 6 до 14 символов"`
	Confirm  string `json:"сonfirm,required" description:"Подтверждение пароля" valid:"required~Не заполнено подтверждение пароля"`
}