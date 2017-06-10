
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE publishers DROP INDEX idx_ero;
ALTER TABLE publishers CHANGE ero r18 BOOLEAN NOT NULL DEFAULT 0;
ALTER TABLE publishers ADD INDEX idx_r18(r18);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE publishers DROP INDEX idx_r18;
ALTER TABLE publishers CHANGE r18 ero BOOLEAN NOT NULL DEFAULT 0;
ALTER TABLE publishers ADD INDEX idx_ero(ero);
