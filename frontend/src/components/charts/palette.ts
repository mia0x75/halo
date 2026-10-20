/**
 * Apex renders SVG, so it needs concrete colors instead of CSS variables.
 * Every chart in the app pulls from this module — never hard-code hex values
 * in feature files.
 */
export const chartPalette = {
  primary: "#066fd1",
  success: "#2fb344",
  warning: "#f59f00",
  danger: "#d63939",
  neutral: "#a3a3a3",
  text: "#404040",
  textMuted: "#737373",
} as const
