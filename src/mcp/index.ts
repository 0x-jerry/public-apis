import './html-to-markdown.ts'
import './img-to-ascii.ts'
import './ocr.ts'
import { app, server, transport } from './_app.ts'

await server.connect(transport)

export { app }
