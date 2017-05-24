CREATE DATABASE manganow CHARACTER SET utf8mb4;
USE manganow;

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
    tree_type ENUM('main','sub') NOT NULL,
    asin CHAR(191) NOT NULL UNIQUE,
    sub_asins CHAR(255),
    publish_type CHAR(255),
    title TEXT NOT NULL,
    region CHAR(255) NOT NULL,
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
    publisher_id BIGINT UNSIGNED,
    author_id BIGINT UNSIGNED,
    updated_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    INDEX(tree_type),
    INDEX(date_publish),
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
