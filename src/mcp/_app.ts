import { createMcpHonoApp } from '@modelcontextprotocol/hono'
import { McpServer, WebStandardStreamableHTTPServerTransport } from '@modelcontextprotocol/server'

const app = createMcpHonoApp({
  host: '0.0.0.0'
})

const server = new McpServer({
  name: 'public-apis',
  version: '1.0.0',
})

const transport = new WebStandardStreamableHTTPServerTransport()

app.all('/', async (c) => {
  return transport.handleRequest(c.req.raw)
})

export { app, server, transport }
