import { IconChevronDown } from "@tabler/icons-react"

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { periodLabels, periods, type Period } from "@/features/dashboard/periods"

export interface PeriodSelectProps {
  value: Period
  onChange: (period: Period) => void
}

/**
 * 时段菜单。触发器是设计稿里的 dropdown-toggle：透明无边框的次要色文字 +
 * 小箭头，选中项高亮。
 */
export function PeriodSelect({ value, onChange }: PeriodSelectProps) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <button
          type="button"
          aria-label={`时间范围：${periodLabels[value]}`}
          className="inline-flex items-center gap-1 rounded-sm border-0 bg-transparent p-0 text-sm leading-4 text-muted-foreground transition-colors hover:text-foreground focus-visible:outline-none focus-visible:ring-3 focus-visible:ring-ring/50"
        >
          {periodLabels[value]}
          <IconChevronDown size={16} stroke={1.5} />
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {periods.map((period) => (
          <DropdownMenuItem
            key={period}
            onSelect={() => onChange(period)}
            className={period === value ? "bg-primary/5 font-medium" : undefined}
          >
            {periodLabels[period]}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
