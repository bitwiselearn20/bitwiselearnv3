from fastapi import FastAPI

app = FastAPI(title="BitwiseLearn Notification Service")

from routers.contact import router as contact_router

app.include_router(contact_router)


@app.get("/health")
async def health():
    return {"status": "ok", "service": "notification"}
