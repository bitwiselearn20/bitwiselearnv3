import React from "react";
import V1VendorInstitutions from "./v1/V1VendorInstitutions";

type VendorInstitutionsProps = {
  vendorId: string;
};

function VendorInstitutions({ vendorId }: VendorInstitutionsProps) {
  return (
    <div>
      <V1VendorInstitutions vendorId={vendorId} />
    </div>
  );
}

export default VendorInstitutions;
