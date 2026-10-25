export type NoteTone = "success" | "danger"

// Long feedback can wrap the action row, so it belongs below the card.
export const NOTE_INLINE_MAX = 24

export function isInlineNote(text: string): boolean {
  return text.length <= NOTE_INLINE_MAX
}
