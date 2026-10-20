import { IconChevronDown } from "@tabler/icons-react"

import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { periodLabels, periods, type Period } from "@/lib/mock/data"

export interface PeriodSelectProps {
  value: Period
  onChange: (period: Period) => void
}

/** 卡片右上角的时段菜单：只留一个下拉小箭头，当前时段显示在卡片描述里。 */
export function PeriodSelect({ value, onChange }: PeriodSelectProps) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          aria-label={`时间范围：${periodLabels[value]}`}
          title={periodLabels[value]}
        >
          <IconChevronDown size={16} stroke={1.5} />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {periods.map((period) => (
          <DropdownMenuItem
            key={period}
            onSelect={() => onChange(period)}
            className={period === value ? "font-medium text-primary" : undefined}
          >
            {periodLabels[period]}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
