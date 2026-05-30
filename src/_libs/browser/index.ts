import puppeteer, { type Browser } from 'puppeteer-core'

let browserPromise: Promise<Browser> | null = null

function isBrowserEnabled() {
  const v = process.env.BROWSER_WS_ENABLED
  return v === 'true' || v === '1'
}

export function getBrowser(): Promise<Browser> | null {
  if (!isBrowserEnabled()) {
    return null
  }
  if (!browserPromise) {
    const wsEndpoint = process.env.BROWSER_WS || 'ws://127.0.0.1:9222'
    browserPromise = puppeteer.connect({ browserWSEndpoint: wsEndpoint })
  }
  return browserPromise
}

export async function fetchHtmlWithBrowser(url: string): Promise<string | null> {
  const browser = await getBrowser()
  if (!browser) return null

  const page = await browser.newPage()
  try {
    await page.goto(url, { waitUntil: 'networkidle0', timeout: 30_000 })
    return await page.content()
  } finally {
    await page.close()
  }
}
