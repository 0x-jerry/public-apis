import TurndownService from "turndown";
import { Readability } from "@mozilla/readability";
import { JSDOM } from "jsdom";

const turndown = new TurndownService({
  headingStyle: "atx",
  codeBlockStyle: "fenced",
});

export function htmlToMarkdown(html: string, url?: string, limit = 50000) {
  const dom = new JSDOM(html);
  const doc = dom.window.document;

  const yaml = extractMetadata(url, doc);

  const article = new Readability(doc).parse();

  const markdown = turndown.turndown(article?.content || doc.body?.innerHTML || "");

  let result = yaml + "\n" + markdown;

  if (result.length > limit) {
    result = result.slice(0, limit) + "\n\n... (truncated)";
  }

  return result;
}

function extractMetadata(url: string | undefined, doc: Document) {
  const getMetaContent = (attr: string, value: string) =>
    doc.querySelector(`meta[${attr}="${value}"]`)?.getAttribute("content") || "";

  const metadata = {
    url: url ?? "",
    title: doc.querySelector("title")?.textContent?.trim() || "",
    description: getMetaContent('name', 'description'),
    ogTitle: getMetaContent('property', 'og:title'),
    ogDescription: getMetaContent('property', 'og:description'),
    ogImage: getMetaContent('property', 'og:image'),
    language: doc.documentElement?.getAttribute("lang") || "",
  };

  const yaml = [
    "---",
    `url: "${metadata.url}"`,
    `title: "${metadata.title}"`,
    metadata.description ? `description: "${metadata.description}"` : "",
    metadata.ogTitle ? `og:title: "${metadata.ogTitle}"` : "",
    metadata.ogDescription ? `og:description: "${metadata.ogDescription}"` : "",
    metadata.ogImage ? `og:image: "${metadata.ogImage}"` : "",
    metadata.language ? `language: "${metadata.language}"` : "",
    "---",
  ]
    .filter(Boolean)
    .join("\n");

  return yaml;
}