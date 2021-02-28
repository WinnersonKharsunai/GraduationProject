CREATE TABLE `Message` (
  `messageId` varchar(45) NOT NULL,
  `data` varchar(500) DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT NULL,
  `expiredAt` timestamp NULL DEFAULT NULL,
  `pubId` int(10) DEFAULT NULL,
  `topicId` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`messageId`),
  KEY `publisherID_idx` (`pubId`),
  KEY `topicID_idx` (`topicId`),
  CONSTRAINT `message_publisher` FOREIGN KEY (`pubId`) REFERENCES `Publisher` (`publisherid`),
  CONSTRAINT `message_topic` FOREIGN KEY (`topicId`) REFERENCES `Topic` (`topicid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
