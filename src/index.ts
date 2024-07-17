import { app } from './_app.ts'
import './_index.tsx'

import { app as rssRouter } from './rss/index.ts'
import { app as qrRouter } from './qr/index.ts'
import { app as svgRouter } from './svg/index.ts'

app.route('/rss', rssRouter)
app.route('/qr', qrRouter)
app.route('/svg', svgRouter)

export { app }
