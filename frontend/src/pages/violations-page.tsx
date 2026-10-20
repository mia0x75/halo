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
import { violations } from "@/lib/mock/data"

const levelTone = {
  错误: "danger",
  警告: "warning",
  提示: "info",
} as const

export function ViolationsPage() {
  return (
    <>
      <PageHeader title="命中记录" description="规则命中明细，可按工单与规则追溯" />

      <SectionCard title="命中明细" description={`共 ${violations.length} 条`} bodyClassName="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>规则</TableHead>
              <TableHead>级别</TableHead>
              <TableHead>工单</TableHead>
              <TableHead>提示</TableHead>
              <TableHead>时间</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {violations.map((violation) => (
              <TableRow key={violation.id}>
                <TableCell className="font-mono text-xs font-medium">{violation.rule}</TableCell>
                <TableCell>
                  <Badge variant={levelTone[violation.level]}>{violation.level}</Badge>
                </TableCell>
                <TableCell className="text-muted-foreground">{violation.ticket}</TableCell>
                <TableCell>{violation.message}</TableCell>
                <TableCell className="text-muted-foreground">{violation.createdAt}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </SectionCard>
    </>
  )
}
