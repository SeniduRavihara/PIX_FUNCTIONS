"use client";

import dynamic from "next/dynamic";
import { useState } from "react";

// Dynamically import Monaco Editor to avoid SSR issues
const MonacoEditor = dynamic(() => import("@monaco-editor/react"), {
  ssr: false,
  loading: () => (
    <div className="w-full h-full flex items-center justify-center bg-gray-900 text-white">
      Loading editor...
    </div>
  ),
});

interface CodeEditorProps {
  value: string;
  onChange?: (value: string | undefined) => void;
  language?: string;
  readOnly?: boolean;
  height?: string;
}

export default function CodeEditor({
  value,
  onChange,
  language = "javascript",
  readOnly = false,
  height = "500px",
}: CodeEditorProps) {
  const [mounted, setMounted] = useState(false);

  // Fallback to textarea if Monaco fails to load
  const [useMonaco, setUseMonaco] = useState(true);

  if (!useMonaco) {
    return (
      <textarea
        value={value}
        onChange={(e) => onChange?.(e.target.value)}
        readOnly={readOnly}
        className="w-full font-mono text-sm p-4 bg-gray-900 text-white border border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
        style={{ height }}
      />
    );
  }

  return (
    <div className="border border-gray-700 rounded-md overflow-hidden">
      <MonacoEditor
        height={height}
        language={language}
        value={value}
        onChange={onChange}
        theme="vs-dark"
        options={{
          readOnly,
          minimap: { enabled: false },
          fontSize: 14,
          lineNumbers: "on",
          scrollBeyondLastLine: false,
          automaticLayout: true,
          tabSize: 2,
        }}
        onMount={() => setMounted(true)}
        loading={
          <div className="w-full h-full flex items-center justify-center bg-gray-900 text-white">
            Loading editor...
          </div>
        }
      />
    </div>
  );
}
