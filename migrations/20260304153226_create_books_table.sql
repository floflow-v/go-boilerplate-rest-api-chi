-- +goose Up
CREATE TABLE
    IF NOT EXISTS books (
        `id` CHAR(36) NOT NULL UNIQUE,
        `title` VARCHAR(255) NOT NULL UNIQUE,
        `description` VARCHAR(255) NOT NULL,
        `author_id` CHAR(36) NOT NULL,
        PRIMARY KEY (`id`),
        FOREIGN KEY (author_id) REFERENCES authors(id)
    );

-- +goose Down
DROP TABLE IF EXISTS books
