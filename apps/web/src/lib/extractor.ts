import type { KnowledgeCategory } from "./types";

export interface ExtractedPageData {
  title: string;
  url: string;
  domain: string;
  author: string;
  category: KnowledgeCategory;
  summary: string;
  tags: string[];
}

export function detectDomain(url: string): string {
  try {
    const parsed = new URL(url);
    return parsed.hostname.replace(/^www\./, "");
  } catch {
    return "web";
  }
}

export function detectCategory(
  url: string,
  title: string,
  text: string,
): KnowledgeCategory {
  const lowerUrl = url.toLowerCase();
  const lowerContent = `${title} ${text}`.toLowerCase();

  // 1. Research (ArXiv, papers, breakthroughs, studies)
  if (
    lowerUrl.includes("arxiv.org") ||
    lowerUrl.includes("biorxiv.org") ||
    lowerUrl.includes("openreview.net") ||
    lowerUrl.includes("paperswithcode.com") ||
    lowerContent.includes("paper") ||
    lowerContent.includes("preprint") ||
    lowerContent.includes("benchmark") ||
    lowerContent.includes("ablation") ||
    lowerContent.includes("architecture") ||
    lowerContent.includes("abstract")
  ) {
    return "research";
  }

  // 2. Models (Hugging Face, LLMs, weights, vision, checkpoints)
  if (
    lowerUrl.includes("huggingface.co") ||
    lowerUrl.includes("civitai.com") ||
    lowerContent.includes("llm") ||
    lowerContent.includes("checkpoint") ||
    lowerContent.includes("weights") ||
    lowerContent.includes("parameters") ||
    lowerContent.includes("transformer") ||
    lowerContent.includes("diffusion model") ||
    lowerContent.includes("fine-tune") ||
    lowerContent.includes("quantization")
  ) {
    return "model";
  }

  // 3. Agents (Autonomous frameworks, agentic workflows, multi-agent)
  if (
    lowerContent.includes("agent") ||
    lowerContent.includes("agentic") ||
    lowerContent.includes("autonomous") ||
    lowerContent.includes("crewai") ||
    lowerContent.includes("autogen") ||
    lowerContent.includes("langgraph") ||
    lowerContent.includes("swarm")
  ) {
    return "agent";
  }

  // 4. Skills (antigravity skills, system prompts, function calling specs)
  if (
    lowerContent.includes("skill") ||
    lowerContent.includes("prompt engineering") ||
    lowerContent.includes("mcp server") ||
    lowerContent.includes("system prompt") ||
    lowerContent.includes("tool use") ||
    lowerContent.includes("function call")
  ) {
    return "skill";
  }

  // 5. Default to Tool (developer tools, libraries, repos, CLI)
  return "tool";
}

export function extractAuthorFromUrl(url: string): string {
  try {
    const parsed = new URL(url);
    const parts = parsed.pathname.split("/").filter(Boolean);

    // X / Twitter: https://x.com/username/...
    if (
      parsed.hostname.includes("x.com") ||
      parsed.hostname.includes("twitter.com")
    ) {
      if (parts.length > 0) return `@${parts[0]}`;
    }

    // Instagram: https://instagram.com/username/...
    if (parsed.hostname.includes("instagram.com")) {
      if (parts.length > 0 && parts[0] !== "p" && parts[0] !== "reel") {
        return `@${parts[0]}`;
      }
    }

    // GitHub: https://github.com/owner/repo
    if (parsed.hostname.includes("github.com")) {
      if (parts.length > 0) return parts[0];
    }

    // Substack: https://author.substack.com
    if (parsed.hostname.includes("substack.com")) {
      const sub = parsed.hostname.split(".")[0];
      if (sub && sub !== "www") return sub;
    }

    return parsed.hostname.replace(/^www\./, "");
  } catch {
    return "unknown";
  }
}

export function extractPageMetadata(
  url: string,
  documentTitle: string,
  pageTextSnippet: string = "",
): ExtractedPageData {
  const domain = detectDomain(url);
  const author = extractAuthorFromUrl(url);
  const category = detectCategory(url, documentTitle, pageTextSnippet);

  // Generate initial tags based on domain and category
  const tags: string[] = [category];
  if (domain.includes("github")) tags.push("opensource", "code");
  if (domain.includes("arxiv")) tags.push("paper", "academic");
  if (domain.includes("huggingface")) tags.push("ai-model");
  if (domain.includes("x.com") || domain.includes("twitter"))
    tags.push("social", "discussion");

  const summary = pageTextSnippet.trim()
    ? pageTextSnippet.slice(0, 320)
    : `Captured from ${domain} (${author}). Auto-indexed for developer memory.`;

  return {
    title: documentTitle || domain || "Web Note",
    url,
    domain,
    author,
    category,
    summary,
    tags,
  };
}
