export type Testimonial = {
  id: string;
  name: string;
  role: string;
  company: string;
  quote: string;
  image: string;
};

export const testimonials: Testimonial[] = [
  {
    id: "1",
    name: "Arjun Singh",
    role: "SDE",
    company: "TCS",
    quote:
      "The DSA roadmap on BitwiseLearn is what got me through three technical rounds. I walked in having already solved variations of every question they asked.",
    image: "https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?w=200&h=200&fit=crop&crop=face",
  },
  {
    id: "2",
    name: "Kavya Reddy",
    role: "Frontend Developer",
    company: "Swiggy",
    quote:
      "Mock interviews felt uncomfortable at first — which is exactly why they worked. By the real thing, nothing caught me off guard.",
    image: "https://images.unsplash.com/photo-1580489944761-15a19d654956?w=200&h=200&fit=crop&crop=face",
  },
  {
    id: "3",
    name: "Vikram Kumar",
    role: "Data Engineer",
    company: "Infosys",
    quote:
      "I came in strong on theory and weak on practice. The company-specific training closed that gap fast, with problems shaped like what I'd actually be asked.",
    image: "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=200&h=200&fit=crop&crop=face",
  },
  {
    id: "4",
    name: "Ananya Desai",
    role: "Cloud Engineer",
    company: "Wipro",
    quote:
      "The AWS certification track gave me something to point to in interviews beyond my coursework — it's what separated my resume from the pile.",
    image: "https://images.unsplash.com/photo-1573497019940-1c28c88b4f3e?w=200&h=200&fit=crop&crop=face",
  },
  {
    id: "5",
    name: "Rohan Mehta",
    role: "Full Stack Developer",
    company: "Zomato",
    quote:
      "My mentor connect sessions were the difference between memorizing answers and actually understanding the systems I was being asked about.",
    image: "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=200&h=200&fit=crop&crop=face",
  },
];
