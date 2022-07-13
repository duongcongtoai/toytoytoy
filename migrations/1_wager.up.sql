CREATE TABLE wagers (
	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	total_wager_value INT NOT NULL,
	odds INT NOT NULL,
	selling_percentage INT NOT NULL,
	selling_price DECIMAL(20, 2) NOT NULL,
	current_selling_price DECIMAL(20, 2) NOT NULL,
	percentage_sold INT NULL,
	amount_sold DECIMAL(20, 2) NULL,
	placed_at TIMESTAMP NOT NULL
);

