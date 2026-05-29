import { initWasm, Resvg } from '@resvg/resvg-wasm'
import { app } from './_app.ts'

app.post('/png', async (ctx) => {
  const data = await ctx.req.json()

  const { svg, url } = data

  const svgData = url ? await getSvg(url) : svg

  const buf = await svgToPng(svgData)

  return new Response(buf, {
    headers: { 'content-type': 'image/png' },
  })
})

async function getSvg(url: string) {
  const svg = await (await fetch(url)).text()

  return svg
}

async function svgToPng(svg: string) {
  // @ts-ignore
  if (!svgToPng.inited) {
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
