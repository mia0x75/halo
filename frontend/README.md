# Halo 前端

SQL 变更工单与审计平台的 Web 前端。技术栈与设计规范见仓库根目录
[FRONTEND-DECISIONS.md](../FRONTEND-DECISIONS.md)。

## 命令

```bash
bun install          # 安装依赖
bun run dev          # 开发服务器（默认 http://localhost:5173）
bun run build        # 类型检查（tsc -b）+ 打包
bun run lint         # oxlint
bun run preview      # 预览打包产物
```

仓库根目录亦提供 `make frontend-install / frontend-build / frontend-lint / frontend-dev`。

## 目录结构

```
src/
  components/
    ui/          # 自持的基础组件（button/card/badge/input/label/table/separator）
    layout/      # 应用外壳：顶栏、页头、容器、主题切换
    charts/      # 图表调色板与 ApexCharts 封装
    shared/      # 跨业务组件：section-card、status-badge、empty-state
  features/
    dashboard/   # 业务域组件（概览）
  hooks/         # use-theme
  lib/
    api/         # types.ts（接口契约）/ mock.ts（类型化样例）/ client.ts（GraphQL 传输）
    format.ts    # 纯格式化函数
    utils.ts     # cn 再导出
  pages/         # 页面级组合
  App.tsx
  main.tsx
  index.css      # 设计令牌 + 主题 + Tailwind 入口
```

约定：`features/*` 可依赖 `components/*` 与 `lib/*`；`components/ui` 必须保持通用、
不掺业务语义。组件不得直接发起网络请求，一律经由 `lib/api`。

## 主题

深色模式为 class 方案（`<html class="dark">`），选择持久化在 `localStorage`
的 `halo-theme` 键，缺省跟随 `prefers-color-scheme`。`index.html` 中的内联脚本在
bundle 之前执行，避免主题闪烁。令牌集中在 `src/index.css`，自定义变量统一使用
`--halo-*` 前缀。

## 数据接入

当前页面数据来自 `src/lib/api/mock.ts` 的类型化样例。接入后端时，把
`lib/api/client.ts` 的 `listTickets()` 接到页面组件即可，组件本身无需改动。
开发服务器已把 `/api` 代理到后端 `https://127.0.0.1:443`（GraphQL 入口
`/api/query`）。
