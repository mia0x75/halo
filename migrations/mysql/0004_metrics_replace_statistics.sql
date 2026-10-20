-- 指标表取代统计表。
--
-- 旧表用「分组 + 键」表达指标，counter 与 gauge 混在一张表里：所有写入都走
-- value = value + delta，导致 latest-login 这类 gauge 指标被累加。新表把维度
-- 拆成「名称 + 作用域 + 作用域实体 + 时间桶」，并显式区分 counter 与 gauge。
--
-- 本脚本可重复执行：目标库若已按权威 schema 建好（mm_metrics 存在、
-- mm_statistics 已删除），数据搬迁会跳过，只补台账。

CREATE TABLE IF NOT EXISTS `mm_metrics` (
  `metric_id`   INT UNSIGNED
                NOT NULL
                AUTO_INCREMENT
                COMMENT '自增主键',
  `uuid`        CHAR(36)
                NOT NULL
                COMMENT 'UUID',
  `name`        VARCHAR(75)
                NOT NULL
                COMMENT '指标名称',
  `scope`       TINYINT UNSIGNED
                NOT NULL DEFAULT 0
                COMMENT '作用域：0 全局 1 按天 2 用户 3 群集 4 分组',
  `scope_id`    CHAR(36)
                NOT NULL DEFAULT ''
                COMMENT '作用域实体 UUID，全局作用域为空串',
  `type`        TINYINT UNSIGNED
                NOT NULL DEFAULT 0
                COMMENT '类型：0 counter 累加 1 gauge 覆盖',
  `bucket`      VARCHAR(20)
                NOT NULL DEFAULT ''
                COMMENT '时间桶，如 2026-10-20',
  `value`       DECIMAL(18,4)
                NOT NULL DEFAULT 0
                COMMENT '指标值',
  `version`     INT UNSIGNED
                NOT NULL DEFAULT 0
                COMMENT '版本',
  `update_at`   DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
                COMMENT '修改时间',
  `create_at`   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
                COMMENT '创建时间',

  PRIMARY KEY (`metric_id`),
  UNIQUE KEY `unique_1` (`uuid`),
  UNIQUE KEY `unique_2` (`name`, `scope`, `scope_id`, `bucket`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '指标表'
;

-- 旧数据搬迁：分组决定作用域与指标名，键在按天分组里是时间桶。
SET @halo_move_metrics = (
    SELECT IF(
        EXISTS(SELECT 1 FROM information_schema.tables
               WHERE table_schema = DATABASE() AND table_name = 'mm_statistics'),
        'INSERT INTO `mm_metrics` (`uuid`, `name`, `scope`, `scope_id`, `type`, `bucket`, `value`, `version`, `update_at`, `create_at`)
         SELECT `uuid`,
                CASE WHEN `group` = ''tickets-daily'' THEN ''total-tickets''
                     WHEN `group` = ''queries-daily'' THEN ''total-queries''
                     ELSE `key` END,
                CASE WHEN `group` IN (''tickets-daily'', ''queries-daily'') THEN 1
                     WHEN `group` REGEXP ''^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$'' THEN 2
                     ELSE 0 END,
                CASE WHEN `group` REGEXP ''^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$'' THEN `group`
                     ELSE '''' END,
                CASE WHEN `key` = ''latest-login'' THEN 1 ELSE 0 END,
                CASE WHEN `group` IN (''tickets-daily'', ''queries-daily'') THEN `key`
                     ELSE '''' END,
                `value`, `version`, `update_at`, `create_at`
         FROM `mm_statistics`
         ON DUPLICATE KEY UPDATE `value` = `mm_metrics`.`value`',
        'SELECT 1')
);
PREPARE halo_move_metrics_stmt FROM @halo_move_metrics;
EXECUTE halo_move_metrics_stmt;
DEALLOCATE PREPARE halo_move_metrics_stmt;

DROP TABLE IF EXISTS `mm_statistics`;
