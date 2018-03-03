DROP TABLE IF EXISTS `abuses`;

CREATE TABLE `abuses` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `createrId` int(11) NOT NULL,
  `abuserId` int(11) NOT NULL,
  `comment` varchar(45) DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `resolved` enum('true','false') DEFAULT 'false',
  `resolvedAt` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
