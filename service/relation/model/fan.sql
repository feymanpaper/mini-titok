CREATE TABLE `fan` (
                       `id`   bigint(64) unsigned NOT NULL AUTO_INCREMENT,
                       `from_id` bigint(64) unsigned NOT NULL,
                       `to_id` bigint(64) unsigned NOT NULL,
                       `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                       `delete_time` timestamp NULL DEFAULT NULL,
                       PRIMARY KEY (`id`),
                       INDEX `idx_from` (`from_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4;

insert into titok_relation.fan (id, from_id, to_id) values (1, 4, 3);