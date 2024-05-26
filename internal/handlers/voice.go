package handlers

import (
	"io/ioutil"
	"log"
	"myapp/internal/logic"
	"myapp/internal/voice"
	"myapp/pkg/io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
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

func CompareAudio(audioData1, audioData2 []byte) float64 {
	// Создаем временные файлы для .wav файлов
	tempFile1, err := os.CreateTemp("", "audio1_*.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tempFile1.Name())
	if _, err := tempFile1.Write(audioData1); err != nil {
		log.Fatal(err)
	}

	tempFile2, err := os.CreateTemp("", "audio2_*.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tempFile2.Name())
	if _, err := tempFile2.Write(audioData2); err != nil {
		log.Fatal(err)
	}

	// Вызываем питоновскую функцию для сравнения аудиофайлов
	cmd := exec.Command("python", "compare_audio.py", tempFile1.Name(), tempFile2.Name())
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	// Преобразуем вывод в число с плавающей точкой
	similarity, err := strconv.ParseFloat(string(output), 64)
	if err != nil {
		log.Fatal(err)
	}

	return similarity
}

// handleVoiceVerification обрабатывает запрос на проверку голоса
func (h *Handlers) VoiceVerificationHandler(w http.ResponseWriter, r *http.Request) {
	// Читаем файл из тела запроса
	audioFile, _, err := r.FormFile("audio1")
	if err != nil {
		io.SendError(w, "Ошибка при чтении данных из тела запроса", http.StatusBadRequest)
		log.Println(err)
		return
	}
	defer audioFile.Close()

	// Читаем содержимое аудиофайла в []byte
	audioBytes, err := ioutil.ReadAll(audioFile)
	if err != nil {
		http.Error(w, "Failed to read audio file", http.StatusInternalServerError)
		return
	}

	// Получаем список пользователей и их голосовых записей
	users, err := h.logic.GetUsersVoices()
	if err != nil {
		io.SendError(w, "Ошибка при получении данных о пользователях", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Проходим по каждому пользователю
	for _, user := range users {
		// Сравниваем бинарные данные с первой голосовой записью
		similarity1, err := h.logic.CompareUserVoices(audioBytes, user.Voice_sample1)
		if err != nil {
			io.SendError(w, "Ошибка при сравнении с первой голосовой записью", http.StatusInternalServerError)
			log.Println("Ошибка при сравнении с первой голосовой записью:", err)
			continue
		}

		// Если сходство с первой записью меньше 60%, переходим к следующему пользователю
		if similarity1 < 60 {
			continue
		}

		// Сравниваем бинарные данные с остальными голосовыми записями
		similarity2, err := h.logic.CompareUserVoices(audioBytes, user.Voice_sample2)
		if err != nil {
			io.SendError(w, "Ошибка при сравнении с второй голосовой записью", http.StatusInternalServerError)
			log.Println("Ошибка при сравнении с второй голосовой записью:", err)
			continue
		}

		similarity3, err := h.logic.CompareUserVoices(audioBytes, user.Voice_sample3)
		if err != nil {
			io.SendError(w, "Ошибка при сравнении с третьей голосовой записью", http.StatusInternalServerError)
			log.Println("Ошибка при сравнении с третьей голосовой записью:", err)
			continue
		}

		// Если среднее значение сходства больше 75%, возвращаем user_id
		averageSimilarity := (similarity1 + similarity2 + similarity3) / 3
		if averageSimilarity > 75 {
			response := struct {
				user_id string `json:"user_id"`
			}{user_id: user.UserID.String()}
			io.WriteJSON(w, http.StatusOK, response)
			return
		}
	}

	// Если ни один пользователь не найден, возвращаем пустой ответ
	w.WriteHeader(http.StatusNotFound)
}
