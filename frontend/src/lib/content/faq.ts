export type FAQItem = {
  id: string;
  question: string;
  answer: string;
};

export const faqItems: FAQItem[] = [
  {
    id: "1",
    question: "What types of courses do you offer?",
    answer:
      "Our process is structured and collaborative. We begin with a consultation, followed by concept development, moodboards, and 3D visualizations. After approvals, we proceed with detailed drawings, material selection, and execution. We keep clients updated at every stage to ensure clarity and confidence.",
  },
  {
    id: "2",
    question: "Are the courses beginner-friendly?",
    answer:
      "Yes, our courses are designed to be beginner-friendly while also offering value to advanced learners. Each course starts with core concepts and gradually progresses to more advanced topics.",
  },
  {
    id: "3",
    question: "How do I access the courses after enrolling?",
    answer:
      "Once you enroll, you'll get instant access to the course through your account dashboard, where you can stream lessons anytime from any device.",
  },
  {
    id: "4",
    question: "Are the courses self-paced?",
    answer:
      "Yes, all courses are fully self-paced, allowing you to learn at your own speed and fit your studies around your schedule.",
  },
  {
    id: "5",
    question: "Do I receive a certificate after completion?",
    answer:
      "Yes, you will receive a certificate of completion after finishing the course, which you can use to showcase your skills and achievements.",
  },
];
