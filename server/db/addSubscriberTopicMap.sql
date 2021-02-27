CREATE TABLE `SubscriberTopicMap` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `subscriberId` int(11) NOT NULL,
  `topicId` varchar(30) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `subID_idx` (`subscriberId`),
  KEY `topicID_idx` (`topicId`),
  CONSTRAINT `submap_subscriber` FOREIGN KEY (`subscriberId`) REFERENCES `Subscriber` (`subscriberid`),
  CONSTRAINT `topic_ID` FOREIGN KEY (`topicId`) REFERENCES `Topic` (`topicid`)
) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4;