-- --------------------------------------------------------
-- 主机:                           127.0.0.1
-- 服务器版本:                        10.6.4-MariaDB-1:10.6.4+maria~focal - mariadb.org binary distribution
-- 服务器操作系统:                      debian-linux-gnu
-- HeidiSQL 版本:                  11.3.0.6295
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

-- 导出  表 obs.user 结构
DROP TABLE IF EXISTS `user`;
CREATE TABLE IF NOT EXISTS `user` (
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
) ENGINE=InnoDB AUTO_INCREMENT=373746917676511096 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='账户信息表';

-- 正在导出表  obs.user 的数据：~1 rows (大约)
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` (`id`, `account_id`, `account`, `user_id`, `pass_word`, `email`, `permission`, `verify`, `desc`, `deleted`, `created_at`, `updated_at`, `deleted_at`) VALUES
	(373746917676511095, '373746917676445559', 'crochee', '373746917676445559', '123456', 'crochee@139.com', '{"*":4}', 0, '{"detail":"some unknown"}', 0, '2021-09-22 08:42:02.647', '2021-09-22 08:42:02.647', NULL);
/*!40000 ALTER TABLE `user` ENABLE KEYS */;

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
