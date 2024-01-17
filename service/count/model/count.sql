CREATE DATABASE titok_count;
use titok_count;
CREATE TABLE `count` (
	`id`   bigint(64) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
	`count_key` varchar(128) NOT NULL DEFAULT '' COMMENT 'count的类型',
	`count_val` bigint(64) unsigned NOT NULL DEFAULT '0' COMMENT '计数值',
	`create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	`update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`),
	UNIQUE INDEX unique_countkey (count_key)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4;

insert into titok_count.count (id, count_key, count_val) values (1, 'follower:1', 22);