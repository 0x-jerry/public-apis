import puppeteer, { type Browser } from 'puppeteer-core'
import { sleep } from '@0x-jerry/utils'

let browser: Browser | null = null
let browserPromise: Promise<Browser> | null = null

function isBrowserEnabled() {
  const v = process.env.BROWSER_WS_ENABLED
  return v === 'true' || v === '1'
}

export function getBrowser(): Promise<Browser> | Browser | null {
  if (!isBrowserEnabled()) {
    return null
  }
  if (browser?.connected) {
    return browser
  }
  browser = null
  if (!browserPromise) {
    const wsEndpoint = process.env.BROWSER_WS || 'ws://127.0.0.1:9222'
    browserPromise = puppeteer
      .connect({ browserWSEndpoint: wsEndpoint })
      .then((b) => {
        browser = b
        browserPromise = null
        return b
      })
      .catch((e) => {
        browserPromise = null
        throw e
      })
  }
  return browserPromise
}

export async function fetchHtmlWithBrowser(url: string): Promise<string | null> {
  const browser = await getBrowser()
  if (!browser) return null

  const page = await browser.newPage()
  try {
    await Promise.race([
      page.goto(url, { waitUntil: 'networkidle0' }),
      sleep(10_000)
    ])
    
    return await page.content()
  } finally {
    await page.close()
  }
}
