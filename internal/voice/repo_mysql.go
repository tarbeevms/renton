package voice

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// MySQLRepository реализация интерфейса UserRepository для работы с MySQL
type MySQLRepository struct {
	db *sql.DB
}

// NewMySQLRepository создает новый экземпляр MySQLRepository
func NewMySQLRepository(db *sql.DB) UserRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) CreateUser(user *User) error {
	_, err := r.db.Exec("INSERT INTO Credentials (user_id, phone_number, password, firstname, surname) VALUES (?, ?, ?, ?, ?)", user.UserID, user.Phone, user.Password, user.Firstname, user.Surname)
	if err != nil {
		return err
	}
	return nil
}

func (r *MySQLRepository) GetUserByPhone(phone string) (*User, error) {
	var user User
	err := r.db.QueryRow("SELECT user_id, phone_number, password FROM Credentials WHERE phone_number = ?", phone).Scan(&user.UserID, &user.Phone, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *MySQLRepository) SaveVoiceRecordings(recordings *VoiceRecording) error {
	_, err := r.db.Exec("INSERT INTO UsersVoices (user_id, voice_sample1, voice_sample2, voice_sample3) VALUES (?, ?, ?, ?)",
		recordings.UserID, recordings.Voice_sample1, recordings.Voice_sample2, recordings.Voice_sample3)
	if err != nil {
		return err
	}

	return nil
}

func (r *MySQLRepository) UpdateVoiceRecordings(recordings *VoiceRecording) error {
	_, err := r.db.Exec("UPDATE UsersVoices SET voice_sample1 = ?, voice_sample2 = ?, voice_sample3 = ? WHERE user_id = ?",
		recordings.Voice_sample1, recordings.Voice_sample2, recordings.Voice_sample3, recordings.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (r *MySQLRepository) GenerateBankAccount() string {
	return fmt.Sprintf("%020d", uuid.New().ID())
}

// SaveBankAccount сохраняет банковский счет пользователя в базе данных
func (r *MySQLRepository) SaveBankAccount(account *BankAccount) error {
	_, err := r.db.Exec("INSERT INTO BankAccounts (user_id, account) VALUES (?, ?)", account.UserID, account.Account)
	if err != nil {
		return err
	}
	return nil
}

// GetUsersVoices получает данные о пользователях из базы данных
func (r *MySQLRepository) GetUsersVoices() ([]*VoiceRecording, error) {
	rows, err := r.db.Query("SELECT user_id, voice_sample1, voice_sample2, voice_sample3 FROM UsersVoices")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*VoiceRecording
	for rows.Next() {
		var user VoiceRecording
		if err := rows.Scan(&user.UserID, &user.Voice_sample1, &user.Voice_sample2, &user.Voice_sample3); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
