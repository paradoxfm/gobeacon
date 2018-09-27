package code

// Пароль не совпадает с подтверждением
const UnknownErr = -1
// Пароль не совпадает с подтверждением
const PasswordsNotEqual = 2000003
//Неправильный формат email
const InvalidEmailFormat = 2000004
//не заполнен email
const EmailRequired = 2000005
//не указан пароль
const PasswordRequired = 2000006
//Длина пароля от 6 до 14 символов
const PasswordSizeInvalid = 2000007
//Пользователь уже существует
const UserAlreadyExist = 2000008
//Ошибка сохранения нового пользователя
const UserCreateUnknownError = 2000009
//Пользователья с таким email не существует
const UserWithEmailNotFound = 2000010
//Ошибка обновления пароля пользователя
const UserUpdatePwdUnknownError = 2000011
//Ошибка обправки по email нового пароля пользователя
const EmailSendError = 2000012
//Ошибка БД при обновлении push id
const DbErrorUpdateUserPush = 2000013
//Ошибка БД не удалось достать трекер из базы
const DbErrorGetTracker = 2000014
// Пароль не совпадает с подтверждением
const InavlidCurrentPasswords = 2000015
// Пароль не совпадает с подтверждением
const ZoneCrossing = 2000016
// Невозможно открыть файл
const CantOpenFile = 2000017
// Невозможно Прочитать файл
const CantReadFile = 2000018
// Неверный формат изображения, ожидается jpeg
const InvalidImage = 2000019
// Размер изображения должен быть 250 на 250
const InvalidImageSize = 2000020
// Попытка добавить уже существующий трекер
const TrackForUserExist = 2000021
//Ошибка БД
const DbError = 2000100
