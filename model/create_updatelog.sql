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
