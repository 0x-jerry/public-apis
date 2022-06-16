export function sleep(ts = 0) {
  return new Promise<void>((resolve) => setTimeout(resolve, ts))
}
