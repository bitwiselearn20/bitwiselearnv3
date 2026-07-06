"use client";

import { useEffect, useMemo, useState, useRef } from "react";
import { useRouter, usePathname } from "next/navigation";
import { motion } from "framer-motion";
import { ChevronDown, ChevronUp, SearchX, Search } from "lucide-react";

type Status = "all" | "pending" | "submitted";

const filterLabel = (v: Status) => {
  if (v === "all") return "All Levels";
  if (v === "pending") return "Pending";
  return "Submitted";
};

const filters = [
  { label: "All Levels", value: "all" },
  { label: "Pending", value: "pending" },
  { label: "Submitted", value: "submitted" },
];

const assignmentsData = [
  {
    id: 1,
    name: "<Assignment_Name>",
    description:
      "Proident minim mollit irure est et Lorem deserunt ipsum nisi irure qui.",
    questions: 10,
    duration: "60 min",
    issuedBy: "John Doe",
    dueIn: "9 Days",
    status: "submitted" as Status,
  },
  {
    id: 2,
    name: "<Assignment_Name>",
    description:
      "Proident minim mollit irure est et Lorem deserunt ipsum nisi irure qui.",
    questions: 10,
    duration: "60 min",
    issuedBy: "John Doe",
    dueIn: "9 Days",
    status: "pending" as Status,
  },
  {
    id: 3,
    name: "<Assignment_Name>",
    description:
      "Proident minim mollit irure est et Lorem deserunt ipsum nisi irure qui.",
    questions: 10,
    duration: "60 min",
    issuedBy: "John Doe",
    dueIn: "9 Days",
    status: "pending" as Status,
  },
];

export default function Assignments() {
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [filter, setFilter] = useState<Status>("all");
  const [open, setOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const timer = setTimeout(() => setLoading(false), 800);
    return () => clearTimeout(timer);
  }, []);

  // overlay: click outside to close
  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(e.target as Node)
      ) {
        setOpen(false);
      }
    };
    document.addEventListener("mousedown", handler);
    return () => document.removeEventListener("mousedown", handler);
  }, []);

  const filteredAssignments = useMemo(() => {
    return assignmentsData.filter((a) => {
      const matchesSearch = a.name.toLowerCase().includes(search.toLowerCase());
      const matchesFilter = filter === "all" ? true : a.status === filter;
      return matchesSearch && matchesFilter;
    });
  }, [search, filter]);

  const router = useRouter();
  const pathname = usePathname();

  const handleRoute = () => {
    router.push(`${pathname}/assignment/assignment-1`);
  };

  return (
    <div className="h-full w-full p-4 space-y-6 overflow-y-auto">
      {/* SEARCH + FILTER */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div className="relative sm:w-185">
          <Search
            size={16}
            className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"
          />
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search"
            className="bg-[#121313] w-full pl-9 pr-4 sm:w-1/2 py-2 rounded-md outline-none text-sm"
          />
        </div>
        <div className="relative w-fit hidden sm:block" ref={dropdownRef}>
          <button
            onClick={() => setOpen(!open)}
            className="
              bg-[#0F0F10]
              border border-[#1f1f1f]
              px-4 py-2
              rounded-xl
              flex items-center gap-2
              text-sm text-white
              cursor-pointer
              transition-all
              hover:border-gray-500
            "
          >
            {filterLabel(filter)}
            {open ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
          </button>

          {open && (
            <div
              className="
                absolute right-0 mt-2 w-44
                bg-[#0f0f0f]
                rounded-xl
                border border-white/10
                shadow-lg
                overflow-hidden
                z-50
              "
            >
              {(["all", "pending", "submitted"] as Status[]).map((f) => (
                <button
                  key={f}
                  onClick={() => {
                    setFilter(f);
                    setOpen(false);
                  }}
                  className="
                    w-full text-left
                    px-4 py-2
                    text-sm text-gray-200
                    hover:bg-white/5
                    transition-colors
                  "
                >
                  {filterLabel(f)}
                </button>
              ))}
            </div>
          )}
        </div>
      </div>

      {!loading && filteredAssignments.length === 0 && (
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.3, ease: "easeOut" }}
          className="flex flex-col items-center justify-center h-72 text-center gap-4"
        >
          <motion.div
            animate={{ y: [0, -12, 0] }}
            transition={{ repeat: Infinity, duration: 1.8, ease: "easeInOut" }}
            className="p-4 rounded-full"
          >
            <SearchX size={62} className="text-blue-400" />
          </motion.div>

          <div className="space-y-1">
            <p className="text-xl text-gray-200">No Assignments Found</p>
            <p className="text-sm text-gray-500">
              No assignments match your current search.
            </p>
          </div>
        </motion.div>
      )}

      {/* CARDS */}
      {filteredAssignments.length > 0 && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {loading
            ? Array.from({ length: 3 }).map((_, i) => (
                <div
                  key={i}
                  className="bg-[#121313] rounded-xl p-4 h-64 animate-pulse"
                />
              ))
            : filteredAssignments.map((a) => (
                <motion.div
                  key={a.id}
                  whileHover={{ y: -6 }}
                  transition={{ ease: "easeOut", duration: 0.2 }}
                  className="bg-[#121313] rounded-xl p-4 flex flex-col gap-4 shadow-md"
                >
                  <div className="flex items-center justify-between">
                    <h3 className="text-sm text-blue-400 font-mono">
                      {a.name}
                    </h3>

                    <span
                      className={`text-[10px] px-2 py-1 rounded-full ${
                        a.status === "submitted"
                          ? "bg-gray-600"
                          : "bg-linear-to-r from-red-500 to-red-600 text-white"
                      }`}
                    >
                      {a.status}
                    </span>
                  </div>

                  <p className="text-xs text-gray-400 leading-relaxed">
                    {a.description}
                  </p>

                  <div className="text-xs text-gray-300 space-y-1">
                    <div className="flex justify-between">
                      <span>Questions:</span>
                      <span>{a.questions}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Duration:</span>
                      <span>{a.duration}</span>
                    </div>
                  </div>

                  <div className="flex justify-between text-xs text-gray-400">
                    <span>Issued By:</span>
                    <span className="text-white">{a.issuedBy}</span>
                  </div>

                  <div className="flex justify-between text-xs text-gray-400">
                    <span>Due in:</span>
                    <span>{a.dueIn}</span>
                  </div>

                  {a.status === "submitted" ? (
                    <button
                      disabled
                      className="mt-2 bg-gray-700 text-xs py-2 rounded-md cursor-not-allowed"
                    >
                      Already Submitted
                    </button>
                  ) : (
                    <button
                      className="mt-2 bg-[#64ACFF] text-black text-xs py-2 rounded-md cursor-pointer"
                      onClick={handleRoute}
                    >
                      Start Now
                    </button>
                  )}
                </motion.div>
              ))}
        </div>
      )}
    </div>
  );
}
