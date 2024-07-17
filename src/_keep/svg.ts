import { Router } from 'oak'
import { initWasm, Resvg } from 'https://cdn.skypack.dev/@resvg/resvg-wasm'

export const router = new Router()

router.post('/png', async (ctx) => {
  const body = ctx.request.body()

  const data =
    body.type === 'text' ? JSON.parse(await body.value) : await body.value

  const { svg, url } = data

  const svgData = url ? await getSvg(url) : svg

  const buf = await svgToPng(svgData)

  ctx.response.headers.set('Content-Type', 'image/png')
  ctx.response.body = buf
})

router.get('/png', async (ctx) => {
  const params = ctx.request.url.searchParams
  const url = params.get('url')!

  const svgData = await getSvg(url)

  const buf = await svgToPng(svgData)

  ctx.response.headers.set('Content-Type', 'image/png')
  ctx.response.body = buf
})

async function getSvg(url: string) {
  const svg = await (await fetch(url)).text()

  return svg
}

async function svgToPng(svg: string) {
  // deno-lint-ignore ban-ts-comment
  // @ts-ignore
  if (!svgToPng.inited) {
    // deno-lint-ignore ban-ts-comment
    // @ts-ignore
    svgToPng.inited = true
    await initWasm(fetch('https://unpkg.com/@resvg/resvg-wasm/index_bg.wasm'))
  }

  const resvg = new Resvg(svg, {
    fitTo: {
      mode: 'width', // If you need to change the size
      value: 800,
    },
  })

  const pngData = resvg.render()

  return pngData.asPng()
}
