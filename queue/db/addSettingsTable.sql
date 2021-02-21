-- queue: Settings table 

CREATE TABLE `Settings` (
  `topic` varchar(20) NOT NULL,
  `timeCount` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`topic`),
  UNIQUE KEY `topic_UNIQUE` (`topic`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4