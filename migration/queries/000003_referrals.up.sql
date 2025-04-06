-- Таблица с рефералами. У каждого юзера может быть только 1 реферал (либо не быть)
-- При удалении юзера из Credential из этой таблицы также удаляются все поля с referral_id и id
CREATE TABLE IF NOT EXISTS Referrals (
    id integer PRIMARY KEY,
    referrer_id integer,
    CONSTRAINT creds_fk
    FOREIGN KEY (id) REFERENCES Credentials(id) ON DELETE CASCADE,
    CONSTRAINT creds_referrals_fk
    FOREIGN KEY (referrer_id)  REFERENCES Credentials(id) ON DELETE SET NULL 
);

-- Таблица с доступными мессенджерами
CREATE TABLE IF NOT EXISTS Messengers (
    id serial PRIMARY KEY,
    name character varying(255) NOT NULL
);

-- Таблица с доступными действиями в мессенджерах
CREATE TABLE IF NOT EXISTS Actions (
    id serial PRIMARY KEY,
    name character varying(255) NOT NULL
);

-- Таблица с доступными задачами
CREATE TABLE IF NOT EXISTS Tasks (
    id serial PRIMARY KEY,
    messenger_id integer NOT NULL,
    action_id integer NOT NULL,
    task_data character varying(255) NOT NULL,
    award numeric NOT NULL,
    CONSTRAINT messenger_fk
    FOREIGN KEY (messenger_id) REFERENCES Messengers(id) ON DELETE CASCADE,
    CONSTRAINT action_fk
    FOREIGN KEY (action_id) REFERENCES Actions(id) ON DELETE CASCADE
);

-- Таблица с данными о действиях пользователей
CREATE TABLE IF NOT EXISTS TaskLogs (
    id serial PRIMARY KEY,
    user_id integer NOT NULL,
    task_id integer NOT NULL,
    award numeric NOT NULL,
    CONSTRAINT user_fk
    FOREIGN KEY (user_id) REFERENCES Credentials(id),
    CONSTRAINT task_fk
    FOREIGN KEY (task_id) REFERENCES Tasks(id)
);
