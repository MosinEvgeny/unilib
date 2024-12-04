CREATE TABLE books (
    book_id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn VARCHAR(20) UNIQUE,
    publisher VARCHAR(255),
    publication_year INT,
    total_copies INT DEFAULT 1 CHECK (total_copies >= 0),
    category VARCHAR(100),
    description TEXT
);

CREATE TABLE copies (
    copy_id SERIAL PRIMARY KEY,
    book_id INT REFERENCES books(book_id) ON DELETE CASCADE, 
    inventory_number VARCHAR(20) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'Доступен' CHECK (status IN ('Доступен', 'Выдан', 'Списан')),
    acquisition_date DATE
);

CREATE TABLE readers (
    reader_id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    faculty VARCHAR(100),
    course INT CHECK (course BETWEEN 1 AND 6),
    student_id VARCHAR(20) UNIQUE,
    phone_number VARCHAR(255),
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL DEFAULT 'reader',
    registration_date DATE
);

CREATE TABLE issue (
    issue_id SERIAL PRIMARY KEY,
    copy_id INT REFERENCES copies(copy_id) ON DELETE RESTRICT,
    reader_id INT REFERENCES readers(reader_id) ON DELETE RESTRICT,
    issue_date DATE NOT NULL,
    due_date DATE NOT NULL,
    return_date DATE
);

CREATE TABLE reports (
    report_id SERIAL PRIMARY KEY,
    report_date DATE NOT NULL,
    report_content TEXT
);

-- Функция для проверки уникальности логина и студенческого билета
CREATE OR REPLACE FUNCTION check_unique_reader()
RETURNS TRIGGER AS $$
BEGIN
    -- Проверка на уникальность логина
    IF EXISTS (SELECT 1 FROM readers WHERE username = NEW.username AND reader_id != NEW.reader_id) THEN
        RAISE EXCEPTION 'Пользователь с таким логином уже существует';
    END IF;

    -- Проверка на уникальность студенческого билета (если он указан)
    IF NEW.student_id IS NOT NULL AND EXISTS (SELECT 1 FROM readers WHERE student_id = NEW.student_id AND reader_id != NEW.reader_id) THEN
        RAISE EXCEPTION 'Читатель с таким студенческим билетом уже существует';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- Триггер для проверки перед вставкой или обновлением
CREATE TRIGGER trigger_check_unique_reader
BEFORE INSERT OR UPDATE ON readers
FOR EACH ROW
EXECUTE PROCEDURE check_unique_reader();