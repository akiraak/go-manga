CREATE DATABASE manganow CHARACTER SET utf8mb4;
USE manganow;

CREATE TABLE manganow.titles (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME NOT NULL
);

CREATE TABLE manganow.publishers (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name CHAR(191) BINARY NOT NULL UNIQUE,
    created_at DATETIME NOT NULL
);

CREATE TABLE manganow.authors (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name CHAR(191) BINARY NOT NULL UNIQUE,
    created_at DATETIME NOT NULL
);

CREATE TABLE manganow.books (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    asin CHAR(191) NOT NULL UNIQUE,
    publish_type CHAR(255),
    name TEXT NOT NULL,
    region CHAR(255) NOT NULL,
    date_publish CHAR(8) NOT NULL,
    image_s_url CHAR(255) NOT NULL,
    image_s_width INT NOT NULL,
    image_s_height INT NOT NULL,
    image_m_url CHAR(255) NOT NULL,
    image_m_width INT NOT NULL,
    image_m_height INT NOT NULL,
    image_l_url CHAR(255) NOT NULL,
    image_l_width INT NOT NULL,
    image_l_height INT NOT NULL,
    xml TEXT,
    title_id BIGINT UNSIGNED,
    publisher_id BIGINT UNSIGNED,
    author_id BIGINT UNSIGNED,
    updated_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    INDEX(date_publish),
    FOREIGN KEY(title_id) REFERENCES titles(id),
    FOREIGN KEY(publisher_id) REFERENCES publishers(id),
    FOREIGN KEY(author_id) REFERENCES authors(id)
);

CREATE TABLE manganow.update_logs (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    date DATETIME NOT NULL,
    fetch_asin_count INT NOT NULL,
	update_asin_count INT NOT NULL,
	updated_book_count INT NOT NULL,
    updated_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    INDEX(date)
);

CREATE TABLE manganow.one (
    last_update_book_page INT NOT NULL
);
INSERT INTO manganow.one (last_update_book_page) VALUES (0);
