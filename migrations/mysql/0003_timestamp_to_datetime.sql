-- 时间列由 UNIX 秒（INT UNSIGNED）改为 DATETIME，统一按 UTC 存取。
-- 存量值用 FROM_UNIXTIME 转换，故脚本先将会话时区置为 UTC。
--
-- 注意：MySQL 的 DDL 会隐式提交，无法整体回滚；执行前必须备份，
-- 且应用侧需暂停写入（见 AGENTS「Schema 迁移」）。
--
-- 本变更已包含在基线迁移 0001_init.sql 中（项目尚未发布，没有任何数据库
-- 执行过本文件），因此迁移器只登记台账、不执行语句。
-- migrate:folded-into-baseline

SET time_zone = '+00:00';

-- mm_avatars
ALTER TABLE `mm_avatars`
    ADD COLUMN `update_at_tmp` DATETIME NULL AFTER `version`,
    ADD COLUMN `create_at_tmp` DATETIME NULL AFTER `update_at_tmp`;
UPDATE `mm_avatars` SET `update_at_tmp` = FROM_UNIXTIME(NULLIF(`update_at`, 0)), `create_at_tmp` = COALESCE(FROM_UNIXTIME(NULLIF(`create_at`, 0)), UTC_TIMESTAMP());
ALTER TABLE `mm_avatars`
    DROP COLUMN `update_at`,
    DROP COLUMN `create_at`,
    CHANGE COLUMN `update_at_tmp` `update_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    CHANGE COLUMN `create_at_tmp` `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';

-- mm_comments
ALTER TABLE `mm_comments`
    ADD COLUMN `update_at_tmp` DATETIME NULL AFTER `version`,
    ADD COLUMN `create_at_tmp` DATETIME NULL AFTER `update_at_tmp`;
UPDATE `mm_comments` SET `update_at_tmp` = FROM_UNIXTIME(NULLIF(`update_at`, 0)), `create_at_tmp` = COALESCE(FROM_UNIXTIME(NULLIF(`create_at`, 0)), UTC_TIMESTAMP());
ALTER TABLE `mm_comments`
    DROP COLUMN `update_at`,
    DROP COLUMN `create_at`,
    CHANGE COLUMN `update_at_tmp` `update_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    CHANGE COLUMN `create_at_tmp` `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';

-- mm_crons
ALTER TABLE `mm_crons`
    ADD COLUMN `update_at_tmp` DATETIME NULL AFTER `version`,
    ADD COLUMN `create_at_tmp` DATETIME NULL AFTER `update_at_tmp`;
UPDATE `mm_crons` SET `update_at_tmp` = FROM_UNIXTIME(NULLIF(`update_at`, 0)), `create_at_tmp` = COALESCE(FROM_UNIXTIME(NULLIF(`create_at`, 0)), UTC_TIMESTAMP());
ALTER TABLE `mm_crons`
    DROP COLUMN `update_at`,
    DROP COLUMN `create_at`,
    CHANGE COLUMN `update_at_tmp` `update_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    CHANGE COLUMN `create_at_tmp` `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';

-- mm_glossaries
ALTER TABLE `mm_glossaries`
    ADD COLUMN `update_at_tmp` DATETIME NULL AFTER `version`,
    ADD COLUMN `create_at_tmp` DATETIME NULL AFTER `update_at_tmp`;
UPDATE `mm_glossaries` SET `update_at_tmp` = FROM_UNIXTIME(NULLIF(`update_at`, 0)), `create_at_tmp` = COALESCE(FROM_UNIXTIME(NULLIF(`create_at`, 0)), UTC_TIMESTAMP());
ALTER TABLE `mm_glossaries`
    DROP COLUMN `update_at`,
    DROP COLUMN `create_at`,
    CHANGE COLUMN `update_at_tmp` `update_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    CHANGE COLUMN `create_at_tmp` `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';

-- mm_clusters
ALTER TABLE `mm_clusters`
    ADD COLUMN `update_at_tmp` DATETIME NULL AFTER `version`,
    ADD COLUMN `create_at_tmp` DATETIME NULL AFTER `update_at_tmp`;
UPDATE `mm_clusters` SET `update_at_tmp` = FROM_UNIXTIME(NULLIF(`update_at`, 0)), `create_at_tmp` = COALESCE(FROM_UNIXTIME(NULLIF(`create_at`, 0)), UTC_TIMESTAMP());
ALTER TABLE `mm_clusters`
    DROP COLUMN `update_at`,
    DROP COLUMN `create_at`,
    CHANGE COLUMN `update_at_tmp` `update_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    CHANGE COLUMN `create_at_tmp` `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';

-- mm_logs（无 update_at）
ALTER TABLE `mm_logs`
    ADD COLUMN `create_at_tmp` DATETIME NULL AFTER `version`;
UPDATE `mm_logs` SET `create_at_tmp` = COALESCE(FROM_UNIXTIME(NULLIF(`create_at`, 0)), UTC_TIMESTAMP());
ALTER TABLE `mm_logs`
    DROP COLUMN `create_at`,
    CHANGE COLUMN `create_at_tmp` `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';

-- mm_options
ALTER TABLE `mm_options`
    ADD COLUMN `update_at_tmp` DATETIME NULL AFTER `version`,
    ADD COLUMN `create_at_tmp` DATETIME NULL AFTER `update_at_tmp`;
UPDATE `mm_options` SET `update_at_tmp` = FROM_UNIXTIME(NULLIF(`update_at`, 0)), `create_at_tmp` = COALESCE(FROM_UNIXTIME(NULLIF(`create_at`, 0)), UTC_TIMESTAMP());
ALTER TABLE `mm_options`
    DROP COLUMN `update_at`,
    DROP COLUMN `create_at`,
    CHANGE COLUMN `update_at_tmp` `update_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    CHANGE COLUMN `create_at_tmp` `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';

-- mm_queries
ALTER TABLE `mm_queries`
    ADD COLUMN `update_at_tmp` DATETIME NULL AFTER `version`,
    ADD COLUMN `create_at_tmp` DATETIME NULL AFTER `update_at_tmp`;
UPDATE `mm_queries` SET `update_at_tmp` = FROM_UNIXTIME(NULLIF(`update_at`, 0)), `create_at_tmp` = COALESCE(FROM_UNIXTIME(NULLIF(`create_at`, 0)), UTC_TIMESTAMP());
ALTER TABLE `mm_queries`
    DROP COLUMN `update_at`,
    DROP COLUMN `create_at`,
    CHANGE COLUMN `update_at_tmp` `update_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    CHANGE COLUMN `create_at_tmp` `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';

-- mm_edges
ALTER TABLE `mm_edges`
    ADD COLUMN `update_at_tmp` DATETIME NULL AFTER `version`,
    ADD COLUMN `create_at_tmp` DATETIME NULL AFTER `update_at_tmp`;
UPDATE `mm_edges` SET `update_at_tmp` = FROM_UNIXTIME(NULLIF(`update_at`, 0)), `create_at_tmp` = COALESCE(FROM_UNIXTIME(NULLIF(`create_at`, 0)), UTC_TIMESTAMP());
ALTER TABLE `mm_edges`
    DROP COLUMN `update_at`,
    DROP COLUMN `create_at`,
    CHANGE COLUMN `update_at_tmp` `update_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    CHANGE COLUMN `create_at_tmp` `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';

-- mm_roles
ALTER TABLE `mm_roles`
    ADD COLUMN `update_at_tmp` DATETIME NULL AFTER `version`,
    ADD COLUMN `create_at_tmp` DATETIME NULL AFTER `update_at_tmp`;
UPDATE `mm_roles` SET `update_at_tmp` = FROM_UNIXTIME(NULLIF(`update_at`, 0)), `create_at_tmp` = COALESCE(FROM_UNIXTIME(NULLIF(`create_at`, 0)), UTC_TIMESTAMP());
ALTER TABLE `mm_roles`
    DROP COLUMN `update_at`,
    DROP COLUMN `create_at`,
    CHANGE COLUMN `update_at_tmp` `update_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    CHANGE COLUMN `create_at_tmp` `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';

-- mm_rules
ALTER TABLE `mm_rules`
    ADD COLUMN `update_at_tmp` DATETIME NULL AFTER `version`,
    ADD COLUMN `create_at_tmp` DATETIME NULL AFTER `update_at_tmp`;
UPDATE `mm_rules` SET `update_at_tmp` = FROM_UNIXTIME(NULLIF(`update_at`, 0)), `create_at_tmp` = COALESCE(FROM_UNIXTIME(NULLIF(`create_at`, 0)), UTC_TIMESTAMP());
ALTER TABLE `mm_rules`
    DROP COLUMN `update_at`,
    DROP COLUMN `create_at`,
    CHANGE COLUMN `update_at_tmp` `update_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    CHANGE COLUMN `create_at_tmp` `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';

-- mm_statements
ALTER TABLE `mm_statements`
    ADD COLUMN `update_at_tmp` DATETIME NULL AFTER `version`,
    ADD COLUMN `create_at_tmp` DATETIME NULL AFTER `update_at_tmp`;
UPDATE `mm_statements` SET `update_at_tmp` = FROM_UNIXTIME(NULLIF(`update_at`, 0)), `create_at_tmp` = COALESCE(FROM_UNIXTIME(NULLIF(`create_at`, 0)), UTC_TIMESTAMP());
ALTER TABLE `mm_statements`
    DROP COLUMN `update_at`,
    DROP COLUMN `create_at`,
    CHANGE COLUMN `update_at_tmp` `update_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    CHANGE COLUMN `create_at_tmp` `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';

-- mm_statistics
ALTER TABLE `mm_statistics`
    ADD COLUMN `update_at_tmp` DATETIME NULL AFTER `version`,
    ADD COLUMN `create_at_tmp` DATETIME NULL AFTER `update_at_tmp`;
UPDATE `mm_statistics` SET `update_at_tmp` = FROM_UNIXTIME(NULLIF(`update_at`, 0)), `create_at_tmp` = COALESCE(FROM_UNIXTIME(NULLIF(`create_at`, 0)), UTC_TIMESTAMP());
ALTER TABLE `mm_statistics`
    DROP COLUMN `update_at`,
    DROP COLUMN `create_at`,
    CHANGE COLUMN `update_at_tmp` `update_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    CHANGE COLUMN `create_at_tmp` `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';

-- mm_tickets
ALTER TABLE `mm_tickets`
    ADD COLUMN `update_at_tmp` DATETIME NULL AFTER `version`,
    ADD COLUMN `create_at_tmp` DATETIME NULL AFTER `update_at_tmp`;
UPDATE `mm_tickets` SET `update_at_tmp` = FROM_UNIXTIME(NULLIF(`update_at`, 0)), `create_at_tmp` = COALESCE(FROM_UNIXTIME(NULLIF(`create_at`, 0)), UTC_TIMESTAMP());
ALTER TABLE `mm_tickets`
    DROP COLUMN `update_at`,
    DROP COLUMN `create_at`,
    CHANGE COLUMN `update_at_tmp` `update_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    CHANGE COLUMN `create_at_tmp` `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';

-- mm_users
ALTER TABLE `mm_users`
    ADD COLUMN `update_at_tmp` DATETIME NULL AFTER `version`,
    ADD COLUMN `create_at_tmp` DATETIME NULL AFTER `update_at_tmp`;
UPDATE `mm_users` SET `update_at_tmp` = FROM_UNIXTIME(NULLIF(`update_at`, 0)), `create_at_tmp` = COALESCE(FROM_UNIXTIME(NULLIF(`create_at`, 0)), UTC_TIMESTAMP());
ALTER TABLE `mm_users`
    DROP COLUMN `update_at`,
    DROP COLUMN `create_at`,
    CHANGE COLUMN `update_at_tmp` `update_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    CHANGE COLUMN `create_at_tmp` `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';

-- mm_templates
ALTER TABLE `mm_templates`
    ADD COLUMN `update_at_tmp` DATETIME NULL AFTER `version`,
    ADD COLUMN `create_at_tmp` DATETIME NULL AFTER `update_at_tmp`;
UPDATE `mm_templates` SET `update_at_tmp` = FROM_UNIXTIME(NULLIF(`update_at`, 0)), `create_at_tmp` = COALESCE(FROM_UNIXTIME(NULLIF(`create_at`, 0)), UTC_TIMESTAMP());
ALTER TABLE `mm_templates`
    DROP COLUMN `update_at`,
    DROP COLUMN `create_at`,
    CHANGE COLUMN `update_at_tmp` `update_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    CHANGE COLUMN `create_at_tmp` `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';
