import { Application } from 'oak'
import { router } from './api/index.ts'

const app = new Application()

app.use(async (ctx, next) => {
  const url = ctx.request.url
  console.log(`[${ctx.request.method}] ${url.pathname}`)
  await next()
})

app.use((ctx, next) => {
  ctx.response.headers.set('Access-Control-Allow-Origin', '*')
  ctx.response.headers.set('Access-Control-Allow-Headers', '*')

  return next()
})

app.use(router.routes(), router.allowedMethods())

console.log('http://localhost:8000')
await app.listen({ port: 8000 })
