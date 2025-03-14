CREATE TABLE IF NOT EXISTS content (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL, -- Оставляем, но без внешнего ключа
    file_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    start_time TEXT NOT NULL,
    end_time TEXT NOT NULL,
    macaddress VARCHAR(17) NOT NULL,
    latest_history INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS content_history (
    id SERIAL PRIMARY KEY,
    content_id INTEGER NOT NULL,
    status_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    user_id INTEGER NOT NULL, -- Убрали внешний ключ
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS status (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

INSERT INTO status (id, name) VALUES
(1, 'Created'),
(2, 'Under Review'),
(3, 'Approved'),
(4, 'Rejected')
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS monitors (
    id SERIAL PRIMARY KEY,
    building VARCHAR NOT NULL,
    floor INTEGER NOT NULL,
    notes TEXT,
    monitor_resolution VARCHAR NOT NULL,
    ip TEXT NOT NULL,
    status BOOLEAN NOT NULL DEFAULT TRUE,
    macaddress VARCHAR(17) UNIQUE NOT NULL
);


INSERT INTO monitors (building, floor, monitor_resolution, macaddress, status, notes) 
VALUES 
('УЛК', 2, '1920x1080', 'AA:BB:CC:DD:EE:01', FALSE, 'Новый монитор на 2 этаже'),
('УЛК', 3, '1920x1080', 'AA:BB:CC:DD:EE:02', TRUE, 'Рабочий монитор в аудитории 305'),
('УЛК', 4, '1280x720', 'AA:BB:CC:DD:EE:03', FALSE, 'Требуется проверка соединения'),

('ФИТ', 1, '1920x1080', 'AA:BB:CC:DD:EE:04', TRUE, 'Установлен в холле 1 этажа'),
('ФИТ', 2, '2560x1440', 'AA:BB:CC:DD:EE:05', FALSE, 'Монитор в лаборатории 207'),
('ФИТ', 3, '1920x1080', 'AA:BB:CC:DD:EE:06', TRUE, 'Экран для презентаций'),

('ЦИСИ', 1, '1920x1080', 'AA:BB:CC:DD:EE:07', FALSE, 'Проблемы с отображением'),
('ЦИСИ', 2, '1280x720', 'AA:BB:CC:DD:EE:08', TRUE, 'Используется для информационных объявлений'),
('ЦИСИ', 3, '1920x1080', 'AA:BB:CC:DD:EE:09', FALSE, 'Ожидает настройки'),
('ЦИСИ', 4, '2560x1440', 'AA:BB:CC:DD:EE:10', TRUE, 'Монитор в конференц-зале');
