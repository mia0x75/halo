-- Test clusters for integration tests. Idempotent.
-- Target ports: MySQL 3306, PG 5432 (host).
-- Column note: the live schema names the engine column `type` (VARCHAR), not
-- `kind`; values are 'mysql' / 'postgres'. Password is a placeholder (0x00) --
-- the integration test encrypts the real target password at runtime.

INSERT INTO mm_clusters
    (cluster_id, uuid, host, ip, port, alias, user, type, password, fingerprint, status, version, update_at, create_at)
VALUES
    (9999, 'fd577136-0000-0000-0000-000000000001', 'localhost', '127.0.0.1', 3306, 'test-mysql-target',
     'root', 'mysql', 0x00, 0x00, 1, 0, NULL, UTC_TIMESTAMP()),
    (9998, 'fd577137-0000-0000-0000-000000000002', 'localhost', '127.0.0.1', 5432, 'test-pg-target',
     'postgres', 'postgres', 0x00, 0x00, 1, 0, NULL, UTC_TIMESTAMP())
ON DUPLICATE KEY UPDATE
    uuid   = VALUES(uuid),
    host   = VALUES(host),
    ip     = VALUES(ip),
    port   = VALUES(port),
    alias  = VALUES(alias),
    user   = VALUES(user),
    type   = VALUES(type);
