import { fetchHtmlWithBrowser } from '../_libs/browser/index.ts'
import {
  htmlToMarkdown as convertToMarkdown,
  type HtmlToMarkdownOptions,
} from '../_libs/html-to-markdown/index.ts'

export async function htmlToMarkdown(url: string, options?: HtmlToMarkdownOptions) {
  const resp = await fetch(url, { headers: { Accept: 'text/markdown' } })

  if (resp.ok && resp.headers.get('content-type')?.includes('text/markdown')) {
    return await resp.text()
  }

  let html = await fetchHtmlWithBrowser(url)

  if (!html) {
    const resp = await fetch(url, { headers: { Accept: 'text/html' } })
    html = await resp.text()
  }

  return convertToMarkdown(html, options)
}
