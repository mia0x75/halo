-- 站内消息表。
--
-- mm_logs 是全局操作日志，没有收件箱语义与已读状态，顶栏的未读数没有数据来源。
-- 本表按接收人存放站内消息，支持未读计数与标记已读。

CREATE TABLE IF NOT EXISTS `mm_messages` (
  `message_id` INT UNSIGNED
                 NOT NULL
                 AUTO_INCREMENT
                 COMMENT '自增主键',
  `uuid`       CHAR(36)
                 NOT NULL
                 COMMENT 'UUID',
  `user_id`    INT UNSIGNED
                 NOT NULL
                 COMMENT '接收人',
  `category`   VARCHAR(30)
                 NOT NULL DEFAULT ''
                 COMMENT '消息分类',
  `title`      VARCHAR(75)
                 NOT NULL
                 COMMENT '标题',
  `content`    TEXT
                 NOT NULL
                 COMMENT '内容',
  `link`       VARCHAR(150)
                 NOT NULL DEFAULT ''
                 COMMENT '跳转链接',
  `read_at`    DATETIME NULL DEFAULT NULL
                 COMMENT '已读时间',
  `create_at`  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
                 COMMENT '创建时间',

  PRIMARY KEY (`message_id`),
  UNIQUE KEY `unique_1` (`uuid`),
  KEY `idx_user_read` (`user_id`, `read_at`),
  KEY `idx_user_create` (`user_id`, `create_at`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '站内消息表'
;
