import React, { useState } from "react";
import Image from "next/image";
import Picture_Angad from "@/app/images/Picture_Angadveer.jpeg";
import Link from "next/link";
import { Github, Linkedin, Star } from "lucide-react";

export const ProfileCard = ({                                                                                      
  profilePicture = Picture_Angad,
  name = "Angadveer Singh",
  role = "Frontend Dev",
  github_id = "AadarshVerma7",
  linkedinUrl = "https://www.linkedin.com/in/aadarsh-verma-59323134a",
  description = "Passionate frontend developer with expertise in React and modern web technologies. Creating beautiful, responsive interfaces with attention to detail and user experience.",
}) => {
  const [isFlipped, setIsFlipped] = useState(false);

  return (
    <div
      className="group relative w-[18.5rem] h-96 cursor-pointer transition-transform duration-300 hover:scale-110 hover:shadow-2xl"
      style={{ perspective: "1000px" }}
      onClick={() => setIsFlipped((prev) => !prev)}
    >
      <div
        className="relative w-full h-full transition-transform duration-700"
        style={{
          transformStyle: "preserve-3d",
          transform: isFlipped ? "rotateY(180deg)" : "rotateY(0deg)",
        }}
      >
        {/* Front */}
        <div
          className="absolute inset-0 rounded-3xl overflow-hidden shadow-lg"
          style={{ backfaceVisibility: "hidden" }}
        >
          <div
            className="absolute inset-0 bg-cover bg-center"
            style={{ backgroundImage: `url(${profilePicture.src})` }}
          />
          <div className="absolute inset-0 bg-gradient-to-t from-black/70 to-transparent" />

          <div className="relative px-2 py-1 rounded-lg top-6 mx-auto bg-[#9a9a98] w-fit z-10">
            <div className="flex gap-1 items-center">
              <Star className="h-4 w-4 text-yellow-400" />
              <p className="text-white text-xs">{role}</p>
            </div>
          </div>

          <div className="relative flex justify-center mt-56 z-10">
            <h1 className="text-white font-medium text-3xl">{name}</h1>
          </div>

          {/* Social Popup */}
          <div className="absolute bottom-12 left-1/2 -translate-x-1/2 opacity-0 translate-y-6 group-hover:opacity-100 group-hover:translate-y-0 transition-all duration-700">
            <div className="bg-white/90 backdrop-blur-sm rounded-full px-4 py-2 shadow-lg flex gap-5">
              <Link href={linkedinUrl} target="_blank">
                <Linkedin className="h-6 w-6 text-black hover:scale-110 transition-transform" />
              </Link>
              <Link href={`https://github.com/${github_id}`} target="_blank">
                <Github className="h-6 w-6 text-black hover:scale-110 transition-transform" />
              </Link>
            </div>
          </div>
        </div>

        {/* Back */}
        <div
          className="absolute inset-0 rounded-3xl overflow-hidden"
          style={{
            transform: "rotateY(180deg)",
            backfaceVisibility: "hidden",
          }}
        >
          <Image
            src={profilePicture}
            alt="Profile"
            fill
            className="object-cover blur-sm"
          />
          <div className="absolute inset-0 bg-black/60" />

          <div className="relative h-full p-6 flex flex-col text-white">
            <h2 className="text-2xl font-bold mb-2">{name}</h2>
            <p className="text-sm font-semibold underline underline-offset-2 mb-3">
              {role}
            </p>
            <div className="flex-grow overflow-y-auto text-sm">
              {description}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProfileCard;
