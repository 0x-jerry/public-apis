import { htmlToMarkdown } from './convert.ts'
import { app } from './_app.ts'

app.get('/to-markdown', async (ctx) => {
  const url = ctx.req.query('url') || ''
  const result = await htmlToMarkdown(url)
  return ctx.text(result)
})
