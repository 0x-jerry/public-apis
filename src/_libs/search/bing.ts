import { getBrowser } from '../browser'
import type { SearchResponse } from './types'

export async function searchBing(query: string): Promise<SearchResponse | null> {
  const browser = await getBrowser()
  if (!browser) return null

  const page = await browser.newPage()
  try {
    const url = `https://www.bing.com/search?q=${encodeURIComponent(query)}`
    await page.goto(url, { waitUntil: 'domcontentloaded' })
    await page.waitForSelector('[aria-label="Search Results"]')

    return await page.evaluate(`(() => {
      const results = [];
      const items = document.querySelectorAll('li.b_algo');
      for (const item of items) {
        const titleEl = item.querySelector('h2 a');
        const title = titleEl?.textContent?.trim() ?? '';
        const url = titleEl?.href ?? '';
        const snippetEl = item.querySelector('.b_caption p, .b_lineclamp2');
        const snippet = snippetEl?.textContent?.trim() ?? '';
        if (title && url) results.push({ title, url, snippet });
      }

      let prev, next;
      const nav = document.querySelector('nav[role="navigation"]');
      if (nav) {
        const prevLink = nav.querySelector('a.sb_pagP');
        const nextLink = nav.querySelector('a.sb_pagN');
        if (prevLink) prev = prevLink.href;
        if (nextLink) next = nextLink.href;
      }

      return { results, prev, next };
    })()`) as SearchResponse
  } finally {
    await page.close()
  }
}
