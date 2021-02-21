-- queue: Publisher table

CREATE TABLE `Publisher` (
  `publisherId` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `createdAt` timestamp NOT NULL,
  `topicName` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`publisherId`),
  KEY `name_idx` (`topicName`),
  CONSTRAINT `pubTopic` FOREIGN KEY (`topicName`) REFERENCES `Settings` (`topic`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4