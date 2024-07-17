import { Hono } from 'hono'
import { app as rssRouter } from './rss/index.ts'

export const app = new Hono()

app.route('/rss', rssRouter)
