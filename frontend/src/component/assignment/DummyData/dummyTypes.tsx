export type Question = {
  id: string;
  question: string;
  choices: string[]; // exactly 4
  correctAnswer: string;
};

export type Assignment = {
  id: string;
  name: string;
  durationInMinutes: number;
  totalQuestions: number;
  questions: Question[];
};
