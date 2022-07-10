import { apiRoot } from './setup.ts'
import { assertEquals } from 'https://deno.land/std@0.96.0/testing/asserts.ts'

Deno.test('transform json to typescript', async () => {
  const test = {
    a: 1,
  }

  const res = await fetch(apiRoot + '/transform/json', {
    method: 'post',
    body: JSON.stringify({
      lang: 'typescript',
      json: JSON.stringify(test),
    }),
  })

  const ts = await res.text()

  assertEquals(ts, 'interface RootInterface {\n  a: number;\n}\n\n')
})

Deno.test('transform json to rust serde', async () => {
  const test = {
    a: 1,
  }

  const res = await fetch(apiRoot + '/transform/json', {
    method: 'post',
    body: JSON.stringify({
      lang: 'rust-serde',
      json: JSON.stringify(test),
    }),
  })

  const rust = await res.text()

  assertEquals(
    rust,
    '#[derive(Serialize, Deserialize)]\nstruct RootInterface {\n  a: i64,\n}\n\n'
  )
})

Deno.test('transform json use get', async () => {
  const test = {
    a: 1,
  }

  const url = new URL(apiRoot + '/transform/json?')
  url.searchParams.set('lang', 'typescript')
  url.searchParams.set('json', JSON.stringify(test))

  const res = await fetch(url)

  const ts = await res.text()

  assertEquals(ts, 'interface RootInterface {\n  a: number;\n}\n\n')
})
