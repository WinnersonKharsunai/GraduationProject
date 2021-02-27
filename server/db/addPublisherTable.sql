CREATE TABLE `Publisher` (
  `publisherId` int(11) NOT NULL,
  `topicId` varchar(30) DEFAULT NULL,
  PRIMARY KEY (`publisherId`),
  KEY `topicId_idx` (`topicId`),
  CONSTRAINT `publisher_topic` FOREIGN KEY (`topicId`) REFERENCES `Topic` (`topicid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;