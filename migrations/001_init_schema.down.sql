-- migrations/001_init_schema.down.sql

-- Сначала удаляем таблицы, которые ссылаются на другие (дочерние)
DROP TABLE IF EXISTS appointments CASCADE;
DROP TABLE IF EXISTS slots CASCADE;

-- Потом основную таблицу
DROP TABLE IF EXISTS users CASCADE;