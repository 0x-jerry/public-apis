import { join } from 'node:path'

export function checkEnvInfos() {
  console.log('Auth: ', process.env.AUTH_ENABLED ?? false)
  console.log('Browser WS: ', process.env.BROWSER_WS_ENABLED ?? false)
  console.log('Storage: ', process.env.STORAGE_PATH || join(process.cwd(), 'files'))
}