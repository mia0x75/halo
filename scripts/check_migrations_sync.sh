#!/usr/bin/env bash
# 校验权威 schema 脚本与迁移链是否收敛到同一批表。
#
# 权威脚本（scripts/halodb_schema_<engine>.sql）是当前结构的快照，迁移链
# （migrations/<engine>/*.sql）是它的增量历史。逐个版本重放迁移链的
# CREATE/DROP，得到链尾的表集合，与快照比对，避免「快照有、迁移没有」这类漂移。
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
status=0

normalize() {
    sed -E "s/CREATE TABLE (IF NOT EXISTS )?//" | tr -d "\"\`"
}

tables_in() {
    grep -hoE "CREATE TABLE (IF NOT EXISTS )?[\"\`]?[A-Za-z_][A-Za-z0-9_]*[\"\`]?" "$@" | normalize | sort -u || true
}

dropped_in() {
    grep -hoE "DROP TABLE (IF EXISTS )?[\"\`]?[A-Za-z_][A-Za-z0-9_]*[\"\`]?" "$@" \
        | sed -E "s/DROP TABLE (IF EXISTS )?//" | tr -d "\"\`" | sort -u || true
}

check() {
    local engine="$1" schema="$2" dir="$3"
    local schema_tables chain_tables created dropped

    schema_tables=$(tables_in "$root/$schema")

    chain_tables=""
    for file in $(ls "$root/$dir"/*.sql | sort); do
        created=$(tables_in "$file")
        dropped=$(dropped_in "$file")
        chain_tables=$(echo "$chain_tables
$created" | grep -v '^$' | sort -u)
        if [ -n "$dropped" ]; then
            chain_tables=$(echo "$chain_tables" | grep -vxF "$dropped" || true)
        fi
    done
    # mm_migrations 由迁移器自身创建，不在迁移文件里。
    chain_tables=$(echo "$chain_tables
mm_migrations" | grep -v '^$' | sort -u)

    if [ "$schema_tables" == "$chain_tables" ]; then
        echo "OK   $engine 权威 schema 与迁移链收敛到同一批表"
        return
    fi

    echo "FAIL $engine 权威 schema 与迁移链不一致："
    diff <(echo "$schema_tables") <(echo "$chain_tables") || true
    status=1
}

check mysql scripts/halodb_schema_mysql.sql migrations/mysql
check postgres scripts/halodb_schema_postgres.sql migrations/postgres

exit $status
