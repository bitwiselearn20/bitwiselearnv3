export default function VideoSection({ loading }: { loading: boolean }) {
  if (loading) {
    return (
      <div className="w-full h-55 sm:h-75 lg:h-full bg-[#121313] animate-pulse flex items-center justify-center"></div>
    );
  }

  return (
    <div className="w-full h-55 sm:h-75 lg:h-full bg-black">
      <iframe
        src="https://www.youtube.com/embed/xR3V5Ow2dTI"
        className="w-full h-full"
        allow="autoplay; encrypted-media"
        allowFullScreen
      />
    </div>
  );
}
