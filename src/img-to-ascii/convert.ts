import { imgToAscii } from '../_libs/img-to-ascii/index.ts'
import { app } from './_app.ts'

app.get('/convert', async (ctx) => {
  const url = ctx.req.query('url') || ''
  if (!url) throw new Error('Missing url parameter')

  const maxSize = parseInt(ctx.req.query('maxSize') || '') || 80

  const resp = await fetch(url)
  if (!resp.ok) {
    throw new Error(`Failed to fetch image: ${resp.status} ${resp.statusText}`)
  }
  const buf = new Uint8Array(await resp.arrayBuffer())
  const ascii = await imgToAscii(buf, { maxWidth: maxSize, maxHeight: maxSize })
  return ctx.text(ascii)
})
