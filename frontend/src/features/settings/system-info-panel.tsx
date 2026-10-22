import { useState } from "react"
import { IconInfoCircle, IconMail } from "@tabler/icons-react"

import { SectionCard } from "@/components/shared/section-card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { useGraphQL } from "@/hooks/use-graphql"
import { graphqlRequest } from "@/lib/api/client"
import { SEND_TEST_MAIL_MUTATION, SYSTEM_INFO_QUERY, type SystemInfoNode } from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"

/** 平台运行信息：版本、元数据库与规模概览。 */
export function SystemInfoPanel() {
  const [testEmail, setTestEmail] = useState("")
  const [sending, setSending] = useState(false)
  const [testResult, setTestResult] = useState<{ ok: boolean; text: string } | null>(null)
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
        ["构建时间", info.BuildAt],
        ["Go 版本", info.GoVer],
        ["构建分支", info.Branch],
        ["元数据库", info.Database],
        ["邮件通道", info.MailEnabled ? "已启用" : "未启用"],
        ["群集数量", String(info.ClusterCount)],
        ["账号数量", String(info.UserCount)],
      ]
    : []

  const sendTestMail = async () => {
    setSending(true)
    setTestResult(null)
    try {
      await graphqlRequest(SEND_TEST_MAIL_MUTATION, { email: testEmail }, getToken())
      setTestResult({ ok: true, text: "测试邮件已交给 SMTP 通道发送。" })
    } catch (cause) {
      setTestResult({
        ok: false,
        text: cause instanceof Error ? cause.message : String(cause),
      })
    } finally {
      setSending(false)
    }
  }

  return (
    <SectionCard bodyClassName="p-5">
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

      <div className="mt-6 flex flex-col gap-2 border-t border-border/60 pt-5">
        <p className="text-sm font-medium text-strong">发送测试邮件</p>
        <p className="text-xs text-muted-foreground">
          SMTP 配置在「平台参数」里修改，保存后立即生效。
        </p>
        <div className="flex flex-wrap items-center gap-2">
          <Input
            value={testEmail}
            onChange={(event) => setTestEmail(event.target.value)}
            placeholder="收件人邮箱"
            className="h-8 max-w-64"
            type="email"
          />
          <Button
            size="sm"
            disabled={sending || testEmail.trim() === ""}
            onClick={() => void sendTestMail()}
          >
            <IconMail size={16} stroke={1.5} />
            发送测试邮件
          </Button>
        </div>
        {testResult ? (
          <p className={testResult.ok ? "text-sm text-success" : "text-sm text-danger"}>
            {testResult.text}
          </p>
        ) : null}
      </div>
    </SectionCard>
  )
}
