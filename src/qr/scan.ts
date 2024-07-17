import { decode, GIF, Image } from 'imagescript'
import jsqr from 'jsqr'
import { app } from './_app.ts'

app.post('/scan', async (ctx) => {
  const form = await ctx.req.formData()
  const file = form.get('file') as File
  if (!file?.name) {
    throw Error('Not found file')
  }

  const imgInfo = await decode(await file.arrayBuffer())

  const imgBuffer = (imgInfo as Image).bitmap || (imgInfo as GIF).at(0)?.bitmap

  const code = jsqr.default(imgBuffer, imgInfo.width, imgInfo.height)

  return ctx.json({
    version: code?.version,
    data: code?.data,
    location: code?.location,
  })
})
