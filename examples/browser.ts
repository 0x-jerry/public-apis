import { getBrowser } from "../src/_libs/browser"
import { htmlToMarkdown } from "../src/_libs/html-to-markdown"
import { sleep } from "@0x-jerry/utils"

const browser = await getBrowser()
if (!browser) {
  throw new Error("Browser not connected!")
}

const page = await browser.newPage()

const url = "https://news.ycombinator.com/news";

await Promise.race([
  // 
  page.goto(url, { waitUntil: "networkidle0" }),
  sleep(20_000),
])

const content = await page.content()

const md = htmlToMarkdown(content)

console.log(md)

await page.close()
await browser.close()