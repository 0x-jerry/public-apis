import { Router } from 'oak'
import { qrcode } from 'https://deno.land/x/qrcode@v2.0.0/mod.ts'

export const router = new Router()

interface QrScanParams {
  file: Blob
}

router.post('/qr/scan', async (ctx) => {
  // const body: FormDataReader = ctx.request.body().value
  // body
})

router.get('/qr/generate', async (ctx) => {
  const content = ctx.request.url.searchParams.get('c') || ''

  const code = await qrcode(content)

  ctx.response.body = code.createDataURL()
})
