
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `bbs_comments` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` char(255) NOT NULL,
  `comment` text NOT NULL,
  `ip_hash` char(255) NOT NULL,
  `updated_at` datetime NOT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `bbs_comments`;
