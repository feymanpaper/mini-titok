CREATE DATABASE titok_user;
use titok_user;
CREATE TABLE `user` (
    `id`   bigint(64) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(32) NOT NULL DEFAULT '' COMMENT '用户名称',
    `avatar` varchar(128) NOT NULL DEFAULT '' COMMENT '用户头像',
    `background_image` varchar(128) NOT NULL DEFAULT '' COMMENT '用户个人页顶部大图',
    `signature` varchar(128) NOT NULL DEFAULT '' COMMENT '用户简介',
    `password` varchar(64)  NOT NULL DEFAULT '' COMMENT '用户密码',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_name` (`name`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4;

