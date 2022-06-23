# Deno Deploy API

A list of useful serverless API.

API endpoint: https://0x-jerry-dd-api.deno.dev/

## /qr/generate

Generate QRCode image.

Usage:

```ts
const root = 'https://0x-jerry-dd-api.deno.dev'
const r = await fetch(root + '/qr/generate?c=1234')
const file = await r.blob()

assertEquals(file.type, 'image/gif')
```

## /qr/scan

Scan QRCode.

Usage:

```ts
const root = 'https://0x-jerry-dd-api.deno.dev'
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
```
