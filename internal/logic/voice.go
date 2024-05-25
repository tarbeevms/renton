package logic

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"myapp/internal/voice"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// VoiceLogic содержит логику для обработки голосовых записей
type VoiceLogic struct {
	repository voice.UserRepository
}

// NewVoiceLogic создает новый экземпляр VoiceLogic
func NewVoiceLogic(repository voice.UserRepository) *VoiceLogic {
	return &VoiceLogic{repository: repository}
}

// CreateVoiceRecordings создает записи голосовых записей пользователя
func (vl *VoiceLogic) CreateVoiceRecordings(userID uuid.UUID, audio1, audio2, audio3 []byte) error {
	recordings := &voice.VoiceRecording{
		UserID: userID,
		Audio1: audio1,
		Audio2: audio2,
		Audio3: audio3,
	}

	err := vl.repository.SaveVoiceRecordings(recordings)
	if err != nil {
		return err
	}

	return nil
}

// UpdateVoiceRecordings обновляет записи голосовых записей пользователя
func (vl *VoiceLogic) UpdateVoiceRecordings(userID uuid.UUID, audio1, audio2, audio3 []byte) error {
	recordings := &voice.VoiceRecording{
		UserID: userID,
		Audio1: audio1,
		Audio2: audio2,
		Audio3: audio3,
	}

	err := vl.repository.UpdateVoiceRecordings(recordings)
	if err != nil {
		return err
	}

	return nil
}

// CreateSessionCookie создает куки-файл для сессии пользователя
func (vl *VoiceLogic) CreateSessionCookie(userID string) (*http.Cookie, error) {
	expiration := time.Now().Add(24 * time.Hour)
	cookie := &http.Cookie{
		Name:    "user_id",
		Value:   userID,
		Expires: expiration,
		Path:    "/",
	}
	return cookie, nil
}

func generateAccountNumber() string {
	// Генерируем случайный номер счета длиной 20 символов
	return uuid.New().String()[:20]
}

// RegisterUser регистрирует нового пользователя
func (vl *VoiceLogic) RegisterUser(phone, password, firstname, surname string) (string, error) {
	// Проверяем, существует ли пользователь с таким номером телефона
	existingUser, err := vl.repository.GetUserByPhone(phone)
	if err != sql.ErrNoRows {
		return "", err
	}
	if existingUser != nil {
		return "", errors.New("user with this phone number already exists")
	}

	// Генерируем уникальный идентификатор пользователя
	userID := uuid.New().String()

	// Создаем учетные данные пользователя
	user := &voice.User{
		UserID:    userID,
		Phone:     phone,
		Password:  password,
		Firstname: firstname,
		Surname:   surname,
	}

	// Сохраняем учетные данные в репозитории
	err = vl.repository.CreateUser(user)
	if err != nil {
		return "", err
	}

	// Создаем банковский счет для пользователя
	accountNumber := generateAccountNumber()
	bankAccount := &voice.BankAccount{
		UserID:  userID,
		Account: accountNumber,
	}

	// Сохраняем банковский счет в репозитории
	err = vl.repository.SaveBankAccount(bankAccount)
	if err != nil {
		return "", err
	}

	return userID, nil
}

// Authenticate аутентифицирует пользователя по номеру телефона и паролю
func (vl *VoiceLogic) Authenticate(phone, password string) (string, error) {
	user, err := vl.repository.GetUserByPhone(phone)
	if err != nil {
		return "", err
	}
	if user == nil || user.Password != password {
		return "", errors.New("invalid phone or password")
	}
	return user.UserID, nil
}

// GetUserByPhone возвращает пользователя по его номеру телефона
func (vl *VoiceLogic) GetUserByPhone(phone string) (*voice.User, error) {
	user, err := vl.repository.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ProcessVoiceRecordingsRequest обрабатывает запрос на создание или обновление голосовых записей
func (vl *VoiceLogic) ProcessVoiceRecordingsRequest(r *http.Request) (uuid.UUID, []byte, []byte, []byte, error) {
	// Парсим multipart/form-data
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		return uuid.UUID{}, nil, nil, nil, err
	}

	// Получаем идентификатор пользователя из URL
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["user-id"])
	if err != nil {
		return uuid.UUID{}, nil, nil, nil, errors.New("invalid user ID")
	}

	// Читаем файлы из формы
	audio1, err := readFileFromRequest(r, "audio1")
	if err != nil {
		return uuid.UUID{}, nil, nil, nil, err
	}

	audio2, err := readFileFromRequest(r, "audio2")
	if err != nil {
		return uuid.UUID{}, nil, nil, nil, err
	}

	audio3, err := readFileFromRequest(r, "audio3")
	if err != nil {
		return uuid.UUID{}, nil, nil, nil, err
	}

	return userID, audio1, audio2, audio3, nil
}

// readFileFromRequest читает файл из запроса
func readFileFromRequest(r *http.Request, fieldName string) ([]byte, error) {
	file, _, err := r.FormFile(fieldName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}
