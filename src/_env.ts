export function checkEnvInfos() {
  console.log('Auth: ', process.env.AUTH_ENABLED ?? false)
  console.log('Browser WS: ', process.env.BROWSER_WS_ENABLED ?? false)
}