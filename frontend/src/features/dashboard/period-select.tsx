import { IconCalendar } from "@tabler/icons-react"

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

/** 时段选择：7d / 30d / 180d / 365d。 */
export function PeriodSelect({ value, onChange }: PeriodSelectProps) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="sm" className="gap-1.5">
          <IconCalendar size={14} stroke={1.5} />
          {periodLabels[value]}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {periods.map((period) => (
          <DropdownMenuItem key={period} onSelect={() => onChange(period)}>
            {periodLabels[period]}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
