import QRCode from 'qrcode'
import { app } from './_app.ts'

app.get('/generate', async (ctx) => {
  const content = ctx.req.query('c') || ''

  const code = await QRCode.toString(content, { type: 'svg' })

  ctx.header('content-type', 'image/svg+xml')

  return ctx.body(code)
})
