-- +goose Up
CREATE TABLE `vote` (
  `id` int NOT NULL AUTO_INCREMENT,
  `from_user` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `to_user` varchar(255) DEFAULT NULL,
  `vote` tinyint DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `from_user_id` (`from_user`),
  KEY `to_user` (`to_user`),
  CONSTRAINT `vote_ibfk_1` FOREIGN KEY (`from_user`) REFERENCES `users` (`uuid`),
  CONSTRAINT `vote_ibfk_2` FOREIGN KEY (`to_user`) REFERENCES `users` (`uuid`),
  CONSTRAINT `check_vote` CHECK ((`vote` in (-(1),1))),
  CONSTRAINT one_vote_per_user UNIQUE (from_user, to_user),
  CONSTRAINT user_cant_vote_oneself CHECK (from_user <> to_user)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- +goose Down
DROP table `vote`
