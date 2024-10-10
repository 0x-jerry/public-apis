import { render } from '@deno/gfm'
import { html } from 'hono/html'
import { app } from './_app.ts'

app.get('/', async (ctx) => {
  const intro = await Deno.readTextFile('./readme.md')

  const app = <RenderMarkdownPage markdown={intro}></RenderMarkdownPage>

  return ctx.render(app)
})

interface SiteData {
  title: string
  style?: string
  children?: unknown
}

const Layout = (props: SiteData) =>
  html`<!DOCTYPE html>
    <html>
      <head>
        <title>${props.title}</title>
        <link
          rel="stylesheet"
          href="https://cdn.jsdelivr.net/npm/github-markdown-css/github-markdown.css"
        />
        <style>
          ${props.style}
        </style>
      </head>
      <body>
        ${props.children}
      </body>
    </html>`

function RenderMarkdownPage(props: { markdown: string }) {
  const md = render(props.markdown)

  return (
    <Layout title='Some API'>
      <div data-color-mode='auto' data-dark-theme='auto' class='markdown-body'>
        <div
          style={{ width: '768px', margin: 'auto' }}
          dangerouslySetInnerHTML={{ __html: md }}
        >
        </div>
      </div>
    </Layout>
  )
}
