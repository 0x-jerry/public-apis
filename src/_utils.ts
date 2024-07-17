export async function get(url: URL | string, query?: Record<string, unknown>) {
  const u = new URL(url)

  if (query) {
    for (const [key, value] of Object.entries(query)) {
      u.searchParams.set(key, String(value))
    }
  }

  const resp = await fetch(url)

  return resp.text()
}
