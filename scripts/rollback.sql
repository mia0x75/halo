USE `rollback`;

DROP TABLE IF EXISTS `$_$halo_rollbacks$_$`;
CREATE TABLE `$_$halo_rollbacks$_$` (
   `opid_time`         VARCHAR(50)
                       NOT NULL
                       COMMENT '执行操作ID，格式为时间戳+线程号+执行序号',
   `start_binlog_file` VARCHAR(25)
                       NOT NULL
                       COMMENT '起始日志文件',
   `start_binlog_pos`  INT UNSIGNED
                       NOT NULL
                       COMMENT '起始位置',
   `end_binlog_file`   VARCHAR(25)
                       NOT NULL
                       COMMENT '终止日志文件',
   `end_binlog_pos`    INT UNSIGNED
                       NOT NULL
                       COMMENT '终止位置',
   `content`           TEXT
                       NOT NULL
                       COMMENT '执行语句',
   `host`              VARCHAR(75)
                       NOT NULL
                       COMMENT '执行主机',
   `database`          VARCHAR(75)
                       NOT NULL
                       COMMENT '执行库名',
   `table`             VARCHAR(75)
                       NOT NULL
                       COMMENT '执行表名',
   `port`              SMALLINT UNSIGNED
                       NOT NULL
                       COMMENT '执行端口',
   `duration`          TIMESTAMP
                       NOT NULL
                       COMMENT '执行耗时',
   `type`              VARCHAR(20)
                       NOT NULL
                       COMMENT '操作类型'
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '回滚信息记录表'
;
