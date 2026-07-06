import React from "react";
import AttemptV1 from "./v1/AttemptV1";
import { AttemptMode } from "./v1/types";

type AttemptProps = {
  id: string;
  mode: AttemptMode;
};


export default function Attempt(data : AttemptProps) {
  return <AttemptV1 id={data.id} mode={data.mode} />;
}
