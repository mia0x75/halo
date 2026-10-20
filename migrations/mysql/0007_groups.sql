-- 群集分组表。
--
-- 用户的组用角色（mm_roles）表达，群集服务器的组用本表表达：群集通过
-- mm_clusters.group_id 归属分组，角色通过 mm_edges 关联分组，实现组间权限隔离。

CREATE TABLE IF NOT EXISTS `mm_groups` (
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

-- 群集归属分组，0 表示未分组。
ALTER TABLE `mm_clusters`
    ADD COLUMN `group_id` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '所属分组，0 表示未分组';
