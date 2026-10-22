-- PostgreSQL 的 BYTEA 没有长度限制，这里只保持迁移编号对齐。
ALTER TABLE mm_clusters ALTER COLUMN fingerprint TYPE BYTEA;
