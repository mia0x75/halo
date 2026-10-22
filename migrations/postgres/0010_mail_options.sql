-- 邮件通道选项。
--
-- SMTP 配置从配置文件迁到 mm_options，改完立即生效，不必重启服务。
-- mail.password 以 AES 加密后 base64 存储；配置文件里的 mail 段只作为缺省值。

INSERT INTO mm_options (name, uuid, parameter, description, element, writable, version, create_at, update_at) VALUES
('mail.enabled','9d2f1b10-0001-4b10-8b01-000000000001','false','邮件通知开关','switch',1,1, NOW(), NOW()),
('mail.addr','9d2f1b10-0002-4b10-8b01-000000000002','','SMTP 服务器（host:port）','input',1,1, NOW(), NOW()),
('mail.user','9d2f1b10-0003-4b10-8b01-000000000003','','SMTP 账号','input',1,1, NOW(), NOW()),
('mail.password','9d2f1b10-0004-4b10-8b01-000000000004','','SMTP 密码（加密存储）','password',1,1, NOW(), NOW()),
('mail.insecure-skip-verify','9d2f1b10-0005-4b10-8b01-000000000005','false','跳过 SMTP 证书校验','switch',1,1, NOW(), NOW())
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description, element = EXCLUDED.element;
