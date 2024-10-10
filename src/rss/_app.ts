import { Hono } from 'hono'

export const app = new Hono()

app.use('/*', async (ctx, next) => {
    ctx.header('content-type', 'text/xml')
    ctx.header('charset', 'utf-8')

    await next()
})
