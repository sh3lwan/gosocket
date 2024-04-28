CREATE TABLE IF NOT EXISTS `messages` (
  `id` int  NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `message` text NOT NULL,
  `username` varchar(191) NOT NULL,
  `is_new` tinyint(1) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
