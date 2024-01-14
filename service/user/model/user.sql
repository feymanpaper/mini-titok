CREATE DATABASE titok_user;
use titok_user;
CREATE TABLE `user` (
    `id`   bigint(64) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(32) NOT NULL DEFAULT '' COMMENT '用户名称',
    `follow_count` bigint(64) unsigned NOT NULL DEFAULT '0' COMMENT '用户关注的人个数',
    `follower_count` bigint(64) unsigned NOT NULL DEFAULT '0' COMMENT '用户的粉丝个数',
    `avatar` varchar(128) NOT NULL DEFAULT '' COMMENT '用户头像',
    `background_image` varchar(128) NOT NULL DEFAULT '' COMMENT '用户个人页顶部大图',
    `signature` varchar(128) NOT NULL DEFAULT '' COMMENT '用户简介',
    `total_favorited` bigint(64) unsigned NOT NULL DEFAULT '0' COMMENT '用户获赞数量',
    `work_count` bigint(64) unsigned NOT NULL DEFAULT '0' COMMENT '用户获赞数量',
    `favorite_count` bigint(64) unsigned NOT NULL DEFAULT '0' COMMENT '用户点赞数量',
    `password` varchar(64)  NOT NULL DEFAULT '' COMMENT '用户密码',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_name` (`name`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4;

