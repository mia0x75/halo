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

### 字体

界面默认字体是**寒蝉团圆体 Sans（Chill Reunion Sans）**，自持于
`src/assets/fonts/ChillReunionSans.woff2`。原始字体为 2.1MB 的 OTF，已按界面
实际用到的字符（含 ASCII 与全部中文文案）子集化到 154KB；后端返回的生僻字会
回落到系统字体。字体权重声明为 100–900（源文件只有一档字重）。

### 字号基准

以卡片为准：**标题 16px/500**（`text-base font-medium`）、**副标题 14px**
（`text-sm`）。其余按此基准对齐：

| 用途 | 字号 |
|------|------|
| 页面标题 | 20px/600 |
| 卡片标题 / 副标题 | 16px/500 / 14px |
| 正文、次要文字、表单标签 | 14px |
| 按钮 | 14px（小号 12px） |
| 表头 | 12px 大写、字距 0.04em |
| 徽标、侧栏计数 | 10px |

指标卡按设计稿的 Revenue 卡片组织三行：标签行（右侧可放操作）、数值行（数值与
环比同行、基线对齐）、图表行（独占一块贴卡片底部）。

其余关键尺寸与设计稿一致：卡片内边距 1.25rem、按钮内边距 0.5625rem/1rem、
徽标内边距 0.25rem/0.125rem、导航项 0.5rem/1rem、进度条 8px 高。

## 分页

统一使用**游标分页**：游标是 `base64("cursor<id>")`，与后端 `EncodeCursor` 一致，
只做前后翻页、不算偏移。工具在 `src/lib/mock/cursor.ts`，状态封装在
`hooks/use-cursor-pagination.ts`（过滤条件变化时 `reset()` 回到第一页）。

## 深色模式

class 方案（`<html class="dark">`），选择持久化在 `localStorage` 的 `halo-theme`
键，缺省跟随 `prefers-color-scheme`。`index.html` 中的内联脚本在 bundle 之前执行，
避免主题闪烁。

## 数据

数据有���个来源：

- **接口**（已登录时）：`src/lib/api/` 下的客户端与查询文档，目前接入
  `me`、`metrics`、`ruleHitLevels`、`ruleHitRules`、`unreadMessageCount`、
  `groups`、`tickets`；
- **演示数据**：`src/lib/mock/data.ts`，接口未登录、超时或返回空时自动回落。

顶栏右侧的「演示数据 / 已连接」按钮可以打开登录框，登录后令牌存在
`localStorage`（`halo-token`），请求带在 `Authentication` 头里（对应后端配置
`secret.jwt.token_name`）。

后端地址默认 `https://127.0.0.1:443`，本地测试可用环境变量覆盖：

```bash
HALO_API=https://127.0.0.1:8443 bun run dev
```
