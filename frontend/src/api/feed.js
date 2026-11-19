import axios from 'axios'

export function validatePage(data, cursor) {
  if (!data || !Array.isArray(data.items) || typeof data.has_more !== 'boolean' || !Number.isSafeInteger(data.next_cursor) || data.next_cursor < 0) throw new Error('Invalid feed response')
  for (const video of data.items) {
    if (!Number.isSafeInteger(video.id) || video.id <= 0 || typeof video.title !== 'string' || typeof video.video_url !== 'string') throw new Error('Invalid video')
  }
  if (data.has_more && (!data.items.length || data.next_cursor <= 0 || (cursor > 0 && data.next_cursor >= cursor))) throw new Error('Invalid feed cursor')
  return data
}

export async function fetchFeedPage(cursor = 0, limit = 8, signal) {
  const response = await axios.get('/api/v1/feed', { params: { cursor, limit }, timeout: 10000, signal })
  if (response.data.code !== 0) throw new Error('Feed request failed')
  return validatePage(response.data.data, cursor)
}

export function playbackURL(value) {
  try {
    const url = new URL(value, 'http://localhost:8000')
    if (!['http:', 'https:'].includes(url.protocol)) return ''
    if (['localhost', '127.0.0.1'].includes(url.hostname) && url.pathname.startsWith('/uploads/videos/')) return url.pathname + url.search
    return url.href
  } catch { return '' }
}
