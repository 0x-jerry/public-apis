import { Hono } from 'hono'

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
