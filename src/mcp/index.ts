import './html-to-markdown.ts'
import './ocr.ts'
import { app, server, transport } from './_app.ts'

await server.connect(transport)

export { app }
