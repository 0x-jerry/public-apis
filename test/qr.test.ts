import { apiRoot } from './setup.ts'
import { assertEquals } from 'https://deno.land/std@0.96.0/testing/asserts.ts'

Deno.test('qr scan', async () => {
  const img = await fetch(apiRoot + '/qr/generate?c=1234')
  const imgBlob = await img.blob()

  const form = new FormData()
  form.set('file', imgBlob)

  const r = await fetch(apiRoot + '/qr/scan', {
    method: 'post',
    body: form,
  })

  const txt = await r.text()
  assertEquals(txt, '1234')
})

Deno.test('qr generate', async () => {
  const r = await fetch(apiRoot + '/qr/generate?c=1234')
  const file = await r.blob()

  assertEquals(file.type, 'image/gif')
})
