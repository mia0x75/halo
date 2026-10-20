import { IconChevronLeft, IconChevronRight } from "@tabler/icons-react"

import { Button } from "@/components/ui/button"

export interface PaginationProps {
  page: number
  totalPages: number
  total: number
  onPageChange: (page: number) => void
}

/** 分页：页码方块 2rem，边框 6px 圆角，当前页为主色。 */
export function Pagination({ page, totalPages, total, onPageChange }: PaginationProps) {
  const pages = Array.from({ length: totalPages }, (_, index) => index + 1)

  return (
    <div className="flex items-center justify-between gap-4 border-t border-border px-5 py-3">
      <span className="text-xs text-muted-foreground">共 {total} 条</span>
      <div className="flex items-center gap-1">
        <Button
          variant="outline"
          size="icon-sm"
          aria-label="上一页"
          disabled={page <= 1}
          onClick={() => onPageChange(page - 1)}
        >
          <IconChevronLeft size={16} stroke={1.5} />
        </Button>
        {pages.map((item) => (
          <Button
            key={item}
            variant={item === page ? "default" : "outline"}
            size="icon-sm"
            aria-label={`第 ${item} 页`}
            aria-current={item === page ? "page" : undefined}
            onClick={() => onPageChange(item)}
          >
            {item}
          </Button>
        ))}
        <Button
          variant="outline"
          size="icon-sm"
          aria-label="下一页"
          disabled={page >= totalPages}
          onClick={() => onPageChange(page + 1)}
        >
          <IconChevronRight size={16} stroke={1.5} />
        </Button>
      </div>
    </div>
  )
}
