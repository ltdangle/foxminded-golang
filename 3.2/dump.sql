-- -------------------------------------------------------------
-- TablePlus 5.6.6(520)
--
-- https://tableplus.com/
--
-- Database: fox_rest
-- Generation Time: 2023-12-06 13:53:56.4200
-- -------------------------------------------------------------


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


DROP TABLE IF EXISTS `user_role`;
CREATE TABLE `user_role` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `title` (`title`),
  KEY `idx_role_title` (`title`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `uuid` varchar(255) NOT NULL,
  `email` varchar(255) DEFAULT NULL,
  `password` longtext,
  `first_name` longtext,
  `last_name` longtext,
  `nickname` longtext,
  `role_id` int DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`uuid`),
  UNIQUE KEY `email` (`email`),
  KEY `idx_users_email` (`email`),
  KEY `role_id` (`role_id`),
  CONSTRAINT `users_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `user_role` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;;

DROP TABLE IF EXISTS `vote`;
CREATE TABLE `vote` (
  `id` int NOT NULL AUTO_INCREMENT,
  `from_user` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `to_user` varchar(255) DEFAULT NULL,
  `vote` tinyint DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `one_vote_per_user` (`from_user`,`to_user`),
  KEY `from_user_id` (`from_user`),
  KEY `to_user` (`to_user`),
  CONSTRAINT `vote_ibfk_1` FOREIGN KEY (`from_user`) REFERENCES `users` (`uuid`),
  CONSTRAINT `vote_ibfk_2` FOREIGN KEY (`to_user`) REFERENCES `users` (`uuid`),
  CONSTRAINT `check_vote` CHECK ((`vote` in (-(1),1))),
  CONSTRAINT `user_cant_vote_oneself` CHECK ((`from_user` <> `to_user`))
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;;

INSERT INTO `user_role` (`id`, `title`, `description`) VALUES
(1, 'admin', 'system admin'),
(2, 'user', 'regular user'),
(3, 'moderator', 'moderator');

INSERT INTO `users` (`uuid`, `email`, `password`, `first_name`, `last_name`, `nickname`, `role_id`, `created_at`, `updated_at`) VALUES
('e0e5ba28-19fc-4c65-8692-f61266608d4n', 'user2@example2.com', 'password1234', 'FirstName2', 'LastName2', 'kicko', 2, '2023-11-24 10:19:30.583', '2023-11-28 16:44:11.833'),
('e0e5ba28-19fc-4c65-8692-f61266608e3e', 'user@example.com', 'password123', 'FirstName1', 'LastName1', 'johnd', 1, '2023-11-24 10:19:30.583', '2023-11-28 16:44:11.833');

INSERT INTO `vote` (`id`, `from_user`, `to_user`, `vote`) VALUES
(4, 'e0e5ba28-19fc-4c65-8692-f61266608e3e', 'e0e5ba28-19fc-4c65-8692-f61266608d4n', 1),
(5, 'e0e5ba28-19fc-4c65-8692-f61266608d4n', 'e0e5ba28-19fc-4c65-8692-f61266608e3e', -1);



/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;