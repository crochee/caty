CREATE TABLE `user` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `account_id` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '主账号ID',
    `account` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '账号名',
    `user_id` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '用户ID',
    `pass_word` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '密码',
    `email` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '邮箱',
    `permission` longtext COLLATE utf8mb4_bin NOT NULL COMMENT '权限文本' CHECK (json_valid(`permission`)),
    `verify` tinyint(3) unsigned NOT NULL COMMENT '身份认证',
    `desc` text COLLATE utf8mb4_bin NOT NULL COMMENT '详细描述',
    `deleted` bigint(20) unsigned NOT NULL COMMENT '软删除标记',
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_id_deleted` (`user_id`,`deleted`),
    KEY `idx_user_deleted_at` (`deleted_at`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='账户信息表';
