-- Добавляем поле price в таблицу seats
ALTER TABLE seats ADD COLUMN price DECIMAL(10,2) NOT NULL DEFAULT 0.00;

-- Обновляем существующие записи, устанавливая цену по умолчанию
-- В реальном приложении здесь должна быть логика установки правильных цен
UPDATE seats SET price = 1000.00 WHERE price = 0.00;

-- Убираем значение по умолчанию после обновления
ALTER TABLE seats ALTER COLUMN price DROP DEFAULT;
