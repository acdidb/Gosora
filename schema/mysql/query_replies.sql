CREATE TABLE `replies` (
	`rid` int not null AUTO_INCREMENT,
	`tid` int not null,
	`content` text not null,
	`parsed_content` text not null,
	`createdAt` datetime not null,
	`createdBy` int not null,
	`lastEdit` int DEFAULT 0 not null,
	`lastEditBy` int DEFAULT 0 not null,
	`lastUpdated` datetime not null,
	`ip` varchar(200) DEFAULT '' not null,
	`likeCount` int DEFAULT 0 not null,
	`attachCount` int DEFAULT 0 not null,
	`words` int DEFAULT 1 not null,
	`actionType` varchar(20) DEFAULT '' not null,
	`poll` int DEFAULT 0 not null,
	primary key(`rid`),
	fulltext key(`content`)
) CHARSET=utf8mb4 COLLATE utf8mb4_general_ci;