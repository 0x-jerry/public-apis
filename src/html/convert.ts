import { JSDOM } from 'jsdom'
import { fetchHtmlWithBrowser } from '../_libs/browser/index.ts'
import {
  htmlToMarkdown as convertToMarkdown,
  type HtmlToMarkdownOptions,
} from '../_libs/html-to-markdown/index.ts'
import { isProbablyReaderable } from '@mozilla/readability'

const ACCEPT_MARKDOWN = 'text/markdown'
const ACCEPT_HTML = 'text/html'

export async function htmlToMarkdown(url: string, options?: HtmlToMarkdownOptions) {
  const resp = await fetch(url, { headers: { Accept: ACCEPT_MARKDOWN } })

  let fetchedHtml: string | null = null

  if (resp.ok) {
    if (resp.headers.get('content-type')?.includes(ACCEPT_MARKDOWN)) {
      return await resp.text()
    }

    if (resp.headers.get('content-type')?.includes(ACCEPT_HTML)) {
      const html = await resp.text()

      const jsdom = new JSDOM(html)
      const doc = jsdom.window.document

      if (isProbablyReaderable(doc)) {
        return convertToMarkdown(doc, options)
      }

      fetchedHtml = html
    }
  }

  let html = (await fetchHtmlWithBrowser(url)) || fetchedHtml

  if (!html) {
    const resp = await fetch(url, { headers: { Accept: ACCEPT_HTML } })
    html = await resp.text()
  }

  return convertToMarkdown(html, options)
}
