import praw
from dotenv import load_dotenv
from pathlib import Path
import os
from datetime import datetime, timezone

# Load env vars from .env
load_dotenv(dotenv_path=Path(__file__).resolve().parent / ".env")

# Set up Reddit client
reddit = praw.Reddit(
    client_id=os.getenv("REDDIT_CLIENT_ID"),
    client_secret=os.getenv("REDDIT_CLIENT_SECRET"),
    user_agent=os.getenv("REDDIT_USER_AGENT")
)

def fetch_reddit_posts(coins, questions, limit=10, time_filter="all"):
    posts = []

    for coin in coins:
        coin_id = coin["id"]
        subreddit = reddit.subreddit(coin["subreddit"])
        for question in questions:
            label = question["label"]
            question_id = question["id"]
            query = question["query"]
            full_query = f"{query} {coin['code']}"
            for post in subreddit.search(full_query, sort="new" if time_filter != "all" else "relevance", time_filter=time_filter, limit=limit):
                posts.append({
                    "coin": coin["code"],
                    "coin_id": coin_id,
                    "subreddit": subreddit.display_name,
                    "category": label,
                    "question_id": question_id,
                    "title": post.title,
                    "text": post.selftext,
                    "timestamp": datetime.fromtimestamp(post.created_utc, tz=timezone.utc).isoformat(),
                    "author": post.author.name if post.author else "unknown",
                    "score": post.score,
                    "url": post.url,
                    "num_comments": post.num_comments
                })
    return posts