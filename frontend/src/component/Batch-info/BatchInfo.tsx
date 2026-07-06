"use client";
import React from "react";
import AdminBatchInfo from "./v1/V1BatchInfo";

type BatchInfoProps = {
  batch: any;
  institutionId: any
};

function BatchInfo({ batch, institutionId }: BatchInfoProps) {
  return (
    <div>
      <AdminBatchInfo batch={batch} institutionId={institutionId} />
    </div>
  );
}

export default BatchInfo;
