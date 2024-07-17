import { Router } from 'oak'
import { getLinkPreview } from 'https://cdn.skypack.dev/link-preview-js?dts'

export const router = new Router()

router.post('/metadata', async (ctx) => {
  const body = ctx.request.body()

  const data =
    body.type === 'text' ? JSON.parse(await body.value) : await body.value

  const { url } = data

  console.log(url, data)

  const res = await getLinkPreview(url)

  ctx.response.body = res
})

router.get('/metadata', async (ctx) => {
  const params = ctx.request.url.searchParams
  const url = params.get('url')!

  const res = await getLinkPreview(url)

  ctx.response.body = res
})
