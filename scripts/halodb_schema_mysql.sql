DROP TABLE IF EXISTS `mm_avatars`;
CREATE TABLE `mm_avatars` (
   `avatar_id` INT UNSIGNED
               NOT NULL
               AUTO_INCREMENT
               COMMENT '自增主键',
   `uuid`      CHAR(36)
               NOT NULL
               COMMENT 'UUID',
   `url`       VARCHAR(75)
               NOT NULL
               COMMENT '头像位置',
   `version`   INT UNSIGNED
               NOT NULL
               COMMENT '版本',
   `update_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
               COMMENT '修改时间',
   `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
               COMMENT '创建时间',

   PRIMARY KEY (`avatar_id`),
   UNIQUE KEY `unique_1` (`uuid`),
   UNIQUE KEY `unique_2` (`url`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '头像表'
;

DROP TABLE IF EXISTS `mm_comments`;
CREATE TABLE `mm_comments` (
  `comment_id` INT UNSIGNED
               NOT NULL
               AUTO_INCREMENT
               COMMENT '自增主键',
  `uuid`       CHAR(36)
               NOT NULL
               COMMENT 'UUID',
  `content`    TINYTEXT
               NOT NULL
               COMMENT '意见建议',
  `ticket_id`  INT UNSIGNED
               NOT NULL
               COMMENT '对应工单',
  `user_id`    INT UNSIGNED
               NOT NULL
               COMMENT '发起人',
  `version`    INT UNSIGNED
               NOT NULL
               COMMENT '版本',
  `update_at`  DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
               COMMENT '修改时间',
  `create_at`  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
               COMMENT '创建时间',

  PRIMARY KEY (`comment_id`),
  UNIQUE KEY `unique_1` (`uuid`),
  KEY `index_1` (`ticket_id`),
  KEY `index_2` (`user_id`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '审核意见表'
;

DROP TABLE IF EXISTS `mm_crons`;
CREATE TABLE `mm_crons` (
  `cron_id`   INT UNSIGNED
              NOT NULL
              AUTO_INCREMENT
              COMMENT '自增主键',
  `uuid`      CHAR(36)
              NOT NULL
              COMMENT 'UUID',
  `status`    CHAR(1)
              NOT NULL
              COMMENT '执行状态',
  `name`      VARCHAR(75)
              NOT NULL
              COMMENT '任务名称',
  `cmd`       VARCHAR(75)
              NOT NULL
              COMMENT '函数名称',
  `params`    VARCHAR(150)
              NOT NULL
              COMMENT '运行参数',
  `last_run`  CHAR(25)
              COMMENT '上一次运行时间',
  `next_run`  CHAR(25)
              COMMENT '下一次运行时间',
  `interval`  VARCHAR(20)
              COMMENT '执行间隔',
  `duration`  VARCHAR(20)
              COMMENT '执行耗时',
  `user_id`   INT UNSIGNED
              NOT NULL DEFAULT 0
              COMMENT '创建人',
  `cluster_id` INT UNSIGNED
              NOT NULL DEFAULT 0
              COMMENT '目标群集，0 表示不限定',
  `recurrent` TINYINT UNSIGNED
              COMMENT '是否周期运行',
  `hash`      VARCHAR(60)
              COMMENT '哈希值',
  `version`   INT UNSIGNED
              NOT NULL
              COMMENT '版本',
  `update_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
              COMMENT '修改时间',
  `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
              COMMENT '创建时间',

  PRIMARY KEY (`cron_id`),
  UNIQUE KEY `unique_1` (`uuid`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '计划任务表'
;

DROP TABLE IF EXISTS `mm_glossaries`;
CREATE TABLE `mm_glossaries` (
  `group`       VARCHAR(25)
                NOT NULL
                COMMENT '分组',
  `key`         TINYINT UNSIGNED
                NOT NULL
                COMMENT '键',
  `value`       VARCHAR(50)
                NOT NULL
                COMMENT '值',
  `uuid`        CHAR(36)
                NOT NULL
                COMMENT 'UUID',
  `description` VARCHAR(150)
                NOT NULL
                COMMENT '值描述',
  `version`     INT UNSIGNED
                NOT NULL
                COMMENT '版本',
  `update_at`   DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
                COMMENT '修改时间',
  `create_at`   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
                COMMENT '创建时间',

  PRIMARY KEY (`group`,`key`),
  UNIQUE KEY `unique_1` (`group`,`value`),
  UNIQUE KEY `unique_2` (`uuid`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '字典表'
;

DROP TABLE IF EXISTS `mm_clusters`;
CREATE TABLE `mm_clusters` (
  `cluster_id`  INT UNSIGNED
                NOT NULL
                AUTO_INCREMENT
                COMMENT '自增主键',
  `uuid`        CHAR(36)
                NOT NULL
                COMMENT 'UUID',
  `host`        VARCHAR(150)
                NOT NULL
                COMMENT '主机名称',
  `alias`       VARCHAR(75)
                NOT NULL
                COMMENT '主机别名',
  `ip`          VARCHAR(15)
                NOT NULL
                COMMENT '主机地址',
  `port`        INT UNSIGNED
                NOT NULL
                DEFAULT 3306
                COMMENT '端口',
  `user`        VARCHAR(50)
                NOT NULL
                COMMENT '连接用户',
  `type`        VARCHAR(20)
                NOT NULL
                DEFAULT 'mysql'
                COMMENT '数据库类型',
  `group_id`    INT UNSIGNED
                NOT NULL
                DEFAULT 0
                COMMENT '所属分组，0 表示未分组',
  `password`    VARBINARY(48)
                NOT NULL
                COMMENT '密码',
  `fingerprint` VARBINARY(32)
                NOT NULL
                COMMENT '指纹',
  `status`      TINYINT UNSIGNED
                NOT NULL
                DEFAULT 1
                COMMENT '状态',
  `version`     INT UNSIGNED
                NOT NULL
                COMMENT '版本',
  `update_at`   DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
                COMMENT '修改时间',
  `create_at`   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
                COMMENT '创建时间',

  PRIMARY KEY (`cluster_id`),
  UNIQUE KEY `unique_1` (`uuid`),
  UNIQUE KEY `unique_2` (`alias`),
  UNIQUE KEY `unique_3` (`host`,`port`),
  UNIQUE KEY `unique_4` (`ip`,`port`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '群集表'
;

DROP TABLE IF EXISTS `mm_logs`;
CREATE TABLE `mm_logs` (
  `log_id`    INT UNSIGNED
              NOT NULL
              AUTO_INCREMENT
              COMMENT '自增主键',
  `uuid`      CHAR(36)
              NOT NULL
              COMMENT 'UUID',
  `user_id`   INT UNSIGNED
              NOT NULL
              DEFAULT 0
              COMMENT '操作员',
  `operation` TINYTEXT
              NOT NULL
              COMMENT '操作记录',
  `version`   INT UNSIGNED
              NOT NULL
              COMMENT '版本',
  `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
              COMMENT '创建时间',

  PRIMARY KEY (`log_id`),
  UNIQUE KEY `unique_1` (`uuid`),
  KEY `index_1` (`user_id`,`create_at`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '日志表'
;

DROP TABLE IF EXISTS `mm_options`;
CREATE TABLE `mm_options` (
  `name`        VARCHAR(50)
                NOT NULL
                COMMENT '配置项',
  `uuid`        CHAR(36)
                NOT NULL
                COMMENT 'UUID',
  `parameter`   TINYTEXT
                NOT NULL
                COMMENT '配置值',
  `description` VARCHAR(75)
                NOT NULL
                COMMENT '描述',
  `element`     VARCHAR(15)
                NOT NULL
                DEFAULT '-'
                COMMENT '展现组件类型',
  `writable`    TINYINT UNSIGNED
                NOT NULL
                COMMENT '是否可写',
  `version`     INT UNSIGNED
                NOT NULL
                COMMENT '版本',
  `update_at`   DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
                COMMENT '修改时间',
  `create_at`   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
                COMMENT '创建时间',
  `group_name`  VARCHAR(50)
                NOT NULL
                DEFAULT ''
                COMMENT '所属分组，同组共用一个标题',
  `sort_order`  INT
                NOT NULL
                DEFAULT 0
                COMMENT '展示顺序：组内按此排序，跨组时也决定分组的先后',
  `depends_on`  VARCHAR(50)
                NOT NULL
                DEFAULT ''
                COMMENT '依赖的选项名，为空表示无依赖',
  `min_value`   INT NULL
                COMMENT '数值下界',
  `max_value`   INT NULL
                COMMENT '数值上界',

  PRIMARY KEY (`name`),
  UNIQUE KEY `unique_1` (`uuid`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '系统选项表'
;

DROP TABLE IF EXISTS `mm_queries`;
CREATE TABLE `mm_queries` (
  `query_id`   INT UNSIGNED
               NOT NULL
               AUTO_INCREMENT
               COMMENT '主键',
  `uuid`       CHAR(36)
               NOT NULL
               COMMENT 'UUID',
  `type`       TINYINT UNSIGNED
               NOT NULL
               COMMENT '查询类型',
  `cluster_id` INT UNSIGNED
               NOT NULL
               COMMENT '目标群集',
  `database`   VARCHAR(75)
               NOT NULL
               COMMENT '目标库',
  `content`    TEXT
               NOT NULL
               COMMENT '执行语句',
  `plan`       TEXT
               NOT NULL
               COMMENT '执行计划',
  `user_id`    INT UNSIGNED
               NOT NULL
               COMMENT '发起人',
  `version`    INT UNSIGNED
               NOT NULL
               COMMENT '版本',
  `update_at`  DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
               COMMENT '修改时间',
  `create_at`  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
               COMMENT '创建时间',

  PRIMARY KEY (`query_id`),
  UNIQUE KEY `unique_1` (`uuid`),
  KEY `index_1` (`user_id`),
  KEY `index_2` (`cluster_id`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '查询表'
;

DROP TABLE IF EXISTS `mm_edges`;
CREATE TABLE `mm_edges` (
  `edge_id`       INT UNSIGNED
                  NOT NULL
                  AUTO_INCREMENT
                  COMMENT '自增主键',
  `uuid`          CHAR(36)
                  NOT NULL
                  COMMENT 'UUID',
  `type`          INT UNSIGNED
                  NOT NULL
                  COMMENT '分类标识',
  `ancestor_id`   INT UNSIGNED
                  NOT NULL
                  COMMENT '先代',
  `descendant_id` INT UNSIGNED
                  NOT NULL
                  COMMENT '后代',
  `version`       INT UNSIGNED
                  NOT NULL
                  COMMENT '版本',
  `update_at`     DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
                  COMMENT '修改时间',
  `create_at`     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
                  COMMENT '创建时间',

  PRIMARY KEY (`edge_id`),
  UNIQUE KEY `unique_1` (`uuid`),
  UNIQUE KEY `unique_2` (`type`,`ancestor_id`,`descendant_id`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '多对多关联表'
;

DROP TABLE IF EXISTS `mm_roles`;
CREATE TABLE `mm_roles` (
  `role_id`     INT UNSIGNED
                NOT NULL
                AUTO_INCREMENT
                COMMENT '自增主键',
  `uuid`        CHAR(36)
                NOT NULL
                COMMENT 'UUID',
  `name`        VARCHAR(25)
                NOT NULL
                COMMENT '角色名称',
  `description` VARCHAR(75)
                NOT NULL
                COMMENT '描述',
  `version`     INT UNSIGNED
                NOT NULL
                COMMENT '版本',
  `update_at`   DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
                COMMENT '修改时间',
  `create_at`   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
                COMMENT '创建时间',

  PRIMARY KEY (`role_id`),
  UNIQUE KEY `unique_1` (`uuid`),
  UNIQUE KEY `unique_2` (`name`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '角色表'
;

DROP TABLE IF EXISTS `mm_rules`;
CREATE TABLE `mm_rules` (
  `name`        CHAR(10)
                NOT NULL
                COMMENT '规则名称',
  `uuid`        CHAR(36)
                NOT NULL
                COMMENT 'UUID',
  `group`       TINYINT UNSIGNED
                NOT NULL
                COMMENT '规则分组',
  `description` VARCHAR(75)
                NOT NULL
                COMMENT '规则描述',
  `level`       TINYINT UNSIGNED
                NOT NULL
                COMMENT '严重级别',
  `vldr_group`  TINYINT UNSIGNED
                NOT NULL
                COMMENT '审核分组',
  `operator`    VARCHAR(10)
                NOT NULL
                COMMENT '比较符',
  `parameters`  VARCHAR(150)
                NOT NULL
                COMMENT '有效值',
  `bitwise`     TINYINT UNSIGNED
                NOT NULL
                COMMENT '是否可用',
  `func`        VARCHAR(75)
                NOT NULL
                DEFAULT 'nil'
                COMMENT '处理函数',
  `message`     VARCHAR(150)
                NOT NULL
                COMMENT '错误提示',
  `element`     VARCHAR(50)
                NOT NULL
                COMMENT '展现组件类型',
  `version`     INT UNSIGNED
                NOT NULL
                COMMENT '版本',
  `update_at`   DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
                COMMENT '修改时间',
  `create_at`   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
                COMMENT '创建时间',

  PRIMARY KEY (`name`),
  UNIQUE KEY `unique_1` (`uuid`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '规则表'
;

DROP TABLE IF EXISTS `mm_statements`;
CREATE TABLE `mm_statements` (
  `ticket_id`     INT UNSIGNED
                  NOT NULL
                  COMMENT '所属工单',
  `sequence`      SMALLINT UNSIGNED
                  NOT NULL
                  COMMENT '分解序号',
  `uuid`          CHAR(36)
                  NOT NULL
                  COMMENT 'UUID',
  `content`       TEXT
                  NOT NULL
                  COMMENT '单独语句',
  `type`          TINYINT UNSIGNED
                  NOT NULL
                  COMMENT '类型',
  `status`        TINYINT UNSIGNED
                  NOT NULL
                  COMMENT '审核状态',
  `report`        TEXT
                  NOT NULL
                  COMMENT '审核结果',
  `plan`          TEXT
                  COMMENT '执行计划',
  `results`       TEXT
                  COMMENT '执行结果',
  `rows_affected` INT UNSIGNED
                  COMMENT '在服务器正确执行后影响的行数',
  `version`       INT UNSIGNED
                  NOT NULL
                  COMMENT '版本',
  `update_at`     DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
                  COMMENT '修改时间',
  `create_at`     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
                  COMMENT '创建时间',

  PRIMARY KEY (`ticket_id`,`sequence`),
  UNIQUE KEY `unique_1` (`uuid`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '工单分解表'
;

DROP TABLE IF EXISTS `mm_metrics`;
CREATE TABLE `mm_metrics` (
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


DROP TABLE IF EXISTS `mm_tickets`;
CREATE TABLE `mm_tickets` (
  `ticket_id`   INT UNSIGNED
                NOT NULL
                AUTO_INCREMENT
                COMMENT '自增主键',
  `uuid`        CHAR(36)
                NOT NULL
                COMMENT 'UUID',
  `cluster_id` INT UNSIGNED
                NOT NULL
                COMMENT '目标群集',
  `database`    VARCHAR(75)
                NOT NULL
                COMMENT '目标库',
  `subject`     VARCHAR(50)
                NOT NULL
                COMMENT '主题',
  `content`     TEXT
                NOT NULL
                COMMENT '更新语句',
  `status`      TINYINT UNSIGNED
                NOT NULL
                COMMENT '状态',
  `user_id`     INT UNSIGNED
                NOT NULL
                COMMENT '申请人',
  `reviewer_id` INT UNSIGNED
                NOT NULL
                COMMENT '审核人',
  `cron_id`     INT UNSIGNED
                COMMENT '预约信息',
  `version`     INT UNSIGNED
                NOT NULL
                COMMENT '版本',
  `update_at`   DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
                COMMENT '修改时间',
  `create_at`   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
                COMMENT '创建时间',

  PRIMARY KEY (`ticket_id`),
  UNIQUE KEY `unique_1` (`uuid`),
  UNIQUE KEY `unique_2` (`cron_id`),
  KEY `index_1` (`user_id`),
  KEY `index_2` (`cluster_id`),
  KEY `index_3` (`reviewer_id`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '工单表'
;

DROP TABLE IF EXISTS `mm_users`;
CREATE TABLE `mm_users` (
  `user_id`   INT UNSIGNED
              NOT NULL
              AUTO_INCREMENT
              COMMENT '自增主键',
  `uuid`      CHAR(36)
              NOT NULL
              COMMENT 'UUID',
  `email`     VARCHAR(75)
              NOT NULL
              COMMENT '电子邮件',
  `password`  CHAR(60)
              NOT NULL
              COMMENT '密码',
  `status`    TINYINT UNSIGNED
              NOT NULL
              COMMENT '状态',
  `name`      VARCHAR(15)
              NOT NULL
              COMMENT '真实名称',
  `phone`     BIGINT UNSIGNED
              NOT NULL
              COMMENT '电话号码',
  `avatar_id` INT UNSIGNED
              NOT NULL
              COMMENT '头像',
  `version`   INT UNSIGNED
              NOT NULL
              COMMENT '版本',
  `update_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
              COMMENT '修改时间',
  `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
              COMMENT '创建时间',

  PRIMARY KEY (`user_id`),
  UNIQUE KEY `unique_1` (`uuid`),
  UNIQUE KEY `unique_2` (`email`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '用户表'
;

DROP TABLE IF EXISTS `mm_templates`;
CREATE TABLE `mm_templates` (
  `uuid`        CHAR(36)
                NOT NULL
                COMMENT 'UUID',
  `subject`     CHAR(100)
                NOT NULL
                COMMENT '主题',
  `body`        TEXT
                NOT NULL
                COMMENT '内容',
  `description` VARCHAR(50)
                NOT NULL
                COMMENT '模板说明',
  `version`     INT UNSIGNED
                NOT NULL
                COMMENT '版本',
  `update_at`   DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
                COMMENT '修改时间',
  `create_at`   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
                COMMENT '创建时间',

  PRIMARY KEY (`uuid`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '邮件模板表'
;

DROP TABLE IF EXISTS `mm_migrations`;
CREATE TABLE `mm_migrations` (
  `version`     INT UNSIGNED
                NOT NULL
                COMMENT '迁移版本号',
  `name`        VARCHAR(75)
                NOT NULL
                COMMENT '迁移名称',
  `checksum`    CHAR(64)
                NOT NULL
                COMMENT '迁移文件内容 sha256',
  `duration_ms` INT UNSIGNED
                NOT NULL DEFAULT 0
                COMMENT '执行耗时（毫秒）',
  `success`     TINYINT(1)
                NOT NULL DEFAULT 1
                COMMENT '是否执行成功',
  `apply_at`    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
                COMMENT '应用时间',

  PRIMARY KEY (`version`),
  UNIQUE KEY `unique_1` (`name`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '迁移版本台账'
;

DROP TABLE IF EXISTS `mm_violations`;
CREATE TABLE `mm_violations` (
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

DROP TABLE IF EXISTS `mm_messages`;
CREATE TABLE `mm_messages` (
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

DROP TABLE IF EXISTS `mm_groups`;
CREATE TABLE `mm_groups` (
  `group_id`    INT UNSIGNED
                NOT NULL
                AUTO_INCREMENT
                COMMENT '自增主键',
  `uuid`        CHAR(36)
                NOT NULL
                COMMENT 'UUID',
  `name`        VARCHAR(75)
                NOT NULL
                COMMENT '分组名称',
  `description` VARCHAR(150)
                NOT NULL DEFAULT ''
                COMMENT '分组说明',
  `status`      TINYINT UNSIGNED
                NOT NULL DEFAULT 1
                COMMENT '状态：1 正常 2 停用',
  `version`     INT UNSIGNED
                NOT NULL DEFAULT 0
                COMMENT '版本',
  `update_at`   DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
                COMMENT '修改时间',
  `create_at`   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
                COMMENT '创建时间',

  PRIMARY KEY (`group_id`),
  UNIQUE KEY `unique_1` (`uuid`),
  UNIQUE KEY `unique_2` (`name`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '群集分组表'
;

DROP TABLE IF EXISTS `mm_tasks`;
CREATE TABLE `mm_tasks` (
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
