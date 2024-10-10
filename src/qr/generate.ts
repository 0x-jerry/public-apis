import { qrcode } from '@libs/qrcode'
import { app } from './_app.ts'

app.get('/generate', (ctx) => {
  const content = ctx.req.query('c') || ''

  const code = qrcode(content, { output: 'svg' })

  ctx.header('content-type', 'image/svg+xml')

  return ctx.body(code)
})
