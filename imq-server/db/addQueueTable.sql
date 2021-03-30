CREATE TABLE `Queue` (
  `queueId` varchar(45) NOT NULL,
  `topicId` varchar(45) DEFAULT NULL,
  `messageId` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`queueId`),
  KEY `queue_topic_idx` (`topicId`),
  KEY `queue_message_idx` (`messageId`),
  CONSTRAINT `queue_message` FOREIGN KEY (`messageId`) REFERENCES `Message` (`messageid`),
  CONSTRAINT `queue_topic` FOREIGN KEY (`topicId`) REFERENCES `Topic` (`topicid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;