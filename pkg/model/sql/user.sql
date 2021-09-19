CREATE TABLE `user`
(
    `id`         BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `account_id` VARCHAR(250)        NOT NULL COMMENT '主账号ID' COLLATE 'utf8mb4_bin',
    `user_id`    VARCHAR(250)        NOT NULL COMMENT '用户ID' COLLATE 'utf8mb4_bin',
    `nick`       VARCHAR(50)         NOT NULL COMMENT '昵称' COLLATE 'utf8mb4_bin',
    `pass_word`  VARCHAR(50)         NOT NULL COMMENT '密码' COLLATE 'utf8mb4_bin',
    `email`      VARCHAR(50)         NOT NULL COMMENT '邮箱' COLLATE 'utf8mb4_bin',
    `permission` TEXT                NOT NULL COMMENT '权限文本' COLLATE 'utf8mb4_bin',
    `verify`     TINYINT(3)          NOT NULL COMMENT '身份认证',
    `deleted`    BIGINT(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '软删除标记',
    `created_at` DATETIME            NOT NULL DEFAULT current_timestamp() COMMENT '创建时间',
    `updated_at` DATETIME            NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '更新时间',
    `deleted_at` DATETIME            NULL     DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `idx_user_id_deleted` (`user_id`, `deleted`) USING BTREE,
    INDEX `idx_deleted_at` (`deleted_at`) USING BTREE
)
    COMMENT ='账户信息表'
    COLLATE = 'utf8mb4_bin'
    ENGINE = InnoDB
;
