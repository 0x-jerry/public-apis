import { assertEquals } from 'https://deno.land/std/testing/asserts.ts'

import { join as pJoin, fromFileUrl } from 'https://deno.land/std/path/mod.ts'

const root = 'http://localhost:8000'

const join = (...args: string[]) =>
  pJoin(fromFileUrl(import.meta.url), '..', ...args)

Deno.test('qr scan', async () => {
  const img = await fetch(root + '/qr/generate?c=1234')
  const imgBlob = await img.blob()

  const form = new FormData()
  form.set('file', imgBlob)

  const r = await fetch(root + '/qr/scan', {
    method: 'post',
    body: form,
  })

  const txt = await r.text()
  assertEquals(txt, '1234')
})

Deno.test('qr generate', async () => {
  const r = await fetch(root + '/qr/generate?c=1234')
  const file = await r.blob()

  assertEquals(file.type, 'image/gif')
})
