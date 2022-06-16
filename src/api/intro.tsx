import { h, renderSSR } from 'https://deno.land/x/nano_jsx@v0.0.32/mod.ts'

interface APIItemProps {
  path: string
  children?: APIItemProps[]
}

function APIItem({ path, children }: APIItemProps) {
  if (children?.length) {
    const item = children.map((api) => (
      <APIItem path={path + '/' + api.path} children={api.children}></APIItem>
    ))
    return <ul>{item}</ul>
  }

  return <li>{path}</li>
}

function App() {
  const api: APIItemProps[] = [
    {
      path: 'qr',
      children: [
        {
          path: 'generate',
        },
        {
          path: 'scan',
        },
      ],
    },
  ]

  return (
    <html>
      <head>
        <title>Deno Deploy API</title>
      </head>
      <body>
        <h1> Deno Deploy API</h1>
        <h3>
          Github:&nbsp;
          <a href="https://github.com/0x-jerry/dd-api" target="_blank">
            https://github.com/0x-jerry/dd-api
          </a>
        </h3>
        <div>
          <h3>API List:</h3>
          <div>
            {api.map((item) => (
              <APIItem {...item}></APIItem>
            ))}
          </div>
        </div>
      </body>
    </html>
  )
}

export const html = renderSSR(<App />)
