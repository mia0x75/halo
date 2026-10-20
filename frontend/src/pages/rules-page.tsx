import { IconShieldCheck } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { Badge } from "@/components/ui/badge"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { rules } from "@/lib/mock/data"

const levelTone = {
  错误: "danger",
  警告: "warning",
  提示: "info",
} as const

export function RulesPage() {
  return (
    <>
      <PageHeader
        title="规则库"
        description={`共 ${rules.length} 条规则，命中明细按规则与级别聚合`}
        action={
          <Badge variant="default" className="px-2 py-1">
            <IconShieldCheck size={14} stroke={1.5} />
            248 条已启用
          </Badge>
        }
      />

      <SectionCard title="规则清单" bodyClassName="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>规则名</TableHead>
              <TableHead>分组</TableHead>
              <TableHead>级别</TableHead>
              <TableHead>状态</TableHead>
              <TableHead>说明</TableHead>
              <TableHead className="text-right">近 7 天命中</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {rules.map((rule) => (
              <TableRow key={rule.name}>
                <TableCell className="font-mono text-xs font-medium">{rule.name}</TableCell>
                <TableCell className="text-muted-foreground">{rule.group}</TableCell>
                <TableCell>
                  <Badge variant={levelTone[rule.level]}>{rule.level}</Badge>
                </TableCell>
                <TableCell className="text-muted-foreground">
                  {rule.enabled ? "已启用" : "已停用"}
                </TableCell>
                <TableCell className="text-muted-foreground">{rule.description}</TableCell>
                <TableCell className="text-right font-medium">{rule.hitCount}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </SectionCard>
    </>
  )
}
