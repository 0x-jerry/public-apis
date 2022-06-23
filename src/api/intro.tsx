import { h, renderSSR } from 'https://deno.land/x/nano_jsx@v0.0.32/mod.ts'

import * as gfm from 'https://deno.land/x/gfm@0.1.20/mod.ts'
import 'https://esm.sh/prismjs@1.28.0/components/prism-typescript?no-check'

function App() {
  const intro = Deno.readTextFileSync('./readme.md')
  const md = gfm.render(intro)

  return (
    <html>
      <head>
        <title>Deno Deploy API</title>
      </head>
      <style>
        {`html,body {padding: 0;margin: 0;}`}

        {gfm.CSS}
      </style>
      <body>
        <div
          data-color-mode="auto"
          data-dark-theme="auto"
          class="markdown-body"
        >
          <div
            style={{ width: '768px', margin: 'auto' }}
            dangerouslySetInnerHTML={{ __html: md }}
          ></div>
        </div>
      </body>
    </html>
  )
}

export const html = renderSSR(<App />)
