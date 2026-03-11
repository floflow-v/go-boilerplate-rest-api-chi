-- +goose Up
CREATE TABLE
	IF NOT EXISTS authors (
		`id` CHAR(36) NOT NULL UNIQUE,
		`name` VARCHAR(255) NOT NULL UNIQUE,
		PRIMARY KEY (`id`)
	);

-- +goose Down
DROP TABLE IF EXISTS authors
