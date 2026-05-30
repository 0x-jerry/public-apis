import { load } from 'cheerio'
import TurndownService from 'turndown'
import puppeteer, { type Browser } from 'puppeteer-core'
import QuickLRU from 'quick-lru'

const htmlCache = new QuickLRU<string, string>({ maxSize: 100, maxAge: 60_000 })

let browserPromise: Promise<Browser> | null = null

function getBrowser(): Promise<Browser> {
  if (!browserPromise) {
    const wsEndpoint = process.env.BROWSER_WS || 'ws://127.0.0.1:9222'
    browserPromise = puppeteer.connect({ browserWSEndpoint: wsEndpoint })
  }
  return browserPromise
}

export async function fetchHtml(url: string) {
  const cached = htmlCache.get(url)
  if (cached !== undefined) {
    return cached
  }

  const enabled = process.env.BROWSER_WS_ENABLED === 'true' || process.env.BROWSER_WS_ENABLED === '1'
  let html: string
  if (enabled) {
    const browser = await getBrowser()
    const page = await browser.newPage()

    try {
      await page.goto(url, { waitUntil: 'networkidle0', timeout: 30_000 })
      html = await page.content()
    } finally {
      await page.close()
    }
  } else {
    const resp = await fetch(url, {
      headers: { Accept: 'text/html' },
    })
    html = await resp.text()
  }

  htmlCache.set(url, html)
  return html
}

export async function htmlToMarkdown(url: string, limit = 50000) {
  const raw = await fetchHtml(url)

  const $ = load(raw)

  const metadata = {
    url,
    title: $('title').text().trim(),
    description: $('meta[name="description"]').attr('content') || '',
    ogTitle: $('meta[property="og:title"]').attr('content') || '',
    ogDescription: $('meta[property="og:description"]').attr('content') || '',
    ogImage: $('meta[property="og:image"]').attr('content') || '',
    language: $('html').attr('lang') || '',
  }

  const yaml = [
    '---',
    `url: "${metadata.url}"`,
    `title: "${metadata.title}"`,
    metadata.description ? `description: "${metadata.description}"` : '',
    metadata.ogTitle ? `og:title: "${metadata.ogTitle}"` : '',
    metadata.ogDescription ? `og:description: "${metadata.ogDescription}"` : '',
    metadata.ogImage ? `og:image: "${metadata.ogImage}"` : '',
    metadata.language ? `language: "${metadata.language}"` : '',
    '---',
  ].filter(Boolean).join('\n')

  $('style, script, link, meta, head').remove()

  const turndown = new TurndownService()
  const markdown = turndown.turndown($.html())

  let result = yaml + '\n' + markdown

  if (result.length > limit) {
    result = result.slice(0, limit) + '\n\n... (truncated)'
  }

  return result
}
