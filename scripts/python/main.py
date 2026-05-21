import asyncio
import httpx
import os
import argparse
from faker import Faker
from rich.console import Console
from rich.progress import Progress, SpinnerColumn, TextColumn, BarColumn, TaskProgressColumn

console = Console()
faker = Faker()

BASE_URL = os.environ.get("API_URL", "http://localhost:8080")

async def register_user(client: httpx.AsyncClient) -> dict | None:
    email = faker.unique.email()
    password = "password"
    username = faker.name()

    resp = await client.post(f"{BASE_URL}/auth/signup", json={
        "email":email,
        "username":username,
        "password":password
    })

    if resp.status_code != 201 :
        console.print(f"[yellow]register failed[/yellow] {email} — {resp.status_code} {resp.text}")
        return None
    
    body = resp.json()
    data = body["data"]
    return {
        "email":email,
        "password":password,
        "access_token":data["access_token"]
    }


STATUSES = ["pending", "in_progress", "done"]

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


async def update_task_status(client: httpx.AsyncClient, token: str, task_id: str) -> None:
    # randomly move some tasks past pending so the data looks realistic
    import random
    status = random.choice(STATUSES)
    if status == "pending":
        return   # already pending by default, skip

    await client.put(
        f"{BASE_URL}/task/{task_id}",
        headers={"Authorization": f"Bearer {token}"},
        json={"is_completed": True},
    )


async def seed_tasks_for_user(client: httpx.AsyncClient, user: dict, count: int) -> int:
    created = 0
    for _ in range(count):
        ok = await create_task(client, user["access_token"])
        if ok:
            created += 1
    return created   
    

async def seed(user_count:int, tasks_per_user:int) -> None:
    total_users_created = 0
    total_tasks_created = 0
    total_failed = 0

    async with httpx.AsyncClient(timeout=10.0) as client:
        try:
            health = await client.get(f"{BASE_URL}/health")
            health.raise_for_status()
        except Exception as e:
            console.print(f"[red]cannot reach API at {BASE_URL}[/red] - {e}")
            console.print("is the server running? is API_URL set correctly")
            return
        
        console.print(f"[green]API reachable[/green] at {BASE_URL}")
        console.print(f"seeding {user_count} users * {tasks_per_user} tasks each\n")

        with Progress(
            SpinnerColumn(),
            TextColumn("[progess.description]{task.description}"),
            BarColumn(),
            TaskProgressColumn(),
            console=console,
        ) as progress:
            user_task = progress.add_task("Creating users...",  total=user_count)
            task_task = progress.add_task("Creating tasks...", total=user_count * tasks_per_user)

            for _ in range(user_count):
                user = await register_user(client)

                if user is None:
                    total_failed += 1
                    progress.update(user_task, advance=1)
                    continue

                total_users_created += 1
                progress.update(user_task, advance=1)

                created = await seed_tasks_for_user(client, user, tasks_per_user)
                total_tasks_created += created
                progress.update(task_task, advance=tasks_per_user)

        console.print("\n[bold]Done[/bold]")
        console.print(f"  Users created:  {total_users_created}")
        console.print(f"  Users failed:   {total_failed}")
        console.print(f"  Tasks created:  {total_tasks_created}")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="seed forge with fake users and tasks")
    parser.add_argument("--users",          type=int, default=10,  help="number of users to create (default: 10)")
    parser.add_argument("--tasks-per-user", type=int, default=5,   help="tasks per user (default: 5)")
    return parser.parse_args()


if __name__ == "__main__":
    args = parse_args()
    asyncio.run(seed(args.users, args.tasks_per_user))