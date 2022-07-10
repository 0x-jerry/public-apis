# Deno Deploy API

A list of useful serverless API.

API endpoint: https://0x-jerry-dd-api.deno.dev/

## /qr/generate

Generate QRCode image.

Usage:

```ts
const apiRoot = 'https://0x-jerry-dd-api.deno.dev'
const r = await fetch(apiRoot + '/qr/generate?c=1234')
const file = await r.blob()

assertEquals(file.type, 'image/gif')
```

## /qr/scan

Scan QRCode.

Usage:

```ts
const apiRoot = 'https://0x-jerry-dd-api.deno.dev'
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
```

## /transform/json

Use [quicktype-core](https://github.com/quicktype/quicktype) to transform json to types.

Usage:

```ts
const apiRoot = 'https://0x-jerry-dd-api.deno.dev'

const test = {
  a: 1,
}

const url = new URL(apiRoot + '/transform/json?')
url.searchParams.set('lang', 'typescript')
url.searchParams.set('json', JSON.stringify(test))

const res = await fetch(url)

const ts = await res.text()

assertEquals(ts, 'interface RootInterface {\n  a: number;\n}\n\n')
```
