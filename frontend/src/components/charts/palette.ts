/**
 * 图表渲染为 SVG，需要具体色值；所有图表从此处取色，不要在业务组件里写死颜色。
 */
export const chartPalette = {
  primary: "#066fd1",
  success: "#2fb344",
  info: "#4299e1",
  warning: "#f59f00",
  danger: "#d63939",
  secondary: "#737373",
  neutral: "#a3a3a3",
  grid: "#d4d4d4",
  text: "#262626",
  textMuted: "#737373",
} as const
