import WhatYouWillLearn from "./WhatYoullLearn";

const colors = {
  primary_Bg: "bg-[#121313]",
  secondary_Bg: "bg-[#1E1E1E]",
  primary_Hero: "bg-[#129274]",
  primary_Hero_Faded: "bg-[rgb(18, 146, 116, 0.24)]",
  secondary_Hero: "bg-[#64ACFF]",
  secondary_Hero_Faded: "bg-[rgb(100, 172, 255, 0.56)]",
  primary_Font: "text-[#FFFFFF]",
  secondary_Font: "text-[#B1AAA6]",
  special_Font: "text-[#64ACFF]",
  accent: "#B1AAA6",
  accent_Faded: "bg-[rgb(177, 170, 166, 0.41)]",
  primary_Icon: "white",
  secondary_Icon: "black",
  special_Icon: "#64ACFF",
  border_Accent: "border-[rgb(177, 170, 166, 0.41)]",
};

export default function Study() {
  return (
    <div className="h-full w-full grid grid-cols-1 lg:grid-cols-3 grid-rows-[auto_1fr] gap-5 overflow-hidden">
      {/* VIDEO */}
      <div className="lg:col-span-2 rounded-xl overflow-hidden bg-[#121313] aspect-video">
        <iframe
          src="https://www.youtube.com/embed/xR3V5Ow2dTI"
          className="w-full h-full"
          allow="autoplay; encrypted-media"
          allowFullScreen
        />
      </div>

      {/* WHAT YOU'LL LEARN */}
      <div className="lg:col-span-1">
        <WhatYouWillLearn />
      </div>

      {/* DESCRIPTION */}
      <div className="lg:col-span-3 text-white font-mono overflow-hidden">
        <h2 className="text-xl mb-2">Getting Started with HTML</h2>

        <div className="flex items-center gap-3 mb-3 text-sm text-gray-400">
          <span>John Doe</span>
          <span>•</span>
          <span>Beginner</span>
        </div>

        <p className="text-sm leading-relaxed text-gray-400 pb-2 border-b border-gray-600">
          Id sint voluptate incididunt occaecat qui mollit quis sint Lorem anim
          magna deserunt est anim velit...
        </p>

        <button className="mt-4 px-3 py-1 bg-black rounded-md text-sm w-fit">
          🔗 Resources
        </button>
      </div>
    </div>
  );
}
