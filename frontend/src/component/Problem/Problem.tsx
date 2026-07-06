import V1Problem from "./v1/V1Problem";

function Problem({ data }: any) {
  return (
    <div className="w-full h-screen">
      <V1Problem data={data} />
    </div>
  );
}

export default Problem;
