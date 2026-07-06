"use client";

import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import remarkMath from "remark-math";
import rehypeRaw from "rehype-raw";
import rehypeSanitize from "rehype-sanitize";
import rehypeKatex from "rehype-katex";
import rehypeStringify from "rehype-stringify";

type MarkdownComponentProps = {
  content: string;
};

export default function MarkdownComponent({ content }: MarkdownComponentProps) {
  return (
    <div className="prose prose-invert max-w-none">
      <ReactMarkdown
        remarkPlugins={[remarkGfm, remarkMath]}
        rehypePlugins={[
          rehypeRaw,
          rehypeSanitize,
          rehypeKatex,
          rehypeStringify,
        ]}
      >
        {content}
      </ReactMarkdown>
    </div>
  );
}
