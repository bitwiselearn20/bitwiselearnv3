"use client";
import { useEffect } from "react";

const PDFViewer = ({ url }: { url: string }) => {
  useEffect(() => {
    const script = document.createElement("script");
    script.src = "https://acrobatservices.adobe.com/view-sdk/viewer.js";
    script.async = true;
    document.body.append(script);

    document.addEventListener("adobe_dc_view_sdk.ready", function () {
      const adobeDCView = new AdobeDC.View({
        clientId: "<Your Adobe Client ID>",
        divId: "pdf-view",
      });

      adobeDCView.previewFile(
        {
          content: {
            location: {
              url: `${url}`,
            },
          },
          metaData: { fileName: `${url.split("/").slice(-1)[0]}` },
        },
        {
          embedMode: "FULL_WINDOW",
          showAnnotationTools: false,
          showPrintPDF: true,
          showDownloadPDF: true,
          defaultViewMode: "FIT_WIDTH",
          enableLinearization: true,
        },
      );

      // Listen to events
      const eventsOptions = {
        listenOn: [
          AdobeDC.View.Enum.Events.APP_RENDERING_DONE,
          AdobeDC.View.Enum.Events.APP_RENDERING_FAILED,
        ],
        enableFilePreviewEvents: true,
      };

      adobeDCView.registerCallback(
        AdobeDC.View.Enum.CallbackType.EVENT_LISTENER,
        function (event: any) {
          switch (event.type) {
            case "APP_RENDERING_DONE":
              break;
            case "APP_RENDERING_FAILED":
              break;
          }
        },
        eventsOptions,
      );
    });
  }, []);

  return <div id="pdf-view" className="w-[40%] h-screen"></div>;
};
export default PDFViewer;
