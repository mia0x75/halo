-- 规则命中明细表。
--
-- 命中结果原先只以文本形式写进 mm_statements.report，无法按级别或规则聚合，
-- Dashboard 的「命中分级」「命中 Top 规则」没有数据来源。本表按工单、语句、
-- 规则保留每一条命中明细，聚合类统计由 mm_metrics 承担。

CREATE TABLE IF NOT EXISTS `mm_violations` (
  `violation_id` INT UNSIGNED
                  NOT NULL
                  AUTO_INCREMENT
                  COMMENT '自增主键',
  `uuid`         CHAR(36)
                  NOT NULL
                  COMMENT 'UUID',
  `ticket_id`    INT UNSIGNED
                  NOT NULL
                  COMMENT '对应工单',
  `sequence`     SMALLINT UNSIGNED
                  NOT NULL DEFAULT 0
                  COMMENT '语句序号',
  `rule_name`    VARCHAR(75)
                  NOT NULL
                  COMMENT '命中规则名',
  `level`        TINYINT UNSIGNED
                  NOT NULL
                  COMMENT '严重程度：0 提示 1 错误 2 警告',
  `message`      VARCHAR(150)
                  NOT NULL
                  COMMENT '命中提示',
  `create_at`    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
                  COMMENT '创建时间',

  PRIMARY KEY (`violation_id`),
  UNIQUE KEY `unique_1` (`uuid`),
  KEY `idx_ticket` (`ticket_id`),
  KEY `idx_rule` (`rule_name`),
  KEY `idx_level_create` (`level`, `create_at`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '规则命中明细表'
;
