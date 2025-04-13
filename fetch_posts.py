import praw
from datetime import datetime
from dotenv import load_dotenv
from pathlib import Path
import os

load_dotenv(dotenv_path=Path(__file__).resolve().parent / ".env")

# Reddit API credentials (replace with your values)
reddit = praw.Reddit(
    client_id=os.getenv("REDDIT_CLIENT_ID"),
    client_secret=os.getenv("REDDIT_CLIENT_SECRET"),
    user_agent=os.getenv("REDDIT_USER_AGENT")
)

# Target subreddit
subreddit = reddit.subreddit("Bitcoin")

# Fetch top 10 hot posts
for post in subreddit.hot(limit=10):
    data = {
        "platform": "reddit",
        "title": post.title,
        "text": post.selftext,
        "timestamp": datetime.utcfromtimestamp(post.created_utc).isoformat(),
        "author": post.author.name if post.author else "unknown"
    }
    print(data)