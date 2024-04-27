CREATE Table IF NOT EXISTS User  (
    id INT Primary Key AUTO_INCREMENT,
    email VARCHAR(255) unique,
    password VARCHAR(255),
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE Table IF NOT EXISTS Profile (
    id INT Primary Key AUTO_INCREMENT,
    avatar VARCHAR(255),
    displayName VARCHAR(255),
    userName VARCHAR(255) unique,
    phoneNumber VARCHAR(255) unique,
    status ENUM('online', 'away', 'offline') DEFAULT 'offline',
    customStatus VARCHAR(255),
    clearAfter DATETIME,
    userId INT unique NOT NULL,
    FOREIGN KEY (userId) REFERENCES User(id)
)