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

.PHONY: all
all: | fmt build

# TODO: 依赖检查
.PHONY: deps
deps:
	@echo "$(CGREEN)Dependency check$(CEND)"
	@bash ./deps.sh
	# The retool tools.json is setup from retool-install.sh
	retool sync
	retool do gometalinter.v2 intall

# 代码格式化
.PHONY: fmt
fmt:
	@echo "$(CGREEN)Run gofmt on all source files ...$(CEND)"
	@echo "gofmt -l -s -w ..."
	@ret=0 && for d in $$(go list -f '{{.Dir}}' ./... | grep -v /vendor/ | grep -v '/frontend'); do \
		gofmt -l -s -w $$d/*.go || ret=$$? ; \
	done ; exit $$ret

# 运行全部测试用例（无数据库时 DB 相关用例自动跳过）
.PHONY: test
test:
	@echo "$(CGREEN)Run all test cases ...$(CEND)"
	go test ./...
	@echo "test Success!"

# 启动测试用元数据库（MySQL 3307 / PG 5433），等待 schema+seed 加载完成
.PHONY: test-db-up
test-db-up: test-config
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
test-db-down:
	@echo "$(CGREEN)Stopping test databases ...$(CEND)"
	$(COMPOSE) -f $(COMPOSE_FILE) down -v

# 从示例生成测试配置（已存在则不覆盖）
.PHONY: test-config
test-config:
	@test -f bin/config.test.json || cp config.test.json.example bin/config.test.json

# 起测试库 → 跑全部测试（含 DB 集成）→ 拆测试库；无库时不依赖此目标
.PHONY: test-all
test-all: test-db-up
	@echo "$(CGREEN)Run all test cases against test databases ...$(CEND)"
	@HALO_TEST_DSN='$(TEST_MYSQL_DSN)' \
	 HALO_TEST_PG_DSN='$(TEST_PG_DSN)' \
	 HALO_TEST_CONFIG='$(TEST_CONFIG)' \
	 go test ./...; status=$$?; \
	$(COMPOSE) -f $(COMPOSE_FILE) down -v; \
	exit $$status

# 测试代码覆盖率
.PHONY: cover
cover: test
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
build: fmt
	@echo "$(CGREEN)Building ...$(CEND)"
	@bash ./genver.sh $(GO_VERSION_MIN)
	@ret=0 && for pkg in $$(go list -f '{{if (eq .Name "main")}}{{.ImportPath}}{{end}}' ./...); do \
		bin=$$(basename $$pkg); \
		echo "building $$pkg -> bin/$$bin"; \
		GOOS=linux GOARCH=amd64 go build -trimpath -mod=readonly -ldflags="-s -w" -race -o bin/$$bin $$pkg || ret=$$? ; \
	done ; exit $$ret
	@echo "build Success!"

# 安装
install: build
	@echo "$(CGREEN)Install ...$(CEND)"
	go install ./...
	@echo "install Success!"

# gometalinter
# 如果有不想改的lint问题可以使用metalinter.sh加黑名单
#@bash doc/example/metalinter.sh
.PHONY: lint
lint: build
	@echo "$(CGREEN)Run linter check ...$(CEND)"
	CGO_ENABLED=0 retool do gometalinter.v2 -j 1 --config doc/example/metalinter.json ./...
	retool do revive -formatter friendly --exclude vendor/... -config doc/example/revive.toml ./...
	retool do golangci-lint --tests=false run
	@echo "gometalinter check your code is pretty good"

# 清理
.PHONY: clean
clean:
	@echo "$(CGREEN)Cleanup ...$(CEND)"
	go clean
	rm -f ${BINARY}
	rm -f ${BINARY} coverage.*
	find . -name "*.log" -delete
	git clean -fi
