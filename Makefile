BINARY=halo
GOPATH ?= $(shell go env GOPATH)
ifeq "$(GOPATH)" ""
  $(error 运行`make`命令前请先设置好GOPATH环境变量。)
endif
PATH := ${GOPATH}/bin:$(PATH)
CRED:=$(shell tput setaf 1 2>/dev/null)
CGREEN:=$(shell tput setaf 2 2>/dev/null)
CYELLOW:=$(shell tput setaf 3 2>/dev/null)
CEND:=$(shell tput sgr0 2>/dev/null)
GO_VERSION_MIN=1.26

# Test service compose: prefer podman-compose (rootless, SELinux-friendly),
# fall back to docker compose. TEST_* can be overridden on the command line.
COMPOSE ?= $(shell command -v podman-compose 2>/dev/null || (command -v docker-compose 2>/dev/null || echo docker compose))
ENGINE ?= $(shell command -v podman >/dev/null 2>&1 && echo podman || echo docker)
COMPOSE_FILE ?= docker-compose.test.yml
TEST_MYSQL_DSN ?= root:test@tcp(127.0.0.1:3307)/halo_test?charset=utf8mb4&parseTime=true&loc=UTC&time_zone=%27%2B00%3A00%27
TEST_PG_DSN ?= postgres://postgres:test@127.0.0.1:5433/halo_test?sslmode=disable
TEST_CONFIG ?= $(CURDIR)/bin/config.test.json


# 帮助：列出可用目标（目标后的 ## 注释即是说明）
.PHONY: help
help:
	@echo "用法: make <target>"
	@echo ""
	@grep -hE '^[a-zA-Z0-9_.-]+:.*## ' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*## "}; {printf "  $(CGREEN)%-20s$(CEND) %s\n", $$1, $$2}'

.PHONY: all
all: | fmt build ## 格式化并构建全部二进制

# TODO: 依赖检查
.PHONY: deps
deps: ## 安装 Go 依赖
	@echo "$(CGREEN)Dependency check$(CEND)"
	@bash ./deps.sh
	# The retool tools.json is setup from retool-install.sh
	retool sync
	retool do gometalinter.v2 intall

# 代码格式化
.PHONY: fmt
fmt: ## 格式化 Go 代码
	@echo "$(CGREEN)Run gofmt on all source files ...$(CEND)"
	@echo "gofmt -l -s -w ..."
	@ret=0 && for d in $$(go list -f '{{.Dir}}' ./... | grep -v /vendor/ | grep -v '/frontend'); do \
		gofmt -l -s -w $$d/*.go || ret=$$? ; \
	done ; exit $$ret

# 运行全部测试用例（无数据库时 DB 相关用例自动跳过）
.PHONY: test
test: ## 运行 Go 测试
	@echo "$(CGREEN)Run all test cases ...$(CEND)"
	go test ./...
	@echo "test Success!"

# 启动测试用元数据库（MySQL 3307 / PG 5433），等待 schema+seed 加载完成
.PHONY: test-db-up
test-db-up: test-config ## 拉起测试用 MySQL/PostgreSQL
	@echo "$(CGREEN)Starting test databases ...$(CEND)"
	$(COMPOSE) -f $(COMPOSE_FILE) up -d
	@for c in halo-test-mysql halo-test-postgres; do \
		for i in $$(seq 1 90); do \
			[ "$$($(ENGINE) inspect -f '{{.State.Health.Status}}' $$c 2>/dev/null)" = "healthy" ] && break; \
			sleep 2; \
		done; \
		[ "$$($(ENGINE) inspect -f '{{.State.Health.Status}}' $$c 2>/dev/null)" = "healthy" ] \
			|| { echo "$(CRED)timeout waiting for $$c$(CEND)"; exit 1; }; \
		echo "$$c healthy"; \
	done

# 停止并清除测试数据库（含数据卷）
.PHONY: test-db-down
test-db-down: ## 停止测试数据库
	@echo "$(CGREEN)Stopping test databases ...$(CEND)"
	$(COMPOSE) -f $(COMPOSE_FILE) down -v

# 从示例生成测试配置（已存在则不覆盖）
.PHONY: test-config
test-config: ## 生成测试配置
	@test -f bin/config.test.json || cp config.test.json.example bin/config.test.json

# 元数据库连接读取的配置文件：优先 bin/config.json，回退 bin/config.mysql.json
META_CONFIG ?= $(shell test -f bin/config.json && echo bin/config.json || echo bin/config.mysql.json)

# 向元数据库灌入测试 cluster（幂等）。连接信息取自 $(META_CONFIG) 的 database.addr。
.PHONY: test-seed-meta
test-seed-meta: ## 写入测试种子数据
	@dsn=$$(grep -m1 '"addr"' $(META_CONFIG) | sed -E 's/.*"addr"[[:space:]]*:[[:space:]]*"([^"]*)".*/\1/'); \
	if [ -z "$$dsn" ]; then echo "$(CRED)cannot parse database.addr from $(META_CONFIG)$(CEND)"; exit 1; fi; \
	user=$${dsn%%:*}; rest=$${dsn#*:}; pass=$${rest%%@tcp(*}; \
	hp=$${dsn#*@tcp(}; hp=$${hp%%)*}; host=$${hp%%:*}; port=$${hp##*:}; \
	dbpart=$${dsn#*)/}; db=$${dbpart%%\?*}; \
	echo "$(CGREEN)seeding test clusters into $$host:$$port/$$db ...$(CEND)"; \
	mysql --protocol=tcp -h "$$host" -P "$$port" -u "$$user" -p"$$pass" "$$db" < scripts/test/seed_meta_test.sql

# 起测试库 → 跑全部测试（含 DB 集成）→ 拆测试库；无库时不依赖此目标
.PHONY: test-all
test-all: test-db-up ## 拉起测试库并跑全量测试
	@echo "$(CGREEN)Run all test cases against test databases ...$(CEND)"
	@HALO_TEST_DSN='$(TEST_MYSQL_DSN)' \
	 HALO_TEST_PG_DSN='$(TEST_PG_DSN)' \
	 HALO_TEST_CONFIG='$(TEST_CONFIG)' \
	 go test ./...; status=$$?; \
	$(COMPOSE) -f $(COMPOSE_FILE) down -v; \
	exit $$status

# 测试代码覆盖率
.PHONY: cover
cover: test ## 运行测试并生成覆盖率报告
	@echo "$(CGREEN)Run test cover check ...$(CEND)"
	go test -coverpkg=./... -coverprofile=coverage.data ./... | column -t
	go tool cover -html=coverage.data -o coverage.html
	go tool cover -func=coverage.data -o coverage.txt
	@tail -n 1 coverage.txt | awk '{sub(/%/, "", $$NF); \
		if($$NF < 80) \
			{print "\033[91m"$$0"%\033[0m"} \
		else if ($$NF >= 90) \
			{print "\033[92m"$$0"%\033[0m"} \
		else \
			{print "\033[93m"$$0"%\033[0m"}}'

# 项目构建
build: fmt frontend-build ## 构建前端并编译全部二进制（前端打入 server）
	@echo "$(CGREEN)Building ...$(CEND)"
	@P=github.com/mia0x75/halo/internal/app; \
	git="$$(git describe --tags --always 2>/dev/null)"; \
	[ -n "$$git" ] || git="not a git repo"; \
	version="$$(cat VERSION 2>/dev/null || echo 0.0.0)"; \
	build="$$(cat BUILD 2>/dev/null || echo 0)"; \
	echo $$((build + 1)) > BUILD; \
	branch="$$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo unknown)"; \
	kernel="$$(uname -r 2>/dev/null || echo unknown)"; \
	[ -r /etc/os-release ] && . /etc/os-release; \
	distro="$${PRETTY_NAME:-Unknown}"; \
	buildAt="$$(date '+%F %T %z')"; \
	gover="$$(go env GOVERSION | sed -E 's/^go//; s/[[:space:]-].*//')"; \
	[ -n "$$gover" ] || gover="unknown"; \
	ldflags="-X \"$$P.Version=$$version.$$build\""; \
	for kv in "Git:$$git" "BuildAt:$$buildAt" "Branch:$$branch" "Distro:$$distro" "Kernel:$$kernel" "GoVer:$$gover"; do \
		ldflags="$$ldflags -X \"$$P.$${kv%%:*}=$${kv#*:}\""; \
	done; \
	ret=0; \
	for pkg in $$(go list -f '{{if (eq .Name "main")}}{{.ImportPath}}{{end}}' ./...); do \
		bin=$$(basename $$pkg); \
		GOOS=linux GOARCH=amd64 go build -trimpath -mod=readonly -ldflags "$$ldflags -s -w" -race -o bin/$$bin $$pkg || ret=$$? ; \
	done; exit $$ret
	@echo "$(CGREEN)build Success!$(CEND)"

# 安装
install: build ## 构建并安装二进制
	@echo "$(CGREEN)Install ...$(CEND)"
	go install ./...
	@echo "install Success!"

# gometalinter
# 如果有不想改的lint问题可以使用metalinter.sh加黑名单
#@bash doc/example/metalinter.sh
.PHONY: lint
lint: build ## 运行 Go 静态检查
	@echo "$(CGREEN)Run linter check ...$(CEND)"
	CGO_ENABLED=0 retool do gometalinter.v2 -j 1 --config doc/example/metalinter.json ./...
	retool do revive -formatter friendly --exclude vendor/... -config doc/example/revive.toml ./...
	retool do golangci-lint --tests=false run
	@echo "gometalinter check your code is pretty good"

# 重新生成 sqlc 与 GraphQL 服务端代码
.PHONY: generate
generate: sqlc gqlgen ## 重新生成 sqlc 与 gqlgen 代码

.PHONY: sqlc
sqlc: ## 重新生成 sqlc 数据访问代码
	@echo "$(CGREEN)Generating sqlc code ...$(CEND)"
	sqlc generate

.PHONY: gqlgen
gqlgen: ## 重新生成 gqlapi/generated.go
	@echo "$(CGREEN)Generating GraphQL server code ...$(CEND)"
	go tool gqlgen generate

# 校验迁移基线与权威 schema 脚本是否同步
.PHONY: check-migrations
check-migrations: ## 校验迁移链与权威 schema 是否同步
	@echo "$(CGREEN)Checking migration baselines against schema scripts ...$(CEND)"
	@bash scripts/check_migrations_sync.sh

# 执行元数据库 schema 迁移（需先设置 HALO_CFG）
.PHONY: migrate-up
migrate-up: ## 执行元数据库迁移
	@echo "$(CGREEN)Applying metadata migrations ...$(CEND)"
	@go run ./cmd/halocli migrate up

# 查看元数据库迁移状态
.PHONY: migrate-status
migrate-status: ## 查看元数据库迁移状态
	@echo "$(CGREEN)Migration status ...$(CEND)"
	@go run ./cmd/halocli migrate status

# 前端依赖安装
.PHONY: frontend-install
frontend-install: ## 安装前端依赖
	@echo "$(CGREEN)Installing frontend dependencies ...$(CEND)"
	cd frontend && bun install

# 前端构建（类型检查 + 打包）
.PHONY: frontend-build
frontend-build: ## 构建前端并复制进 embed 目录
	@echo "$(CGREEN)Building frontend ...$(CEND)"
	cd frontend && bun run build
	@echo "$(CGREEN)Copying frontend into embed dir ...$(CEND)"
	@mkdir -p internal/transport/web/dist
	@find internal/transport/web/dist -mindepth 1 ! -name .gitkeep -delete
	@cp -r frontend/dist/. internal/transport/web/dist/
	@echo "$(CGREEN)frontend build Success!$(CEND)"

# 前端静态检查
.PHONY: frontend-lint
frontend-lint: ## 前端静态检查
	@echo "$(CGREEN)Run frontend linter ...$(CEND)"
	cd frontend && bun run lint
	@echo "$(CGREEN)frontend lint Success!$(CEND)"

# 前端开发服务器
.PHONY: frontend-dev
frontend-dev: ## 启动前端开发服务器
	@echo "$(CGREEN)Starting frontend dev server ...$(CEND)"
	cd frontend && bun run dev

# 清理
.PHONY: clean
clean:
	@echo "$(CGREEN)Cleanup ...$(CEND)"
	go clean
	rm -f ${BINARY}
	rm -f ${BINARY} coverage.*
	find . -name "*.log" -delete
	git clean -fi
