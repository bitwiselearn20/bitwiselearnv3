import StudentSideBar from "@/component/general/StudentSidebar";
import React from "react";
import CodeCompiler from "./CodeCompiler";

function CompilerV1() {
  return (
    <div className="flex h-screen w-full bg-[#0f0f0f] text-white">
      <StudentSideBar />
      <div className="w-full">
        <CodeCompiler />
      </div>
    </div>
  );
}

export default CompilerV1;
