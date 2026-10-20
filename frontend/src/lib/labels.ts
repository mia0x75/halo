/**
 * 显示值映射。接口只给编码，中文名统一放在前端，
 * 后端改动不会牵动界面文案。
 */

/** 规则严重级别：0 提示 1 错误 2 警告。 */
export const ruleLevelLabels: Record<number, string> = {
  0: "提示",
  1: "错误",
  2: "警告",
}

/** 规则分组，按规则名的三字前缀归类。 */
export const ruleGroupLabels: Record<string, string> = {
  CDB: "数据库",
  DDB: "数据库",
  MDB: "数据库",
  CTB: "表",
  MTB: "表",
  RTB: "表",
  DTB: "表",
  CIX: "索引",
  RIX: "索引",
  DEL: "增删改",
  INS: "增删改",
  RPL: "增删改",
  UPD: "增删改",
  SEL: "查询",
  CVW: "视图",
  DVW: "视图",
  MVW: "视图",
  CFU: "函数",
  DFU: "函数",
  MFU: "函数",
  CTG: "触发器",
  DTG: "触发器",
  MTG: "触发器",
  CEV: "事件",
  DEV: "事件",
  MEV: "事件",
  CSP: "存储过程",
  DSP: "存储过程",
  MSP: "存储过程",
  MSC: "杂项",
}

/** 取规则名里的分组前缀。 */
export function ruleGroupOf(name: string): string {
  return ruleGroupLabels[name.slice(0, 3)] ?? "其它"
}

/** 规则的适用范围：位掩码的低三位表示读/写/执行。 */
export function ruleScopeOf(bitwise: number): string {
  const parts: string[] = []
  if (bitwise & 4) parts.push("读")
  if (bitwise & 2) parts.push("写")
  if (bitwise & 1) parts.push("执行")
  return parts.join(" / ")
}

const logObjectLabels: Record<string, string> = {
  ticket: "工单",
  statement: "语句",
  comment: "审核意见",
  cluster: "群集",
  group: "分组",
  user: "用户",
  role: "角色",
  edge: "关联",
  rule: "规则",
  option: "系统选项",
  avatar: "头像",
  cron: "定时任务",
  task: "任务",
  glossary: "术语",
  metric: "指标",
  query: "语句模板",
  template: "模板",
}

const logVerbLabels: Record<string, string> = {
  create: "新建",
  update: "修改",
  delete: "删除",
  execute: "执行",
  replace: "替换",
}

export interface LogDisplay {
  /** 操作描述，例如「修改工单」 */
  label: string
  /** 操作对象，例如「工单 #00099054」 */
  target: string
}

/**
 * 解析日志里的结构化内容。审计日志是 "操作|对象|详情" 三段，
 * 登录日志之类的自由文本没有分隔符，原样展示。
 */
export function parseLogOperation(operation: string): LogDisplay {
  const parts = operation.split("|")
  if (parts.length < 2) {
    return { label: operation, target: "" }
  }
  const [code, rawTarget] = parts
  const [object, verb] = code.split(".")
  const objectLabel = logObjectLabels[object ?? ""] ?? object ?? "对象"
  const verbLabel = logVerbLabels[verb ?? ""] ?? verb ?? "变更"
  const target = rawTarget ?? ""
  const [targetType, targetId] = target.split(":")
  const targetLabel = logObjectLabels[targetType ?? ""] ?? targetType ?? ""
  const shortId = (targetId ?? "").slice(0, 8)
  return {
    label: `${verbLabel}${objectLabel}`,
    target: targetId ? `${targetLabel} #${shortId}` : targetLabel,
  }
}
