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

INSERT INTO books (title, author, isbn, publisher, publication_year, total_copies, category, description) VALUES
('Преступление и наказание', 'Ф.М. Достоевский', '978-5-17-098231-7', 'Азбука', 2017, 10, 'Классическая литература', 'Роман о преступлении и его последствиях.'),
('Война и мир', 'Л.Н. Толстой', '978-5-389-14036-0', 'Эксмо', 2018, 5, 'Классическая литература', 'Эпический роман о России в эпоху наполеоновских войн.'),
('Мастер и Маргарита', 'М.А. Булгаков', '978-5-04-118828-8', 'Эксмо', 2021, 8, 'Классическая литература', 'Мистический роман о добре и зле.'),
('Евгений Онегин', 'А.С. Пушкин', '978-5-699-99585-4', 'Эксмо', 2019, 12, 'Классическая литература', 'Роман в стихах о жизни русской аристократии.'),
('Мертвые души', 'Н.В. Гоголь', '978-5-17-114057-3', 'АСТ', 2020, 6, 'Классическая литература', 'Сатирический роман о похождениях Чичикова.'),
('Анна Каренина', 'Л.Н. Толстой', '978-5-389-15122-9', 'Эксмо', 2019, 4, 'Классическая литература', 'Роман о трагической любви замужней женщины.'),
('Горе от ума', 'А.С. Грибоедов', '978-5-17-102604-8', 'АСТ', 2019, 7, 'Классическая литература', 'Комедия в стихах о нравах московского дворянства.'),
('Отцы и дети', 'И.С. Тургенев', '978-5-17-110583-1', 'АСТ', 2020, 9, 'Классическая литература', 'Роман о конфликте поколений.'),
('История государства Российского', 'Н.М. Карамзин', '978-5-4453-0186-5', 'Лабиринт', 2016, 3, 'История', 'Фундаментальный труд по истории России с древнейших времен.'),
('Курс русской истории', 'В.О. Ключевский', '978-5-4468-1875-2', 'АСТ', 2018, 2, 'История', 'Классический курс лекций по истории России.'),
('Философия права', 'Г.В.Ф. Гегель', '978-5-4444-8157-5', 'Издательский дом ''Территория будущего''', 2010, 4, 'Философия', 'Фундаментальный труд по философии права.'),
('Так говорил Заратустра', 'Ф. Ницше', '978-5-17-098093-1', 'АСТ', 2018, 6, 'Философия', 'Философский роман о преодолении человеком себя.'),
('Бытие и ничто', 'Ж.-П. Сартр', '978-5-91579-102-2', 'Академический Проект', 2015, 2, 'Философия', 'Фундаментальный труд по экзистенциализму.'),
('Чистое искусство', 'В.В. Кандинский', '978-5-903060-68-2', 'ГИТИС', 2016, 3, 'Искусство', 'Теоретический труд о абстрактном искусстве.'),
('Черный квадрат', 'К.С. Малевич', '978-5-9909896-1-4', 'Азбука', 2019, 5, 'Искусство', 'Сборник статей об искусстве.'),
('Психология интеллекта', 'Ж. Пиаже', '978-5-496-01300-5', 'Питер', 2016, 4, 'Психология', 'Классический труд по психологии развития интеллекта.'),
('Психология масс и анализ человеческого "Я"', 'З. Фрейд', '978-5-17-100445-9', 'АСТ', 2019, 7, 'Психология', 'Сборник работ по социальной психологии.'),
('1984', 'Джордж Оруэлл', '978-5-17-090411-1', 'АСТ', 2018, 10, 'Фантастика', 'Антиутопический роман о тоталитарном обществе будущего.'),
('451 градус по Фаренгейту', 'Рэй Брэдбери', '978-5-17-114679-7', 'АСТ', 2020, 8, 'Фантастика', 'Роман о тоталитарном обществе, где книги находятся под запретом.'),
('Солярис', 'Станислав Лем', '978-5-17-113862-4', 'АСТ', 2020, 6, 'Фантастика', 'Философский роман о контакте с внеземным разумом.');

INSERT INTO copies (book_id, inventory_number, status, acquisition_date)
SELECT 
    b.book_id, 
    b.book_id || '-' || generate_series(1, b.total_copies),  -- Генерация инвентарного номера 
    'Доступен', 
    CURRENT_DATE
FROM books AS b;