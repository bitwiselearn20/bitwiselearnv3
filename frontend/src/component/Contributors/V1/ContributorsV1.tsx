"use client";

import { desc } from "framer-motion/client";
import ProfileCard from "./ProfileCard";
import Picture_Angadveer from "@/app/images/Picture_Angadveer.jpeg";
import Picture_Angad from "@/app/images/Picture_Angad.jpg";
import Picture_Adheesh from "@/app/images/Picture_Adheesh.jpg";
import Picture_Aadarsh from "@/app/images/Picture_Aadarsh.jpg";
import Picture_Aayush from "@/app/images/Picture_Aayush.jpg";

const contributors = [
  {
    name: "Aadarsh Verma",
    role: "FullStack SWE",
    Picture: Picture_Aadarsh,
    description:
      "Full-stack SWE, building features end to end using React, Next.js, Node.js, and Express. I enjoy turning ideas into clean, scalable code, refining product flows, and solving tricky problems while keeping performance and usability front and center.",
    linkedinUrl: "https://www.linkedin.com/in/aadarsh-verma-59323134a",
    github_id: "AadarshVerma7",
  },
  {
    name: "Angad Sudan",
    role: "Technical Lead & Full Stack SWE",
    Picture: Picture_Angad,
    description:
      "I love to work with systems, checking how it works under load. Designing system, Solving problems and innovating with ideas that genuinely brings meaning and improves productivity. Not a vibe-coder but i vibe with the code and i dont slay bugs they slay me.",
    linkedinUrl: "https://www.linkedin.com/in/angadsudan",
    github_id: "AngadSudan",
  },
  {
    name: "Angadveer Singh",
    role: "DevOps & Frontend SWE",
    Picture: Picture_Angadveer,
    description:
      "In the tech industry, I have a strong interest in Web Development and DevOps, and I enjoy exploring how modern web applications are built, and deployed. I am particularly curious about creating scalable, efficient solutions and understanding the tools and workflows that support smooth development and deployment processes.",
    linkedinUrl: "https://www.linkedin.com/in/angadveer-singh-1751842b2",
    github_id: "Angadveer185",
  },
  {
    name: "Adheesh Verma",
    role: "FullStack SWE",
    Picture: Picture_Adheesh,
    description:
      "I’m Adheesh, a curious builder and problem-solver who loves turning ideas into real projects. I explore electronics  & programming, system design, enjoy experimenting with Linux, and constantly push myself to learn deeper, create smarter solutions, and build things that actually make an impact",
    linkedinUrl: "https://www.linkedin.com/in/adheesh-verma-177538324/",
    github_id: "AdheeshVerma",
  },
  {
    name: "Aayush Vats",
    role: "Full Stack SWE",
    Picture: Picture_Aayush,
    description:
      "I am a second-year computer science student who loves building real-world tech. I work with JavaScript, Java, and Next.js, explore backend systems, and contribute to open source. Hackathons, problem-solving, and impactful projects in education, healthcare, and safety drive my learning journey. I enjoy mentoring peers and continuously pushing my limits forward.",
    linkedinUrl: "https://www.linkedin.com/in/theaayushvats/",
    github_id: "Aayush-0821",
  },
];

export default function ContributorsV1() {
  return (
    <div className="w-full overflow-hidden py-20 bg-black">
      <h2 className="text-white text-4xl font-bold text-center mb-12">
        Contributors
      </h2>

      <div className="relative">
        <div className="animate-scroll flex gap-10 px-10">
          {[...contributors, ...contributors].map((c, index) => (
            <ProfileCard
              key={index}
              name={c.name}
              role={c.role}
              profilePicture={c.Picture}
              description={c.description}
              linkedinUrl={c.linkedinUrl}
              github_id={c.github_id}
            />
          ))}
        </div>
      </div>
    </div>
  );
}
