import { getLinkPreview } from 'link-preview-js'
import { app } from './_app.ts'

app.get('/metadata', async (ctx) => {
  const url = ctx.req.query('url') || ''

  const res = await getLinkPreview(url)

  ctx.json(res)
})
