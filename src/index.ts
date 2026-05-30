import { app } from './_app.ts'
import './_index.tsx'

import { app as rssRouter } from './rss/index.ts'
import { app as qrRouter } from './qr/index.ts'
import { app as svgRouter } from './svg/index.ts'
import { app as htmlRouter } from './html/index.ts'
import { app as imgToAsciiRouter } from './img-to-ascii/index.ts'
import { app as mcpRouter } from './mcp/index.ts'

app.route('/rss', rssRouter)
app.route('/qr', qrRouter)
app.route('/svg', svgRouter)
app.route('/html', htmlRouter)
app.route('/img-to-ascii', imgToAsciiRouter)
app.route('/mcp', mcpRouter)

export { app }
