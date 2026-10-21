const eventName = "halo:messages-changed"

/** 消息已读状态变化后通知顶栏与侧栏刷新未读数。 */
export function notifyMessagesChanged(): void {
  window.dispatchEvent(new Event(eventName))
}

export function subscribeMessagesChanged(listener: () => void): () => void {
  window.addEventListener(eventName, listener)
  return () => window.removeEventListener(eventName, listener)
}
