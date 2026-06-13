import { fetchHtmlWithBrowser } from '../_libs/browser/index.ts'
import { htmlToMarkdown as convertToMarkdown, type HtmlToMarkdownOptions } from '../_libs/html-to-markdown/index.ts'

async function fetchHtml(url: string) {
  let html: string | null = await fetchHtmlWithBrowser(url)
  if (!html) {
    const resp = await fetch(url, { headers: { Accept: 'text/html' } })
    html = await resp.text()
  }

  return html
}

export async function htmlToMarkdown(url: string, options?: HtmlToMarkdownOptions) {
  const resp = await fetch(url, { headers: { Accept: 'text/markdown' } })
  if (resp.ok) {
    return await resp.text()
  }

  const raw = await fetchHtml(url)
  return convertToMarkdown(raw, options)
}
