import asyncio
import os 
import httpx
from faker import Faker

faker = Faker()

BASE_URL = os.environ.get("API_URL", "http://localhost:8080")

async def create_task(client: httpx.AsyncClient, token:str) -> bool:
    resp = await client.post(
        f"{BASE_URL}/task",
        headers={"Authorization":f"Bearer {token}"},
        json={
            "title": faker.sentence(nb_words=5).rstrip("."),
            "description": faker.paragraph(nb_sentences=2)
        }
    )

    return resp.status_code == 201


async def main():
    cleint = httpx.AsyncClient(timeout=10)
    token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODAwMDM4ODYsInVzZXJJZCI6IkNBbDFwaHdWVWwtdHBRUjNocUdpIn0.Aej-h4t7XF4SnC4K0p43DYHHsIxZFxJPp7l1f7j_-5U"
    ok = await create_task(cleint, token)
    if not ok:
        print("failed")
        return
    print("done")    
        

if __name__ == "__main__":
    asyncio.run(main())