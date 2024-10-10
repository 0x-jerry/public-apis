import { app } from './_app.ts'
import { get } from './_utils.ts'
import { load } from 'cheerio'
import { Feed } from 'feed'
import dayjs from 'dayjs'

const baseUrl = 'https://ollama.com'

app.get('/ollama.com/blog', async (ctx) => {
  const rssUrl = `${baseUrl}/blog`
  const content = await get(rssUrl)

  const $ = load(content)

  const feed = new Feed({
    id: baseUrl,
    title: 'Ollama BLog',
    language: 'en',
    copyright: 'Ollama',
    link: rssUrl,
    description: $("meta[name='description']").attr('content'),
  })

  const items = $('a.group.border-b.py-10')
    .toArray()

  for (const item of items) {
    feed.addItem({
      title: $(item).children('h2').first().text(),
      link: baseUrl + $(item).attr('href'),
      date: dayjs(
        $(item).children('h3').first().text().trim(),
        'MMMM D, YYYY',
      ).toDate(),
      description: $(item).children('p').first().text(),
    })
  }

  const rss = feed.rss2()

  return ctx.body(rss)
})
