/** 概览与报表共用的时段选项。 */
export const periods = [7, 30, 180, 365] as const

export type Period = (typeof periods)[number]

export const periodLabels: Record<Period, string> = {
  7: "近 7 天",
  30: "近 30 天",
  180: "近 180 天",
  365: "近 365 天",
}
