"use client";

import { useEffect, useState } from "react";
import axios from "axios";

interface CodeQuestion {
  id: string;
  title: string;
  difficulty: string;
}

const AddCodeQuestion = () => {
  const [query, setQuery] = useState("");
  const [questions, setQuestions] = useState<CodeQuestion[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchQuestions = async (search: string) => {
    try {
      setLoading(true);
      const res = await axios.get(`/api/code-questions?search=${search}`);
      setQuestions(res.data);
    } catch (err) {
      // console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (query.length > 2) {
      fetchQuestions(query);
    }
  }, [query]);

  return (
    <div className="space-y-4">
      <input
        type="text"
        placeholder="Search coding questions..."
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        className="w-full border p-2 rounded"
      />

      {loading && <p>Loading...</p>}

      <ul className="space-y-2">
        {questions.map((q) => (
          <li key={q.id} className="border p-3 rounded flex justify-between">
            <span>{q.title}</span>
            <span className="text-sm text-gray-500">{q.difficulty}</span>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default AddCodeQuestion;
