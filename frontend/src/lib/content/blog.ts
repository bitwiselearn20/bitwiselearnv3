export type BlogPost = {
  id: string;
  title: string;
  category: string;
  date: string;
  image: string;
  slug: string;
};

export const blogPosts: BlogPost[] = [
  {
    id: "1",
    title: "The skills employers are looking for in 2025",
    category: "AI Trends",
    date: "December 12, 2025",
    image: "https://images.unsplash.com/photo-1498050108023-c5249f4df085?w=400&h=240&fit=crop",
    slug: "s3",
  },
  {
    id: "2",
    title: "How to choose the right course for your career path",
    category: "Career",
    date: "December 12, 2025",
    image: "https://images.unsplash.com/photo-1522071820081-009f0129c71c?w=400&h=240&fit=crop",
    slug: "choose-right-course",
  },
  {
    id: "3",
    title: "The secret to staying motivated during online courses",
    category: "Design",
    date: "December 13, 2025",
    image: "https://images.unsplash.com/photo-1434030216411-0b793f4b4173?w=400&h=240&fit=crop",
    slug: "staying-motivated-online",
  },
  {
    id: "4",
    title: "Building a learning routine that actually sticks",
    category: "Productivity",
    date: "December 15, 2025",
    image: "https://images.unsplash.com/photo-1484480974693-6ca0a78fb36b?w=400&h=240&fit=crop",
    slug: "learning-routine-that-sticks",
  },
  {
    id: "5",
    title: "Why soft skills matter more than ever in tech",
    category: "Career",
    date: "December 16, 2025",
    image: "https://images.unsplash.com/photo-1522202176988-66273c2fd55f?w=400&h=240&fit=crop",
    slug: "soft-skills-in-tech",
  },
  {
    id: "6",
    title: "Getting the most out of video-based learning",
    category: "Learning",
    date: "December 18, 2025",
    image: "https://images.unsplash.com/photo-1516321318423-f06f85e504b3?w=400&h=240&fit=crop",
    slug: "video-based-learning",
  },
];
