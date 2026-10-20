-- 异步任务执行表。
--
-- 定时任务（mm_crons）存的是调度定义，internal/worker 的任务实例只存在内存里，
-- 重启即丢，也无法追查执行历史。本表记录每一次执行：状态、重试次数、起止时间
-- 与失败原因，是「执行必经审计」的落点。

CREATE TABLE IF NOT EXISTS `mm_tasks` (
  `task_id`      INT UNSIGNED
                 NOT NULL
                 AUTO_INCREMENT
                 COMMENT '自增主键',
  `uuid`         CHAR(36)
                 NOT NULL
                 COMMENT 'UUID',
  `name`         VARCHAR(75)
                 NOT NULL
                 COMMENT '任务名称',
  `type`         VARCHAR(30)
                 NOT NULL DEFAULT ''
                 COMMENT '任务类型',
  `payload`      TEXT
                 NOT NULL
                 COMMENT '任务参数（JSON）',
  `status`       TINYINT UNSIGNED
                 NOT NULL DEFAULT 0
                 COMMENT '状态：0 待执行 1 执行中 2 成功 3 失败 4 取消',
  `priority`     TINYINT UNSIGNED
                 NOT NULL DEFAULT 0
                 COMMENT '优先级，数值越小越优先',
  `attempt`      TINYINT UNSIGNED
                 NOT NULL DEFAULT 0
                 COMMENT '已尝试次数',
  `max_attempt`  TINYINT UNSIGNED
                 NOT NULL DEFAULT 3
                 COMMENT '最大尝试次数',
  `cron_id`      INT UNSIGNED
                 NOT NULL DEFAULT 0
                 COMMENT '来源定时任务，0 表示无',
  `scheduled_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
                 COMMENT '计划执行时间',
  `started_at`   DATETIME NULL DEFAULT NULL
                 COMMENT '开始执行时间',
  `finished_at`  DATETIME NULL DEFAULT NULL
                 COMMENT '结束执行时间',
  `error`        TEXT
                 NOT NULL
                 COMMENT '失败原因',
  `version`      INT UNSIGNED
                 NOT NULL DEFAULT 0
                 COMMENT '版本',
  `update_at`    DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
                 COMMENT '修改时间',
  `create_at`    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
                 COMMENT '创建时间',

  PRIMARY KEY (`task_id`),
  UNIQUE KEY `unique_1` (`uuid`),
  KEY `idx_status_scheduled` (`status`, `scheduled_at`),
  KEY `idx_type` (`type`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '异步任务执行表'
;

-- 定时任务补齐审计维度：谁配置的、跑在哪个群集。
ALTER TABLE `mm_crons`
    ADD COLUMN `user_id` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建人',
    ADD COLUMN `cluster_id` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '目标群集，0 表示不限定';
