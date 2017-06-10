
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE publishers ADD ero BOOLEAN NOT NULL DEFAULT 0;
ALTER TABLE publishers ADD INDEX idx_ero(ero);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE publishers DROP INDEX idx_ero;
ALTER TABLE publishers DROP ero;
