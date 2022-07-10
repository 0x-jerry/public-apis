import { apiRoot } from './setup.ts'
import { assertSnapshot } from 'https://deno.land/std@0.147.0/testing/snapshot.ts'

Deno.test('transform json to typescript', async (t) => {
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

  assertSnapshot(t, ts)
})

Deno.test('transform json to rust serde', async (t) => {
  const test = {
    a: 1,
  }

  const res = await fetch(apiRoot + '/transform/json', {
    method: 'post',
    body: JSON.stringify({
      lang: 'rust',
      json: JSON.stringify(test),
    }),
  })

  const rust = await res.text()

  assertSnapshot(t, rust)
})

Deno.test('transform json use get', async (t) => {
  const test = {
    a: 1,
  }

  const url = new URL(apiRoot + '/transform/json?')
  url.searchParams.set('lang', 'typescript')
  url.searchParams.set('json', JSON.stringify(test))

  const res = await fetch(url)

  const ts = await res.text()

  assertSnapshot(t, ts)
})
