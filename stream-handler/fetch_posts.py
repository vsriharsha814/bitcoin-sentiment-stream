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
    "BTC": (1, "Bitcoin"),
    "ETH": (2, "ethereum"),
    "USDT": (3, "Tether+CryptoCurrency"),
    "XRP": (4, "Ripple"),
    "BNB": (5, "binance"),
    "SOL": (6, "solana"),
    "USDC": (7, "CryptoCurrency"),
    "TRX": (8, "Tronix"),
    "DOGE": (9, "dogecoin"),
    "ADA": (10, "cardano")
}

def fetch_reddit_posts(limit=10, time_filter="all"):
    posts = []
    queries = {
        "features": (101, '"new features" OR "use cases"'),
        "leadership": (102, '"founder" OR "CEO" OR "leadership"'),
        "security": (103, '"hack" OR "security breach" OR "exploit"'),
        "market": (104, '"price prediction" OR "market trend"'),
        "regulations": (105, '"regulation" OR "government policy"'),
        "community": (106, '"adoption" OR "community sentiment"'),
        "partnerships": (107, '"partnership" OR "integration"'),
        "staking": (108, '"mining" OR "staking" OR "validator"')
    }

    for coin, (coin_id, subreddit_name) in COIN_SUBREDDITS.items():
        subreddit = reddit.subreddit(subreddit_name)
        for label, (question_id, query) in queries.items():
            full_query = f"{query} {coin}"
            for post in subreddit.search(full_query, sort="new" if time_filter != "all" else "relevance", time_filter=time_filter, limit=limit):
                posts.append({
                    "coin": coin,
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