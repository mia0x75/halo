CREATE DATABASE IF NOT EXISTS `halodb` DEFAULT CHARACTER SET utf8mb4;

USE `halodb`;

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
   `update_at` INT UNSIGNED
               COMMENT '修改时间',
   `create_at` INT UNSIGNED
               NOT NULL
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

LOCK TABLES `mm_avatars` WRITE;
INSERT INTO `mm_avatars` VALUES
(0,'fec01539-1116-462c-91f6-22fe9e9fdf3b','/assets/images/avatars/01.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'f9e06623-3a1b-4221-8b27-e5461467edc1','/assets/images/avatars/02.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'11ae1ddc-e0d3-4c59-8fd4-c781522813fe','/assets/images/avatars/03.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'a57b0095-c398-4107-bf8c-056f6b3043d9','/assets/images/avatars/04.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'5af668d9-a5b1-46d8-bd09-1221d3b5bd69','/assets/images/avatars/05.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'494e55b3-a327-4d61-980a-a3d656d08146','/assets/images/avatars/06.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'93813cf7-9db3-4bf3-81cb-b06019aa206b','/assets/images/avatars/07.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'d108e7cd-f704-48af-8db8-ddae4c025183','/assets/images/avatars/08.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'2fee40ef-d15f-44a1-bc38-a738675438b4','/assets/images/avatars/09.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'b35d5d56-37c6-4846-8011-285c6ffd45e7','/assets/images/avatars/10.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'d4836d39-9dbf-492c-9846-b0f442de3c93','/assets/images/avatars/11.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'3d657f77-c08e-4c80-92ad-608a453c85a3','/assets/images/avatars/12.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'242bc8f7-3298-48ff-972e-e7ca313d9695','/assets/images/avatars/13.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'6ba5f281-77c9-4c90-a25d-b1368075782b','/assets/images/avatars/14.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'f53a9b19-996b-428b-ab32-83dd7609a378','/assets/images/avatars/15.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'b997e733-744d-438d-a99b-a0c8f6516fcf','/assets/images/avatars/16.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'334a0b7a-025a-4a26-a09b-f1354f22e08c','/assets/images/avatars/17.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'0a9a4366-3a04-4d3b-b595-b8f16b6892fe','/assets/images/avatars/18.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'a8aff0bc-440c-4f7c-af22-102adddf3af1','/assets/images/avatars/19.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'86a67d30-458c-4f5b-9646-50deee256754','/assets/images/avatars/20.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'8890af26-5675-4939-9fb6-a0d974cc45a6','/assets/images/avatars/21.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'c95496b9-b05b-45ba-b9e9-f3b07ca80d00','/assets/images/avatars/22.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'7005d057-0b0a-485e-8fb3-15ccdf56b4e6','/assets/images/avatars/23.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'f17584ae-3d23-46d4-abb8-5bff1ac8fcc4','/assets/images/avatars/24.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'35cc2cda-6f36-4c3f-b221-a6e5009abbf5','/assets/images/avatars/25.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'0ad5d283-87b9-491d-b54e-0244fd8732fa','/assets/images/avatars/26.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'9d45a30f-dbc8-443a-b09e-455fec2f0948','/assets/images/avatars/27.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'0e93f5c5-b38a-46c9-bc7d-60c0bf46e114','/assets/images/avatars/28.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'10ab8966-3742-49a3-b181-03e5cd93f722','/assets/images/avatars/29.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'9ca8e9d6-48f3-408e-82f5-874a5094fb38','/assets/images/avatars/30.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'bc6d8884-fa2a-4cf4-acab-d05f54190a24','/assets/images/avatars/31.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'5a1b2ada-c27f-490e-b858-10835eeb4127','/assets/images/avatars/32.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'b40af0f2-5998-4456-a310-021f8ae20e0c','/assets/images/avatars/33.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'ef09727c-f24b-4362-b3a8-118a878b614c','/assets/images/avatars/34.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'b4725215-a634-4df0-aa9e-6eff3571efe0','/assets/images/avatars/35.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'82a8a876-cc4b-40a0-8500-bfce5e95dce0','/assets/images/avatars/36.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'91dbcbb4-f29d-4167-8363-508b9ea7238a','/assets/images/avatars/37.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'a61ae780-8e58-4052-8158-2ec861a813fd','/assets/images/avatars/38.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'ad8cb15a-3e89-4d6c-830f-60d64922b049','/assets/images/avatars/39.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'82a9766f-c30f-426a-bc1e-d09c94ccd740','/assets/images/avatars/40.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'d02822aa-c1df-4e8e-96f5-32818e7aeccb','/assets/images/avatars/41.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'fff6ddf8-a461-4fe1-804f-11973de20c9d','/assets/images/avatars/42.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'cfa70a71-97dc-4585-9493-5e1e2aff16d4','/assets/images/avatars/43.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'e3c8549d-7cef-4f78-ac3e-795a48c9d8b6','/assets/images/avatars/44.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'d09316ff-9659-48e2-a5b5-3676fac6ecd8','/assets/images/avatars/45.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'2c6e0f48-a5e3-49e6-8484-b275b7796a79','/assets/images/avatars/46.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'cb2337ca-ad6b-46ff-8f8c-9ef1ff12ee70','/assets/images/avatars/47.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'7ff53eb4-231f-49fe-a85b-5ee315fcefc3','/assets/images/avatars/48.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'f60e1de7-972a-4c2f-9686-b1a6f814d391','/assets/images/avatars/49.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'6a669df1-53cd-44de-8138-8ff64995aa56','/assets/images/avatars/50.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'7e119b36-df2b-4f4c-880b-c30b1c12d974','/assets/images/avatars/51.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'86fa371d-038b-4553-8000-435cd8b428f7','/assets/images/avatars/52.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'6b2fb09a-f131-44a0-9ffb-a82a1a43e355','/assets/images/avatars/53.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'6512daf8-f3f6-4574-b3dc-b413482215a9','/assets/images/avatars/54.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'a308b254-64b5-41c7-b241-288aba391454','/assets/images/avatars/55.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'92237ff1-e1e4-4004-967b-fbb5709232f7','/assets/images/avatars/56.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'3a570e84-e82d-4e5e-8c0b-233f63e27332','/assets/images/avatars/57.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'60f1036c-2643-4561-8f28-89f98df03cf0','/assets/images/avatars/58.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'5b55ca9f-0569-428d-b24d-f0bc3efd0e5a','/assets/images/avatars/59.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'31a7bb40-0738-4dc1-81e6-ea520800f141','/assets/images/avatars/60.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'fe178ab8-a97b-46fb-9bdd-2cdab2ecb2ed','/assets/images/avatars/61.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'46530de5-03f0-4f3c-8f1b-38597cea4abb','/assets/images/avatars/62.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'28002cc0-3ca8-484d-852d-d442af81d275','/assets/images/avatars/63.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'280e40f0-0306-451b-8066-380c001f60c2','/assets/images/avatars/64.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'c1c5bf8e-f305-49a6-9e13-ae300297847f','/assets/images/avatars/65.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'82a4f20b-6c04-4229-bb8c-2614cad1432d','/assets/images/avatars/66.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'6a762662-670f-42c5-9d89-3217786c3673','/assets/images/avatars/67.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'af4bc4e7-6cc0-4309-8004-70eb64784eba','/assets/images/avatars/68.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'f70f8703-73f6-4a80-a967-0fc4039558d4','/assets/images/avatars/69.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'5298263c-4b9b-44ee-a879-40e75bb69e05','/assets/images/avatars/70.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'69a0963c-74b5-4a20-830f-700af9d90b26','/assets/images/avatars/71.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'c73e539a-16e9-49af-95e8-e6dec1e0407c','/assets/images/avatars/72.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'b21d1709-4db2-4c60-a9d7-22860440276c','/assets/images/avatars/73.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'3687403a-80ac-46e8-95e1-66f0ad9f1f37','/assets/images/avatars/74.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'dee08364-3951-4e43-8f95-d3b8473139ab','/assets/images/avatars/75.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'7fd7041f-8c89-423d-b1b7-5c8d17a1c072','/assets/images/avatars/76.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'cf44be0a-820c-4540-8d43-3c9af7cf39a6','/assets/images/avatars/77.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'493d8674-60cd-4add-be22-6d9b3d006f7c','/assets/images/avatars/78.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'3d725f8d-66c5-43b0-8a69-4ee8ffbb5e11','/assets/images/avatars/79.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'1beb5cae-5824-4639-8a7a-a7e1f169bdd9','/assets/images/avatars/80.png',1,1549962182, UNIX_TIMESTAMP()),
(0,'9d0db5e9-e9d1-4d61-a56e-0a4d216d0f63','/assets/images/avatars/81.png',1,1549962182, UNIX_TIMESTAMP());
UNLOCK TABLES;

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
  `update_at`  INT UNSIGNED
               COMMENT '修改时间',
  `create_at`  INT UNSIGNED
               NOT NULL
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
  `recurrent` TINYINT UNSIGNED
              COMMENT '是否周期运行',
  `hash`      VARCHAR(60)
              COMMENT '哈希值',
  `version`   INT UNSIGNED
              NOT NULL
              COMMENT '版本',
  `update_at` INT UNSIGNED
              COMMENT '修改时间',
  `create_at` INT UNSIGNED
              NOT NULL
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
  `update_at`   INT UNSIGNED
                COMMENT '修改时间',
  `create_at`   INT UNSIGNED
                NOT NULL
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

LOCK TABLES `mm_glossaries` WRITE;
INSERT INTO `mm_glossaries` VALUES
('2b45a7ea-e3f6-432e-a606-87964696c9ff','charsets',1,'big5','',1,1552357173, UNIX_TIMESTAMP()),
('8c3a0561-0d5f-4e18-a055-9f0cb670cabb','charsets',6,'latin1','',1,1552357173, UNIX_TIMESTAMP()),
('06ce1b18-dead-4be0-bcc2-2bfe1882241e','charsets',7,'latin2','',1,1552357173, UNIX_TIMESTAMP()),
('4245b14a-0853-439c-9b75-883e065d29b6','charsets',9,'ascii','',1,1552357173, UNIX_TIMESTAMP()),
('eded3f08-0193-40b9-8890-13a955aab2c5','charsets',16,'gb2312','',1,1552357173, UNIX_TIMESTAMP()),
('99f43671-1882-480b-98c1-a027f65d7136','charsets',19,'gbk','',1,1552357173, UNIX_TIMESTAMP()),
('9d101bb8-eee6-4356-a5ce-54fd4d4b55a0','charsets',20,'latin5','',1,1552357173, UNIX_TIMESTAMP()),
('29e5d2b3-26f4-4684-977e-f976b41fb9c0','charsets',22,'utf8','',1,1552357173, UNIX_TIMESTAMP()),
('d9d55fee-ee68-43ed-a06a-3221447f40ec','charsets',29,'latin7','',1,1552357173, UNIX_TIMESTAMP()),
('a870dbce-8ddd-4e54-bbf6-036c730dbac3','charsets',30,'utf8mb4','',1,1552357173, UNIX_TIMESTAMP()),
('12ac3ba1-d808-4b68-99ea-eb28fff279e6','charsets',37,'binary','',1,1552357173, UNIX_TIMESTAMP()),
('0cd6501f-e9fe-4658-b690-f00150e77168','collates',1,'big5_chinese_ci','',1,1552357173, UNIX_TIMESTAMP()),
('50915ad7-819b-4221-9f61-8ecbc7dd2acb','collates',9,'latin2_general_ci','',1,1552357173, UNIX_TIMESTAMP()),
('7b2cdb22-45ad-4545-bd58-f347d80bb725','collates',11,'ascii_general_ci','',1,1552357173, UNIX_TIMESTAMP()),
('849ffa37-d853-4629-8b60-9feee7e05939','collates',24,'gb2312_chinese_ci','',1,1552357173, UNIX_TIMESTAMP()),
('4dc2884a-bd6b-4e5c-b2d8-346ec19850a2','collates',28,'gbk_chinese_ci','',1,1552357173, UNIX_TIMESTAMP()),
('05126857-f916-49aa-8086-ce49cfe3c485','collates',33,'utf8_general_ci','',1,1552357173, UNIX_TIMESTAMP()),
('5344c214-12f2-42f4-a546-c3dba69a7b83','collates',41,'latin7_general_ci','',1,1552357173, UNIX_TIMESTAMP()),
('52c1c9a7-4a56-4ca4-80f1-cb15be2103db','collates',42,'latin7_general_cs','',1,1552357173, UNIX_TIMESTAMP()),
('63283439-1e38-4107-8768-0b72689f8797','collates',45,'utf8mb4_general_ci','',1,1552357173, UNIX_TIMESTAMP()),
('512720e6-d0c0-43af-9164-72704e0417d6','collates',46,'utf8mb4_bin','',1,1552357173, UNIX_TIMESTAMP()),
('fbe8d595-dd8f-452e-aa32-b75dd9a7bd24','collates',47,'latin1_bin','',1,1552357173, UNIX_TIMESTAMP()),
('7d9b0978-59e8-4b92-b138-900ad8fd071d','collates',48,'latin1_general_ci','',1,1552357173, UNIX_TIMESTAMP()),
('22068a13-77db-4f84-adda-623f314b05a4','collates',49,'latin1_general_cs','',1,1552357173, UNIX_TIMESTAMP()),
('5cf1df88-78eb-4dbe-8d46-84746910f51d','collates',63,'binary','',1,1552357173, UNIX_TIMESTAMP()),
('18da1c87-2413-4f49-9b51-bf2e942222ef','collates',65,'ascii_bin','',1,1552357173, UNIX_TIMESTAMP()),
('30455150-25d7-4903-8a98-defc092917e0','collates',77,'latin2_bin','',1,1552357173, UNIX_TIMESTAMP()),
('ba7a8501-408b-470c-bf16-e4564f0459be','collates',78,'latin5_bin','',1,1552357173, UNIX_TIMESTAMP()),
('079210d8-f029-4ece-bb0e-762fbb641dd6','collates',79,'latin7_bin','',1,1552357173, UNIX_TIMESTAMP()),
('f5760bb3-8f95-4c6b-bb72-a8c84c2edc9e','collates',83,'utf8_bin','',1,1552357173, UNIX_TIMESTAMP()),
('47b726b5-6fb2-47c2-9b61-bc92c758c0ea','collates',84,'big5_bin','',1,1552357173, UNIX_TIMESTAMP()),
('dd160bb0-6447-4981-9f95-94cd1914e4dd','collates',86,'gb2312_bin','',1,1552357173, UNIX_TIMESTAMP()),
('bfc7d81c-e4fe-4da1-bb79-d43b4c1e80b3','collates',87,'gbk_bin','',1,1552357173, UNIX_TIMESTAMP()),
('1a632227-bb99-4309-b493-afebe5b34b33','collates',192,'utf8_unicode_ci','',1,1552357173, UNIX_TIMESTAMP()),
('67015690-52f0-4b2a-8e08-66480c2915dc','collates',224,'utf8mb4_unicode_ci','',1,1552357173, UNIX_TIMESTAMP()),
('ff2b7083-8a15-41e3-b0a9-76fc777d7733','data-types',1,'bit','',1,1552357173, UNIX_TIMESTAMP()),
('c6f5448e-3964-435f-8f65-a4c321e183b8','data-types',2,'boolean','',1,1552357173, UNIX_TIMESTAMP()),
('7a5e07d2-91ef-4006-b4b2-41ab1e723c1d','data-types',3,'tinyint','',1,1552357173, UNIX_TIMESTAMP()),
('c1f2dcbb-a4c3-4091-8e26-ee277920f4ff','data-types',4,'smallint','',1,1552357173, UNIX_TIMESTAMP()),
('f00019b5-ec19-4c83-91cf-31e462e49c09','data-types',5,'mediumint','',1,1552357173, UNIX_TIMESTAMP()),
('95decec8-23c3-46d8-837f-255d77014da2','data-types',6,'int','',1,1552357173, UNIX_TIMESTAMP()),
('1f646b75-9d2d-4d6c-aa44-829d568c55ef','data-types',7,'bigint','',1,1552357173, UNIX_TIMESTAMP()),
('1bd6cef3-68d4-42f4-8e8c-0dbe65879e00','data-types',8,'decimal','',1,1552357173, UNIX_TIMESTAMP()),
('eaa68ac7-2f70-4d1f-ac88-d1fa9b45d533','data-types',9,'float','',1,1552357173, UNIX_TIMESTAMP()),
('6ff2cf69-a701-426b-bde2-4a0bc3c9bfae','data-types',10,'double','',1,1552357173, UNIX_TIMESTAMP()),
('307d7774-2b55-4b39-a9d4-c1297943f10e','data-types',11,'real','',1,1552357173, UNIX_TIMESTAMP()),
('6f04f90f-1c63-4f1f-83f8-a523e96586e0','data-types',12,'timestamp','',1,1552357173, UNIX_TIMESTAMP()),
('de02df71-1652-4e57-bed0-e90a04487da0','data-types',13,'date','',1,1552357173, UNIX_TIMESTAMP()),
('b2d46779-0276-4cd8-be10-90da8b192a32','data-types',14,'time','',1,1552357173, UNIX_TIMESTAMP()),
('6df81c48-781b-42ae-907f-af72979efa12','data-types',15,'datetime','',1,1552357173, UNIX_TIMESTAMP()),
('ad61b129-46ca-406e-8ddf-dff67a364fda','data-types',16,'year','',1,1552357173, UNIX_TIMESTAMP()),
('0ab0a3fe-80bb-4c21-85d3-0931ebd5ad8d','data-types',17,'char','',1,1552357173, UNIX_TIMESTAMP()),
('6603e15f-1320-4fcb-bfc2-7dedad8e2f11','data-types',18,'varchar','',1,1552357173, UNIX_TIMESTAMP()),
('8f241cd5-b8a2-4cf9-bf32-9056e7d39209','data-types',19,'json','',1,1552357173, UNIX_TIMESTAMP()),
('7cd538ea-baeb-4a45-9179-95392641c1ac','data-types',20,'enum','',1,1552357173, UNIX_TIMESTAMP()),
('38da2334-7104-4799-9110-ce43b4821d06','data-types',21,'set','',1,1552357173, UNIX_TIMESTAMP()),
('722d9c25-fadd-4bcb-adf4-aa423adb2f66','data-types',22,'binary','',1,1552357173, UNIX_TIMESTAMP()),
('27ab9242-2085-47c5-8f2e-2119bc7c2275','data-types',23,'varbinary','',1,1552357173, UNIX_TIMESTAMP()),
('264f47c5-e35d-4e4b-99af-e88b4dcb24be','data-types',24,'tinyblob','',1,1552357173, UNIX_TIMESTAMP()),
('35cf4ab6-5cc7-406d-95be-c43a84bdc776','data-types',25,'blob','',1,1552357173, UNIX_TIMESTAMP()),
('e0a46c77-fafd-4428-b3a8-958dd2095771','data-types',26,'mediumblob','',1,1552357173, UNIX_TIMESTAMP()),
('fa8e7a4b-a224-4b1c-8daf-9cf3a34198a1','data-types',27,'longblob','',1,1552357173, UNIX_TIMESTAMP()),
('81497f01-d756-4b68-b62a-70428ac33bfd','data-types',28,'tinytext','',1,1552357173, UNIX_TIMESTAMP()),
('bd1006da-2cca-4279-8791-fbe8467dbb18','data-types',29,'text','',1,1552357173, UNIX_TIMESTAMP()),
('c74104a0-845d-47c1-97d1-3bb8b8ab29b3','data-types',30,'mediumtext','',1,1552357173, UNIX_TIMESTAMP()),
('12d90114-2c2e-4f62-9635-4aa178f50048','data-types',31,'longtext','',1,1552357173, UNIX_TIMESTAMP()),
('22645b4d-3630-448c-80d3-192f71278e7c','engines',1,'innodb','',1,1552357173, UNIX_TIMESTAMP()),
('b557b62d-2c4d-4ff5-975e-dd81747b3cbd','engines',2,'myisam','',1,1552357173, UNIX_TIMESTAMP()),
('66a045fd-f3d1-44f3-9603-d33300ddabc8','engines',3,'csv','',1,1552357173, UNIX_TIMESTAMP()),
('5c073a3a-c5b1-4580-91e3-247a049b14a1','engines',4,'memory','',1,1552357173, UNIX_TIMESTAMP()),
('45f38653-a147-4003-9fb4-beed1eb5434f','engines',5,'blackhole','',1,1552357173, UNIX_TIMESTAMP()),
('e14a5a52-3198-49dc-bba2-4c8320e55a79','engines',6,'tokudb','',1,1552357173, UNIX_TIMESTAMP()),
('eb13ec88-e543-40bb-a212-dcca923321f3','engines',7,'rocksdb','',1,1552357173, UNIX_TIMESTAMP()),
('b7ff66ec-88a8-4e60-a8e3-9e12a5787dc8','engines',8,'archive','',1,1552357173, UNIX_TIMESTAMP()),
('3ec651ec-90ab-4b72-8926-ba79b1b11e55','engines',9,'aria','',1,1552357173, UNIX_TIMESTAMP()),
('41ee05aa-469a-4d2b-a0a7-73dcb00ac0a6','engines',10,'cassandra','',1,1552357174, UNIX_TIMESTAMP()),
('8c23a0b5-3b80-436e-8b42-240c76c6d4b4','engines',11,'federated','',1,1552357174, UNIX_TIMESTAMP()),
('b32b3a83-f2ce-42bc-aff4-7b24ae5396d9','clusters.status',1,'正常','正常',1,1552357174, UNIX_TIMESTAMP()),
('204471d6-4947-44c6-9771-a96573fed512','clusters.status',2,'停用','停用',1,1552357174, UNIX_TIMESTAMP()),
('b8876688-eea9-4b0b-9614-97623644edf0','options.group',1,'邮件服务器 配置选项','用于配置邮件发送服务器，通过邮件发送系统通知',1,1552357174, UNIX_TIMESTAMP()),
('86903652-d561-4262-b33a-d3a3d3d34180','options.group',2,'LDAP登录 配置选项','用于配置LDAP登录',1,1552357174, UNIX_TIMESTAMP()),
('1776cec8-6f6f-42ee-adbe-527f7b6eb70d','edges.type',1,'...','...',1,1552357174, UNIX_TIMESTAMP()),
('aaac96e4-f34b-48f0-ad06-ec228ab9df0d','rules.group',1,'DATABASE 规则分组','用于检查数据库相关操作的审核规则',1,1552357174, UNIX_TIMESTAMP()),
('6ff95c2f-d554-4c45-8f82-abeff13789be','rules.group',2,'CREATE-TABLE 规则分组','用于检查创建表的审核规则',1,1552357174, UNIX_TIMESTAMP()),
('c12e0348-4ab9-48e4-b6b6-79a1a53cb76b','rules.group',3,'ALTER-TABLE 规则分组','用于检查修改表的审核规则',1,1552357174, UNIX_TIMESTAMP()),
('2c2ca13f-45e2-44e0-b665-019015e6b003','rules.group',4,'RENAME-TABLE 规则分组','用于检查对表重命名的审核规则',1,1552357174, UNIX_TIMESTAMP()),
('efc3111c-73c6-4055-a42e-63c1d657e5e5','rules.group',5,'DROP-TABLE 规则分组','用于检查删除表的审核规则',1,1552357174, UNIX_TIMESTAMP()),
('92cc9880-c9de-435e-ad0d-bf4fbb156c4e','rules.group',6,'INSERT 规则分组','用于检查插入数据的审核规则',1,1552357174, UNIX_TIMESTAMP()),
('1e938c51-e7f5-4fc1-a0f2-076f2174c2f3','rules.group',7,'UPDATE 规则分组','用于检查更新数据的审核规则',1,1552357174, UNIX_TIMESTAMP()),
('64f48887-66c6-4855-bda8-1e7976187801','rules.group',8,'DELETE 规则分组','用于检查删除数据的审核规则',1,1552357174, UNIX_TIMESTAMP()),
('3d41e8a8-3afd-40d9-9cbe-3c6c04d88141','rules.group',9,'SELECT 规则分组','用于检查查询数据的审核规则',1,1552357174, UNIX_TIMESTAMP()),
('4b47ad17-e40b-49c6-bbb3-8661da79b0b8','rules.group',10,'VIEW 规则分组','用于检查视图相关的审核规则',1,1552357174, UNIX_TIMESTAMP()),
('70d80234-7c41-4f99-be7c-0220763a34f3','rules.group',11,'FUNCTION 规则分组','用于检查函数相关的审核规则',1,1552357174, UNIX_TIMESTAMP()),
('63c0f1cf-4cab-41e6-b075-afd4c86c36f6','rules.group',12,'TRIGGER 规则分组','用于检查触发器相关的审核规则',1,1552357174, UNIX_TIMESTAMP()),
('55621602-4e63-4922-9957-60f0bfeafd6d','rules.group',13,'EVENT 规则分组','用于检查事件相关的审核规则',1,1552357174, UNIX_TIMESTAMP()),
('4afe41e2-1fd2-4eb4-8538-2e49418b3445','rules.group',14,'PROCEDURE 规则分组','用于检查存储过程相关的审核规则',1,1552357174, UNIX_TIMESTAMP()),
('571a3e92-d641-4ca0-b111-4a30b267361c','rules.group',15,'其他 规则分组','用于检查没有分类的其他操作的审核规则',1,1552357174, UNIX_TIMESTAMP()),
('0dd6123b-c517-4302-9aee-c22ced97c7e6','rules.group',16,'合并拆分 规则分组','用于检查没有分类的其他操作的审核规则',1,1552357174, UNIX_TIMESTAMP()),
('664de121-0571-48d4-b839-653dc1821d76','rules.severity',0,'消息','保留',1,1552357174, UNIX_TIMESTAMP()),
('42dfa1f6-5e06-423e-91d7-49220b0c733d','rules.severity',1,'错误','错误级别的审核规则，不允许上线执行',1,1552357174, UNIX_TIMESTAMP()),
('cd36779c-18fe-407a-ab0d-22c0851592ee','rules.severity',2,'警告','警告级别的审核规则，不建议上线执行',1,1552357174, UNIX_TIMESTAMP()),
('283706fd-e149-4033-909f-9410ae163aaf','statements.status',1,'待审核','待审核',1,1552357174, UNIX_TIMESTAMP()),
('f1505ec1-6a5e-40c8-953f-ac692c66b1c6','statements.status',2,'通过','通过',1,1552357174, UNIX_TIMESTAMP()),
('368713c4-18ee-4ab0-ae0d-8db01038430f','statements.status',3,'警告','警告',1,1552357174, UNIX_TIMESTAMP()),
('afbed129-4c11-439d-89f5-688ff535e020','statements.status',4,'错误','错误',1,1552357174, UNIX_TIMESTAMP()),
('2f8e3909-d4b5-4916-95d7-7b6423db0f27','tickets.status',1,'等待系统审核','等待系统审核',1,1552357174, UNIX_TIMESTAMP()),
('89221f30-5ab3-4a63-b933-9f823f551785','tickets.status',2,'等待人工审核','等待人工审核',1,1552357174, UNIX_TIMESTAMP()),
('978aeceb-9134-45ab-a0f5-3b851731e5d2','tickets.status',3,'系统审核失败','系统审核失败',1,1552357174, UNIX_TIMESTAMP()),
('857a3c15-299b-488a-9121-3ced66c6e4f3','tickets.status',4,'warn','warn',1,1552357174, UNIX_TIMESTAMP()),
('8bfe6a5c-2f54-42cd-8654-0dad1ac5b022','tickets.status',5,'人工审核失败','人工审核失败',1,1552357174, UNIX_TIMESTAMP()),
('21131a99-e459-443e-a37a-32f3fa80811b','tickets.status',6,'上线执行完成','上线执行完成',1,1552357174, UNIX_TIMESTAMP()),
('cb6fe0ea-815b-4857-8d3d-52227ca92112','tickets.status',7,'上线执行失败','上线执行失败',1,1552357174, UNIX_TIMESTAMP()),
('cff33cc0-705e-4342-8407-dcc1d1aa8f35','users.status',0,'等待激活','等待激活',1,1552357174, UNIX_TIMESTAMP()),
('64068bae-20b2-4752-bfae-b0806b7aeab0','users.status',1,'正常','正常',1,1552357174, UNIX_TIMESTAMP()),
('f97cb12f-5082-45bb-b0b5-52510859a53e','users.status',2,'禁用','禁用',1,1552357174, UNIX_TIMESTAMP());
UNLOCK TABLES;

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
  `password`    VARBINARY(48)
                NOT NULL
                COMMENT '密码',
  `fingerprint` VARBINARY(20)
                NOT NULL
                COMMENT '指纹',
  `status`      TINYINT UNSIGNED
                NOT NULL
                DEFAULT 1
                COMMENT '状态',
  `version`     INT UNSIGNED
                NOT NULL
                COMMENT '版本',
  `update_at`   INT UNSIGNED
                COMMENT '修改时间',
  `create_at`   INT UNSIGNED
                NOT NULL
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
  `create_at` INT UNSIGNED
              NOT NULL
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
  `value`       TINYTEXT
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
  `update_at`   INT UNSIGNED
                COMMENT '修改时间',
  `create_at`   INT UNSIGNED
                NOT NULL
                COMMENT '创建时间',

  PRIMARY KEY (`name`),
  UNIQUE KEY `unique_1` (`uuid`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '系统选项表'
;

LOCK TABLES `mm_options` WRITE;
INSERT INTO `mm_options` VALUES
('0b66b11d-fc6b-44cc-84e1-852c8cb09b7b','smtp.enabled','127.0.0.1','','-',0,2,1545013963, UNIX_TIMESTAMP()),
('766efa82-6c66-4b5e-b3e4-7d633a09839c','smtp.host','127.0.0.1','','-',0,2,1545013963, UNIX_TIMESTAMP()),
('39cd3a0e-355d-4494-9b9b-e288a26bf3a2','smtp.port','127.0.0.1','','-',0,2,1545013963, UNIX_TIMESTAMP()),
('42c9c43c-5c22-4501-a74a-e9b0e3bac6f3','smtp.user','127.0.0.1','','-',0,2,1545013963, UNIX_TIMESTAMP()),
('5cf28fe0-9d75-4da5-95a8-9bec61985553','smtp.password','127.0.0.1','','-',0,2,1545013963, UNIX_TIMESTAMP()),
('a558ccae-9b61-476c-a3ed-e2b7bcc0ed5b','smtp.encryption','127.0.0.1','','-',0,2,1545013963, UNIX_TIMESTAMP()),
('ff68ad88-1d85-4306-b5f2-20ea1779dc2c','ldap.enabled','127.0.0.1','','-',0,2,1545013963, UNIX_TIMESTAMP()),
('23cd8cd9-0ddd-4a88-b009-484a7fef0b37','ldap.host','127.0.0.1','','-',0,2,1545013963, UNIX_TIMESTAMP()),
('217bd72d-464a-475c-9dd5-3d788ec92f0e','ldap.domain','127.0.0.1','','-',0,2,1545013963, UNIX_TIMESTAMP()),
('caf3cb07-3fa9-41ed-9f30-7fd17342e7c5','ldap.type','127.0.0.1','','-',0,2,1545013963, UNIX_TIMESTAMP()),
('55720a44-f204-4d5f-8ab4-12665d840540','ldap.user','127.0.0.1','','-',0,2,1545013963, UNIX_TIMESTAMP()),
('32f77401-0b1f-40e5-b6a8-bdebfb219ce2','ldap.password','127.0.0.1','','-',0,2,1545013963, UNIX_TIMESTAMP()),
('6e200836-4cb7-4088-8b7f-ea8c110692c3','ldap.sc','127.0.0.1','','-',0,2,1545013963, UNIX_TIMESTAMP()),
('43549a0d-aa4c-43ea-848a-41cc41517277','ldap.ou','127.0.0.1','','-',0,2,1545013963, UNIX_TIMESTAMP());
UNLOCK TABLES;

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
  `update_at`  INT UNSIGNED
               COMMENT '修改时间',
  `create_at`  INT UNSIGNED
               NOT NULL
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
  `update_at`     INT UNSIGNED
                  COMMENT '修改时间',
  `create_at`     INT UNSIGNED
                  NOT NULL
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

LOCK TABLES `mm_edges` WRITE;
INSERT INTO `mm_edges` VALUES
(1, 'dfd02e74-cd54-4ca8-984a-42dcfc2f4d35', 3, 1, 1,  1, 0,  UNIX_TIMESTAMP()),
(2, '6483249c-3fe9-48c1-9ba9-520f89f2db07', 2, 1, 1,  1, 0,  UNIX_TIMESTAMP()),
(3, 'dd98a7d7-3d69-4161-b45b-d090b019acf7', 2, 1, 2,  1, 0,  UNIX_TIMESTAMP()),
(4, 'eea7ef37-47e8-444d-b1bd-459d6bc99354', 2, 1, 3,  1, 0,  UNIX_TIMESTAMP()),
(5, '716989c1-ad05-4d06-98e2-9f59c9cdbd90', 3, 1, 5,  1, 0,  UNIX_TIMESTAMP()),
(6, '45970fc3-bde6-4fa5-a81f-758391bee397', 1, 1, 1,  0, 0,  UNIX_TIMESTAMP()),
(7, '4742b7f2-e270-4906-b2c1-08850e3e7512', 1, 1, 21, 0, 0,  UNIX_TIMESTAMP());
UNLOCK TABLES;

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
  `update_at`   INT UNSIGNED
                COMMENT '修改时间',
  `create_at`   INT UNSIGNED
                NOT NULL
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

LOCK TABLES `mm_roles` WRITE;
INSERT INTO `mm_roles` VALUES
(1, '97452ae1-17b9-4791-bc29-62aad89e2dbb', 'Root',      '系统管理员', 1, 0,  UNIX_TIMESTAMP()),
(2, 'd0c5dbd4-e42e-4958-8606-19ffcc5c7dfd', 'Reviewer',  '普通审核', 1, 0,  UNIX_TIMESTAMP()),
(3, '3d0f426d-101a-410f-b61e-b95133847f49', 'Developer', '数据查询及工单提交',1, 0,  UNIX_TIMESTAMP()),
(4, '914d2121-fc17-4328-92e8-5642d6482fa8', 'Viewer',    '查询用户', 1, 0,  UNIX_TIMESTAMP()),
(5, '0aa345ab-3235-43f0-a4c3-4d7f10b62e2f', 'User',      '注册用户', 1, 0,  UNIX_TIMESTAMP()),
(6, 'c562a523-eb23-44f2-9e07-d73237638ba6', 'Guest',     '来宾账号', 1, 0,  UNIX_TIMESTAMP());
UNLOCK TABLES;

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
  `values`      VARCHAR(150)
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
  `update_at`   INT UNSIGNED
                COMMENT '修改时间',
  `create_at`   INT UNSIGNED
                NOT NULL
                COMMENT '创建时间',

  PRIMARY KEY (`name`),
  UNIQUE KEY `unique_1` (`uuid`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '规则表'
;

LOCK TABLES `mm_rules` WRITE;
INSERT INTO `mm_rules` VALUES
('e4b6f058-f8e3-469d-b754-ac754c9bd07e', 'CDB-L2-001', 10, '新建数据库时允许的字符集', 2, 100, 'in', '["utf8mb4","binary"]', 7, 'DatabaseCreateAvailableCharsets', '建库禁用字符集"%s"，请使用"%s"。', 'checkboxes/key=CHARsets', 1, 0, UNIX_TIMESTAMP()),
('5126367e-19bf-4996-96eb-b92d51860acc', 'CDB-L2-002', 10, '新建数据库时允许的排序规则', 2, 100, 'none', '["utf8mb4_general_ci", "utf8mb4_bin", "utf8mb4_unicode_ci"]', 7, 'DatabaseCreateAvailableCollates', '建库禁用排序规则"%s"，请使用"%s"。', 'checkboxes/key=collates', 1, 0, UNIX_TIMESTAMP()),
('8e549891-ace6-48ba-bab7-6c333851098f', 'CDB-L2-003', 10, '新建数据库时字符集与排序规则必须匹配', 2, 100, 'none', 'nil', 5, 'DatabaseCreateCharsetCollateMatch', '建库使用的字符集"%s"和排序规则"%s"不匹配，请查阅官方文档。', 'none', 1, 0, UNIX_TIMESTAMP()),
('439a8103-7664-48c8-ae3d-227deb057416', 'CDB-L2-004', 10, '库名规则', 2, 100, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 7, 'DatabaseCreateDatabaseNameQualified', '库名"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('ce313eeb-aaa3-4ccc-ab52-6f8694bff63e', 'CDB-L2-005', 10, '库名必须小写', 2, 100, 'none', '^[_a-z0-9]+$', 7, 'DatabaseCreateDatabaseNameLowerCaseRequired', '库名"%s"中含有大写字母。', 'none', 1, 0, UNIX_TIMESTAMP()),
('5575b7a8-fda5-4c07-8fed-5457839334bb', 'CDB-L2-006', 10, '库名最大长度', 2, 100, 'lte', '15', 7, 'DatabaseCreateDatabaseNameMaxLength', '库名"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('d8a0fad2-773a-49d7-b1b9-d1219158222b', 'CDB-L3-001', 10, '新建数据库时目标库必须不存在', 3, 100, 'none', 'nil', 4, 'DatabaseCreateTargetDatabaseExists', '目标库"%s"已存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('1abf3c57-d84d-4295-9cbe-4b95f9e3c9bc', 'DDB-L3-001', 10, '删除数据库时目标库必须已存在', 3, 102, 'none', 'nil', 4, 'DatabaseDropTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('78ab3a56-28dd-463d-b8dd-bcb1929d2314', 'MDB-L2-001', 10, '修改数据库时允许的字符集', 2, 101, 'in', '["utf8mb4"]', 7, 'DatabaseAlterAvailableCharsets', '改库禁用字符集"%s"，请使用"%s"。', 'checkboxes/key=CHARsets', 1, 0, UNIX_TIMESTAMP()),
('a28e9a5d-2d26-433b-b4c6-fb68faa05daa', 'MDB-L2-002', 10, '修改数据库时允许的排序规则', 2, 101, 'in', '["utf8mb4_unicode_ci", "utf8mb4_general_ci", "utf8mb4_bin"]', 7, 'DatabaseAlterAvailableCollates', '改库禁用排序规则"%s"，请使用"%s"。', 'checkboxes/key=collates', 1, 0, UNIX_TIMESTAMP()),
('3c10b2cc-ebfb-4bf6-8543-44a91506d86b', 'MDB-L2-003', 10, '修改数据库时字符集与排序规则必须匹配', 2, 101, 'none', 'nil', 5, 'DatabaseAlterCharsetCollateMatch', '改库使用的字符集"%s"和排序规则"%s"不匹配，请查阅官方文档。', 'none', 1, 0, UNIX_TIMESTAMP()),
('c370fd11-b92a-4619-a64a-2deae0873224', 'MDB-L3-001', 10, '修改数据库时目标库必须已存在', 3, 101, 'none', 'nil', 4, 'DatabaseAlterTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('2e7cbe13-90d3-472a-9da2-4594c8424df0', 'CTB-L2-001', 11, '允许的字符集', 2, 110, 'in', '["utf8mb4"]', 7, 'TableCreateAvailableCharsets', '建表禁用字符集"%s"，请使用"%s"。', 'checkboxes/key=CHARsets', 1, 0, UNIX_TIMESTAMP()),
('5028dfe3-949d-4220-b9ae-a29d8a16de52', 'CTB-L2-002', 11, '允许的排序规则', 2, 110, 'in', '["utf8mb4_unicode_ci", "utf8mb4_general_ci", "utf8mb4_bin"]', 5, 'TableCreateAvailableCollates', '建表禁用排序规则"%s"，请使用"%s"。', 'checkboxes/key=collates', 1, 0, UNIX_TIMESTAMP()),
('81b9dec9-3769-4dff-a29a-8739c3ac4ec9', 'CTB-L2-003', 11, '字符集与排序规则必须匹配', 2, 110, 'none', 'nil', 5, 'TableCreateTableCharsetCollateMatch', '建表使用的字符集"%s"和排序规则"%s"不匹配，请查阅官方文档。', 'none', 1, 0, UNIX_TIMESTAMP()),
('a6dd1817-a41f-4688-924f-df959ffbb2db', 'CTB-L2-004', 11, '允许的存储引擎', 2, 110, 'in', '["innodb", "tokudb", "rocksdb", "archive"]', 7, 'TableCreateAvailableEngines', '建表禁用存储引擎"%s"，请使用"%s"。', 'checkboxes/key=engines', 1, 0, UNIX_TIMESTAMP()),
('0862861b-4328-49eb-86dc-07b15e0a2a6c', 'CTB-L2-005', 11, '表名规则', 2, 110, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 5, 'TableCreateTableNameQualified', '表名"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('f77613a7-6f1c-48ff-8946-e725d5431e72', 'CTB-L2-006', 11, '表名必须小写', 2, 110, 'regexp', '^[_a-z0-9]+$', 5, 'TableCreateTableNameLowerCaseRequired', '表名"%s"含有大写字母。', 'none', 1, 0, UNIX_TIMESTAMP()),
('19630e1b-fee5-4afc-8fd1-28f43efb74c6', 'CTB-L2-007', 11, '表名最大长度', 2, 110, 'lte', '20', 7, 'TableCreateTableNameMaxLength', '表名"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('335bd848-2027-4c66-b5b0-09f225f02bb9', 'CTB-L2-008', 11, '表必须有注释', 2, 110, 'none', 'nil', 7, 'TableCreateTableCommentRequired', '需要为表"%s"需要提供COMMENT注解。', 'none', 1, 0, UNIX_TIMESTAMP()),
('6f945479-4788-4bdb-b1e2-e21dc3b335d6', 'CTB-L2-009', 11, '禁止使用CREATE TABLE ... SELECT ...建表', 2, 110, 'none', 'nil', 5, 'TableCreateUseSelectEnabled', '禁止使用CREATE TABLE AS SELECT的方式建表。', 'none', 1, 0, UNIX_TIMESTAMP()),
('0b0e3799-7f6e-4a46-a054-7437b5485d06', 'CTB-L2-010', 11, '列名规则', 2, 110, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 5, 'TableCreateColumnNameQualified', '列名"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('d99679c9-ffbf-48e6-9dc8-bb1a82d10ff6', 'CTB-L2-011', 11, '列名必须小写', 2, 110, 'regexp', '^[_a-z0-9]+$', 7, 'TableCreateColumnNameLowerCaseRequired', '列名"%s"中含有大写字母。', 'none', 1, 0, UNIX_TIMESTAMP()),
('86530dce-57e4-41b6-a296-0ccbec5036a3', 'CTB-L2-012', 11, '列名最大长度', 2, 110, 'lte', '20', 7, 'TableCreateColumnNameMaxLength', '列名"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('97842c81-7b38-406b-a9a3-88244ea79271', 'CTB-L2-013', 11, '列名是否重复', 2, 110, 'none', 'nil', 5, 'TableCreateColumnNameDuplicate', '表"%s"中的定义了重复的列"%s"。', 'none', 1, 0, UNIX_TIMESTAMP()),
('2c98b6b7-0e0e-438e-8ccd-fa73fbcdbbea', 'CTB-L2-014', 11, '表允许的最大列数', 2, 110, 'lte', '25', 7, 'TableCreateColumnCountLimit', '表"%s"中定义%d个列，数量超出了规则允许的上限%d，请考虑拆分表。', 'number', 1, 0, UNIX_TIMESTAMP()),
('31ab4d7a-4243-4951-bf7e-ddbdbc03254f', 'CTB-L2-015', 11, '列禁用的数据类型', 2, 110, 'not-in', '["bit", "enum", "set", "double", "real", "float"]', 7, 'TableCreateColumnUnwantedTypes', '列"%s"使用了不期望的数据类型"%s"，请避免使用"%s"数据类型。', 'checkboxes/key=data-types', 1, 0, UNIX_TIMESTAMP()),
('b14a529b-bbfc-4d9a-8c1a-067f99449e26', 'CTB-L2-016', 11, '列必须有注释', 2, 110, 'none', 'nil', 5, 'TableCreateColumnCommentRequired', '列"%s"需要提供COMMENT注解。', 'none', 1, 0, UNIX_TIMESTAMP()),
('95aa9816-ed7c-41a7-8f4f-6929757a92d8', 'CTB-L2-017', 11, '列允许的字符集', 2, 110, 'in', '["utf8mb4", "binary"]', 7, 'TableCreateColumnAvailableCharsets', '列"%s"禁用字符集"%s"，请使用"%s"。', 'checkboxes/key=CHARsets', 1, 0, UNIX_TIMESTAMP()),
('858cabba-39f1-4c47-af8c-85c891605400', 'CTB-L2-018', 11, '列允许的排序规则', 2, 110, 'in', '["utf8mb4_unicode_ci", "utf8mb4_general_ci", "utf8mb4_bin", "binary"]', 7, 'TableCreateColumnAvailableCollates', '列"%s"禁用排序规则"%s"，请使用"%s"。', 'checkboxes/key=collates', 1, 0, UNIX_TIMESTAMP()),
('55e72c3d-f0d2-45e9-bd3e-0c172ae6ad69', 'CTB-L2-019', 11, '列字符集与排序规则必须匹配', 2, 110, 'none', 'nil', 5, 'TableCreateColumnCharsetCollateMatch', '列"%s"使用的字符集"%s"和排序规则"%s"不匹配，请查阅官方文档。', 'none', 1, 0, UNIX_TIMESTAMP()),
('dba74873-c113-482a-8552-8d81dd640bd4', 'CTB-L2-020', 11, '非空列必须有默认值', 2, 110, 'none', 'nil', 5, 'TableCreateColumnNotNullWithDefaultRequired', '列"%s"不允许为空，但没有指定默认值。', 'none', 1, 0, UNIX_TIMESTAMP()),
('bc94e7f3-9e3e-48c7-b21b-85593b2c2af7', 'CTB-L2-021', 11, '自增列允许的数据类型', 2, 110, 'in', '["int", "bigint"]', 7, 'TableCreateColumnAutoIncAvailableTypes', '自增列"%s"禁用"%s"类型，请使用"%s"。', 'none', 1, 0, UNIX_TIMESTAMP()),
('d6894864-4d41-4930-b3fd-62da07588239', 'CTB-L2-022', 11, '自增列必须是无符号', 2, 110, 'none', 'nil', 7, 'TableCreateColumnAutoIncIsUnsigned', '自增列"%s"必须使用无符号的整数。', 'none', 1, 0, UNIX_TIMESTAMP()),
('526ac72c-55d4-484b-8bb2-88eda1e67d6c', 'CTB-L2-023', 11, '自增列必须是主键', 2, 110, 'none', 'nil', 5, 'TableCreateColumnAutoIncMustPrimaryKey', '自增列"%s"不是主键。', 'none', 1, 0, UNIX_TIMESTAMP()),
('84fcb061-af77-4a7f-8a0e-dfc166b2c531', 'CTB-L2-024', 11, '仅允许一个时间戳类型的列', 2, 110, 'none', 'nil', 7, 'TableCreateTimestampColumnCountLimit', '表"%s"中的定义了多个时间戳列，请改用DATETIME类型。', 'none', 1, 0, UNIX_TIMESTAMP()),
('0df017c0-d106-4df5-845c-bf50557d497f', 'CTB-L2-025', 11, '单一索引最大列数', 2, 110, 'lte', '3', 7, 'TableCreateIndexMaxColumnLimit', '索引"%s"索引的列数超出了规则允许的上限，请控制在%d个列以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('8d319516-890e-45f2-af90-8061b37e9a04', 'CTB-L2-026', 11, '必须有主键', 2, 110, 'none', 'nil', 5, 'TableCreatePrimaryKeyRequired', '必须为表指定一个主键。', 'none', 1, 0, UNIX_TIMESTAMP()),
('31e193df-bf11-4462-9558-c182eb8366e5', 'CTB-L2-027', 11, '主键是否显式命名', 2, 110, 'none', 'nil', 4, 'TableCreatePrimaryKeyNameExplicit', '主键没有提供名称。', 'none', 1, 0, UNIX_TIMESTAMP()),
('bedd8dfa-b325-451b-981a-a1cede6bc787', 'CTB-L2-028', 11, '主键名规则', 2, 110, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 4, 'TableCreatePrimaryKeyNameQualified', '主键名"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('dee29598-1f59-4d5f-9881-ad16fa4038fd', 'CTB-L2-029', 11, '主键名必须小写', 2, 110, 'regexp', '^[_a-z0-9]+$', 4, 'TableCreatePrimryKeyLowerCaseRequired', '主键名"%s"含有大写字母。', 'none', 1, 0, UNIX_TIMESTAMP()),
('c47768d2-24b3-4e9c-9079-8ef39bbd0a53', 'CTB-L2-030', 11, '主键名最大长度', 2, 110, 'lte', '20', 4, 'TableCreatePrimryKeyMaxLength', '主键名"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('bc1164e0-45ee-4d58-b4c5-cea244ebeb3f', 'CTB-L2-031', 11, '主键名前缀规则', 2, 110, 'regexp', '^pk_[_a-zA-Z0-9]+$', 4, 'TableCreatePrimryKeyPrefixRequired', '主键名"%s"需要满足前缀正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('7f7182f3-3da8-43c0-a43f-dbe44ccf0130', 'CTB-L2-032', 11, '索引必须命名', 2, 110, 'none', 'nil', 5, 'TableCreateIndexNameExplicit', '一个或多个索引没有提供索引名称。', 'none', 1, 0, UNIX_TIMESTAMP()),
('364f58af-29a8-4a49-a3ad-a35ea500352c', 'CTB-L2-033', 11, '索引名规则', 2, 110, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 5, 'TableCreateIndexNameQualified', '索引名"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('1dccb327-2b34-4f3d-a900-9560a7b9cd3b', 'CTB-L2-034', 11, '索引名必须小写', 2, 110, 'regexp', '^[_a-z0-9]+$', 5, 'TableCreateIndexNameLowerCaseRequired', '索引名"%s"含有大写字母。', 'none', 1, 0, UNIX_TIMESTAMP()),
('0e51b2a9-b251-4a25-be86-7895285ad77a', 'CTB-L2-035', 11, '索引名最大长度', 2, 110, 'lte', '10', 7, 'TableCreateIndexNameMaxLength', '索引名"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('abac99f7-ad0a-4d24-bdf4-00439e4abad6', 'CTB-L2-036', 11, '索引名前缀规则', 2, 110, 'regexp', '^index_[1-9][0-9]*$', 5, 'TableCreateIndexNamePrefixRequired', '索引名"%s"需要满足前缀正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('de701fa6-6ed4-4608-a95c-53087c43e177', 'CTB-L2-037', 11, '唯一索引必须命名', 2, 110, 'none', 'nil', 5, 'TableCreateUniqueNameExplicit', '一个或多个唯一索引没有提供索引名称。', 'none', 1, 0, UNIX_TIMESTAMP()),
('26d84e04-25e6-4865-b76b-00995ef56f81', 'CTB-L2-038', 11, '唯一索引索名规则', 2, 110, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 5, 'TableCreateUniqueNameQualified', '唯一索引"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('3b07760f-c055-4d3c-b27a-d23077a0c768', 'CTB-L2-039', 11, '唯一索引名必须小写', 2, 110, 'regexp', '^[_a-z0-9]+$', 5, 'TableCreateUniqueNameLowerCaseRequired', '唯一索引"%s"含有大写字母。', 'none', 1, 0, UNIX_TIMESTAMP()),
('f9a615b4-cbaa-40a7-b6c2-28a0fc78a35e', 'CTB-L2-040', 11, '唯一索引名最大长度', 2, 110, 'lte', '10', 7, 'TableCreateUniqueNameMaxLength', '唯一索引"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('8aab8c07-59d3-43d5-a74f-0005697955ee', 'CTB-L2-041', 11, '唯一索引名前缀规则', 2, 110, 'regexp', '^unique_[1-9][0-9]*$', 5, 'TableCreateUniqueNamePrefixRequired', '唯一索引"%s"需要满足前缀正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('30a7e9c3-0502-4793-9b88-fca8223efa90', 'CTB-L2-042', 11, '禁止外键', 2, 110, 'none', 'nil', 5, 'TableCreateForeignKeyEnabled', '禁止外键。', 'none', 1, 0, UNIX_TIMESTAMP()),
('8c599732-e3d0-432e-9396-4b4ed6b05479', 'CTB-L2-043', 11, '外键是否显式命名', 2, 110, 'none', 'nil', 5, 'TableCreateForeignKeyNameExplicit', '没有为外键指定名称。', 'none', 1, 0, UNIX_TIMESTAMP()),
('77708821-c945-4a6d-aa0d-84329b1b01fe', 'CTB-L2-044', 11, '外键名规则', 2, 110, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 5, 'TableCreateForeignKeyNameQualified', '外键名"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('867fa702-b10f-482d-b1a2-a71d5be32f21', 'CTB-L2-045', 11, '外键名必须小写', 2, 110, 'regexp', '^[_a-z0-9]+$', 5, 'TableCreateForeignKeyNameLowerCaseRequired', '外键名"%s"含有大写字母。', 'none', 1, 0, UNIX_TIMESTAMP()),
('3c3052b4-fc6d-4976-98d1-e2a49ecd3927', 'CTB-L2-046', 11, '外键名最大长度', 2, 110, 'lte', '25', 5, 'TableCreateForeignKeyNameMaxLength', '外键名"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('35207240-5094-4696-b209-1260865fde8a', 'CTB-L2-047', 11, '外键名前缀规则', 2, 110, 'regexp', '^fk_[_a-zA-Z0-9]+$', 5, 'TableCreateForeignKeyNamePrefixRequired', '外键名"%s"需要满足前缀正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('e5c85872-1cff-48c0-b705-0c61b04070b9', 'CTB-L2-048', 11, '表中最多可建多少个索引', 2, 110, 'lte', '5', 7, 'TableCreateIndexCountLimit', '表"%s"中定义了%d个索引，数量超过允许的阈值%d。', 'number', 1, 0, UNIX_TIMESTAMP()),
('7eedb2fb-5e26-41a7-baa7-8c37b8ef3170', 'CTB-L2-049', 11, '禁止使用CREATE TABLE ... LIKE ...建表', 2, 110, 'none', 'nil', 5, 'TableCreateUseLikeEnabled', '禁止使用CREATE TABLE LIKE的方式建表。', 'none', 1, 0, UNIX_TIMESTAMP()),
('5ace374b-0f35-4132-8941-44283bbc87c1', 'CTB-L2-050', 11, '仅允许定义一个自增列', 2, 110, 'nil', 'nil', 5, 'TableCreateAutoIncColumnCountLimit', '表"%s"中定义了多个自增列。', 'none', 1, 0, UNIX_TIMESTAMP()),
('b9dce6ba-5612-4bc4-afdb-730d7b15930a', 'CTB-L2-051', 11, '仅允许定义一个主键', 2, 110, 'nil', 'nil', 5, 'TableCreatePrimaryKeyCountLimit', '表"%s"中定义了多个主键。', 'none', 1, 0, UNIX_TIMESTAMP()),
('887cf290-2288-451a-be70-57394e23dde4', 'CTB-L3-001', 11, '目标库必须已存在', 3, 110, 'none', 'nil', 4, 'TableCreateTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('1e75c5a2-7bcc-42c4-ac3f-5a64ead5b25c', 'CTB-L3-002', 11, '目标表必须不存在', 3, 110, 'none', 'nil', 4, 'TableCreateTargetTableExists', '目标表"%s"已存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('b17d640a-16fd-4628-a845-d62cdb582884', 'MTB-L2-001', 12, '改表允许的字符集', 2, 120, 'in', '["utf8mb4"]', 7, 'TableAlterAvailableCharsets', '表禁用字符集"%s"，请使用"%s"。', 'checkboxes/key=CHARsets', 1, 0, UNIX_TIMESTAMP()),
('11460f45-3517-42d6-91ec-5b36d041cf4e', 'MTB-L2-002', 12, '改表允许的校验规则', 2, 120, 'in', '["utf8mb4_unicode_ci", "utf8mb4_general_ci", "utf8mb4_bin"]', 7, 'TableAlterAvailableCollates', '表禁用排序规则"%s"，请使用"%s"。', 'checkboxes/key=collates', 1, 0, UNIX_TIMESTAMP()),
('b4ae4473-9677-4db5-b2d4-7987fe904e88', 'MTB-L2-003', 12, '表的字符集与排序规则必须匹配', 2, 120, 'none', 'nil', 5, 'TableAlterCharsetCollateMatch', '表字符集"%s"和排序规则"%s"不匹配，请查阅官方文档。', 'none', 1, 0, UNIX_TIMESTAMP()),
('aee962c0-411e-4875-9291-d708473a1286', 'MTB-L2-004', 12, '改表允许的存储引擎', 2, 120, 'in', '["innodb", "tokudb", "rocksdb", "archive"]', 7, 'TableAlterAvailableEngines', '不支持的存储引擎"%s"，请使用"%s"。', 'checkboxes/key=engines', 1, 0, UNIX_TIMESTAMP()),
('07ca1452-fa25-4e33-b2dd-6623215a018e', 'MTB-L2-005', 12, '列名必须符合命名规范', 2, 120, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 5, 'TableAlterColumnNameQualified', '列名"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('8463928d-8d32-47fd-a882-efa93983ea3d', 'MTB-L2-006', 12, '列名必须小写', 2, 120, 'regexp', '^[_a-z0-9]+$', 7, 'TableAlterColumnNameLowerCaseRequired', '列名"%s"中含有除小写字母、数字和下划线以外的字符。', 'none', 1, 0, UNIX_TIMESTAMP()),
('d4a3a9c9-c16f-485a-a741-2c60ffda1b10', 'MTB-L2-007', 12, '列名最大长度', 2, 120, 'lte', '20', 7, 'TableAlterColumnNameMaxLength', '列名"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('8ceb2794-a706-406d-859d-4d7fe67ed970', 'MTB-L2-008', 12, '列禁用的数据类型', 2, 120, 'not-in', '["bit", "enum", "set", "double", "real", "float"]', 7, 'TableAlterColumnUnwantedTypes', '列"%s"使用了不期望的数据类型"%s"，请避免使用"%s"数据类型。', 'checkboxes/key=data-types', 1, 0, UNIX_TIMESTAMP()),
('f5b265d1-30c2-4c98-ab43-8876177c6ccb', 'MTB-L2-009', 12, '列必须有注释', 2, 120, 'none', 'nil', 5, 'TableAlterColumnCommentRequired', '列"%s"需要提供COMMENT注解。', 'none', 1, 0, UNIX_TIMESTAMP()),
('77130452-bf96-46f6-bb3e-8ce779d82dfc', 'MTB-L2-010', 12, '列允许的字符集', 2, 120, 'in', '["utf8mb4", "binary"]', 7, 'TableAlterColumnAvailableCharsets', '列"%s"禁用字符集"%s"，请使用"%s"。', 'checkboxes/key=CHARsets', 1, 0, UNIX_TIMESTAMP()),
('d6401716-b714-42a9-b06c-10f515f69cde', 'MTB-L2-011', 12, '列允许的排序规则', 2, 120, 'in', '["utf8mb4_unicode_ci", "utf8mb4_general_ci", "utf8mb4_bin", "binary"]', 7, 'TableAlterColumnAvailableCollates', '列"%s"禁用排序规则"%s"，请使用"%s"。', 'checkboxes/key=collates', 1, 0, UNIX_TIMESTAMP()),
('9c7f24ec-2d91-435f-b4ee-c31936670e6c', 'MTB-L2-012', 12, '列的字符集与排序规则必须匹配', 2, 120, 'none', 'nil', 5, 'TableAlterColumnCharsetCollateMatch', '列"%s"使用的字符集"%s"和排序规则"%s"不匹配，请查阅官方文档。', 'none', 1, 0, UNIX_TIMESTAMP()),
('76b8f1fd-7ebb-41fc-9b57-763e2590121d', 'MTB-L2-013', 12, '非空列必须有默认值', 2, 120, 'none', 'nil', 5, 'TableAlterColumnNotNullWithDefaultRequired', '列"%s"不允许为空，但没有指定默认值。', 'none', 1, 0, UNIX_TIMESTAMP()),
('90c0dece-d69c-4c56-a1a4-7824fe00fe22', 'MTB-L2-014', 12, '索引必须命名', 2, 120, 'none', 'nil', 7, 'TableAlterIndexNameExplicit', '一个或多个索引没有提供索引名称。', 'none', 1, 0, UNIX_TIMESTAMP()),
('a28f4d53-4c81-4b73-a8a9-2bb441bbfb4b', 'MTB-L2-015', 12, '索引名标识符必须满足规则', 2, 120, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 7, 'TableAlterIndexNameQualified', '索引名"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('b462710c-af50-4c96-85de-e3e6ddb20f47', 'MTB-L2-016', 12, '索引名必须小写', 2, 120, 'regexp', '^[_a-z0-9]+$', 7, 'TableAlterIndexNameLowerCaseRequired', '索引名"%s"含有除小写字母、数字和下划线以外的字符。', 'none', 1, 0, UNIX_TIMESTAMP()),
('543b5732-8a6d-4690-a02b-f3c8055435c5', 'MTB-L2-017', 12, '索引名最大长度', 2, 120, 'lte', '10', 7, 'TableAlterIndexNameMaxLength', '索引名"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('9f752e8b-df62-46b0-bae4-e2c6941f3d88', 'MTB-L2-018', 12, '索引名前缀规则', 2, 120, 'regexp', '^index_[1-9][0-9]*$', 7, 'TableAlterIndexNamePrefixRequired', '索引名"%s"需要满足前缀正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('1885389b-4a1b-475b-a37b-582670a908f8', 'MTB-L2-019', 12, '唯一索引必须命名', 2, 120, 'none', 'nil', 7, 'TableAlterUniqueNameExplicit', '一个或多个唯一索引没有提供索引名称。', 'none', 1, 0, UNIX_TIMESTAMP()),
('66dd81be-fd2c-4304-b782-86006c52fc5d', 'MTB-L2-020', 12, '唯一索引索名标识符必须符合规则', 2, 120, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 7, 'TableAlterUniqueNameQualified', '唯一索引"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('b68ffe6e-8482-475b-afa4-90b26b44aeb4', 'MTB-L2-021', 12, '唯一索引名必须小写', 2, 120, 'regexp', '^[_a-z0-9]+$', 7, 'TableAlterUniqueNameLowerCaseRequired', '唯一索引"%s"含有除小写字母、数字和下划线以外的字符。', 'none', 1, 0, UNIX_TIMESTAMP()),
('9dbce3a5-46b7-4465-af71-e5b6291ebc7a', 'MTB-L2-022', 12, '唯一索引名不能超过最大长度', 2, 120, 'lte', '10', 7, 'TableAlterUniqueNameMaxLength', '唯一索引"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('7d77c21e-e0c8-48c0-a115-fc460e820c2b', 'MTB-L2-023', 12, '唯一索引名前缀必须符合规则', 2, 120, 'regexp', '^unique_[1-9][0-9]*$', 7, 'TableAlterUniqueNamePrefixRequired', '唯一索引"%s"需要满足前缀正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('690b58b4-9be8-4560-ade6-be534a34120b', 'MTB-L2-024', 12, '禁止外键', 2, 120, 'none', 'nil', 7, 'TableAlterForeignKeyEnabled', '禁止外键。', 'none', 1, 0, UNIX_TIMESTAMP()),
('98436226-6cb8-4776-bc0b-6bb53ac8cfb4', 'MTB-L2-025', 12, '外键是否显式命名', 2, 120, 'none', 'nil', 5, 'TableAlterForeignKeyNameExplicit', '没有为外键指定名称。', 'none', 1, 0, UNIX_TIMESTAMP()),
('8ee524ec-2359-4a40-ad91-6c3716dc07b6', 'MTB-L2-026', 12, '外键名标识符规则', 2, 120, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 5, 'TableAlterForeignKeyNameQualified', '外键名"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('b107886d-b28d-41c9-9268-81b6fd631c1a', 'MTB-L2-027', 12, '外键名必须小写', 2, 120, 'regexp', '^[_a-z0-9]+$', 5, 'TableAlterForeignKeyNameLowerCaseRequired', '外键名"%s"含有除小写字母、数字和下划线以外的字符。', 'none', 1, 0, UNIX_TIMESTAMP()),
('f3bfe2da-162e-46db-8be5-6e152a16696b', 'MTB-L2-028', 12, '外键名最大长度', 2, 120, 'lte', '25', 5, 'TableAlterForeignKeyNameMaxLength', '外键名"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('9df7d265-a301-4fa4-9933-87fe3f0fcae2', 'MTB-L2-029', 12, '外键名前缀规则', 2, 120, 'regexp', '^fk_[_a-zA-Z0-9]+$', 5, 'TableAlterForeignKeyNamePrefixRequired', '外键名"%s"需要满足前缀正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('7024faae-3902-4af1-a849-55b09a7cd85a', 'MTB-L2-030', 12, '更名新表规则', 2, 120, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 7, 'TableAlterNewTableNameQualified', '目标表"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('92ee30e7-d138-463e-873f-899fb2be4375', 'MTB-L2-031', 12, '更名新表必须小写', 2, 120, 'regexp', '^[_a-z0-9]+$', 7, 'TableAlterNewTableNameLowerCaseRequired', '目标表"%s"含有除小写字母、数字和下划线以外的字符。', 'none', 1, 0, UNIX_TIMESTAMP()),
('106db8a3-da54-423c-8104-ed92625f80f4', 'MTB-L2-032', 12, '更名新表最大长度', 2, 120, 'lte', '20', 7, 'TableAlterNewTableNameMaxLength', '目标表"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('b9a1a81d-c328-48b9-8919-d47f14f50906', 'MTB-L2-033', 12, '禁用全文索引', 2, 120, 'none', 'nil', 7, 'TableAlterFullTextEnabled', '禁止使用全文索引。', 'none', 1, 0, UNIX_TIMESTAMP()),
('40375d77-dea4-410b-b8de-9fe4960a138d', 'MTB-L2-034', 12, '索引必须命名', 2, 120, 'none', 'nil', 7, 'TableAlterFullTextNameExplicit', '一个或多个全文索引没有提供索引名称。', 'none', 1, 0, UNIX_TIMESTAMP()),
('36d7aa62-fa25-4aa5-a7fa-afb15479b037', 'MTB-L2-035', 12, '索引名标识符规则', 2, 120, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 7, 'TableAlterFullTextNameQualified', '全文索引"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('86dd974c-f6ca-475a-84c8-981d43f489ff', 'MTB-L2-036', 12, '索引名必须小写', 2, 120, 'regexp', '^[_a-z0-9]+$', 7, 'TableAlterFullTextNameLowerCaseRequired', '全文索引"%s"含有除小写字母、数字和下划线以外的字符。', 'none', 1, 0, UNIX_TIMESTAMP()),
('b928f242-aacd-41b9-8347-aa80ab679746', 'MTB-L2-037', 12, '索引名不能超过最大长度', 2, 120, 'lte', '10', 7, 'TableAlterFullTextNameMaxLength', '全文索引"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('f754c334-32cf-4c4c-9c46-92a666d24192', 'MTB-L2-038', 12, '索引名前缀必须匹配规则', 2, 120, 'regexp', '^ft_[1-9][0-9]*$', 7, 'TableAlterFullTextNamePrefixRequired', '全文索引"%s"需要满足前缀正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('9280ff7c-ff2b-4c3e-bf6c-4886d49634d1', 'MTB-L2-039', 12, '单一索引最大列数', 2, 120, 'lte', '3', 7, 'TableAlterIndexMaxColumnLimit', '索引"%s"索引的列数超出了规则允许的上限，请控制在%d个列以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('768fb105-d609-4b7b-8ff9-0d3854cabfff', 'MTB-L3-001', 12, '目标库必须已存在', 3, 120, 'none', 'nil', 4, 'TableAlterEnabled', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('aa482f88-074c-45ab-97b9-81509bd4567f', 'MTB-L3-002', 12, '目标表必须已存在', 3, 120, 'none', 'nil', 4, 'TableAlterEnabled', '目标表"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('1651eb45-bc11-4aa9-81dc-ee34d4f0f861', 'MTB-L3-003', 12, '新增列时目标列必须不存在', 3, 120, 'none', 'nil', 4, 'TableAlterEnabled', '目标列"%s"已存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('cfe017e3-339a-45a5-b821-825d656b85e8', 'MTB-L3-004', 12, '位置标记列必须已存在', 3, 120, 'none', 'nil', 4, 'TableAlterEnabled', '位置标记列"%s"(BEFORE/AFTER)不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('a130f2c4-e6c0-4dbd-9e55-51af3cb2ecab', 'MTB-L3-005', 12, '列名是否重复', 2, 120, 'none', 'nil', 4, 'TableAlterColumnNameDuplicate', '表"%s"中的定义了重复的列"%s"。', 'none', 1, 0, UNIX_TIMESTAMP()),
('6677109a-5798-44b9-9f8e-597cce76f168', 'MTB-L3-006', 12, '表允许的最大列数', 2, 120, 'lte', '25', 6, 'TableAlterColumnCountLimit', '表"%s"中定义%d个列，数量超出了规则允许的上限%d，请考虑拆分表。', 'number', 1, 0, UNIX_TIMESTAMP()),
('5e179511-ae0a-4e55-abf6-23120adda97f', 'MTB-L3-007', 12, '仅允许一个时间戳类型的列', 2, 120, 'none', 'nil', 6, 'TableAlterTimestampColumnCountLimit', '表"%s"中的定义了多个时间戳列，请改用DATETIME类型。', 'none', 1, 0, UNIX_TIMESTAMP()),
('084397d9-86fd-415f-b37e-eaa44ba47be9', 'MTB-L3-008', 12, '删除列时目标列必须已存在', 3, 120, 'none', 'nil', 4, 'TableAlterEnabled', '目标列"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('6a336124-aea6-455e-bc67-a7eb685e80a8', 'MTB-L3-009', 12, '修改列时目标列必须已存在', 3, 120, 'none', 'nil', 4, 'TableAlterEnabled', '目标列"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('8541b5b4-66ab-4aed-a254-73ce37194b34', 'MTB-L3-010', 12, '修改列时目标列必须不存在', 3, 120, 'none', 'nil', 4, 'TableAlterEnabled', '目标列"%s"已存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('8600038b-f6ad-444c-a333-02eda47d75a5', 'MTB-L3-011', 12, '仅允许一个时间戳类型的列', 3, 120, 'none', 'nil', 4, 'TableAlterEnabled', '表"%s"中timestamp类型将超过最大限制"%d"。', 'none', 1, 0, UNIX_TIMESTAMP()),
('49251321-a0ae-428b-a743-026dc548a4ff', 'MTB-L3-012', 12, '列名是否重复', 2, 120, 'none', 'nil', 4, 'TableAlterColumnNameDuplicateForChg', '表"%s"中的定义了重复的列"%s"。', 'none', 1, 0, UNIX_TIMESTAMP()),
('b4fd2895-118e-4e80-b792-da47fdedcb20', 'MTB-L3-013', 12, '添加索引时索引必须不存在', 3, 120, 'none', 'nil', 4, 'TableAlterEnabled', '目标表"%s"上已存在索引"%s"。', 'none', 1, 0, UNIX_TIMESTAMP()),
('08fb5c7f-f4f9-490f-a3a6-a8c4c82bd1a1', 'MTB-L3-014', 12, '覆盖索引检查', 3, 120, 'none', 'nil', 4, 'TableAlterEnabled', '目标表"%s"上已存在索引"%s"。', 'none', 1, 0, UNIX_TIMESTAMP()),
('c49c27b8-5d1d-4675-a734-8670eee65d9b', 'MTB-L3-015', 12, '同名外键检查', 3, 120, 'none', 'nil', 4, 'TableAlterEnabled', '目标表"%s"上已存在外键"%s"。', 'none', 1, 0, UNIX_TIMESTAMP()),
('afcc5ed3-1e4c-47df-b837-b6f3e5dee543', 'MTB-L3-016', 12, '添加外键时外键必须不存在', 3, 120, 'none', 'nil', 4, 'TableAlterEnabled', '目标表"%s"上已存在外键"%s"。', 'none', 1, 0, UNIX_TIMESTAMP()),
('a21ea970-7e26-420d-95f3-e1621f9c5630', 'MTB-L3-017', 12, '启用禁用KEY时KEY必须已存在', 3, 120, 'none', 'nil', 4, 'TableAlterEnabled', '目标KEY"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('061df5d6-cd83-4b1a-abbe-9af10307e9dd', 'MTB-L3-018', 12, '删主键时主键必须存在', 3, 120, 'none', 'nil', 4, 'TableAlterEnabled', '目标表"%s"上未定义主键。', 'none', 1, 0, UNIX_TIMESTAMP()),
('964b019b-771d-4983-b7fa-25b089581323', 'MTB-L3-019', 12, '删索引时索引必须存在', 3, 120, 'none', 'nil', 4, 'TableAlterEnabled', '目标表"%s"上未定义索引"%s"。', 'none', 1, 0, UNIX_TIMESTAMP()),
('7012f584-f7d5-492e-8cd4-37834845558c', 'MTB-L3-020', 12, '删外键时外键必须存在', 3, 120, 'none', 'nil', 4, 'TableAlterEnabled', '目标表"%s"上未定义外键"%s"。', 'none', 1, 0, UNIX_TIMESTAMP()),
('2f2af8e5-2cf1-4ced-a1ec-a47c7c744178', 'MTB-L3-021', 12, '改名时目标表已存在', 3, 120, 'none', 'nil', 4, 'TableAlterEnabled', '目标表"%s"已存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('86df78a8-ce02-4b06-81c0-f56a7d99ae09', 'MTB-L3-022', 12, '全文索引必须不存在', 3, 120, 'none', 'nil', 4, 'TableAlterEnabled', '目标表"%s"上已存在全文索引"%s"。', 'none', 1, 0, UNIX_TIMESTAMP()),
('c74c9e43-4600-42f5-8711-9b2114917438', 'MTB-L3-023', 12, '删全文索引时索引必须存在', 3, 120, 'none', 'nil', 4, 'TableAlterEnabled', '目标表"%s"上未定义全文索引"%s"。', 'none', 1, 0, UNIX_TIMESTAMP()),
('a9aa23c7-f79e-4a79-b5b8-520341992474', 'MTB-L3-024', 12, '合并改表语句', 2, 120, 'none', 'nil', 4, 'TableAlterEnabled', '对同一个表"%s"的多条修改语句需要合并成一条语句。', 'none', 1, 0, UNIX_TIMESTAMP()),
('3ec2eaff-5471-405f-ab34-bc2158b8f884', 'RTB-L2-001', 13, '目标表跟源表是同一个表', 3, 130, 'none', 'nil', 5, 'TableRenameTablesIdentical', '源表"%s"和目标表"%s"相同。', 'none', 1, 0, UNIX_TIMESTAMP()),
('95c374ee-52a4-48f2-b24a-ac273e5e8f36', 'RTB-L2-002', 13, '目标表名规则', 2, 130, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 5, 'TableRenameTargetTableNameQualified', '目标表名"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('7daae444-78ea-46ac-a7e9-18149c725545', 'RTB-L2-003', 13, '目标表名必须小写', 2, 130, 'none', '^[_a-z0-9]+$', 5, 'TableRenameTargetTableNameLowerCaseRequired', '目标表名"%s"含有大写字母。', 'none', 1, 0, UNIX_TIMESTAMP()),
('22e5db11-6afb-4e50-9fa1-c77c09dfb168', 'RTB-L2-004', 13, '目标表名最大长度', 2, 130, 'lte', '20', 7, 'TableRenameTargetTableNameMaxLength', '目标表名"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('8e75b778-e6c5-4d34-9edd-03097a4450d2', 'RTB-L3-001', 13, '源库必须已存在', 3, 130, 'none', 'nil', 4, 'TableRenameSourceTableExists', '源库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('8a1a4f9b-35da-4ff0-9ee2-d190b54f2ba5', 'RTB-L3-002', 13, '源表必须已存在', 3, 130, 'none', 'nil', 4, 'TableRenameSourceDatabaseExists', '源表"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('bb56fcfd-736f-4033-aa6f-2ce9ea1108f8', 'RTB-L3-003', 13, '目标库必须已存在', 3, 130, 'none', 'nil', 4, 'TableRenameTargetTableExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('2affbef0-4856-4eb3-a9a7-a1bf74a531c2', 'RTB-L3-004', 13, '目标表必须不存在', 3, 130, 'none', 'nil', 4, 'TableRenameTargetDatabaseExists', '目标表"%s"已存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('da2d4fd4-af73-4cbf-9b83-ca32d5594396', 'DTB-L3-001', 14, '目标库必须已存在', 3, 140, 'none', 'nil', 4, 'TableDropSourceDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('a2d5b7ec-0fda-42f9-9fea-71687a3066f7', 'DTB-L3-002', 14, '目标表必须已存在', 3, 140, 'none', 'nil', 4, 'TableDropSourceTableExists', '目标表"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('2221dc67-94c6-4a3d-a3b5-638a3d02b70a', 'CIX-L2-001', 15, '组合索引允许的最大列数', 2, 150, 'lte', '3', 7, 'IndexCreateIndexMaxColumnLimit', '索引"%s"中索引列数量超过允许的阈值%d。', 'number', 1, 0, UNIX_TIMESTAMP()),
('0076aae4-6196-4f0d-bfb3-37df11fd45ed', 'CIX-L2-002', 15, '索引名规则', 2, 150, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 5, 'IndexCreateIndexNameQualified', '索引名"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('a316aa97-534b-4c9d-a881-278b9b7bdbdf', 'CIX-L2-003', 15, '索引名必须小写', 2, 150, 'regexp', '^[_a-z0-9]+$', 5, 'IndexCreateIndexNameLowerCaseRequired', '索引名"%s"含有大写字母。', 'none', 1, 0, UNIX_TIMESTAMP()),
('71db9655-eeb1-4be3-8c0a-67cef0709b38', 'CIX-L2-004', 15, '索引名最大长度', 2, 150, 'lte', '10', 5, 'IndexCreateIndexNameMaxLength', '索引名"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('6d676a9c-90a5-4c9c-81e6-00ec8e5ab350', 'CIX-L2-005', 15, '索引名前缀规则', 2, 150, 'regexp', '^index_[1-9][0-9]*$', 5, 'IndexCreateIndexNamePrefixRequired', '索引名"%s"需要满足前缀正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('47ee8f95-9d48-4a4c-a917-c4679d5fe4f0', 'CIX-L2-006', 15, '组合索引中是否有重复列', 2, 150, 'none', 'nil', 5, 'IndexCreateDuplicateIndexColumn', '索引"%s"中索引了重复的列。', 'none', 1, 0, UNIX_TIMESTAMP()),
('034b0bcd-e246-4ade-9ab0-35522b5382f9', 'CIX-L3-001', 15, '目标库必须已存在', 3, 150, 'none', 'nil', 4, 'IndexCreateTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('2d222b20-a4a5-44c4-b7aa-4cda28a46369', 'CIX-L3-002', 15, '目标表必须已存在', 3, 150, 'none', 'nil', 4, 'IndexCreateTargetTableExists', '目标表"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('22b1f933-bce4-4d20-9b6f-7a0282c453c1', 'CIX-L3-003', 15, '索引列必须已存在', 3, 150, 'none', 'nil', 4, 'IndexCreateTargetColumnExists', '目标索引"%s"已存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('2e301212-02af-4b84-955f-3178812ded5a', 'CIX-L3-004', 15, '索引内容是否重复', 3, 150, 'none', 'nil', 4, 'IndexCreateTargetIndexExists', '索引"%s"在已有索引"%s"相同或者存在覆盖关系。', 'none', 1, 0, UNIX_TIMESTAMP()),
('57f92fe9-3540-4541-a53c-0487db7c0702', 'CIX-L3-005', 15, '索引名是否重复', 3, 150, 'none', 'nil', 4, 'IndexCreateTargetNameExists', '索引名"%s"在表"%s"已经存在，请使用另外一个索引名称。', 'none', 1, 0, UNIX_TIMESTAMP()),
('6f69fd99-4d89-4c6f-b8ce-1b9dd437e687', 'CIX-L3-006', 15, '最多能建多少个索引', 3, 150, 'lte', '5', 6, 'IndexCreateIndexCountLimit', '索引数量超过允许的阈值%d。', 'number', 1, 0, UNIX_TIMESTAMP()),
('9525f4ff-efa2-4b09-a5af-ab3cc3deb24c', 'CIX-L3-007', 15, '禁止在BLOB/TEXT列上建索引', 2, 150, 'none', 'nil', 4, 'IndexCreateIndexBlobColumnEnabled', '禁止在BLOB/TEXT类型的列上建立索引。', 'none', 1, 0, UNIX_TIMESTAMP()),
('43181c58-4cba-4209-99ad-4534c38e455a', 'RIX-L3-001', 15, '目标库必须已存在', 1, 151, 'none', 'nil', 4, 'IndexDropTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('342f0510-6e41-4b63-871e-b0c7ae95715e', 'RIX-L3-002', 15, '目标表必须已存在', 1, 151, 'none', 'nil', 4, 'IndexDropTargetTableExists', '目标表"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('6df03432-6248-4ff6-af35-04bd678e8812', 'RIX-L3-003', 15, '目标索引必须已存在', 1, 151, 'none', 'nil', 4, 'IndexDropTargetIndexExists', '目标索引"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('ed326131-edce-456a-a19e-1fdece47060e', 'DEL-L2-001', 16, '禁止没有WHERE的删除', 2, 163, 'none', 'nil', 5, 'DeleteWithoutWhereEnabled', '禁止没有WHERE从句的DELETE语句。', 'none', 1, 0, UNIX_TIMESTAMP()),
('de3855f6-b835-49eb-89ff-bf854c515a02', 'DEL-L3-001', 16, '单次删除的最大行数', 3, 163, 'lte', '1000', 6, 'DeleteRowsLimit', '单条DELETE语句不得操作超过%d条记录。', 'number', 1, 0, UNIX_TIMESTAMP()),
('7ebccfe3-6ad7-4c83-b7a4-9cbfa12e35db', 'DEL-L3-002', 16, '目标库必须已存在', 3, 163, 'none', 'nil', 4, 'DeleteTargetDatabaseExists', 'DELETE语句中指定的库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('0b5dbb5e-8223-42a0-bb76-a31ed16bf537', 'DEL-L3-003', 16, '目标表必须已存在', 3, 163, 'none', 'nil', 4, 'DeleteTargetTableExists', 'DELETE语句中指定的表"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('36b9693f-26e1-47b0-8ebb-e26f6ae9d5fa', 'DEL-L3-004', 16, '条件过滤列必须已存在', 3, 163, 'none', 'nil', 4, 'DeleteFilterColumnExists', 'DELETE语句中条件限定列"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('4a5c2e53-1b50-43ce-954b-ce2d60eeeed7', 'INS-L2-001', 16, 'INSERT时强制显式列申明', 2, 160, 'none', 'nil', 5, 'InsertExplicitColumnRequired', '禁止没有显式提供列列表的INSERT语句。', 'none', 1, 0, UNIX_TIMESTAMP()),
('6454269e-fff7-4f75-857f-f06bbf1f3a73', 'INS-L2-002', 16, '禁止INSERT...SELECT', 2, 160, 'none', 'nil', 7, 'InsertUsingSelectEnabled', '禁止INSERT ... SELECT ...语句。', 'none', 1, 0, UNIX_TIMESTAMP()),
('16d0a0b4-3525-449c-965a-673bd290efc7', 'INS-L2-004', 16, 'INSERT时单语句允许操作的最大行数', 2, 160, 'lte', '10000', 7, 'InsertRowsLimit', '单条INSERT语句不得操作超过%d条记录。', 'number', 1, 0, UNIX_TIMESTAMP()),
('ab221728-8c92-4019-92e4-4cf0fc57941b', 'INS-L2-005', 16, 'INSERT时列类型、值是否匹配', 2, 160, 'none', 'nil', 5, 'InsertColumnValueMatch', 'INSERT语句的列数量和值数量不匹配。', 'none', 1, 0, UNIX_TIMESTAMP()),
('d37860b1-d9d6-4c4f-9af6-7e93dd33c312', 'INS-L3-001', 16, 'INSERT时目标库必须已存在', 3, 160, 'none', 'nil', 4, 'InsertTargetDatabaseExists', 'INSERT语句中指定的库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('75d3d38e-b32b-438c-b01c-761fff3b6786', 'INS-L3-002', 16, 'INSERT时目标表必须已存在', 3, 160, 'none', 'nil', 4, 'InsertTargetTableExists', 'INSERT语句中指定的表"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('355540e6-733a-4b83-88e7-e04e838fa8c2', 'INS-L3-003', 16, 'INSERT时目标列必须已存在', 3, 160, 'none', 'nil', 4, 'InsertTargetColumnExists', 'INSERT语句中插入的列"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('9f6cae35-89b7-46ae-ae5c-a7962f473e44', 'INS-L3-004', 16, 'INSERT时非空列是否有值', 3, 160, 'none', 'nil', 4, 'InsertValueForNotNullColumnRequired', 'INSERT语句没有为非空列"%s"提供值。', 'none', 1, 0, UNIX_TIMESTAMP()),
('f3578823-a5c8-4cdf-be09-55fc899d6875', 'RPL-L2-001', 16, 'REPLACE时强制显式列申明', 2, 161, 'none', 'nil', 5, 'ReplaceExplicitColumnRequired', '禁止没有显式提供列列表的REPLACE语句。', 'none', 1, 0, UNIX_TIMESTAMP()),
('41b437f0-eca7-46c3-9935-69c2bcacae70', 'RPL-L2-002', 16, '禁止REPLACE...SELECT', 2, 161, 'none', 'nil', 5, 'ReplaceUsingSelectEnabled', '禁止REPLACE ... SELECT ...语句。', 'none', 1, 0, UNIX_TIMESTAMP()),
('3901b27f-3fd6-429a-92f3-cbfb327aec6e', 'RPL-L2-004', 16, 'REPLACE时单语句允许操作的最大行数', 2, 161, 'lte', '1000', 7, 'ReplaceRowsLimit', '单条REPLACE语句不得操作超过%d条记录。', 'number', 1, 0, UNIX_TIMESTAMP()),
('eaaa918c-611f-467c-9c65-03ff4124a998', 'RPL-L2-005', 16, 'REPLACE时列类型、值是否匹配', 2, 161, 'none', 'nil', 5, 'ReplaceColumnValueMatch', 'REPLACE语句的列数量和值数量不匹配。', 'none', 1, 0, UNIX_TIMESTAMP()),
('5679453f-d160-4f0a-860f-5f6f5c50d505', 'RPL-L3-001', 16, 'REPLACE时目标库必须已存在', 3, 161, 'none', 'nil', 4, 'ReplaceTargetDatabaseExists', 'REPLACE语句中指定的库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('612b1902-88c1-46e2-b139-5fec2b74c692', 'RPL-L3-002', 16, 'REPLACE时目标表必须已存在', 3, 161, 'none', 'nil', 4, 'ReplaceTargetTableExists', 'REPLACE语句中指定的表"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('d54e672f-e360-4789-bd2d-5ca5f14fdf5a', 'RPL-L3-003', 16, 'REPLACE时目标列必须已存在', 3, 161, 'none', 'nil', 4, 'ReplaceTargetColumnExists', 'REPLACE语句中替换的列"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('68b029cc-8fff-4da7-8a3e-62b124b779a3', 'RPL-L3-004', 16, 'REPLACE时非空列是否有值', 3, 161, 'none', 'nil', 4, 'ReplaceValueForNotNullColumnRequired', 'REPLACE语句没有为非空列"%s"提供值。', 'none', 1, 0, UNIX_TIMESTAMP()),
('67fa90a8-5df6-479c-bfd4-11ab0bb9fa8b', 'UPD-L2-001', 16, '禁止没有WHERE的更新', 3, 162, 'none', 'nil', 5, 'UpdateWithoutWhereEnabled', '禁止没有WHERE从句的UPDATE语句。', 'none', 1, 0, UNIX_TIMESTAMP()),
('022dbea9-af6f-44cb-b0ae-3f4712a014c8', 'UPD-L3-001', 16, '目标库必须已存在', 3, 162, 'none', 'nil', 4, 'UpdateTargetDatabaseExists', 'UPDATE语句中指定的库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('749e0a05-843b-4c8e-8d5c-6ef252b261ad', 'UPD-L3-002', 16, '目标表必须已存在', 3, 162, 'none', 'nil', 4, 'UpdateTargetTableExists', 'UPDATE语句中指定的表"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('bdf11d37-9534-4553-826d-15d3b08f1059', 'UPD-L3-003', 16, '目标列必须已存在', 3, 162, 'none', 'nil', 4, 'UpdateTargetColumnExists', 'UPDATE语句中更新的列"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('a75a4aed-8dde-467c-b85b-5b962a8afa70', 'UPD-L3-004', 16, '条件限定列必须已存在', 3, 162, 'none', 'nil', 4, 'UpdateFilterColumnExists', 'UPDATE语句中条件限定列"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('a72f3d07-7512-46bb-a442-8a505610a72e', 'UPD-L3-005', 16, '允许单次更新的最大行数', 3, 162, 'lte', '1000', 6, 'UpdateRowsLimit', '单条UPDATE语句不得操作超过%d条记录。', 'number', 1, 0, UNIX_TIMESTAMP()),
('1b163a6b-633b-4285-978f-fefd0af3fb9c', 'SEL-L2-001', 17, '禁止没有WHERE的查询', 2, 170, 'none', 'nil', 5, 'SelectWithoutWhereEnabled', '禁止没有WHERE从句的查询语句。', 'none', 1, 0, UNIX_TIMESTAMP()),
('f9d57f6c-2ad0-4bb6-a4f8-a56fbff3cb65', 'SEL-L2-002', 17, '禁止没有LIMIT的查询', 2, 170, 'none', 'nil', 7, 'SelectWithoutLimitEnabled', '禁止没有LIMIT从句的查询语句。', 'none', 1, 0, UNIX_TIMESTAMP()),
('b44f75af-51f0-48f5-8d1c-b8f8b1eb225d', 'SEL-L2-003', 17, '禁止SELECT STAR', 2, 170, 'none', 'nil', 7, 'SelectStarEnabled', '禁止SELECT *操作，需要显式指定需要查询的列。', 'none', 1, 0, UNIX_TIMESTAMP()),
('e1018221-684c-49f0-baf4-2ea77c80fd3d', 'SEL-L2-004', 17, '禁止SELECT FOR UPDATE', 2, 170, 'none', 'nil', 5, 'SelectForUpdateEnabled', '禁止SELECT FOR UPDATE操作。', 'none', 1, 0, UNIX_TIMESTAMP()),
('bbd20e8b-4eeb-4467-a481-75bde843913f', 'SEL-L3-001', 17, '目标数据库必须已存在', 3, 170, 'none', 'nil', 4, 'SelectTargetDatabaseExists', 'SELECT语句中指定的库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('2430a63b-67c1-48f6-8b34-a4f1161a63b1', 'SEL-L3-002', 17, '目标表必须已存在', 3, 170, 'none', 'nil', 4, 'SelectTargetTableExists', 'SELECT语句中指定的表"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('90f2b354-19c5-4d4d-be3f-8361ecbe968d', 'SEL-L3-003', 17, '目标列必须已存在', 3, 170, 'none', 'nil', 4, 'SelectTargetColumnExists', 'SELECT语句中返回的列"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('d59e3dad-cf02-41b6-a306-987e1e07814e', 'SEL-L3-004', 17, '是否允许返回BLOB/TEXT列', 3, 170, 'none', 'nil', 4, 'SelectBlobOrTextEnabled', '查询语句中指的列"%s"是BLOB/TEXT类型。', 'none', 1, 0, UNIX_TIMESTAMP()),
('83120116-cd38-4872-935a-7164fc97e68e', 'SEL-L3-005', 17, '条件过滤列必须已存在', 3, 170, 'none', 'nil', 4, 'SelectFilterColumnExists', 'SELECT语句中条件限定列"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('217771f6-c3de-4956-8656-2fd8e482f5b7', 'CVW-L2-001', 18, '新建视图时视图名规则', 2, 180, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 5, 'ViewCreateViewNameQualified', '视图名"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('8d9a1c43-325a-47ac-9891-c1155697da3a', 'CVW-L2-002', 18, '新建视图时视图名必须小写', 2, 180, 'regexp', '^[_a-z0-9]+$', 5, 'ViewCreateViewNameLowerCaseRequired', '视图名"%s"含有大写字母。', 'none', 1, 0, UNIX_TIMESTAMP()),
('19f69c74-7f29-4eba-8131-dbc5992ae3c4', 'CVW-L2-003', 18, '新建视图时视图名最大长度', 2, 180, 'lte', '25', 7, 'ViewCreateViewNameMaxLength', '视图名"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('a776c690-8576-4664-ab29-5db0885960b3', 'CVW-L2-004', 18, '新建视图时视图名前缀规则', 2, 180, 'regexp', '^vw_[_a-zA-Z0-9]+$', 7, 'ViewCreateViewNamePrefixRequired', '视图名"%s"需要满足前缀正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('dff7c8a6-f9db-4d86-9c7f-df8b1614a307', 'CVW-L3-001', 18, '新建视图时目标库必须已存在', 3, 180, 'none', 'nil', 4, 'ViewCreateTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('5fba8576-d5fb-4135-9dc2-d01578ebfb9c', 'CVW-L3-002', 18, '新建视图时目标视图必须不存在', 3, 180, 'none', 'nil', 4, 'ViewCreateTargetViewExists', '目标视图"%s"已存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('4b1e7450-4662-4d47-948f-7e6d61f9d1f3', 'DVW-L3-001', 18, '删除视图时目标库必须已存在', 3, 182, 'none', 'nil', 4, 'ViewDropTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('fe4af9de-5292-45a5-82e3-16348b1a9da1', 'DVW-L3-002', 18, '删除视图时目标视图必须已存在', 3, 182, 'none', 'nil', 4, 'ViewDropTargetViewExists', '目标视图"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('2da54955-a029-4ae5-8535-f77bb461973d', 'MVW-L3-001', 18, '修改视图时目标库必须已存在', 3, 181, 'none', 'nil', 4, 'ViewAlterTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('50edc670-7d6e-4c6f-92f0-24539074404b', 'MVW-L3-002', 18, '修改视图时目标视图必须已存在', 3, 181, 'none', 'nil', 4, 'ViewAlterTargetViewExists', '目标视图"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('3a2fc192-febc-4ba5-84d4-ed95b963fdfc', 'CFU-L2-001', 19, '新建函数时函数名规则', 2, 190, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 7, 'FuncCreateFuncNameQuilified', '函数名"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('33004fce-753f-4a77-b552-e0cf7bbe4636', 'CFU-L2-002', 19, '新建函数时函数名必须小写', 2, 190, 'regexp', '^[_a-z0-9]+$', 7, 'FuncCreateFuncNameLowerCaseRequired', '函数名"%s"含有大写字母。', 'none', 1, 0, UNIX_TIMESTAMP()),
('a924afeb-8def-4443-ab2a-d8d429e5dd49', 'CFU-L2-003', 19, '新建函数时函数名最大长度', 2, 190, 'lte', '25', 7, 'FuncCreateFuncNameMaxLength', '函数名"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('0a74e7bd-c21d-421e-af12-874a60bc61f2', 'CFU-L2-004', 19, '新建函数时函数名前缀规则', 2, 190, 'regexp', '^fn_[_a-zA-Z0-9]+$', 7, 'FuncCreateFuncNamePrefixRequired', '函数名"%s"需要满足前缀正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('5069d187-8cb1-41b9-8876-14e017daf992', 'CFU-L3-001', 19, '新建函数时目标库必须已存在', 3, 190, 'none', 'nil', 4, 'FuncCreateTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('b134e0b3-a6e5-4404-b845-c3f86712fcb9', 'CFU-L3-002', 19, '新建函数时目标函数必须不存在', 3, 190, 'none', 'nil', 4, 'FuncCreateTargetFuncExists', '目标函数"%s"已存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('dfaf91bb-7353-4af6-b88c-3a6f04ad6ab8', 'DFU-L3-001', 19, '删除函数时目标库必须已存在', 3, 192, 'none', 'nil', 4, 'FuncDropTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('e8f789bd-cf1f-4b3e-aa2e-3359979f5485', 'DFU-L3-002', 19, '删除函数时目标函数必须已存在', 3, 192, 'none', 'nil', 4, 'FuncDropTargetFuncExists', '目标函数"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('d3ed9d76-9825-44f3-af9c-b61b9f704e43', 'MFU-L3-001', 19, '修改函数时目标库必须已存在', 3, 191, 'none', 'nil', 4, 'FuncAlterTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('ddec411d-9b2c-427b-965e-86e65d7edfce', 'MFU-L3-002', 19, '修改函数时目标函数必须已存在', 3, 191, 'none', 'nil', 4, 'FuncAlterTargetFuncExists', '目标函数"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('56661395-36fc-4cdc-9346-8cb4608f6f44', 'CTG-L2-001', 20, '新建触发器时触发器名规则', 2, 200, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 5, 'TriggerCreateTriggerNameQualified', '触发器名"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('0d99687d-810f-4f08-aa28-9738ea29eadc', 'CTG-L2-002', 20, '新建触发器时触发器名必须小写', 2, 200, 'regexp', '^[_a-z0-9]+$', 7, 'TriggerCreateTriggerNameLowerCaseRequired', '触发器名"%s"含有大写字母。', 'none', 1, 0, UNIX_TIMESTAMP()),
('73120f74-fdf3-4230-8f31-3f27e7571005', 'CTG-L2-003', 20, '新建触发器时触发器名最大长度', 2, 200, 'lte', '25', 7, 'TriggerCreateTriggerNameMaxLength', '触发器名"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('447e8fe4-4583-4986-ab11-d62ece063b7b', 'CTG-L2-004', 20, '新建触发器时触发器名前缀规则', 2, 200, 'regexp', '^tg_[_a-zA-Z0-9]+$', 5, 'TriggerCreateTriggerPrefixRequired', '触发器名"%s"需要满足前缀正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('27849545-075f-4fb6-8d47-f947c9894dac', 'CTG-L3-001', 20, '新建触发器时目标库必须已存在', 2, 200, 'none', 'nil', 4, 'TriggerCreateTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('35a926d4-f842-4c91-a109-a22327114c52', 'CTG-L3-002', 20, '新建触发器时目标表必须已存在', 2, 200, 'none', 'nil', 4, 'TriggerCreateTargetTableExists', '目标表"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('f2fabd72-4453-45ac-badf-105f46b6f00a', 'CTG-L3-003', 20, '新建触发器时目标触发器必须不存在', 2, 200, 'none', 'nil', 4, 'TriggerCreateTargetTriggerExists', '目标触发器"%s"已存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('f9e45bbe-528d-4ab3-b6dd-4f79749fc12b', 'DTG-L3-001', 20, '删除触发器时目标库必须已存在', 3, 202, 'none', 'nil', 4, 'TriggerDropTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('d55a27c3-0501-4e73-8d6f-c850bbab411e', 'DTG-L3-002', 20, '删除触发器时目标表必须已存在', 3, 202, 'none', 'nil', 4, 'TriggerDropTargetTableExists', '目标表"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('6607f3d3-74c3-4894-b590-2cb073e983d3', 'DTG-L3-003', 20, '删除触发器时目标触发器必须已存在', 3, 202, 'none', 'nil', 4, 'TriggerDropTargetTriggerExists', '目标触发器"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('e6cbd545-643c-40c2-b2d9-924ef0f2691e', 'MTG-L3-001', 20, '修改触发器时目标库必须已存在', 3, 201, 'none', 'nil', 4, 'TriggerAlterTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('10f58b5f-33ed-4ca6-b2af-ad2e01cf360b', 'MTG-L3-002', 20, '修改触发器时目标表必须已存在', 3, 201, 'none', 'nil', 4, 'TriggerAlterTargetTableExists', '目标表"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('7dbfe6fb-e4eb-460d-af51-0427cfe4a600', 'MTG-L3-003', 20, '修改触发器时目标触发器必须已存在', 3, 201, 'none', 'nil', 4, 'TriggerAlterTargetTriggerExists', '目标触发器"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('c8cf9058-f788-48d7-9d87-44201fd574a4', 'CEV-L2-001', 21, '创建事件时事件名规则', 2, 210, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 5, 'EventCreateEventNameQualified', '事件名"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('bf0dc49f-6c6e-40a8-a716-872ebf3bbd71', 'CEV-L2-002', 21, '创建事件时事件名必须小写', 2, 210, 'regexp', '^[_a-z0-9]+$', 5, 'EventCreateEventNameLowerCaseRequired', '事件名"%s"含有大写字母。', 'none', 1, 0, UNIX_TIMESTAMP()),
('8e4484d1-be13-4307-8293-870407428de9', 'CEV-L2-003', 21, '创建事件时事件名最大长度', 2, 210, 'lte', '25', 7, 'EventCreateEventNameMaxLength', '事件名"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('1f67c0d5-86ca-4f9e-95f6-ba6075b4eefc', 'CEV-L2-004', 21, '创建事件时事件名前缀规则', 2, 210, 'regexp', '^ev_[_a-zA-Z0-9]+$', 5, 'EventCreateEventNamePrefixRequired', '事件名"%s"需要满足前缀正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('65fe735f-306c-463e-90d7-30c293d352b7', 'CEV-L3-001', 21, '创建事件时目标库必须已存在', 3, 210, 'none', 'nil', 4, 'EventCreateTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('23dd98bc-685a-413f-a63f-61485bc392cc', 'CEV-L3-002', 21, '创建事件时目标事件必须已存在', 3, 210, 'none', 'nil', 4, 'EventCreateTargetEventExists', '目标事件"%s"已存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('3f2405f2-4fde-40b3-825f-16cb29dea9b9', 'DEV-L3-001', 21, '删除事件时目标库必须已存在', 3, 212, 'none', 'nil', 4, 'EventDropTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('7d46d4bf-6085-401f-9939-63a5120017ad', 'DEV-L3-002', 21, '删除事件时目标事件必须已存在', 3, 212, 'none', 'nil', 4, 'EventDropTargetEventExists', '目标事件"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('d29700b0-1588-4034-aba9-1c2dbc92be05', 'MEV-L3-001', 21, '修改事件时目标库必须已存在', 3, 211, 'none', 'nil', 4, 'EventAlterTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('65003024-7556-4812-b2f6-8aceb77aa052', 'MEV-L3-002', 21, '修改事件时目标事件必须已存在', 3, 211, 'none', 'nil', 4, 'EventAlterTargetEventExists', '目标事件"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('157bdbad-8862-4497-91aa-83827badae02', 'CSP-L2-001', 22, '新建存储过程时存储过程名规则', 2, 220, 'regexp', '^[a-zA-Z][_a-zA-Z0-9]*$', 5, 'ProcCreateProcNameQualified', '存储过程名"%s"需要满足正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('6e82df00-56bf-44e3-adec-fb08bf4d5f3e', 'CSP-L2-002', 22, '新建存储过程时存储过程名必须小写', 2, 220, 'regexp', '^[_a-z0-9]+$', 5, 'ProcCreateProcNameLowerCaseRequired', '存储过程名"%s"含有大写字母。', 'none', 1, 0, UNIX_TIMESTAMP()),
('a4055ee1-caf6-4a72-828e-abea01123af2', 'CSP-L2-003', 22, '新建存储过程时存储过程名最大长度', 2, 220, 'lte', '25', 5, 'ProcCreateProcNameMaxLength', '存储过程名"%s"的长度超出了规则允许的上限，请控制在%d个字符以内。', 'number', 1, 0, UNIX_TIMESTAMP()),
('150f0a36-d1e9-42b1-a30e-4b5e4b2edd71', 'CSP-L2-004', 22, '新建存储过程时存储过程名前缀规则', 2, 220, 'regexp', '^sp_[_a-zA-Z0-9]+$', 5, 'ProcCreateProcNamePrefixRequired', '存储过程名"%s"需要满足前缀正则"%s"。', 'regexp', 1, 0, UNIX_TIMESTAMP()),
('b467ce86-779c-41e6-9079-50b2c9ff3676', 'CSP-L3-001', 22, '新建存储过程时目标库必须已存在', 3, 220, 'none', 'nil', 4, 'ProcCreateTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('60f6be7b-0115-4e9a-bf30-79be2a7c2c00', 'CSP-L3-002', 22, '新建存储过程时目标存储过程必须不存在', 3, 220, 'none', 'nil', 4, 'ProcCreateTargetProcExists', '目标存储过程"%s"已存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('bafb9ff1-cafb-4178-98c9-823ecd8f6600', 'DSP-L3-001', 22, '删除存储过程时目标库必须已存在', 3, 222, 'none', 'nil', 4, 'ProcDropTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('4d627159-45a5-47d9-8942-153b49ac795f', 'DSP-L3-002', 22, '删除存储过程时目标存储过程必须已存在', 3, 222, 'none', 'nil', 4, 'ProcDropTargetProcExists', '目标存储过程"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('841f30ac-bd86-4088-a958-c0bf2e7d26fd', 'MSP-L3-001', 22, '修改存储过程时目标库必须已存在', 3, 221, 'none', 'nil', 4, 'ProcAlterTargetDatabaseExists', '目标库"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('a46eb708-2fdb-495a-8303-ebe72781c875', 'MSP-L3-002', 22, '修改存储过程时目标存储过程必须已存在', 3, 221, 'none', 'nil', 4, 'ProcAlterTargetProcExists', '目标存储过程"%s"不存在。', 'none', 1, 0, UNIX_TIMESTAMP()),
('ddcee0ba-b76c-4037-8ea8-49616c381f8f', 'MSC-L1-001', 23, '禁止LOCK TABLE', 1, 230, 'none', 'nil', 5, 'LockTableProhibited', '禁止LOCK TABLE。', 'none', 1, 0, UNIX_TIMESTAMP()),
('f64e59c8-327e-4d77-9b60-be99d6978421', 'MSC-L1-002', 23, '禁止FLUSH TABLE', 1, 230, 'none', 'nil', 5, 'FlushTableProhibited', '禁止FLUSH操作。', 'none', 1, 0, UNIX_TIMESTAMP()),
('21194aaf-b646-4817-a30d-39e6d1c655b5', 'MSC-L1-003', 23, '禁止TRUNCATE TABLE', 1, 230, 'none', 'nil', 5, 'TruncateTableProhibited', '禁止TRUNCATE TABLE。', 'none', 1, 0, UNIX_TIMESTAMP()),
('f377fbd7-4256-4807-af22-8f8b69d5f6a5', 'MSC-L1-004', 23, '对同一个表/库的操作需要合并', 1, 230, 'none', 'nil', 5, 'MergeRequired', '对同一个对象"%s"的多个操作需要合并。', 'none', 1, 0, UNIX_TIMESTAMP()),
('364f7fb1-e5b8-4308-8166-7849b9e6a08a', 'MSC-L1-005', 23, '禁止PURGE LOG', 1, 230, 'none', 'nil', 5, 'PurgeLogsProhibited', '禁止PURGE LOGS。', 'none', 1, 0, UNIX_TIMESTAMP()),
('e2e2492e-f0fa-4593-a72b-2a9187bd2005', 'MSC-L1-006', 23, '禁止UNLOCK TABLE', 1, 230, 'none', 'nil', 5, 'UnlockTableProhibited', '禁止UNLOCK TABLES。', 'none', 1, 0, UNIX_TIMESTAMP()),
('6d4bce85-5608-4b1a-a810-36a5db17435d', 'MSC-L1-007', 23, '禁止KILL', 1, 230, 'none', 'nil', 5, 'KillProhibited', '禁止KILL。', 'none', 1, 0, UNIX_TIMESTAMP()),
('404354d1-9592-45d9-9f06-81740c5ad3c3', 'MSC-L1-008', 23, '禁止同时出现DDL、DML', 1, 230, 'none', 'nil', 5, 'SplitRequired', '禁止在一个工单中同时出现DML和DDL操作，请分开多个工单提交。', 'none', 1, 0, UNIX_TIMESTAMP());
UNLOCK TABLES;

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
                  COMMENT '执行结果'
  `rows_affected` INT UNSIGNED
                  COMMENT '在服务器正确执行后影响的行数',
  `version`       INT UNSIGNED
                  NOT NULL
                  COMMENT '版本',
  `update_at`     INT UNSIGNED
                  COMMENT '修改时间',
  `create_at`     INT UNSIGNED
                  NOT NULL
                  COMMENT '创建时间',

  PRIMARY KEY (`ticket_id`,`sequence`),
  UNIQUE KEY `unique_1` (`uuid`)
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '工单分解表'
;

DROP TABLE IF EXISTS `mm_statistics`;
CREATE TABLE `mm_statistics` (
  `group`        VARCHAR(36)
                 NOT NULL
                 COMMENT '分组',
  `key`          VARCHAR(50)
                 NOT NULL
                 COMMENT '键',
  `uuid`         CHAR(36)
                 NOT NULL
                 COMMENT 'UUID',
  `value`        DECIMAL(18,4)
                 NOT NULL
                 COMMENT '值',
  `version`      INT UNSIGNED
                 NOT NULL
                 COMMENT '版本',
  `update_at`    INT UNSIGNED
                 COMMENT '修改时间',
  `create_at`    INT UNSIGNED
                 NOT NULL
                 COMMENT '创建时间',

  PRIMARY KEY (`group`, `key`)
  UNIQUE KEY `unique_1` (`uuid`),
)
ENGINE = InnoDB
CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '统计表'
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
  `version`     INT UNSIGNED
                NOT NULL
                COMMENT '版本',
  `update_at`   INT UNSIGNED
                COMMENT '修改时间',
  `create_at`   INT UNSIGNED
                NOT NULL
                COMMENT '创建时间',

  PRIMARY KEY (`ticket_id`),
  UNIQUE KEY `unique_1` (`uuid`),
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
  `update_at` INT UNSIGNED
              COMMENT '修改时间',
  `create_at` INT UNSIGNED
              NOT NULL
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

LOCK TABLES `mm_users` WRITE;
INSERT INTO `mm_users` VALUES
(0,'00000000-0000-0000-0000-000000000000','系统用户','$2a$10$QJT45HdMQIaEHPCNvqkKeeLZpggFEKKU5SdNl.c3hRSGVGbCcMogS',1,'系统用户',0,1,3,0, UNIX_TIMESTAMP());

UPDATE `mm_users` SET `user_id` = 0;

INSERT INTO `mm_users` VALUES
(1,'e70e78bb-9d08-405d-a0ed-266ec703de19','root@163.com','$2a$10$QJT45HdMQIaEHPCNvqkKeeLZpggFEKKU5SdNl.c3hRSGVGbCcMogS',1,'root',0,1,3,0, UNIX_TIMESTAMP());
UNLOCK TABLES;