import { Router } from 'oak'
import { qrcode } from 'https://deno.land/x/qrcode@v2.0.0/mod.ts'
import { toUint8Array } from 'https://deno.land/x/base64@v0.2.1/mod.ts'
import { html } from './intro.tsx'

export const router = new Router()

interface QrScanParams {
  file: Blob
}

router.get('/', (ctx) => {
  ctx.response.body = html
})

router.post('/qr/scan', async (ctx) => {
  // const body: FormDataReader = ctx.request.body().value
  // body
})

router.get('/qr/generate', async (ctx) => {
  const content = ctx.request.url.searchParams.get('c') || ''

  const code = await qrcode(content)

  ctx.response.headers.set('Content-Type', 'image/gif')
  const b64 = String(code).replace('data:image/gif;base64,', '')
  const buf: Uint8Array = toUint8Array(b64)

  ctx.response.body = buf
})
