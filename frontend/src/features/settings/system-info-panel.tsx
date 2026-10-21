import { IconInfoCircle } from "@tabler/icons-react"

import { SectionCard } from "@/components/shared/section-card"
import { Badge } from "@/components/ui/badge"
import { useGraphQL } from "@/hooks/use-graphql"
import { SYSTEM_INFO_QUERY, type SystemInfoNode } from "@/lib/api/queries"

/** 平台运行信息：版本、元数据库与规模概览。 */
export function SystemInfoPanel() {
  const remote = useGraphQL<{ systemInfo: SystemInfoNode }>(
    SYSTEM_INFO_QUERY,
    undefined,
    { systemInfo: null } as unknown as { systemInfo: SystemInfoNode },
  )

  const info = remote.data?.systemInfo ?? null

  const rows: [string, string][] = info
    ? [
        ["版本", info.Version],
        ["Git 提交", info.Git],
        ["构建时间", info.Compile],
        ["构建分支", info.Branch],
        ["元数据库", info.Database],
        ["邮件通道", info.MailEnabled ? "已启用" : "未启用"],
        ["群集数量", String(info.ClusterCount)],
        ["账号数量", String(info.UserCount)],
      ]
    : []

  return (
    <SectionCard title="运行信息" description="平台版本与规模概览" bodyClassName="p-5">
      {!info ? (
        <p className="flex items-center gap-2 text-sm text-muted-foreground">
          <IconInfoCircle size={16} stroke={1.5} />
          正在读取运行信息…
        </p>
      ) : (
        <dl className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
          {rows.map(([label, value]) => (
            <div key={label} className="flex flex-col gap-1">
              <dt className="text-xs text-muted-foreground">{label}</dt>
              <dd className="flex items-center gap-2 text-sm font-medium text-strong">
                {value}
                {label === "邮件通道" ? (
                  <Badge variant={info.MailEnabled ? "success" : "muted"}>
                    {info.MailEnabled ? "启用" : "关闭"}
                  </Badge>
                ) : null}
              </dd>
            </div>
          ))}
        </dl>
      )}
    </SectionCard>
  )
}
