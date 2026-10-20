import { useMemo, useState } from "react"
import { IconDatabase } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { Badge } from "@/components/ui/badge"
import { Label } from "@/components/ui/label"
import { Select } from "@/components/ui/select"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { useClusters } from "@/hooks/use-clusters"
import { useGraphQL } from "@/hooks/use-graphql"
import { DATABASES_QUERY } from "@/lib/api/queries"

interface DatabaseNode {
  Name: string
  Charset: string
  Collate: string
}

/** 数据库页：按群集列出目标实例上的用户库。 */
export function DatabasesPage() {
  const clusters = useClusters(50)
  const [clusterUUID, setClusterUUID] = useState("")

  const selected = clusterUUID || clusters.clusters[0]?.uuid || ""

  const remote = useGraphQL<{ databases: DatabaseNode[] | null }, { clusterUUID: string }>(
    DATABASES_QUERY,
    { clusterUUID: selected },
    { databases: null },
  )

  const databases = useMemo(() => remote.data?.databases ?? [], [remote.data])

  return (
    <>
      <PageHeader
        title="数据库"
        description="目标实例上的用户库与字符集"
        action={
          <Badge variant="default" className="px-2 py-1">
            <IconDatabase size={14} stroke={1.5} className="mr-1" />
            {databases.length} 个库
          </Badge>
        }
      />

      <SectionCard title="用户库" description="选择群集后读取目标实例的库清单" bodyClassName="p-0">
        <div className="flex flex-wrap items-end gap-3 border-b border-border px-5 py-4">
          <div className="flex flex-col gap-1.5 sm:w-72">
            <Label htmlFor="database-cluster">目标群集</Label>
            <Select
              id="database-cluster"
              value={selected}
              onChange={(event) => setClusterUUID(event.target.value)}
            >
              {clusters.clusters.map((cluster) => (
                <option key={cluster.uuid} value={cluster.uuid}>
                  {cluster.alias} · {cluster.type} ({cluster.ip}:{cluster.port})
                </option>
              ))}
            </Select>
          </div>
        </div>

        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>库名</TableHead>
              <TableHead>字符集</TableHead>
              <TableHead>排序规则</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {databases.length === 0 ? (
              <TableRow>
                <TableCell colSpan={3} className="py-8 text-center text-sm text-muted-foreground">
                  {remote.loading ? "正在读取…" : "该群集没有可见的数据库。"}
                </TableCell>
              </TableRow>
            ) : null}
            {databases.map((database) => (
              <TableRow key={database.Name}>
                <TableCell className="font-medium">{database.Name}</TableCell>
                <TableCell className="text-muted-foreground">{database.Charset}</TableCell>
                <TableCell className="text-muted-foreground">{database.Collate}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </SectionCard>
    </>
  )
}
