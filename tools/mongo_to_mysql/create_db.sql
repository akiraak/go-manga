CREATE DATABASE manganow CHARACTER SET utf8mb4;

CREATE TABLE manganow.publishers (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name CHAR(191) BINARY NOT NULL UNIQUE,
    date_created DATETIME NOT NULL
);

CREATE TABLE manganow.authors (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name CHAR(191) BINARY NOT NULL UNIQUE,
    date_created DATETIME NOT NULL
);

CREATE TABLE manganow.books (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    kindle BOOLEAN NOT NULL,
    date_publish DATETIME NOT NULL,
    image_s_url CHAR(255) NOT NULL,
    image_s_width INT NOT NULL,
    image_s_height INT NOT NULL,
    image_m_url CHAR(255) NOT NULL,
    image_m_width INT NOT NULL,
    image_m_height INT NOT NULL,
    image_l_url CHAR(255) NOT NULL,
    image_l_width INT NOT NULL,
    image_l_height INT NOT NULL,
    asin CHAR(255) NOT NULL,
    title TEXT NOT NULL,
    region CHAR(255) NOT NULL,
    publisher_id BIGINT UNSIGNED,
    author_id BIGINT UNSIGNED,
    date_last_modify DATETIME NOT NULL,
    date_created DATETIME NOT NULL,
    INDEX(date_publish),
    FOREIGN KEY(publisher_id) REFERENCES publishers(id),
    FOREIGN KEY(author_id) REFERENCES authors(id)
);
