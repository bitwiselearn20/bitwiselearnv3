"use client";
import React from "react";
import AdminInstituteInfo from "./v1/V1InstitutionInfo";

type InstitutionInfoProps = {
  institution: any;
};

function InstitutionInfo({ institution }: InstitutionInfoProps) {
  return (
    <div>
      <AdminInstituteInfo institution={institution} />
    </div>
  );
}

export default InstitutionInfo;
