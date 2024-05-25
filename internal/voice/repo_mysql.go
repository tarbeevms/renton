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
	_, err := r.db.Exec("INSERT INTO Credentials (userid, phone_number, password, firstname, surname) VALUES (?, ?, ?, ?, ?)", user.UserID, user.Phone, user.Password, user.Firstname, user.Surname)
	if err != nil {
		return err
	}
	return nil
}

func (r *MySQLRepository) GetUserByPhone(phone string) (*User, error) {
	var user User
	err := r.db.QueryRow("SELECT userid, phone_number, password FROM Credentials WHERE phone_number = ?", phone).Scan(&user.UserID, &user.Phone, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *MySQLRepository) SaveVoiceRecordings(recordings *VoiceRecording) error {
	_, err := r.db.Exec("INSERT INTO UsersVoices (userid, voice_sample1, voice_sample2, voice_sample3) VALUES (?, ?, ?, ?)",
		recordings.UserID, recordings.Audio1, recordings.Audio2, recordings.Audio3)
	if err != nil {
		return err
	}

	return nil
}

func (r *MySQLRepository) UpdateVoiceRecordings(recordings *VoiceRecording) error {
	_, err := r.db.Exec("UPDATE UsersVoices SET voice_sample1 = ?, voice_sample2 = ?, voice_sample3 = ? WHERE userid = ?",
		recordings.Audio1, recordings.Audio2, recordings.Audio3, recordings.UserID)
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
