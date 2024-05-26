package voice

import "github.com/google/uuid"

// User структура для хранения информации о пользователе
type User struct {
	UserID    string `json:"userid"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Surname   string `json:"surname"`
}

// BankAccount структура для хранения информации о банковском счете
type BankAccount struct {
	UserID  string `json:"userid"`
	Account string `json:"account"`
}

// VoiceRecording структура для хранения информации о голосовых записях пользователя
type VoiceRecording struct {
	UserID        uuid.UUID `json:"userid"`
	Voice_sample1 []byte    `json:"-"`
	Voice_sample2 []byte    `json:"-"`
	Voice_sample3 []byte    `json:"-"`
}

// UserRepository представляет интерфейс для работы с пользователями и их банковскими счетами в базе данных
type UserRepository interface {
	CreateUser(user *User) error
	GetUserByPhone(phone string) (*User, error)
	SaveVoiceRecordings(recordings *VoiceRecording) error
	UpdateVoiceRecordings(recordings *VoiceRecording) error
	SaveBankAccount(account *BankAccount) error
	GetUsersVoices() ([]*VoiceRecording, error)
}
