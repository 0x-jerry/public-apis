import { h, html } from 'https://deno.land/x/htm@0.0.7/mod.tsx'

import * as gfm from 'https://deno.land/x/gfm@0.1.20/mod.ts'
import 'https://esm.sh/prismjs@1.28.0/components/prism-typescript?no-check'

export const htmlString = async () => {
  const intro = await Deno.readTextFile('./readme.md')

  const res = await html({
    title: 'Deno Deploy API',
    styles: [`html,body {padding: 0;margin: 0;}`, gfm.CSS],
    headers: [],
    body: <App markdown={intro}></App>,
  })

  return res.text()
}

function App(props: { markdown: string }) {
  const md = gfm.render(props.markdown)

  return (
    <div data-color-mode="auto" data-dark-theme="auto" class="markdown-body">
      <div
        style={{ width: '768px', margin: 'auto' }}
        dangerouslySetInnerHTML={{ __html: md }}
      ></div>
    </div>
  )
}
