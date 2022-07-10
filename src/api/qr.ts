import { FormDataReader, Router } from 'oak'
import { qrcode } from 'https://deno.land/x/qrcode@v2.0.0/mod.ts'
import { toUint8Array } from 'https://deno.land/x/base64@v0.2.1/mod.ts'
import {
  decode,
  GIF,
  Image,
} from 'https://deno.land/x/imagescript@v1.2.13/mod.ts'
import jsqr from 'https://cdn.skypack.dev/jsqr@v1.4.0?dts'

export const router = new Router()

router.post('/scan', async (ctx) => {
  const body = ctx.request.body().value as FormDataReader
  const form = await body.read()
  const file = form.files?.[0]
  if (!file?.filename) {
    throw Error('Not found file')
  }

  const imgData = await Deno.readFile(file.filename)
  const imgInfo = await decode(imgData)

  const imgBuffer = (imgInfo as Image).bitmap || (imgInfo as GIF).at(0)?.bitmap

  const code = jsqr(imgBuffer, imgInfo.width, imgInfo.height)

  ctx.response.body = code?.data
})

router.get('/generate', async (ctx) => {
  const content = ctx.request.url.searchParams.get('c') || ''

  const code = await qrcode(content)

  ctx.response.headers.set('Content-Type', 'image/gif')
  const b64 = String(code).replace('data:image/gif;base64,', '')
  const buf: Uint8Array = toUint8Array(b64)

  ctx.response.body = buf
})
