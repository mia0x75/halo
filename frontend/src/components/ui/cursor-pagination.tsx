import { IconArrowBarToLeft, IconChevronLeft, IconChevronRight } from "@tabler/icons-react"

import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

export interface CursorPaginationProps {
  /** 当前页码，从 1 开始（由会话内的翻页次数得到，不查询总数） */
  page: number
  /** 当前页没有任何记录：只提示无结果，不显示页码与翻页按钮 */
  empty?: boolean
  hasPreviousPage: boolean
  hasNextPage: boolean
  onPrevious: () => void
  onNext: () => void
  /** 回到第一页，仅在非首屏显示 */
  onReset: () => void
  /** 独立展示时（不在卡片内）去掉上边框 */
  divided?: boolean
}

/**
 * 游标分页：前后翻页不重算偏移，也不统计 COUNT(*)。
 * 大表加复杂过滤时计数会拖垮查询，页码只反映当前会话翻了几次，
 * 因此不提供"跳到第 N 页"，只提供回首页。
 */
export function CursorPagination({
  page,
  empty = false,
  hasPreviousPage,
  hasNextPage,
  onPrevious,
  onNext,
  onReset,
  divided = true,
}: CursorPaginationProps) {
  if (empty) {
    return (
      <div
        className={cn(
          "px-5 py-3 text-sm text-muted-foreground",
          divided && "border-t border-border",
        )}
      >
        查询无记录
      </div>
    )
  }

  return (
    <div
      className={cn(
        "flex items-center justify-between gap-4 px-5 py-3",
        divided && "border-t border-border",
      )}
    >
      <span className="text-sm text-muted-foreground">第 {page} 页</span>
      <div className="flex items-center gap-2">
        <Button
          variant="outline"
          size="sm"
          disabled={page <= 1}
          onClick={onReset}
          aria-label="回到首页"
        >
          <IconArrowBarToLeft size={16} stroke={1.5} />
          回到首页
        </Button>
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
    </div>
  )
}
