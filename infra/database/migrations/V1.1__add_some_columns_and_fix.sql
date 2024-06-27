-- V1.1__add_status_column_to_users.sql
ALTER TABLE users
ADD COLUMN status INT DEFAULT 1 NOT NULL;

ALTER TABLE comments
ADD COLUMN status INT DEFAULT 1 NOT NULL;