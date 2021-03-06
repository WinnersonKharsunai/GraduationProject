CREATE TABLE `DLQ` (
  `dlqId` varchar(45) NOT NULL,
  `topicId` varchar(45) DEFAULT NULL,
  `messageId` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`dlqId`),
  KEY `dlq_message_idx` (`messageId`),
  KEY `dlq_topic_idx` (`topicId`),
  CONSTRAINT `dlq_message` FOREIGN KEY (`messageId`) REFERENCES `Message` (`messageid`),
  CONSTRAINT `dlq_topic` FOREIGN KEY (`topicId`) REFERENCES `Topic` (`topicid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;