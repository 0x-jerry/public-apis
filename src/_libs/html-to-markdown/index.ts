import { load } from 'cheerio'
import TurndownService from 'turndown'

export function htmlToMarkdown(html: string, url?: string, limit = 50000) {
  const $ = load(html)

  const metadata = {
    url: url ?? '',
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
