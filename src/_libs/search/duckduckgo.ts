import { load } from 'cheerio'
import type { SearchResponse } from './types'

export async function searchDuckDuckGo(query: string): Promise<SearchResponse> {
  const url = `https://html.duckduckgo.com/html/?q=${encodeURIComponent(query)}`
  const res = await fetch(url)
  const html = await res.text()
  const $ = load(html)

  const results = $('.result')
    .map((_, item) => {
      const title = $(item).find('.result__a').text().trim()
      const url = $(item).find('.result__url').text().trim()
      const snippet = $(item).find('.result__snippet').text().trim()
      if (title && url) return { title, url, snippet }
    })
    .get()

  let prev: string | undefined
  let next: string | undefined

  const navLinks = $('.nav-link a')
  navLinks.each((_, link) => {
    const text = $(link).text().trim()
    const href = $(link).attr('href')
    if (href) {
      if (/prev/i.test(text)) prev = href
      if (/next/i.test(text)) next = href
    }
  })

  const forms = $('.nav-link form')
  forms.each((_, form) => {
    const value = $(form).find('input[type="submit"]').attr('value') ?? ''
    const action = $(form).attr('action') ?? '/html/'
    const params = new URLSearchParams()
    $(form)
      .find('input[type="hidden"]')
      .each((_, input) => {
        params.set($(input).attr('name') ?? '', $(input).attr('value') ?? '')
      })
    const formUrl = `https://html.duckduckgo.com${action}?${params.toString()}`
    if (/prev/i.test(value)) prev = formUrl
    if (/next/i.test(value)) next = formUrl
  })

  return { results, prev, next }
}
