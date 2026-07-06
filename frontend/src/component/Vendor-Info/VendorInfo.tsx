"use client";
import React from "react";
import V1VendorInfo from "./v1/V1VendorInfo";

type VendorInfoProps = {
  vendor: any;
};

function VendorInfo({ vendor }: VendorInfoProps) {
  return (
    <div>
      <V1VendorInfo vendor={vendor} />
    </div>
  );
}

export default VendorInfo;

