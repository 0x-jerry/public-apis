import { h, renderSSR } from 'https://deno.land/x/nano_jsx@v0.0.32/mod.ts'

import * as gfm from 'https://deno.land/x/gfm@0.1.20/mod.ts'

function App() {
  const intro = Deno.readTextFileSync('./readme.md')
  const md = gfm.render(intro)

  return (
    <html>
      <head>
        <title>Deno Deploy API</title>
      </head>
      <style>{gfm.CSS}</style>
      <body>
        <div
          style={{ width: '768px', margin: 'auto' }}
          dangerouslySetInnerHTML={{ __html: md }}
        ></div>
      </body>
    </html>
  )
}

export const html = renderSSR(<App />)
