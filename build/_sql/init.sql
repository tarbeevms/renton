-- Создание таблицы Credentials
CREATE TABLE Credentials (
    user_id CHAR(36) PRIMARY KEY NOT NULL UNIQUE,
    phone_number VARCHAR(20) NOT NULL,
    password VARCHAR(255) NOT NULL,
    firstname VARCHAR(20) NOT NULL,
    surname VARCHAR(20) NOT NULL
);

-- Создание таблицы BankAccounts
CREATE TABLE BankAccounts (
    user_id CHAR(36) NOT NULL,
    account VARCHAR(34) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Credentials(user_id)
);

-- Создание таблицы UsersVoices
CREATE TABLE UsersVoices (
    user_id CHAR(36) NOT NULL,
    voice_sample1 BLOB NOT NULL,
    voice_sample2 BLOB NOT NULL,
    voice_sample3 BLOB NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Credentials(user_id)
);
