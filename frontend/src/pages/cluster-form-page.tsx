import { useMemo } from "react"
import { IconArrowLeft } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { Button } from "@/components/ui/button"
import { ClusterForm, type ClusterFormValues } from "@/features/clusters/cluster-form"
import { useGraphQL } from "@/hooks/use-graphql"
import { CLUSTERS_QUERY, type ClusterNode } from "@/lib/api/queries"

export interface ClusterFormPageProps {
  mode: "create" | "edit"
  /** 编辑时带的群集 UUID */
  uuid?: string
  onBack: () => void
  onSaved: () => void
}

/** 群集接入/编辑页。表单是整页，不是弹窗。 */
export function ClusterFormPage({ mode, uuid, onBack, onSaved }: ClusterFormPageProps) {
  const remote = useGraphQL<
    { clusters: { edges: { node: ClusterNode }[] } | null },
    { first: number }
  >(CLUSTERS_QUERY, { first: 100 }, { clusters: null })

  // 编辑态的初始值由接口数据推导，数据没到位前不渲染表单。
  const initial = useMemo<(ClusterFormValues & { uuid: string }) | undefined>(() => {
    if (mode !== "edit" || !uuid) return undefined
    const node = remote.data?.clusters?.edges.find((edge) => edge.node.UUID === uuid)?.node
    if (!node) return undefined
    return {
      uuid: node.UUID,
      alias: node.Alias,
      host: node.Host,
      ip: node.IP,
      port: node.Port,
      user: node.User,
      password: "",
      status: node.Status === 1 ? 1 : 2,
      type: node.Type === "postgres" ? "postgres" : "mysql",
      groupUUID: node.GroupUUID ?? "",
    }
  }, [mode, uuid, remote.data])

  return (
    <>
      <PageHeader
        title={mode === "create" ? "接入群集" : "编辑群集"}
        description={
          mode === "create"
            ? "登记目标数据库所在的主机，密码由服务端加密保存"
            : "修改群集的接入信息，密码留空表示不修改"
        }
        action={
          <Button variant="outline" onClick={onBack}>
            <IconArrowLeft size={16} stroke={1.5} />
            返回群集列表
          </Button>
        }
      />

      <SectionCard
        description="密码由服务端加密保存；编辑时留空表示不修改原密码"
      >
        {mode === "edit" && !initial ? (
          <p className="text-sm text-muted-foreground">正在读取群集信息…</p>
        ) : (
          <ClusterForm
            mode={mode}
            initial={initial}
            onSaved={onSaved}
            onCancel={onBack}
          />
        )}
      </SectionCard>
    </>
  )
}
