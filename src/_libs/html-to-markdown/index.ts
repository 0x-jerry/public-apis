import TurndownService from "turndown";
import { isProbablyReaderable, Readability } from "@mozilla/readability";
import { JSDOM } from "jsdom";

const turndown = new TurndownService({
  headingStyle: "atx",
  codeBlockStyle: "fenced",
});

export interface HtmlToMarkdownOptions {
  /** Source URL, included in YAML front matter */
  url?: string
  /** Maximum output length in characters (default: 50000) */
  limit?: number
  /** `"readable"` (default) extracts article content, `"full"` converts the entire body */
  mode?: "full" | "readable"
}

/**
 * Convert HTML to Markdown with YAML front matter metadata.
 *
 * By default (`readable` mode), uses Mozilla's Readability to extract the main
 * article content, stripping navigation, sidebars, and other clutter.
 * In `full` mode, the entire HTML body is converted as-is.
 *
 * @param html - Raw HTML string to convert
 * @param options
 * @param options.url - Source URL, included in YAML front matter
 * @param options.limit - Maximum output length in characters (default: 50000)
 * @param options.mode - `"readable"` (default) extracts article content,
 *   `"full"` converts the entire body without filtering
 * @returns Markdown string prefixed with YAML front matter
 */
export function htmlToMarkdown(html: string, options?: HtmlToMarkdownOptions) {
  const { url, limit = 50000, mode = "readable" } = options ?? {};

  const dom = new JSDOM(html);
  const doc = dom.window.document;

  removeInvisibleTags(doc);

  const yaml = extractMetadata(url, doc);

  let article = doc.body.innerHTML;
  
  if (mode === "readable" && isProbablyReaderable(doc)) {
    article = new Readability(doc).parse()?.content || article;
  }

  const markdown = turndown.turndown(article);

  let result = yaml + "\n" + markdown;

  if (result.length > limit) {
    result = result.slice(0, limit) + "\n\n... (truncated)";
  }

  return result;
}

function removeInvisibleTags(doc: Document) {
  const tags = ["script", "style", "noscript", "template"];
  for (const tag of tags) {
    doc.querySelectorAll(tag).forEach((el) => el.remove());
  }
}

function extractMetadata(url: string | undefined, doc: Document) {
  const getMetaContent = (attr: string, value: string) =>
    doc.querySelector(`meta[${attr}="${value}"]`)?.getAttribute("content") || "";

  const metadata = {
    url: url ?? "",
    title: doc.querySelector("title")?.textContent?.trim() || "",
    description: getMetaContent("name", "description"),
    ogTitle: getMetaContent("property", "og:title"),
    ogDescription: getMetaContent("property", "og:description"),
    ogImage: getMetaContent("property", "og:image"),
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