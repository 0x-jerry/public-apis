import { getBrowser } from "../src/_libs/browser";
import { htmlToMarkdown } from "../src/_libs/html-to-markdown";

const browser = await getBrowser();
if (!browser) {
  throw new Error("Browser not connected!");
}

const page = await browser.newPage();

const url = "https://developer.mozilla.org/en-US/docs/Learn_web_development/Howto/Design_and_accessibility/HTML_features_for_accessibility";
await page.goto(url, { waitUntil: "networkidle0" });

const content = await page.content();

const md = htmlToMarkdown(content);

console.log(md);

await page.close();
await browser.close();