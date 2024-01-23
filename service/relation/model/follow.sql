CREATE DATABASE titok_relation;
use titok_relation;
CREATE TABLE `follow` (
                        `id`   bigint(64) unsigned NOT NULL AUTO_INCREMENT,
                        `from_id` bigint(64) unsigned NOT NULL,
                        `to_id` bigint(64) unsigned NOT NULL,
                        `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                        `delete_time` timestamp NULL DEFAULT NULL,
                        PRIMARY KEY (`id`),
                        INDEX `idx_from` (`from_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4;

insert into titok_relation.follow (id, from_id, to_id) values (1, 3, 4);

select to_id from `follow` where from_id=3 and create_time < '2025-01-23 10:52:22' order by create_time desc limit 20

select `id`,`from_id`,`to_id`,`create_time`,`delete_time` from `follow` where from_id=3;

select * from follow;