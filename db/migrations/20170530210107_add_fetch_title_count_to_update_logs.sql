
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE update_logs ADD fetch_title_count INT NOT NULL DEFAULT 0;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE update_logs DROP fetch_title_count;
