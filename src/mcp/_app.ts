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
  const headers = c.req.header();
  console.log('--- mcp request headers ---')
  for (const [key, value] of Object.entries(headers)) {
    console.log(`  ${key}: ${value}`)
  }
  console.log('--- end mcp request headers ---')
  
  return transport.handleRequest(c.req.raw)
})

export { app, server, transport }
