import { IconChevronLeft, IconChevronRight } from "@tabler/icons-react"

import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

export interface CursorPaginationProps {
  total: number
  hasPreviousPage: boolean
  hasNextPage: boolean
  onPrevious: () => void
  onNext: () => void
  /** 当前页第一条之前的游标，翻页时回退用 */
  startCursor?: string
  endCursor?: string
  /** 独立展示时（不在卡片内）去掉上边框 */
  divided?: boolean
}

/** 游标分页：只提供前后翻页，不显示页码。 */
export function CursorPagination({
  total,
  hasPreviousPage,
  hasNextPage,
  onPrevious,
  onNext,
  startCursor,
  endCursor,
  divided = true,
}: CursorPaginationProps) {
  return (
    <div
      className={cn(
        "flex items-center justify-between gap-4 px-5 py-3",
        divided && "border-t border-border",
      )}
    >
      <span className="text-xs text-muted-foreground">共 {total} 条</span>
      <div className="flex items-center gap-2">
        {startCursor ? (
          <span className="hidden font-mono text-[0.625rem] text-muted-foreground sm:inline">
            {startCursor}
          </span>
        ) : null}
        <Button
          variant="outline"
          size="sm"
          disabled={!hasPreviousPage}
          onClick={onPrevious}
          aria-label="上一页"
        >
          <IconChevronLeft size={16} stroke={1.5} />
          上一页
        </Button>
        <Button
          variant="outline"
          size="sm"
          disabled={!hasNextPage}
          onClick={onNext}
          aria-label="下一页"
        >
          下一页
          <IconChevronRight size={16} stroke={1.5} />
        </Button>
      </div>
      {endCursor ? (
        <span className="hidden font-mono text-[0.625rem] text-muted-foreground sm:inline">
          {endCursor}
        </span>
      ) : null}
    </div>
  )
}
