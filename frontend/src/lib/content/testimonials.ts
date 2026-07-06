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
    name: "Mark Manhold",
    role: "UI Designer",
    company: "Google",
    quote: "Bitwise Learn's UX Design track was a complete game-changer",
    image: "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=200&h=200&fit=crop&crop=face",
  },
  {
    id: "2",
    name: "James Hook",
    role: "UI Designer",
    company: "Olax",
    quote: "The lessons were clear, practical, and easy to apply right away.",
    image: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=200&h=200&fit=crop&crop=face",
  },
  {
    id: "3",
    name: "Jhon Lee",
    role: "Jr Marketer",
    company: "Dalax",
    quote: "Bitwise Learn gave me confidence to try things on my own.",
    image: "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=200&h=200&fit=crop&crop=face",
  },
  {
    id: "4",
    name: "Brock Clinton",
    role: "AI Engineer",
    company: "Frebrik",
    quote: "Everything was explained in a simple, no-fluff way.",
    image: "https://images.unsplash.com/photo-1519345182560-3f2917c472ef?w=200&h=200&fit=crop&crop=face",
  },
  {
    id: "5",
    name: "Jonathan Wisley",
    role: "SEO Expert",
    company: "Frebrik",
    quote: "Bitwise Learn's SEO program changed the way I think about organic growth",
    image: "https://images.unsplash.com/photo-1507591064344-4c6ce005b128?w=200&h=200&fit=crop&crop=face",
  },
  {
    id: "6",
    name: "Michel Clerk",
    role: "AI Engineer",
    company: "Dalax",
    quote: "Bitwise Learn's AI & ML course track reshaped how I approach",
    image: "https://images.unsplash.com/photo-1560250097-0b93528c311a?w=200&h=200&fit=crop&crop=face",
  },
];
