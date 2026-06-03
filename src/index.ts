import { logger } from 'hono/logger'
import { app } from './_app.ts'
import './_index.tsx'
import { authRequired } from './_auth.ts'
import { app as qrRouter } from './qr/index.ts'
import { app as svgRouter } from './svg/index.ts'
import { app as htmlRouter } from './html/index.ts'
import { app as imgToAsciiRouter } from './img-to-ascii/index.ts'
import { app as linkRouter } from './link/index.ts'
import { app as mcpRouter } from './mcp/index.ts'
import { app as uploadRouter } from './upload/index.ts'
import { trimTrailingSlash } from 'hono/trailing-slash'

app.use(logger())
app.use(trimTrailingSlash())

app.use(async (ctx, next) => {
  await next()

  if (ctx.error) {
    console.error(ctx.error)

    ctx.status(400)
    ctx.res = ctx.json({
      message: String(ctx.error),
    })
  }
})

app.route('/qr', qrRouter)
app.route('/svg', svgRouter)

// app.use('/html/*', auth)
app.route('/html', htmlRouter)

app.route('/img-to-ascii', imgToAsciiRouter)
app.route('/link', linkRouter)

// app.use('/mcp/*', auth)
app.route('/mcp', mcpRouter)

app.use('/upload/file', authRequired)
app.route('/upload', uploadRouter)

export { app }
