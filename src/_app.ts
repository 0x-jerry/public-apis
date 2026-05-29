import { Hono } from 'hono'

export const app = new Hono()

app.use(async (ctx, next) => {
  const start = Date.now()
  await next()
  const ms = Date.now() - start
  console.log(`${ctx.req.method} ${ctx.req.path} ${ctx.res.status} ${ms}ms`)
})

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
