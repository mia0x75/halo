-- 选项重排：历史配置拆成 SMTP 与 LDAP 两组，组的先后与组内顺序统一由
-- sort_order 决定，各组的开关排在首位。

-- 邮件通知：总开关在前。
UPDATE mm_options SET group_name = '邮件通知', sort_order = 1 WHERE name = 'mail.enabled';
UPDATE mm_options SET group_name = '邮件通知', sort_order = 2 WHERE name = 'mail.addr';
UPDATE mm_options SET group_name = '邮件通知', sort_order = 3 WHERE name = 'mail.user';
UPDATE mm_options SET group_name = '邮件通知', sort_order = 4 WHERE name = 'mail.password';
UPDATE mm_options SET group_name = '邮件通知', sort_order = 5 WHERE name = 'mail.insecure-skip-verify';

-- 数据保留期。
UPDATE mm_options SET group_name = '数据保留期', sort_order = 6 WHERE name = 'retention.metrics-days';
UPDATE mm_options SET group_name = '数据保留期', sort_order = 7 WHERE name = 'retention.hits-days';
UPDATE mm_options SET group_name = '数据保留期', sort_order = 8 WHERE name = 'retention.logs-days';
UPDATE mm_options SET group_name = '数据保留期', sort_order = 9 WHERE name = 'retention.messages-days';
UPDATE mm_options SET group_name = '数据保留期', sort_order = 10 WHERE name = 'retention.tasks-days';

-- 历史配置：SMTP 与 LDAP 各一组，开关排在本组首位。
UPDATE mm_options SET group_name = 'SMTP（历史配置）', sort_order = 11 WHERE name = 'smtp.enabled';
UPDATE mm_options SET group_name = 'SMTP（历史配置）', sort_order = 12 WHERE name = 'smtp.host';
UPDATE mm_options SET group_name = 'SMTP（历史配置）', sort_order = 13 WHERE name = 'smtp.port';
UPDATE mm_options SET group_name = 'SMTP（历史配置）', sort_order = 14 WHERE name = 'smtp.user';
UPDATE mm_options SET group_name = 'SMTP（历史配置）', sort_order = 15 WHERE name = 'smtp.password';
UPDATE mm_options SET group_name = 'SMTP（历史配置）', sort_order = 16 WHERE name = 'smtp.encryption';

UPDATE mm_options SET group_name = 'LDAP（历史配置）', sort_order = 17 WHERE name = 'ldap.enabled';
UPDATE mm_options SET group_name = 'LDAP（历史配置）', sort_order = 18 WHERE name = 'ldap.host';
UPDATE mm_options SET group_name = 'LDAP（历史配置）', sort_order = 19 WHERE name = 'ldap.domain';
UPDATE mm_options SET group_name = 'LDAP（历史配置）', sort_order = 20 WHERE name = 'ldap.type';
UPDATE mm_options SET group_name = 'LDAP（历史配置）', sort_order = 21 WHERE name = 'ldap.user';
UPDATE mm_options SET group_name = 'LDAP（历史配置）', sort_order = 22 WHERE name = 'ldap.password';
UPDATE mm_options SET group_name = 'LDAP（历史配置）', sort_order = 23 WHERE name = 'ldap.sc';
UPDATE mm_options SET group_name = 'LDAP（历史配置）', sort_order = 24 WHERE name = 'ldap.ou';
