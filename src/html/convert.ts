import QuickLRU from 'quick-lru'
import { fetchHtmlWithBrowser } from '../_libs/browser/index.ts'
import { htmlToMarkdown as convertToMarkdown, type HtmlToMarkdownOptions } from '../_libs/html-to-markdown/index.ts'

const htmlCache = new QuickLRU<string, string>({ maxSize: 100, maxAge: 60_000 })

async function fetchHtml(url: string) {
  const cached = htmlCache.get(url)
  if (cached !== undefined) {
    return cached
  }

  let html: string | null = await fetchHtmlWithBrowser(url)
  if (!html) {
    const resp = await fetch(url, { headers: { Accept: 'text/html' } })
    html = await resp.text()
  }

  if (html) {
    htmlCache.set(url, html)
  }

  return html
}

export async function htmlToMarkdown(url: string, options?: HtmlToMarkdownOptions) {
  const raw = await fetchHtml(url)
  return convertToMarkdown(raw, options)
}
