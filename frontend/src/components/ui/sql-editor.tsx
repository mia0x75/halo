import { useMemo } from "react"
import { syntaxHighlighting } from "@codemirror/language"
import { MySQL, PostgreSQL, sql } from "@codemirror/lang-sql"
import CodeMirror, { EditorView, oneDarkHighlightStyle } from "@uiw/react-codemirror"

import { useTheme } from "@/hooks/use-theme"
import { cn } from "@/lib/utils"

export type SqlDialect = "mysql" | "postgres"

export interface SqlEditorProps {
  value: string
  onChange: (value: string) => void
  /** 目标库类型决定关键字与函数补全，默认按 MySQL 处理。 */
  dialect?: SqlDialect
  placeholder?: string
  className?: string
}

const monoFont = 'ui-monospace, SFMono-Regular, Menlo, Consolas, "Liberation Mono", monospace'

/**
 * 编辑器底色交给外层卡片，避免 oneDark 自带背景和暗色主题打架；
 * 深色下只借用它的语法高亮配色，浅色沿用 CodeMirror 默认高亮。
 */
function editorTheme(dark: boolean) {
  return EditorView.theme(
    {
      "&": { backgroundColor: "transparent", color: "var(--halo-text)", fontSize: "13px" },
      "&.cm-focused": { outline: "none" },
      ".cm-content": {
        fontFamily: monoFont,
        padding: "10px 0",
        caretColor: "var(--halo-primary)",
      },
      ".cm-cursor, .cm-dropCursor": { borderLeftColor: "var(--halo-primary)" },
      ".cm-scroller": { lineHeight: "1.6" },
      ".cm-gutters": {
        backgroundColor: "transparent",
        border: "none",
        color: "var(--halo-text-muted)",
      },
      ".cm-activeLine": { backgroundColor: "var(--halo-surface-secondary)" },
      ".cm-activeLineGutter": { backgroundColor: "transparent" },
      ".cm-selectionBackground, &.cm-focused .cm-selectionBackground, .cm-content ::selection": {
        backgroundColor: "color-mix(in oklab, var(--halo-primary) 25%, transparent)",
      },
      ".cm-placeholder": { color: "var(--halo-text-muted)" },
      ".cm-tooltip": {
        backgroundColor: "var(--halo-surface)",
        border: "1px solid var(--halo-border)",
        color: "var(--halo-text)",
      },
      ".cm-tooltip-autocomplete > ul > li[aria-selected]": {
        backgroundColor: "var(--halo-surface-secondary)",
        color: "var(--halo-text)",
      },
      ".cm-panels": { backgroundColor: "var(--halo-surface)", color: "var(--halo-text)" },
    },
    { dark },
  )
}

/** 工单 SQL 编辑器：行号、括号匹配、方言高亮，跟随全局浅色/深色主题。 */
export function SqlEditor({ value, onChange, dialect = "mysql", placeholder, className }: SqlEditorProps) {
  const { theme } = useTheme()
  const dark = theme === "dark"

  const extensions = useMemo(
    () => [
      sql({ dialect: dialect === "postgres" ? PostgreSQL : MySQL }),
      EditorView.lineWrapping,
      editorTheme(dark),
      ...(dark ? [syntaxHighlighting(oneDarkHighlightStyle)] : []),
    ],
    [dialect, dark],
  )

  return (
    <div
      className={cn(
        "overflow-hidden rounded-sm border border-border bg-card transition-colors focus-within:border-ring focus-within:ring-3 focus-within:ring-ring/50",
        className,
      )}
    >
      <CodeMirror
        value={value}
        onChange={onChange}
        theme="none"
        extensions={extensions}
        placeholder={placeholder}
        aria-label="SQL 内容"
        basicSetup={{
          foldGutter: false,
          highlightActiveLineGutter: true,
          autocompletion: true,
          bracketMatching: true,
          closeBrackets: true,
        }}
        minHeight="14rem"
      />
    </div>
  )
}
