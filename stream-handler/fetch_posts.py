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

COIN_SUBREDDITS = {
    "BTC": "Bitcoin",
    "ETH": "ethereum",
    "USDT": "Tether+CryptoCurrency",
    "XRP": "Ripple",
    "BNB": "binance",
    "SOL": "solana",
    "USDC": "CryptoCurrency",
    "TRX": "Tronix",
    "DOGE": "dogecoin",
    "ADA": "cardano"
}

def fetch_reddit_posts(limit=10):
    posts = []
    for coin, subreddits in COIN_SUBREDDITS.items():
        subreddit = reddit.subreddit(subreddits)
        for post in subreddit.hot(limit=limit):
            posts.append({
                "coin": coin,
                "subreddit": subreddit.display_name,
                "title": post.title,
                "text": post.selftext,
                "timestamp": datetime.fromtimestamp(post.created_utc, tz=timezone.utc).isoformat(),
                "author": post.author.name if post.author else "unknown",
                "score": post.score,
                "url": post.url,
                "num_comments": post.num_comments
            })
    return posts