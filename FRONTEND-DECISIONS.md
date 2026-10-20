# Frontend Stack & Design System — Decision Record

> **Scope:** project-agnostic frontend playbook. It contains no references to any
> specific product or repository and is intended to be copied into new projects.
>
> **Status:** accepted defaults. Deviate deliberately, not accidentally.
> **Pinned:** 2026-10 (adjust versions when adopting; pin them in `package.json`).

---

## 0. How to use this document

- Treat each section as a **default decision**, not a law. If you deviate, record
  why in the decision log at the end.
- All code snippets are copy-paste ready. Replace the placeholder project prefix
  `app-` with your own 2–3 character prefix (see §6).
- The goal is a **small, owned, type-safe, themeable** frontend — not a bag of
  runtime dependencies.

---

## 1. Principles

1. **Own the component code.** UI primitives are copied into the repo (shadcn
   model), not hidden behind an opaque runtime. You can read, patch, and delete
   them.
2. **Tokens first.** Components reference semantic tokens, never raw hex values.
   The only exception is chart SVG colors (§9), which need concrete values.
3. **Type-safe by default.** `strict` TypeScript, `noUnusedLocals`,
   `noUnusedParameters`, no `any` in application code.
4. **Dark mode is first-class.** Both themes ship from day one; no afterthought
   retrofit.
5. **Mock-first.** Build the UI against typed fixtures, then swap in the API.
6. **Accessible by construction.** Keyboard reachable, labeled icon buttons,
   visible focus rings, sufficient contrast.
7. **Deterministic tailwind.** Prefer explicit layout wrappers over fighting
   utility-ordering with `!important`.

---

## 2. Stack decisions

| Area | Decision | Version (reference) | Rationale | Rejected alternatives |
|------|----------|--------------------|-----------|-----------------------|
| Runtime | **React** | 19.x | Concurrent features, mature ecosystem, `use`/transitions | Vue, Svelte (viable, but team fluency) |
| Language | **TypeScript** | 7.x | Strict typing; TS 7 is faster | Plain JS, Flow |
| Build tool | **Vite** | 8.x (Rolldown) | Fast dev server, native CSS/vars, simple config | Webpack, Parcel |
| Package manager | **Bun** | 1.4.x | Fast install + script runner, single binary | npm/pnpm/yarn (work fine; Bun is the default) |
| Styling | **Tailwind CSS** | 4.x via `@tailwindcss/vite` | Utility-first, CSS-variable theming, no PostCSS needed | CSS-in-JS, Sass modules, vanilla CSS at scale |
| Components | **shadcn/ui** | `radix-nova` preset | Copy-in ownership, Radix a11y primitives | MUI, Ant Design, Chakra (heavier, opinionated runtime) |
| Primitives | **Radix UI** (`radix-ui`) | 1.x | Unstyled, accessible, composable | Headless UI (smaller surface) |
| Icons | **`@tabler/icons-react`** | 3.x | Huge consistent set, `stroke` control | Lucide (fine; keep one only) |
| Charts | **ApexCharts** (`apexcharts`) | 8.x | Rich dashboard charts, sparklines, no React wrapper lock-in | Recharts, Chart.js, ECharts (all viable; Apex chosen) |
| Font | **Inter** via `@fontsource/inter` | 5.x | Self-hosted, no FOUT from Google Fonts | System font stack (acceptable fallback) |
| Lint | **oxlint** | 1.x | Very fast, low config | ESLint + plugins (more setup, slower) |
| Class merge | **`cn`** (clsx + tailwind-merge) | — | Conditional classes without conflicts | Manual template strings |
| Animations | **`tw-animate-css`** | 1.x | Tailwind-native animation utilities | Framer Motion (add only when needed) |
| Routing | **None by default** | — | Single-page dashboards don't need it | React Router (add when multi-page) |
| State | **Local/`useState` by default** | — | Avoid premature global stores | Redux/Zustand (add when shared) |

**Rule:** exactly **one** icon library and **one** chart library per project.

---

## 3. Bootstrap

```bash
bun create vite@latest my-app --template react-ts
cd my-app

# runtime
bun add react react-dom \
  @tabler/icons-react apexcharts \
  class-variance-authority clsx tailwind-merge cn tw-animate-css \
  radix-ui @fontsource/inter

# dev
bun add -d tailwindcss @tailwindcss/vite oxlint typescript vite \
  @vitejs/plugin-react @types/react @types/react-dom @types/node
```

Recommended `package.json` scripts:

```json
{
  "scripts": {
    "dev": "vite",
    "build": "tsc -b && vite build",
    "lint": "oxlint",
    "preview": "vite preview"
  }
}
```

`vite.config.ts`:

```ts
import tailwindcss from "@tailwindcss/vite"
import react from "@vitejs/plugin-react"
import { fileURLToPath, URL } from "node:url"
import { defineConfig } from "vite"

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: { "@": fileURLToPath(new URL("./src", import.meta.url)) },
  },
})
```

> **No PostCSS config is required.** Tailwind v4 is wired through the Vite
> plugin. Do not add `postcss.config.js` / `tailwind.config.js` unless you have a
> concrete reason; configuration lives in CSS.

---

## 4. Directory layout

```
src/
  components/
    ui/          # copied shadcn primitives (owned)
    layout/      # app shell: navbar, footer, containers
    charts/      # palette + chart wrappers (ApexCharts)
    shared/      # cross-feature pieces: section-card, status-badge, ...
  features/
    <domain>/    # feature-scoped components (e.g. dashboard, settings)
  hooks/         # use-theme, etc.
  lib/
    api/         # types.ts (API contract), mock.ts (fixtures), client.ts (later)
    format.ts    # date/relative-time/number formatters
    utils.ts     # cn re-export (if not using an external `cn` package)
  pages/         # route-level composition
  App.tsx
  main.tsx
  index.css      # tokens + theme + Tailwind entry
```

**Convention:** `features/*` may import from `components/*` and `lib/*`;
`components/ui` must stay generic and domain-agnostic.

---

## 5. TypeScript configuration

Recommended `tsconfig.app.json` (bundler mode, no emit):

```json
{
  "compilerOptions": {
    "target": "ES2023",
    "lib": ["ES2023", "DOM"],
    "module": "ESNext",
    "moduleResolution": "bundler",
    "types": ["vite/client"],
    "jsx": "react-jsx",
    "noEmit": true,
    "verbatimModuleSyntax": true,
    "moduleDetection": "force",
    "allowImportingTsExtensions": true,
    "skipLibCheck": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "paths": { "@/*": ["./src/*"] }
  },
  "include": ["src"]
}
```

> **TypeScript 7 removed `baseUrl`** (`TS5102`). Declare `paths` directly, as
> above; do not add `baseUrl`.

---

## 6. Design tokens

### 6.1 Naming

- Use a **short project prefix** for every custom CSS variable: `--app-*`
  (replace `app`). Never reuse a UI vendor's prefix (e.g. `--tblr-*`); it
  collides on upgrades and reads as a dependency you don't have.
- Map shadcn's semantic tokens (`--background`, `--foreground`, `--card`, …)
  **onto** your tokens. Components consume the shadcn names; your theme owns the
  values.
- Expose semantic status colors as first-class tokens: `success`, `warning`,
  `danger`, `strong`.

### 6.2 Token sheet (copy-paste)

```css
@import "tailwindcss";
@import "tw-animate-css";
@import "@fontsource/inter/400.css";
@import "@fontsource/inter/500.css";
@import "@fontsource/inter/600.css";
@import "@fontsource/inter/700.css";

@custom-variant dark (&:is(.dark *));

@theme inline {
  --font-sans: "Inter", system-ui, -apple-system, "Segoe UI", Roboto, sans-serif;
  --font-heading: var(--font-sans);

  --color-background: var(--background);
  --color-foreground: var(--foreground);
  --color-card: var(--card);
  --color-card-foreground: var(--card-foreground);
  --color-popover: var(--popover);
  --color-popover-foreground: var(--popover-foreground);
  --color-primary: var(--primary);
  --color-primary-foreground: var(--primary-foreground);
  --color-secondary: var(--secondary);
  --color-secondary-foreground: var(--secondary-foreground);
  --color-muted: var(--muted);
  --color-muted-foreground: var(--muted-foreground);
  --color-accent: var(--accent);
  --color-accent-foreground: var(--accent-foreground);
  --color-destructive: var(--destructive);
  --color-border: var(--border);
  --color-input: var(--input);
  --color-ring: var(--ring);

  --color-success: var(--app-success);
  --color-warning: var(--app-warning);
  --color-danger: var(--app-danger);
  --color-strong: var(--app-text-strong);

  --radius-sm: 4px;
  --radius-md: 6px;
  --radius-lg: 8px;
  --radius-xl: 12px;
  --radius-2xl: 16px;
  --radius-3xl: 24px;
  --radius-4xl: 32px;
}

:root {
  /* surfaces & text */
  --app-canvas: #fafafa;
  --app-surface: #ffffff;
  --app-text: oklch(37.15% 0 0);          /* ~#404040 */
  --app-text-strong: #171717;
  --app-text-muted: oklch(55.55% 0 0);    /* ~#737373 */
  --app-border: color-mix(in oklab, oklch(26.86% 0 0) 11.9%, transparent);

  /* brand & status */
  --app-primary: #066fd1;
  --app-success: #2fb344;
  --app-warning: #f59f00;
  --app-danger: #d63939;

  /* geometry */
  --app-radius-card: 8px;
  --app-radius-control: 6px;

  /* shadcn semantic mapping */
  --background: var(--app-canvas);
  --foreground: var(--app-text);
  --card: var(--app-surface);
  --card-foreground: var(--app-text);
  --popover: var(--app-surface);
  --popover-foreground: var(--app-text);
  --primary: var(--app-primary);
  --primary-foreground: #ffffff;
  --secondary: oklch(96.7% 0 0);
  --secondary-foreground: var(--app-text);
  --muted: oklch(96.7% 0 0);
  --muted-foreground: var(--app-text-muted);
  --accent: oklch(96.7% 0 0);
  --accent-foreground: var(--app-text);
  --destructive: var(--app-danger);
  --border: var(--app-border);
  --input: var(--app-border);
  --ring: var(--app-primary);
  --radius: var(--app-radius-control);

  color-scheme: light;
}

.dark {
  --app-canvas: oklch(20.46% 0 0);        /* ~neutral-900 */
  --app-surface: oklch(26.86% 0 0);       /* ~neutral-800 */
  --app-text: oklch(92.19% 0 0);
  --app-text-strong: #ffffff;
  --app-text-muted: oklch(71.55% 0 0);
  --app-border: rgba(128, 150, 172, 0.2);

  --background: var(--app-canvas);
  --foreground: var(--app-text);
  --card: var(--app-surface);
  --card-foreground: var(--app-text);
  --popover: var(--app-surface);
  --popover-foreground: var(--app-text);
  --primary: var(--app-primary);
  --primary-foreground: #ffffff;
  --secondary: oklch(32% 0 0);
  --secondary-foreground: var(--app-text);
  --muted: oklch(32% 0 0);
  --muted-foreground: var(--app-text-muted);
  --accent: oklch(32% 0 0);
  --accent-foreground: var(--app-text);
  --destructive: var(--app-danger);
  --border: var(--app-border);
  --input: var(--app-border);
  --ring: var(--app-primary);
  --radius: var(--app-radius-control);

  color-scheme: dark;
}

@layer base {
  * { @apply border-border outline-ring/50; }
  html { @apply font-sans; }
  body { @apply bg-background text-foreground; }
}
```

### 6.3 Color reference

| Token | Light | Dark |
|-------|-------|------|
| canvas | `#fafafa` | `oklch(20.46% 0 0)` |
| surface | `#ffffff` | `oklch(26.86% 0 0)` |
| text | `oklch(37.15% 0 0)` | `oklch(92.19% 0 0)` |
| text-strong | `#171717` | `#ffffff` |
| text-muted | `oklch(55.55% 0 0)` | `oklch(71.55% 0 0)` |
| border | `color-mix(oklch(26.86%) 11.9%, transparent)` | `rgba(128,150,172,.2)` |
| primary | `#066fd1` | same |
| success | `#2fb344` | same |
| warning | `#f59f00` | same |
| danger | `#d63939` | same |

Status tints in components use opacity: `bg-success/10 text-success`,
`bg-warning/10 text-warning`, `bg-danger/10 text-danger`.

---

## 7. Dark mode strategy

- **Class-based**, toggled on `<html>`: `.dark`.
- Declare the variant once: `@custom-variant dark (&:is(.dark *));`
- Persist the choice in `localStorage`; fall back to
  `prefers-color-scheme: dark`.
- Inject a tiny inline script in `index.html` **before** the app bundle to avoid
  a flash of the wrong theme (FOUC).

`index.html` (head):

```html
<script>
  try {
    var t = localStorage.getItem("app-theme")
    if (t === "dark" || (t === null && matchMedia("(prefers-color-scheme: dark)").matches)) {
      document.documentElement.classList.add("dark")
    }
  } catch (e) {}
</script>
```

`src/hooks/use-theme.ts`:

```ts
import { useCallback, useEffect, useState } from "react"

export type Theme = "light" | "dark"
const STORAGE_KEY = "app-theme"

function readTheme(): Theme {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored === "light" || stored === "dark") return stored
  return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light"
}

export function useTheme() {
  const [theme, setThemeState] = useState<Theme>(readTheme)

  useEffect(() => {
    document.documentElement.classList.toggle("dark", theme === "dark")
    localStorage.setItem(STORAGE_KEY, theme)
  }, [theme])

  const toggleTheme = useCallback(
    () => setThemeState((p) => (p === "dark" ? "light" : "dark")),
    [],
  )

  return { theme, setTheme: setThemeState, toggleTheme }
}
```

---

## 8. Visual language (enterprise dashboard)

The aesthetic is a **dense, neutral, enterprise admin UI** (inspired by the
Tabler dashboard language): flat surfaces, hairline borders, one saturated
primary, and restrained status colors.

- **Cards:** 8px radius, `1px` border (`--app-border`), white/dark surface,
  `shadow-xs` (`0 1px 2px 0 #0000000d`). No heavy shadows.
- **Controls:** 6px radius; inputs/dialogs follow the same border token.
- **Section header (subheader):** `0.75rem`, weight `500`, uppercase,
  `letter-spacing: 0.04em`, muted color.
- **Icons:** one set, `stroke={1.5}` consistently (most icon libs default to 2).
  Size 16 for inline/buttons, 20–24 for navigation.
- **Layout:** centered container `max-w-[1320px]`, `px-4`, vertical rhythm
  `space-y-4`.
- **Density:** tables use compact rows; avoid oversized padding.
- **Typography scale (Tailwind defaults):** `xs 12 / sm 14 / base 16 / lg 18 /
  2xl 24`. Keep headings semibold (`font-semibold tracking-tight`).
- **Elevation:** prefer borders over shadows; reserve shadows for popovers.
- **Motion:** short (100–150ms) fades/zooms; respect
  `prefers-reduced-motion`.

> ⚠️ If migrating from a Tailwind v3 project that **remapped** the font scale
> (e.g. `xs=11, sm=13, base=14`), re-map deliberately: old `base → text-sm`,
> old `lg → text-base`. Do not assume `text-base` means the same pixels.

---

## 9. Charts (ApexCharts)

- Build **thin wrappers** (`Sparkline`, `TrendChart`) instead of scattering chart
  options through features.
- Apex renders SVG and needs **concrete colors**. Keep a single palette module:

  ```ts
  export const chartPalette = {
    primary: "#066fd1",
    success: "#2fb344",
    warning: "#f59f00",
    danger: "#d63939",
    neutral: "#a3a3a3",
    text: "#404040",
    textMuted: "#737373",
  }
  ```

- Theme grid/axis/tooltip with **CSS overrides** so they follow light/dark
  automatically:

  ```css
  @layer components {
    .apexcharts-gridline,
    .apexcharts-xaxis line,
    .apexcharts-yaxis line { stroke: var(--app-border); }
    .apexcharts-text { fill: var(--app-text-muted); }
    .apexcharts-tooltip {
      background: var(--app-surface);
      border-color: var(--app-border);
      color: var(--app-text);
    }
    .apexcharts-tooltip-title {
      background: var(--app-surface);
      border-color: var(--app-border);
    }
  }
  ```

- Use `sparkline` for KPI cards, `area` for coverage/trends, `bar` with
  `stacked: true` for severity breakdowns.
- Pass `colors` from the palette; never hard-code colors in feature files.

---

## 10. shadcn/ui operational notes

Initialize with `bunx shadcn@latest init`. Recommended `components.json`:

```json
{
  "style": "radix-nova",
  "rsc": false,
  "tsx": true,
  "tailwind": {
    "config": "",
    "css": "src/index.css",
    "baseColor": "neutral",
    "cssVariables": true,
    "prefix": ""
  },
  "iconLibrary": "lucide",
  "aliases": {
    "components": "@/components",
    "utils": "@/lib/utils",
    "ui": "@/components/ui",
    "lib": "@/lib",
    "hooks": "@/hooks"
  }
}
```

Gotchas learned the hard way:

1. **`iconLibrary` is sticky.** `components.json` says `lucide`, and every
   `shadcn add` injects `lucide-react` imports even if you use another icon set.
   Either change the field (if supported) or **fix imports immediately** after
   adding a component.
2. **Adding over existing files prompts for overwrite.** `--yes` does not answer
   it. Pipe the answer:
   ```bash
   printf 'n\n' | bunx shadcn@latest add dialog
   ```
   Do this when you have locally patched `button.tsx` etc., to avoid clobbering.
3. **Generated code is yours.** Patch icon sizes/props and imports as needed;
   keep a note so a future `add` doesn't silently revert them.
4. **Do not import `cn` from `lucide`/shadcn assumptions.** This playbook uses
   `import { cn } from "cn"`. If you prefer the conventional helper, keep
   `src/lib/utils.ts` exporting `cn` and update `components.json.aliases.utils`.
5. **Card spacing is CSS-variable driven.** shadcn cards use
   `py-(--card-spacing)` / `gap-(--card-spacing)`. Trying to override with
   utilities can lose the cascade battle (see §11). For table/chart cards, use a
   purpose-built shell:
   ```tsx
   <div className="flex flex-col overflow-hidden rounded-lg border border-border bg-card shadow-xs">
     <div className="flex min-h-14 items-center justify-between gap-2 border-b border-border px-5">
       <h3 className="text-base font-semibold">{title}</h3>
       {action}
     </div>
     <div className={bodyClassName}>{children}</div>
   </div>
   ```
   This is deterministic and avoids utility-order surprises.

---

## 11. Tailwind v4 gotchas

- **Spacing scale:** `p-4 = 1rem`, `gap-4 = 1rem`, `size-4 = 1rem` (0.25rem
  base). Re-check any pixel values copied from older setups.
- **Utility ordering:** utilities of the same property can appear in an
  unexpected order. Prefer semantic wrappers and explicit values over stacking
  `p-0` against `px-(--var)`.
- **Theme in CSS:** use `@theme inline` for token exposure; there is no
  `tailwind.config.js` in v4.
- **No PostCSS config:** `@tailwindcss/vite` handles it.
- **Content detection is automatic;** avoid dynamic class construction — use
  explicit maps for status → class.
- **Arbitrary properties** (`[--card-spacing:0]`) are fragile when duplicate
  declarations exist; prefer a component wrapper.

---

## 12. Data layer

- **Mock-first.** Ship a typed fixture module (`lib/api/mock.ts`) with realistic
  shapes, then replace it with a client when the backend is ready.
- **Mirror the API contract** in `lib/api/types.ts` (camelCase on the client;
  match the server's JSON tags). Keep one source of truth per entity.
- Keep pure formatters in `lib/format.ts` (relative time, byte sizes, numbers).
  Example relative-time contract: `just now`, `N mins ago`, `N hours ago`, then
  `YYYY-MM-DD HH:mm`.
- Never let components reach into the network directly; route through
  `lib/api`.

---

## 13. Accessibility & UX baseline

- Every icon-only button gets an accessible name (`aria-label` or `sr-only`).
- Visible focus ring: `focus-visible:ring-3 focus-visible:ring-ring/50`.
- Maintain contrast in **both** themes for text and status tints.
- Use semantic elements (`button`, `table`, `label htmlFor`) before ARIA.
- Respect `prefers-reduced-motion` for transitions and chart animations.
- Ensure the shell uses `min-h-svh` so the footer stays at the bottom.

---

## 14. Quality gates (pre-merge)

1. `tsc -b` (or `bun run build`) passes with zero type errors.
2. `oxlint` passes (treat warnings as debt; fix new ones).
3. No new `lucide-react` / duplicate icon library imports.
4. No raw hex colors outside `charts/palette.ts` and `index.css`.
5. Both light and dark themes visually checked.
6. Icon strokes are consistently `1.5`.
7. Keyboard traversal works for dialogs, menus, and tables.

---

## 15. Decision log template

For each meaningful deviation, append:

```md
### ADR-00X — <decision>
- Status: proposed | accepted | superseded
- Date: YYYY-MM-DD
- Context: <forces at play>
- Decision: <what we chose>
- Consequences: <trade-offs, follow-ups>
```

---

## 16. Deliberate non-decisions

These are intentionally **out of scope** until a concrete need appears:

- Router / multi-page navigation
- Global state management
- Server state / query cache
- Form validation library
- i18n framework
- Unit/component test runner
- SSR
- Design-token JSON pipeline (Figma)

Adding any of these should be its own ADR, not an implicit dependency.

---

## 17. 项目决策日志（Halo）

> 本节记录 Halo 项目对上文默认决策的取舍。通用默认仍以 §1–§16 为准。

### ADR-001 — 设计令牌前缀改为 `--halo-*`
- Status: accepted
- Date: 2026-10-20
- Context: §6.1 要求使用 2–3 字符项目前缀替换占位的 `app-`。
- Decision: 全部自定义变量改用 `--halo-*`（如 `--halo-canvas`、`--halo-border`、
  `--halo-danger-foreground`），shadcn 语义变量（`--background`、`--card` 等）保持原名，
  继续映射到 `--halo-*`。
- Consequences: 组件只消费语义变量；新增令牌一律带 `--halo-` 前缀。

### ADR-002 — 包管理器使用 bun
- Status: accepted
- Date: 2026-10-20
- Context: §2 首选 bun 1.4.x；构建环境已安装 bun 1.4.2，且无 node/npm。
- Decision: 使用 bun 作为包管理器与脚本运行器，锁文件为 `frontend/bun.lock`。
- Consequences: 仓库根 Makefile 增加 `frontend-*` 系列目标，内部调用 `bun run`。

### ADR-003 — 路由、全局状态、表单库、表格库、测试运行器暂不引入
- Status: proposed（2026-10-20 人类确认暂不决定，冲突保留至出现真实需求时再处理）
- Date: 2026-10-20
- Context: §16 将 Router / 全局状态 / 服务端状态缓存 / 表单校验 / 组件测试运行器列为
  「有意不做」。而 `AGENTS.md` 技术栈表列出了 urql + Zustand、React Router 7、
  react-hook-form + zod、TanStack Table、Vitest + Testing Library，两者存在冲突。
- Decision: 初始化阶段遵循 §16，不安装上述依赖；页面切换用 `useState` 本地状态，
  GraphQL 请求用 `lib/api/client.ts` 的 fetch 传输层。
- Consequences: 若确认以 `AGENTS.md` 为准，需要补充依赖并新增 ADR 覆盖本条；
  接入真实工单列表、登录态与多页面导航时必须先解决该冲突。

### ADR-004 — 基础组件自研最小集，不直接 `shadcn add`
- Status: accepted
- Date: 2026-10-20
- Context: §10 提到 `shadcn add` 会按 `components.json` 的 `iconLibrary` 注入
  `lucide-react`，与本项目唯一图标库 `@tabler/icons-react` 冲突。
- Decision: 先手写 `button/card/badge/input/label/table/separator` 六个通用组件，
  语义严格对齐 §8 的密集型后台风格；`components.json` 的 `iconLibrary` 设为 `tabler`，
  后续 `shadcn add` 需逐个检查导入。
- Consequences: 组件代码完全自持、可读可改；缺少 dialog/dropdown/select 等组件，
  需要时再从 shadcn 引入并立刻修正图标导入。

### ADR-005 — 字体仅打包 latin 子集
- Status: accepted
- Date: 2026-10-20
- Context: `@fontsource/inter` 默认 CSS 会引入全部语言子集，产物体积偏大。
- Decision: `src/index.css` 只引入 `@fontsource/inter/latin-{400,500,600,700}.css`。
- Consequences: 界面文案保持简体中文/英文；如需其他语种需补引子集。

### ADR-006 — 图表封装自持，ApexCharts 静态引入
- Status: accepted
- Date: 2026-10-20
- Context: §9 要求薄封装；ApexCharts v8 无官方 React 包装器，且体量较大。
- Decision: `components/charts/apex-chart.tsx` 自持实例生命周期，
  `sparkline.tsx` / `trend-chart.tsx` 为业务无关的薄封装；
  调色板集中在 `charts/palette.ts`，主题相关网格/坐标轴/提示框用 CSS 覆盖。
- Consequences: 当前产物约 1.19 MB（gzip 355 KB），ApexCharts 占主要部分；
  若首屏体积成为问题，改为动态 `import("apexcharts")` 拆包。

### ADR-007 — 视觉基准取自 Tabler，令牌统一用 `--halo-*`
- Status: accepted
- Date: 2026-10-20
- Context: 界面需要一套密集型企业后台的视觉规范。Tabler 的设计稿（preview.tabler.io）
  给出了成体系的尺寸、间距、圆角、阴影与配色，但它的实现基于 Bootstrap；
  本项目不用 Bootstrap，也不引入其 CSS 变量命名。
- Decision: **从 Tabler 的样式表提取确切数值，用 Tailwind CSS v4 重新实现组件**，
  变量一律命名为 `--halo-*`。落地要点：
  - 配色取其调色板（primary `oklch(54.6% 0.1724 254.2deg)`、success/info/warning/danger、
    gray-50…900），页面底色 `oklch(98.51% 0 0)`、卡片 `#ffffff`、描边 `oklch(92.19% 0 0)`；
  - 圆角 4/6/8/16px，卡片阴影 `0 1px 3px 0 rgb(18 18 23 / .1), 0 1px 2px 0 rgb(18 18 23 / .06)`；
  - 正文 14px、按钮 12px 字号（小号 0.3125rem/0.5rem）、卡片内边距 1.25rem、
    徽标 10px 字号、导航项 0.5rem/1rem、进度条 8px、表格头 12px 大写。
  - 组件代码自持（shadcn 模式），不复制 Bootstrap 样式表。
- Consequences: 界面观感与设计稿一致而不引入 Bootstrap 依赖；数值全部落在
  `src/index.css` 的令牌里，后续微调只改一处。

### ADR-009 — 先接少量 gqlapi，失败回落假数据
- Status: accepted
- Date: 2026-10-20
- Context: 后端接口已就绪（指标、命中明细、站内消息、任务、分组），但界面需要
  先按设计稿成型并可独立演示。
- Decision: **两者并存**：接口优先，取不到就回落假数据。
  - 传输层 `lib/api/client.ts` 只做一件事——发查询、取 data，不引查询缓存；
  - 查询文档集中在 `lib/api/queries.ts`，字段沿用后端的 PascalCase 契约，
    页面层再转成自己的模型（`lib/api/mappers.ts`）；
  - 令牌存 `localStorage` 的 `halo-token`，请求带 `Authentication` 头
    （后端从 `secret.jwt.token_name` 读取，默认即 Authentication）；
  - 顶栏的「演示数据 / 已连接」按钮提供登录入口，未登录时界面照常可用。
  - 当前已接入：`me`、`metrics`、`ruleHitLevels`、`ruleHitRules`、
    `unreadMessageCount`、`groups`、`tickets`。
- Consequences: 界面既能独立演示，也能逐步换成真实数据；接口不可用时不会白屏；
  参数名大小写必须与 schema 一致（`metrics(Scope:)`、`ruleHitLevels(from:)`），
  写错会在开发期直接报 GraphQL 校验错误。

### ADR-008 — 演示阶段只用假数据，不连接 gqlapi
- Status: accepted
- Date: 2026-10-20
- Context: 后端接口已就绪，但界面需要先按设计稿成型并可独立演示。
- Decision: 数据统一放在 `src/lib/mock/data.ts`，形态按 GraphQL 契约的 camelCase 镜像；
  删掉原 `src/lib/api/client.ts`（gqlapi 传输层），避免留下未接线的死代码。
- Consequences: 页面可离线演示；接入后端时按页替换取值来源，组件结构不变。
