import { apiRoot } from './setup.ts'
import { assertEquals } from 'https://deno.land/std@0.147.0/testing/asserts.ts'

Deno.test('get metadata of a website', async () => {
  const res = await fetch(apiRoot + '/preview/metadata', {
    method: 'post',
    body: JSON.stringify({
      url: 'https://www.bilibili.com/bangumi/play/ep471910',
    }),
  })

  const metadata = await res.json()

  assertEquals(!!metadata.title, true)
})
