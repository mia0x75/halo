-- Target DB for integration tests (MySQL).
-- Port: 3306 (host).
-- Includes both compliant and non-compliant tables for rule hit/miss tests.

CREATE DATABASE IF NOT EXISTS starwars
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_unicode_ci;

USE starwars;

-- Compliant tables (rule miss).
CREATE TABLE IF NOT EXISTS users (
    id     INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'user id',
    name   VARCHAR(50)  NOT NULL COMMENT 'user name',
    email  VARCHAR(75)  NOT NULL COMMENT 'user email',
    status TINYINT      NOT NULL DEFAULT 1 COMMENT '1 active, 0 disabled',
    PRIMARY KEY (id),
    UNIQUE KEY unique_1 (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='users';

CREATE TABLE IF NOT EXISTS orders (
    id         INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'order id',
    user_id    INT UNSIGNED NOT NULL COMMENT 'user id',
    total      DECIMAL(18,4) NOT NULL DEFAULT 0 COMMENT 'total',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created',
    PRIMARY KEY (id),
    KEY index_1 (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='orders';

INSERT INTO users (name, email, status) VALUES
    ('luke', 'luke@example.com', 1),
    ('leia', 'leia@example.com', 1),
    ('han',  'han@example.com',  0);

INSERT INTO orders (user_id, total) VALUES
    (1, 100.00),
    (2, 200.00);
