#!/usr/bin/env bash
# 校验迁移基线与权威 schema 脚本是否同步。
#
# 基线（migrations/<engine>/0001_init.sql）由 scripts/halodb_schema_<engine>.sql
# 去掉 DROP 语句、CREATE 改为 IF NOT EXISTS 得到。两者一旦分叉，迁移历史与
# 权威结构就会出现两套事实，故在提交前用本脚本卡住。
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
status=0

normalize() {
    # 去掉注释与空行、DROP 语句，并把 IF NOT EXISTS 还原成 CREATE，
    # 只比较结构本身。
    grep -vE "^[[:space:]]*--" "$1" \
        | grep -vE "^[[:space:]]*$" \
        | grep -vE "^[[:space:]]*DROP TABLE IF EXISTS" \
        | sed -E "s/^CREATE TABLE IF NOT EXISTS /CREATE TABLE /" \
        | sed -E "s/^CREATE (UNIQUE )?INDEX IF NOT EXISTS /CREATE \1INDEX /" \
        | sed -E "s/[[:space:]]+$//"
}

check() {
    local engine="$1" schema="$2" baseline="$3"
    if [ ! -f "$root/$baseline" ]; then
        echo "缺少基线迁移文件：$baseline"
        status=1
        return
    fi
    if diff -u \
        <(normalize "$root/$schema") \
        <(normalize "$root/$baseline") > /tmp/halo-migrations-sync.diff; then
        echo "OK   $engine 基线与权威 schema 同步"
    else
        echo "FAIL $engine 基线与权威 schema 不一致："
        cat /tmp/halo-migrations-sync.diff
        status=1
    fi
}

check mysql scripts/halodb_schema_mysql.sql migrations/mysql/0001_init.sql
check postgres scripts/halodb_schema_postgres.sql migrations/postgres/0001_init.sql

exit $status
