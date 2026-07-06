import { CheckCircle } from "lucide-react";

const points = [
  "HTML structure",
  "Tags & attributes",
  "Build your first webpage",
];

export default function WhatYouWillLearn() {
  return (
    <div className="bg-[#121313] rounded-xl p-4">
      <h3 className="mb-3 text-[#64ACFF] font-mono">What You’ll Learn</h3>

      <ul className="space-y-3 text-sm">
        {points.map((item) => (
          <li key={item} className="flex gap-2 items-center text-gray-300">
            <CheckCircle size={16} className="text-[#64ACFF]" />
            {item}
          </li>
        ))}
      </ul>
    </div>
  );
}
