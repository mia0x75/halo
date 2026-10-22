-- 保留期选项。
--
-- 指标桶、规则命中、操作日志、站内信与任务执行记录都会持续增长，需要一个
-- 可配置的保留窗口。保留期以「天」为单位，取值来自 mm_options，清理作业
-- (halocli retention cleanup) 按这些值删除过期数据。

INSERT INTO mm_options (name, uuid, parameter, description, element, writable, version, create_at, update_at) VALUES
('retention.metrics-days','7c1f4a10-0001-4a10-9a01-000000000001','90','','input',1,1, NOW(), NOW()),
('retention.hits-days','7c1f4a10-0002-4a10-9a01-000000000002','90','','input',1,1, NOW(), NOW()),
('retention.logs-days','7c1f4a10-0003-4a10-9a01-000000000003','180','','input',1,1, NOW(), NOW()),
('retention.messages-days','7c1f4a10-0004-4a10-9a01-000000000004','90','','input',1,1, NOW(), NOW()),
('retention.tasks-days','7c1f4a10-0005-4a10-9a01-000000000005','30','','input',1,1, NOW(), NOW())
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description, element = EXCLUDED.element;
