'use client';

import { useState, useEffect } from 'react';
import Editor from '@monaco-editor/react';

import { Button } from "@nextui-org/button";
import { Select, SelectItem } from "@nextui-org/select";
import { Input } from "@nextui-org/input";
import { useTheme } from "next-themes";

import { siteConfig } from "@/config/site";
import { title, subtitle } from "@/components/primitives";
import { Skeleton } from "@nextui-org/skeleton";

export default function Home() {
  const [language, setLanguage] = useState('plaintext');
  const [titleText, setTitleText] = useState('');
  const [code, setCode] = useState('');
  const { theme } = useTheme();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const timer = setTimeout(() => {
      setLoading(false);
    }, 2000);

    return () => clearTimeout(timer);
  }, []);

  const programmingLanguages = [
    { value: 'plaintext', label: 'Plain Text' },
    { value: 'javascript', label: 'JavaScript' },
    { value: 'python', label: 'Python' },
    { value: 'java', label: 'Java' },
    { value: 'cpp', label: 'C++' },
    { value: 'csharp', label: 'C#' },
    { value: 'php', label: 'PHP' },
    { value: 'ruby', label: 'Ruby' },
    { value: 'go', label: 'Go' },
    { value: 'rust', label: 'Rust' },
  ];

  const handleSubmit = async () => {
    const jsonData = {
      title: titleText,
      language: language,
      content: code,
    };

    console.log("Generated JSON:", JSON.stringify(jsonData, null, 2));

  };
  const handleEditorChange = (value: any) => {
    console.log(value);
    setCode(value);
  };

  return (
    <section className="flex flex-col items-center justify-center gap-4">
      <div className="inline-block max-w-xl text-center justify-center">
        <span className={title()}>Write your&nbsp;</span>
        <span className={title({ color: "blue" })}>code&nbsp;</span>
      </div>
      <main className="container mx-auto px-4">
        <div className="light:bg-default-100 dark:bg-[#141414] rounded-lg shadow-lg p-6">
          <div className="flex gap-4 mb-4">
            <Input
              type="text"
              label="Title"
              placeholder="Enter paste title..."
              value={titleText}
              onChange={(e) => setTitleText(e.target.value)}
              className="flex-1"
            />
            <Select
              label="Language"
              placeholder="Select language"
              defaultSelectedKeys={[language]}
              onChange={(e: { target: { value: React.SetStateAction<string>; }; }) => setLanguage(e.target.value)}
              className="w-48"
            >
              {programmingLanguages.map((lang) => (
                <SelectItem key={lang.value} value={lang.value}>
                  {lang.label}
                </SelectItem>
              ))}
            </Select>
          </div>

          <div className="relative w-full h-96">
            {loading && (
              <Skeleton className="absolute inset-0 w-full h-full rounded-lg border border-gray-500" />
            )}

            <Editor
              className='w-full dark:bg-[#1e1e1e] light:bg-default-100 h-96 p-4 font-mono rounded-lg border border-gray-500 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none resize-none'
              height="100%"
              language={language}
              value={code}
              onChange={(e) => handleEditorChange(e)}
              theme={theme === "dark" ? "vs-dark" : "vs-light"}
              options={{
                minimap: { enabled: false },
                lineNumbers: "on",
                roundedSelection: true,
                scrollBeyondLastLine: false,
                readOnly: false,
                fontSize: 14,
                contextmenu: false,
                smoothScrolling: true,
                cursorBlinking: "smooth",
                cursorSmoothCaretAnimation: "on",
              }}
            />
          </div>

          <div className="mt-4 flex justify-end">
            <Button color="primary" className="px-6" onClick={handleSubmit}>
              Save Paste
            </Button>
          </div>
        </div>
      </main>
    </section>
  );
}