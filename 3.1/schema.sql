CREATE TABLE `users` (
  `uuid` varchar(191),
  `email` varchar(191) DEFAULT NULL,
  `password` longtext,
  `first_name` longtext,
  `last_name` longtext,
  `nickname` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`uuid`),
  UNIQUE KEY `email` (`email`),
  KEY `idx_users_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci 
