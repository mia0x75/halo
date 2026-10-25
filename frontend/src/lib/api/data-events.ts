/**
 * Shell-wide change notifications. The navigation badges live in the layout,
 * which stays mounted while pages come and go, so a mutation on a page has to
 * announce itself or the badge keeps the count it fetched on mount.
 */
export type DataTopic = "messages" | "tickets"

const eventName = (topic: DataTopic) => `halo:data-changed:${topic}`

/** Announce that data behind a topic changed and subscribed counters are stale. */
export function notifyDataChanged(topic: DataTopic): void {
  window.dispatchEvent(new Event(eventName(topic)))
}

export function subscribeDataChanged(topic: DataTopic, listener: () => void): () => void {
  window.addEventListener(eventName(topic), listener)
  return () => window.removeEventListener(eventName(topic), listener)
}
