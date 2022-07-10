import { Router } from 'oak'
import transform from 'https://cdn.skypack.dev/transform-json-types@v0.7.0'

export const router = new Router()

router.post('/json', async (ctx) => {
  const data = await ctx.request.body().value
  const { json, lang } = JSON.parse(data)

  const result = transform(json, { lang })

  ctx.response.body = result
})

router.get('/json', (ctx) => {
  const params = ctx.request.url.searchParams
  const lang = params.get('lang')
  const json = params.get('json')

  const result = transform(json, { lang })

  ctx.response.body = result
})
