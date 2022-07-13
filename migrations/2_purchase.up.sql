CREATE TABLE purchases (
	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, 
	buying_price DECIMAL(20, 2),
	bought_at TIMESTAMP NOT NULL,
	wager_id BIGINT NOT NULL REFERENCES wagers(id)
);