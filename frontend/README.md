# Halo 前端

SQL 变更工单与审计平台的 Web 前端。技术栈与设计规范见仓库根目录
[FRONTEND-DECISIONS.md](../FRONTEND-DECISIONS.md)，视觉基准（尺寸、间距、圆角、
阴影、配色）取自 Tabler 设计稿，但实现全部为 Tailwind CSS v4 自持组件，
不引入 Bootstrap，也不使用其 CSS 变量命名。

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
    ui/          # 基础组件：button/badge/table/input/select/label
                 # dropdown-menu/tabs/cursor-pagination/progress/avatar/status-dot
    layout/      # 应用外壳：side-nav、top-nav、page-header、nav-items
    charts/      # 图表封装：ApexChart 与调色板
    shared/      # 跨业务组件：section-card、status-badge、empty-state、subheader
  features/
    dashboard/   # 业务域组件：stat-card、top-rules、period-select
  lib/
    mock/        # 演示假数据（data.ts）与游标分页工具（cursor.ts）
    utils.ts     # cn 再导出
  pages/         # 页面级组合：dashboard、tickets、clusters、rules、violations、logs
  App.tsx
  main.tsx
  index.css      # 设计令牌 + 主题 + Tailwind 入口
```

约定：`features/*` 可依赖 `components/*` 与 `lib/*`；`components/ui` 保持通用、
不掺业务语义。组件不得直接发起网络请求，一律经由 `lib`。

## 视觉基准

设计令牌集中在 `src/index.css`，全部使用 `--halo-*` 前缀：

| 令牌 | 值 | 用途 |
|------|-----|------|
| `--halo-primary` | `oklch(54.6% 0.1724 254.2deg)` | 主色 |
| `--halo-success / info / warning / danger` | 对应状态色 | 状态语义 |
| `--halo-canvas` | `oklch(98.51% 0 0)` | 页面底色（深色模式 `oklch(20.46% 0 0)`） |
| `--halo-surface` | `#ffffff` | 卡片与导航底色（深色模式 `oklch(26.86% 0 0)`） |
| `--halo-border` | `oklch(92.19% 0 0)` | 1px 描边 |
| `--halo-radius / -sm / -lg` | `6px / 4px / 8px` | 圆角 |
| `--halo-shadow` | `0 1px 3px 0 rgb(18 18 23 / .1), 0 1px 2px 0 rgb(18 18 23 / .06)` | 卡片投影 |

其余关键尺寸与设计稿一致：正文 14px、按钮 12px 字号（小号 8/4px 内边距）、
卡片内边距 1.25rem、徽标 10px 字号、导航项 0.5rem/1rem 内边距、进度条 8px 高。

## 分页

统一使用**游标分页**：游标是 `base64("cursor<id>")`，与后端 `EncodeCursor` 一致，
只做前后翻页、不算偏移。工具在 `src/lib/mock/cursor.ts`，状态封装在
`hooks/use-cursor-pagination.ts`（过滤条件变化时 `reset()` 回到第一页）。

## 深色模式

class 方案（`<html class="dark">`），选择持久化在 `localStorage` 的 `halo-theme`
键，缺省跟随 `prefers-color-scheme`。`index.html` 中的内联脚本在 bundle 之前执行，
避免主题闪烁。

## 数据

当前全部为 `src/lib/mock/data.ts` 的演示假数据，**不连接 gqlapi**。数据形态按
GraphQL 契约的 camelCase 镜像，接入时把取值换成接口返回即可，页面组件无需改动。
