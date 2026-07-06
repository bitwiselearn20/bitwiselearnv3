import { Assignment } from "./dummyTypes";

export const dummyAssignmentData: Record<string, Assignment> = {
  "assignment-1": {
    id: "assignment-1",
    name: "Introduction to Data Structures",
    durationInMinutes: 30,
    totalQuestions: 5,
    questions: [
      {
        id: "q1",
        question: "Which data structure follows the FIFO principle?",
        choices: ["Stack", "Queue", "Tree", "Graph"],
        correctAnswer: "Queue",
      },
      {
        id: "q2",
        question:
          "What is the time complexity of accessing an element in an array by index?",
        choices: ["O(n)", "O(log n)", "O(1)", "O(n log n)"],
        correctAnswer: "O(1)",
      },
      {
        id: "q3",
        question:
          "Which data structure is best suited for implementing recursion?",
        choices: ["Queue", "Linked List", "Stack", "Heap"],
        correctAnswer: "Stack",
      },
      {
        id: "q4",
        question:
          "Which traversal of a binary search tree gives sorted output?",
        choices: ["Preorder", "Postorder", "Inorder", "Level order"],
        correctAnswer: "Inorder",
      },
      {
        id: "q5",
        question: "Which of the following is NOT a linear data structure?",
        choices: ["Array", "Linked List", "Stack", "Graph"],
        correctAnswer: "Graph",
      },
    ],
  },
};
