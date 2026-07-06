"use client";

import { useEffect, useState } from "react";

export default function NotesPanel({ loading }: { loading: boolean }) {
  return (
    <div className="h-full flex flex-col bg-[#1E1E1E] rounded-xl overflow-hidden">
      {/* NOTES NAVBAR */}
      <div className="flex items-center justify-between px-4 py-3 border-b border-white/10">
        {loading ? (
          <div className="h-6 w-40 bg-[#121313] rounded animate-pulse" />
        ) : (
          <span className="text-xl font-mono text-blue-400">Lecture Notes</span>
        )}

        <div className="flex items-center gap-3 text-sm">
          {loading ? (
            <div className="h-8 w-24 bg-[#121313] rounded animate-pulse" />
          ) : (
            <button className="px-3 py-1.5 bg-[#121313] rounded-md text-gray-300 hover:text-white transition">
              Download
            </button>
          )}
        </div>
      </div>

      {/* NOTES CONTENT */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4 text-sm text-gray-400 leading-relaxed">
        {loading ? (
          <>
            <div className="space-y-2 animate-pulse">
              <div className="h-4 w-2/3 bg-[#121313] rounded" />
              <div className="h-4 w-full bg-[#121313] rounded" />
              <div className="h-4 w-5/6 bg-[#121313] rounded" />
            </div>

            <div className="space-y-2 animate-pulse">
              <div className="h-4 w-1/2 bg-[#121313] rounded" />
              <div className="h-4 w-full bg-[#121313] rounded" />
              <div className="h-4 w-3/4 bg-[#121313] rounded" />
            </div>

            <div className="h-24 bg-[#121313] rounded animate-pulse" />
          </>
        ) : (
          <>
            <div>
              <h3 className="text-gray-200 font-semibold mb-1">
                Introduction to HTML
              </h3>
              <p>
                HTML stands for HyperText Markup Language. It is used to
                structure content on the web using elements such as headings,
                paragraphs, links, images, and more.
              </p>
            </div>

            <div>
              <h3 className="text-gray-200 font-semibold mb-1">
                Why HTML is Important
              </h3>
              <ul className="list-disc list-inside space-y-1">
                <li>Defines the structure of web pages</li>
                <li>Works with CSS and JavaScript</li>
                <li>Supported by all browsers</li>
              </ul>
            </div>

            <div>
              <h3 className="text-gray-200 font-semibold mb-1">
                Basic HTML Document
              </h3>
              <pre className="bg-[#121313] p-3 rounded-md text-xs text-gray-300 overflow-x-auto">
                {`<!DOCTYPE html>
<html>
  <head>
    <title>My First Page</title>
  </head>
  <body>
    <h1>Hello World</h1>
  </body>
</html>`}
              </pre>
            </div>

            <div className="text-xs text-gray-500 italic">
              Notes provided by course coordinator. For detailed explanations,
              refer to the full PDF.
            </div>
          </>
        )}
      </div>
    </div>
  );
}
