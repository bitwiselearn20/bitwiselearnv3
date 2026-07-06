"use client";

import React, { useState } from "react";
import AssignmentInfo from "./AssignmentInfo";
import { addAssignmentToSection } from "@/api/courses/assignment/add-assignment-to-section";
import toast from "react-hot-toast";

type AddAssignmentV1Props = {
  sectionId: string;
  onClose: () => void;
};

export default function AddAssignmentV1({
  sectionId,
  onClose,
}: AddAssignmentV1Props) {
  const [assignment, setAssignment] = useState({
    title: "",
    description: "",
    instructions: "",
    marksPerQuestion: 0,
  });

  const [loading, setLoading] = useState(false);

  const submitAssignment = async () => {
    if (!assignment.title.trim()) {
      toast.error("Assignment title is required");
      return;
    }

    try {
      setLoading(true);
      toast.loading("Creating assignment...", { id: "assignment" });

      await addAssignmentToSection({
        name: assignment.title,
        description: assignment.description,
        instruction: assignment.instructions,
        marksPerQuestion: assignment.marksPerQuestion,
        sectionId,
      });

      toast.success("Assignment created!", { id: "assignment" });
      window.location.reload();
      onClose();
    } catch (error) {
      toast.error("Failed to create assignment", { id: "assignment" });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="my-[-8] flex justify-center">
      <div className={`rounded-xl w-[93%] max-w-3xl p-2`}>
        <AssignmentInfo
          assignment={assignment}
          setAssignment={setAssignment}
          onSubmit={submitAssignment}
          onClose={onClose}
          loading={loading}
        />
      </div>
    </div>
  );
}
