import { htmlToMarkdown } from './convert.ts'
import { app } from './_app.ts'

app.get('/to-markdown', async (ctx) => {
  const url = ctx.req.query('url') || ''
  return ctx.text(await htmlToMarkdown(url))
})
