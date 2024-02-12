-- +goose Up
CREATE TABLE `users` (
  `uuid` varchar(255) NOT NULL,
  `email` varchar(255) DEFAULT NULL,
  `password` longtext,
  `first_name` longtext,
  `last_name` longtext,
  `nickname` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`uuid`),
  UNIQUE KEY `email` (`email`),
  KEY `idx_users_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- +goose Down
DROP TABLE `users`;

