import praw
from dotenv import load_dotenv
from pathlib import Path
import os
from datetime import datetime

# Load env vars from .env
load_dotenv(dotenv_path=Path(__file__).resolve().parent / ".env")

# Set up Reddit client
reddit = praw.Reddit(
    client_id=os.getenv("REDDIT_CLIENT_ID"),
    client_secret=os.getenv("REDDIT_CLIENT_SECRET"),
    user_agent=os.getenv("REDDIT_USER_AGENT")
)

def fetch_reddit_posts(limit=10):
    subreddit = reddit.subreddit("Bitcoin")
    posts = []
    for post in subreddit.hot(limit=limit):
        posts.append({
            "title": post.title,
            "text": post.selftext,
            "timestamp": datetime.fromtimestamp(post.created_utc, tz=datetime.timezone.utc).isoformat(),
            "author": post.author.name if post.author else "unknown"
        })
    return posts