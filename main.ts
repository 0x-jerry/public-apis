import { serve } from '@hono/node-server'
import { checkEnvInfos } from './src/_env.ts'
import { app } from './src/index.ts'

checkEnvInfos()

serve({ fetch: app.fetch, port: 3000 }, (info) => {
  console.log(`http://localhost:${info.port}`)
})
