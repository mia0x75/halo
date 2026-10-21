import { useCallback, useEffect, useState } from "react"

export type Theme = "light" | "dark"

const STORAGE_KEY = "halo-theme"

/**
 * 主题是全局状态：顶栏、登录页、SQL 编辑器各自调用 useTheme，
 * 一处切换必须让所有实例同步，否则页面变深色而编辑器还停在浅色。
 */
let current: Theme | null = null
const listeners = new Set<() => void>()

function readTheme(): Theme {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored === "light" || stored === "dark") return stored
  return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light"
}

function getTheme(): Theme {
  current ??= readTheme()
  return current
}

function applyTheme(theme: Theme): void {
  document.documentElement.classList.toggle("dark", theme === "dark")
  localStorage.setItem(STORAGE_KEY, theme)
  for (const listener of listeners) {
    listener()
  }
}

export function useTheme() {
  const [theme, setThemeState] = useState<Theme>(getTheme)

  useEffect(() => {
    const listener = () => setThemeState(getTheme())
    listeners.add(listener)
    return () => {
      listeners.delete(listener)
    }
  }, [])

  useEffect(() => {
    current = theme
    applyTheme(theme)
  }, [theme])

  const toggleTheme = useCallback(
    () => setThemeState((prev) => (prev === "dark" ? "light" : "dark")),
    [],
  )

  return { theme, setTheme: setThemeState, toggleTheme }
}
