-- 群集表新增 type 列，区分目标库类型。存量记录一律视为 MySQL。
ALTER TABLE mm_clusters ADD COLUMN `type` VARCHAR(20) NOT NULL DEFAULT 'mysql';
