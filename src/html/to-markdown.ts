import { htmlToMarkdown } from './convert.ts'
import { app } from './_app.ts'

app.get('/to-markdown', async (ctx) => {
  const url = ctx.req.query('url') || ''
  const limit = parseInt(ctx.req.query('limit') || '') || 50000
  return ctx.text(await htmlToMarkdown(url, { limit }))
})
