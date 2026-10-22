-- 选项元数据：分组、排序、依赖与数值范围。
--
-- element 与 writable 之前只在 Go 注册表里维护，数据库里的 element 列是历史
-- 遗留的死数据。这次把元数据统一放到数据库，并补上设置页需要的展示信息。

ALTER TABLE mm_options
    ADD COLUMN group_name VARCHAR(50) NOT NULL DEFAULT '',
    ADD COLUMN sort_order INT NOT NULL DEFAULT 0,
    ADD COLUMN depends_on VARCHAR(50) NOT NULL DEFAULT '',
    ADD COLUMN min_value INT NULL,
    ADD COLUMN max_value INT NULL;

-- 历史行的 element 修正为真实控件类型。
UPDATE mm_options SET element = 'switch' WHERE name IN ('smtp.enabled', 'ldap.enabled');
UPDATE mm_options SET element = 'password' WHERE name IN ('smtp.password', 'ldap.password');
UPDATE mm_options SET element = 'input' WHERE name IN (
    'smtp.host', 'smtp.port', 'smtp.user', 'smtp.encryption',
    'ldap.host', 'ldap.domain', 'ldap.type', 'ldap.user', 'ldap.sc', 'ldap.ou');
UPDATE mm_options SET element = 'number' WHERE name LIKE 'retention.%';

-- 分组与组内顺序。
UPDATE mm_options SET group_name = '数据保留期', sort_order = 1 WHERE name = 'retention.metrics-days';
UPDATE mm_options SET group_name = '数据保留期', sort_order = 2 WHERE name = 'retention.hits-days';
UPDATE mm_options SET group_name = '数据保留期', sort_order = 3 WHERE name = 'retention.logs-days';
UPDATE mm_options SET group_name = '数据保留期', sort_order = 4 WHERE name = 'retention.messages-days';
UPDATE mm_options SET group_name = '数据保留期', sort_order = 5 WHERE name = 'retention.tasks-days';

UPDATE mm_options SET group_name = '邮件通知', sort_order = 1 WHERE name = 'mail.enabled';
UPDATE mm_options SET group_name = '邮件通知', sort_order = 2 WHERE name = 'mail.addr';
UPDATE mm_options SET group_name = '邮件通知', sort_order = 3 WHERE name = 'mail.user';
UPDATE mm_options SET group_name = '邮件通知', sort_order = 4 WHERE name = 'mail.password';
UPDATE mm_options SET group_name = '邮件通知', sort_order = 5 WHERE name = 'mail.insecure-skip-verify';

UPDATE mm_options SET group_name = '历史配置（只读）', sort_order = 1 WHERE name = 'smtp.enabled';
UPDATE mm_options SET group_name = '历史配置（只读）', sort_order = 2 WHERE name = 'smtp.host';
UPDATE mm_options SET group_name = '历史配置（只读）', sort_order = 3 WHERE name = 'smtp.port';
UPDATE mm_options SET group_name = '历史配置（只读）', sort_order = 4 WHERE name = 'smtp.user';
UPDATE mm_options SET group_name = '历史配置（只读）', sort_order = 5 WHERE name = 'smtp.password';
UPDATE mm_options SET group_name = '历史配置（只读）', sort_order = 6 WHERE name = 'smtp.encryption';
UPDATE mm_options SET group_name = '历史配置（只读）', sort_order = 7 WHERE name = 'ldap.enabled';
UPDATE mm_options SET group_name = '历史配置（只读）', sort_order = 8 WHERE name = 'ldap.host';
UPDATE mm_options SET group_name = '历史配置（只读）', sort_order = 9 WHERE name = 'ldap.domain';
UPDATE mm_options SET group_name = '历史配置（只读）', sort_order = 10 WHERE name = 'ldap.type';
UPDATE mm_options SET group_name = '历史配置（只读）', sort_order = 11 WHERE name = 'ldap.user';
UPDATE mm_options SET group_name = '历史配置（只读）', sort_order = 12 WHERE name = 'ldap.password';
UPDATE mm_options SET group_name = '历史配置（只读）', sort_order = 13 WHERE name = 'ldap.sc';
UPDATE mm_options SET group_name = '历史配置（只读）', sort_order = 14 WHERE name = 'ldap.ou';

-- 邮件组以 mail.enabled 为总开关，其余项依赖它。
UPDATE mm_options SET depends_on = 'mail.enabled'
    WHERE name IN ('mail.addr', 'mail.user', 'mail.password', 'mail.insecure-skip-verify');

-- 保留期以天为单位，取值 1-3650。
UPDATE mm_options SET min_value = 1, max_value = 3650 WHERE name LIKE 'retention.%';
