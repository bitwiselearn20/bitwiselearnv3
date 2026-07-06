"use client";

import {
  CircleCheck,
  Plus,
  CirclePlus,
  Pencil,
  Trash2,
  AlertTriangle,
  X,
} from "lucide-react";
import "./helper.css";
import { useState } from "react";

/* ================= TYPES ================= */

type TopicContent = {
  video: boolean;
  transcript: boolean;
  file: boolean;
};

type Topic = {
  id: number;
  title: string;
  description: string;
  isCompleted: boolean;
  isEditing: boolean;
  showContentOptions: boolean;
  showDeleteConfirm: boolean;
  contents: TopicContent;
};

type Section = {
  sectionId: number;
  sectionName: string;
  isCompleted: boolean;
  topics: Topic[];
};

type AddSectionProps = {
  sectionNumber: number;
  sectionData?: Section;
};

/* ================= COMPONENT ================= */

const AddSectionV1 = ({ sectionNumber, sectionData }: AddSectionProps) => {
  const [sectionCompleted, setSectionCompleted] = useState(false);
  const [sectionError, setSectionError] = useState<string | null>(null);

  const [sectionName, setSectionName] = useState(
    sectionData?.sectionName ?? "",
  );

  const [topics, setTopics] = useState<Topic[]>(
    sectionData?.topics ?? [
      {
        id: 1,
        title: "",
        description: "",
        isCompleted: false,
        isEditing: true,
        showContentOptions: false,
        showDeleteConfirm: false,
        contents: {
          video: false,
          transcript: false,
          file: false,
        },
      },
    ],
  );

  /* ================= HELPERS ================= */

  const updateTopic = (
    id: number,
    field: "title" | "description",
    value: string,
  ) => {
    setTopics((prev) =>
      prev.map((topic) =>
        topic.id === id ? { ...topic, [field]: value } : topic,
      ),
    );
  };

  const addNewTopic = () => {
    setTopics((prev) => [
      ...prev,
      {
        id: Date.now(),
        title: "",
        description: "",
        isCompleted: false,
        isEditing: true,
        showContentOptions: false,
        showDeleteConfirm: false,
        contents: {
          video: false,
          transcript: false,
          file: false,
        },
      },
    ]);
  };

  /* ---------- SECTION DELETE ---------- */

  const deleteSection = () => {};

  /* ---------- DELETE TOPIC ---------- */

  const askDeleteTopic = (id: number) => {
    setTopics((prev) =>
      prev.map((topic) =>
        topic.id === id ? { ...topic, showDeleteConfirm: true } : topic,
      ),
    );
  };

  const cancelDeleteTopic = (id: number) => {
    setTopics((prev) =>
      prev.map((topic) =>
        topic.id === id ? { ...topic, showDeleteConfirm: false } : topic,
      ),
    );
  };

  const confirmDeleteTopic = (id: number) => {
    if (topics.length === 1) return;
    setTopics((prev) => prev.filter((topic) => topic.id !== id));
  };

  /* ---------- CONTENT ---------- */

  const toggleContentOptions = (id: number) => {
    setTopics((prev) =>
      prev.map((topic) =>
        topic.id === id
          ? { ...topic, showContentOptions: !topic.showContentOptions }
          : topic,
      ),
    );
  };

  const addContent = (id: number, type: keyof TopicContent) => {
    setTopics((prev) =>
      prev.map((topic) =>
        topic.id === id
          ? {
              ...topic,
              contents: { ...topic.contents, [type]: true },
            }
          : topic,
      ),
    );
  };

  /* ---------- COMPLETE / EDIT ---------- */

  const completeTopic = (id: number) => {
    if (sectionCompleted) return;

    setTopics((prev) =>
      prev.map((topic) =>
        topic.id === id
          ? { ...topic, isCompleted: true, isEditing: false }
          : topic,
      ),
    );
  };

  const editTopic = (id: number) => {
    if (sectionCompleted) return;

    setTopics((prev) =>
      prev.map((topic) =>
        topic.id === id
          ? { ...topic, isCompleted: false, isEditing: true }
          : topic,
      ),
    );
  };

  const completeSection = () => {
    const allTopicsCompleted = topics.every((topic) => topic.isCompleted);

    if (!allTopicsCompleted) {
      setSectionError(
        "Please complete all topics before completing the section.",
      );
      return;
    }

    setSectionError(null);
    setSectionCompleted(true);

    const sectionData: Section = {
      sectionId: sectionNumber,
      sectionName,
      isCompleted: true,
      topics,
    };
  };

  const editSection = () => {
    setSectionCompleted(false);
    setSectionError(null);

    setTopics((prev) =>
      prev.map((topic) => ({
        ...topic,
        isCompleted: false,
        isEditing: true,
      })),
    );
  };

  /* ================= UI ================= */

  return (
    <div className="relative text-white bg-divBg h-full w-[90%] rounded-2xl px-8 py-6 shadow-xl border border-white/5">
      {/* DELETE SECTION BUTTON */}
      <button
        onClick={deleteSection}
        disabled={sectionCompleted}
        className="absolute top-4 right-4 p-2 rounded-full hover:bg-red-500/20 transition disabled:opacity-40"
        title="Delete section"
      >
        <X size={18} className="text-red-400" />
      </button>

      {/* SECTION HEADER */}
      <div className="flex items-center justify-between mb-6">
        <div className="flex flex-col gap-3 w-full max-w-md">
          <div className="flex items-center gap-4">
            <h1 className="text-2xl font-semibold">Section {sectionNumber}</h1>
            <span className="text-xs px-3 py-1 rounded-full bg-white/10 text-white/70">
              Curriculum
            </span>
          </div>

          <div className="input-wrapper">
            <CircleCheck size={18} className="input-icon" />
            <input
              type="text"
              value={sectionName}
              disabled={sectionCompleted}
              onChange={(e) => setSectionName(e.target.value)}
              placeholder="Section name (e.g. Introduction to React)"
              className="w-full border border-white/15 bg-transparent rounded-xl px-3 py-2 text-sm placeholder:text-white/40 disabled:opacity-50"
            />
          </div>
        </div>

        {sectionCompleted && (
          <button
            onClick={editSection}
            className="flex items-center gap-1 text-sm text-primaryBlue hover:underline"
          >
            <Pencil size={14} />
            Edit Section
          </button>
        )}
      </div>

      <div className="pl-6 w-full space-y-8">
        {topics.map((topic, index) => (
          <div key={topic.id} className="space-y-4">
            <div className="flex items-center justify-between">
              <h2 className="text-xs uppercase tracking-wider text-white/60">
                Topic {index + 1}
              </h2>

              <div className="flex items-center gap-3">
                {topic.isCompleted && !sectionCompleted && (
                  <button
                    onClick={() => editTopic(topic.id)}
                    className="text-xs text-primaryBlue hover:underline"
                  >
                    Edit
                  </button>
                )}

                {index > 0 && !sectionCompleted && (
                  <button
                    onClick={() => askDeleteTopic(topic.id)}
                    className="p-1 rounded-full hover:bg-red-500/20 transition"
                  >
                    <Trash2 size={16} className="text-red-400" />
                  </button>
                )}

                <button
                  disabled={sectionCompleted}
                  onClick={() => completeTopic(topic.id)}
                  className="p-1 rounded-full hover:bg-white/10 transition"
                >
                  <CircleCheck
                    className={
                      topic.isCompleted ? "text-green-400" : "text-white/30"
                    }
                  />
                </button>
              </div>
            </div>

            {index > 0 && topic.showDeleteConfirm && !sectionCompleted && (
              <div className="flex items-center gap-3 bg-red-500/10 border border-red-500/30 rounded-lg px-3 py-2 text-sm">
                <AlertTriangle size={16} className="text-red-400" />
                <span className="flex-1">Delete this topic?</span>
                <button
                  onClick={() => confirmDeleteTopic(topic.id)}
                  className="text-red-400 font-semibold"
                >
                  Delete
                </button>
                <button
                  onClick={() => cancelDeleteTopic(topic.id)}
                  className="text-white/70"
                >
                  Cancel
                </button>
              </div>
            )}

            <div className="input-wrapper">
              <CircleCheck size={18} className="input-icon" />
              <input
                type="text"
                value={topic.title}
                disabled={!topic.isEditing || sectionCompleted}
                onChange={(e) => updateTopic(topic.id, "title", e.target.value)}
                placeholder="Title: Example Title"
                className="w-full border border-white/15 bg-transparent rounded-xl px-3 py-2 text-sm placeholder:text-white/40 disabled:opacity-50"
              />

              <button
                type="button"
                onClick={() => toggleContentOptions(topic.id)}
                disabled={!topic.isEditing || sectionCompleted}
                className="input-btn bg-primaryBlue text-black px-3 py-1.5 rounded-lg text-xs font-semibold flex items-center gap-1 disabled:opacity-50"
              >
                <Plus size={14} />
                Content
              </button>
            </div>

            <textarea
              value={topic.description}
              disabled={!topic.isEditing || sectionCompleted}
              onChange={(e) =>
                updateTopic(topic.id, "description", e.target.value)
              }
              placeholder="Description: Example abc"
              className="w-full border border-white/15 bg-transparent rounded-xl px-3 pt-2 text-sm resize-none overflow-hidden placeholder:text-white/40 disabled:opacity-50"
              rows={4}
            />

            {topic.showContentOptions &&
              topic.isEditing &&
              !sectionCompleted && (
                <div className="flex gap-3 pt-1">
                  {(["video", "transcript", "file"] as const).map((type) => (
                    <button
                      key={type}
                      onClick={() => addContent(topic.id, type)}
                      className="flex items-center gap-1 bg-primaryBlue/90 text-black px-3 py-1.5 rounded-full text-xs font-semibold hover:bg-primaryBlue transition"
                    >
                      <CirclePlus size={14} />
                      {type}
                    </button>
                  ))}
                </div>
              )}
          </div>
        ))}

        {sectionError && (
          <div className="flex items-center gap-2 bg-red-500/10 border border-red-500/30 rounded-lg px-4 py-3 text-sm text-red-300">
            <AlertTriangle size={16} />
            {sectionError}
          </div>
        )}

        <div className="flex justify-between pt-6 border-t border-white/10">
          <button
            onClick={addNewTopic}
            disabled={sectionCompleted}
            className="px-5 py-2 rounded-full bg-primary-hero text-black text-sm font-semibold hover:scale-105 transition disabled:opacity-50"
          >
            New Topic
          </button>

          <button
            onClick={completeSection}
            className="px-5 py-2 rounded-full bg-primary-hero text-black text-sm font-semibold hover:scale-105 transition"
          >
            Complete Section
          </button>
        </div>
      </div>
    </div>
  );
};

export default AddSectionV1;
