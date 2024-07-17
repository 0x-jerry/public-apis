import { Hono } from 'hono'
import { app as rssRouter } from './rss/index.ts'
import { app as qrRouter } from './qr/index.ts'

export const app = new Hono()

app.use(async (ctx, next) => {
  await next()

  if (ctx.error) {
    ctx.status(400)
    ctx.res = ctx.json({
      message: String(ctx.error),
    })
  }
})

app.route('/rss', rssRouter)
app.route('/qr', qrRouter)
