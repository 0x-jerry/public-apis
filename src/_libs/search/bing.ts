import { getBrowser } from '../browser'

export interface BingSearchResult {
  title: string
  url: string
  snippet: string
}

export interface BingPaginationLink {
  text: string
  url: string
}

export interface BingSearchResponse {
  results: BingSearchResult[]
  pagination: BingPaginationLink[]
}

export async function searchBing(query: string): Promise<BingSearchResponse | null> {
  const browser = await getBrowser()
  if (!browser) return null

  const page = await browser.newPage()
  try {
    const url = `https://www.bing.com/search?q=${encodeURIComponent(query)}`
    await page.goto(url, { waitUntil: 'domcontentloaded' })
    await page.waitForSelector('[aria-label="Search Results"]')
    await page.waitForSelector('[role="navigation"]')

    return await page.evaluate(() => {
      const results: { title: string; url: string; snippet: string }[] = []
      const items = document.querySelectorAll('li.b_algo')

      for (const item of items) {
        const titleEl = item.querySelector('h2 a')
        const title = titleEl?.textContent?.trim() ?? ''
        const url = (titleEl as HTMLAnchorElement)?.href ?? ''
        const snippetEl = item.querySelector('.b_caption p, .b_lineclamp2')
        const snippet = snippetEl?.textContent?.trim() ?? ''

        if (title && url) {
          results.push({ title, url, snippet })
        }
      }

      const pagination: { text: string; url: string }[] = []
      const paginationNav = document.querySelector('nav[role="navigation"]')
      if (paginationNav) {
        const links = paginationNav.querySelectorAll('a')
        for (const link of links) {
          const href = (link as HTMLAnchorElement).href
          const text = link.textContent?.trim() ?? ''
          if (href && href.includes('search?')) {
            pagination.push({ text, url: href })
          }
        }
      }

      return { results, pagination }
    })
  } finally {
    await page.close()
  }
}
