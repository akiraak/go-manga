
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` char(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `updated_at` datetime NOT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `users`;
