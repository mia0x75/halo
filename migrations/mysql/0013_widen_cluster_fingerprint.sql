-- 群集指纹改为 SHA-256 后是 32 字节，旧列按 SHA-1 的 20 字节建，
-- 插入时报 Data too long。这里把列宽补到 32。
ALTER TABLE `mm_clusters`
    MODIFY `fingerprint` VARBINARY(32) NOT NULL COMMENT '指纹';
