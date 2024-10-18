
CREATE TABLE IF NOT EXISTS `product`(
    productId int PRIMARY KEY AUTO_INCREMENT,
    shopId int NOT NULL,
    FOREIGN KEY (shopId) REFERENCES shops(shopId),
    productName VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS `shops`(
    shopId int PRIMARY KEY AUTO_INCREMENT,
    shopname VARCHAR(255) NOT NULL
);
