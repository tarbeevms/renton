package handlers

import (
	"log"
	"myapp/internal/logic"
	"myapp/internal/voice"
	"myapp/pkg/io"
	"net/http"
)

// Handlers содержит обработчики HTTP-запросов
type Handlers struct {
	logic *logic.VoiceLogic
}

// NewHandlers создает новый экземпляр Handlers
func NewHandlers(logic *logic.VoiceLogic) *Handlers {
	return &Handlers{logic: logic}
}

// RegisterHandler обрабатывает запрос на регистрацию пользователя
func (h *Handlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Читаем данные пользователя из тела запроса
	var user voice.User
	err := io.ReadJSON(r, &user)
	if err != nil {
		io.SendError(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	// Регистрируем пользователя в базе данных
	userID, err := h.logic.RegisterUser(user.Phone, user.Password, user.Firstname, user.Surname)
	if err != nil {
		io.SendError(w, "Error registering user", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	// Создаем куки-файл для сессии пользователя
	cookie, err := h.logic.CreateSessionCookie(userID)
	if err != nil {
		io.SendError(w, "Error creating session cookie", http.StatusInternalServerError)
		return
	}
	// Устанавливаем куки-файл в ответе сервера
	http.SetCookie(w, cookie)
	// Устанавливаем заголовок Location на главную страницу
	io.WriteJSON(w, http.StatusOK, map[string]string{"message": "User registered successfully"})
}

// LoginHandler обрабатывает запрос на аутентификацию пользователя
func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Читаем данные аутентификации пользователя из тела запроса
	var parsedUser voice.User
	if err := io.ReadJSON(r, &parsedUser); err != nil {
		io.SendError(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Получаем пользователя по его номеру телефона из базы данных
	user, err := h.logic.GetUserByPhone(parsedUser.Phone)
	if err != nil {
		io.SendError(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Проверяем правильность пароля
	if user.Password != parsedUser.Password {
		io.SendError(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Создаем куки-файл для сессии пользователя
	cookie, err := h.logic.CreateSessionCookie(user.UserID)
	if err != nil {
		io.SendError(w, "Error creating session cookie", http.StatusInternalServerError)
		return
	}

	// Устанавливаем куки-файл в ответе сервера
	http.SetCookie(w, cookie)

	// Возвращаем успешный ответ с данными пользователя
	io.WriteJSON(w, http.StatusOK, map[string]string{"message": "User authenticated successfully"})
}

// VoiceCreationHandler обрабатывает запрос на создание голосовых записей
func (h *Handlers) VoiceCreationHandler(w http.ResponseWriter, r *http.Request) {
	userID, audio1, audio2, audio3, err := h.logic.ProcessVoiceRecordingsRequest(r)
	if err != nil {
		io.SendError(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Создаем голосовые записи
	err = h.logic.CreateVoiceRecordings(userID, audio1, audio2, audio3)
	if err != nil {
		io.SendError(w, "Error creating voice recordings", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	io.SendError(w, "Voice recordings created successfully", http.StatusOK)
}

// VoiceUpdateHandler обрабатывает запрос на перезапись голосовых записей
func (h *Handlers) VoiceUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Реализация обработки перезаписи голосовых записей
}
