[![LICENSE](https://img.shields.io/badge/license-Anti%20996-blue.svg)](https://github.com/996icu/996.ICU/blob/master/LICENSE)
[![Badge](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu/#/zh_CN)
[![](https://img.shields.io/badge/go-1.11-brightgreen.svg)](https://golang.org/dl/)
![](https://img.shields.io/badge/license-GPL--3.0-orange.svg)

基于规则的SQL审核系统

假设下列规则启用（系统共计百余条规则，此处仅例举）：
* [x] 禁止LOCK TABLE
* [x] 禁止FLUSH TABLE
* [x] 禁止TRUNCATE TABLE
* [x] 禁止PURGE LOG
* [x] 禁止UNLOCK TABLE
* [x] 禁止KILL

示例SQL：
```
LOCK TABLES t1 READ;
LOCK TABLES t2 WRITE;
LOCK TABLES t3 READ LOCAL;
LOCK TABLES t4 WRITE;
LOCK TABLES t5 READ, t6 WRITE;
```

分析结果：
```
1 | LOCK TABLES t1 READ;           | [{"Level":1,"Description":"禁止LOCK TABLE。"}]
2 | LOCK TABLES t2 WRITE;          | [{"Level":1,"Description":"禁止LOCK TABLE。"}]
3 | LOCK TABLES t3 READ LOCAL;     | [{"Level":1,"Description":"禁止LOCK TABLE。"}]
4 | LOCK TABLES t4 WRITE;          | [{"Level":1,"Description":"禁止LOCK TABLE。"}]
5 | LOCK TABLES t5 READ, t6 WRITE; | [{"Level":1,"Description":"禁止LOCK TABLE。"}]
```

示例SQL：
```
FLUSH NO_WRITE_TO_BINLOG TABLES t1 WITH READ LOCK;
FLUSH TABLES;
FLUSH TABLES t1;
FLUSH NO_WRITE_TO_BINLOG TABLES t1;
FLUSH TABLES WITH READ LOCK;
FLUSH TABLES t1, t2, t3;
FLUSH TABLES t1, t2, t3 WITH READ LOCK;
FLUSH PRIVILEGES;
FLUSH STATUS;
```

分析结果
```
1 | FLUSH NO_WRITE_TO_BINLOG TABLES `t1` WITH READ LOCK | [{"Level":1,"Description":"禁止FLUSH TABLE。"}]
2 | FLUSH TABLES                                        | [{"Level":1,"Description":"禁止FLUSH TABLE。"}]
3 | FLUSH TABLES `t1`                                   | [{"Level":1,"Description":"禁止FLUSH TABLE。"}]
4 | FLUSH NO_WRITE_TO_BINLOG TABLES `t1`                | [{"Level":1,"Description":"禁止FLUSH TABLE。"}]
5 | FLUSH TABLES WITH READ LOCK                         | [{"Level":1,"Description":"禁止FLUSH TABLE。"}]
6 | FLUSH TABLES `t1`, `t2`, `t3`                       | [{"Level":1,"Description":"禁止FLUSH TABLE。"}]
7 | FLUSH TABLES `t1`, `t2`, `t3` WITH READ LOCK        | [{"Level":1,"Description":"禁止FLUSH TABLE。"}]
8 | FLUSH PRIVILEGES                                    | [{"Level":1,"Description":"禁止FLUSH TABLE。"}]
9 | FLUSH STATUS                                        | [{"Level":1,"Description":"禁止FLUSH TABLE。"}]
```

示例SQL：
```
TRUNCATE TABLE `t1`;
TRUNCATE TABLE `t2`;
```

分析结果：
```
1 | TRUNCATE TABLE `t1` | [{"Level":1,"Description":"禁止TRUNCATE TABLE。"}]
2 | TRUNCATE TABLE `t2` | [{"Level":1,"Description":"禁止TRUNCATE TABLE。"}]
```

示例SQL：
```
PURGE BINARY LOGS TO 'mysql-bin.010';
PURGE BINARY LOGS BEFORE '2008-04-02 22:46:26';
PURGE MASTER LOGS TO 'mysql-bin.010';
PURGE MASTER LOGS BEFORE '2008-04-02 22:46:26';
```

分析结果：
```
1 | PURGE BINARY LOGS TO 'mysql-bin.010';           | [{"Level":1,"Description":"禁止PURGE LOGS。"}]
2 | PURGE BINARY LOGS BEFORE '2008-04-02 22:46:26'; | [{"Level":1,"Description":"禁止PURGE LOGS。"}]
3 | PURGE MASTER LOGS TO 'mysql-bin.010';           | [{"Level":1,"Description":"禁止PURGE LOGS。"}]
4 | PURGE MASTER LOGS BEFORE '2008-04-02 22:46:26'; | [{"Level":1,"Description":"禁止PURGE LOGS。"}]
```

示例SQL：
```
UNLOCK TABLES;
```

分析结果：
```
1 | UNLOCK TABLES; | [{"Level":1,"Description":"禁止UNLOCK TABLES。"}]
```

示例SQL：
```
KILL 1234;
KILL CONNECTION 5678;
KILL QUERY 90;
```

分析结果：

```
1 | KILL 1234 | [{"Level":1,"Description":"禁止KILL。"}]
2 | KILL 5678 | [{"Level":1,"Description":"禁止KILL。"}]
3 | KILL 90   | [{"Level":1,"Description":"禁止KILL。"}]
```

示例SQL：
```
CREATE TABLE `t1` (`id` INT) CHARSET = 'utf8mb4' COLLATE = 'utf8mb4_unicode_ci';
```

分析结果：

```
1 | CREATE TABLE `t1` (`id` INT)   | [{"Level":1,"Description":"需要为表\"t1\"需要提供COMMENT注解。"},
  | CHARSET = 'utf8mb4'            |  {"Level":1,"Description":"必须为表指定一个主键。"},
  | COLLATE = 'utf8mb4_unicode_ci' |  {"Level":1,"Description":"建表禁用存储引擎\"[empty]\"，请使用\"[innodb]\"。"},
  |                                |  {"Level":1,"Description":"列\"id\"需要提供COMMENT注解。"}]
```


邮件通知：

![mail-1](https://s2.ax1x.com/2019/04/18/ESIYvV.png)
![mail-2](https://s2.ax1x.com/2019/04/19/EpOma6.png)


* [ ] 数据备份
* [x] 任务调度
* [x] 事件处理
* [x] 终端程序
* [x] 单元测试
* [ ] Gh-ost整合
* [ ] 用户界面
  * [ ] 数据查询(SOAR整合)
  * [ ] 工单列表
  * [ ] 工单查看
  * [ ] 工单提交
  * [ ] 工单编辑
  * [ ] 群集列表
  * [ ] 群集查看
  * [ ] 群集提交
  * [ ] 群集编辑
  * [ ] 用户列表
  * [ ] 用户查看
  * [ ] 用户提交
  * [ ] 用户编辑
  * [ ] 日志列表
  * [ ] 任务列表
  * [ ] 任务查看
  * [ ] 数据库元数据查看
  * [ ] 规则查看
  * [ ] 系统选项
