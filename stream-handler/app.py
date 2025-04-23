from fetch_tweets import fetch_tweets
from datetime import datetime

from flask import Flask, jsonify, request
from fetch_posts import fetch_reddit_posts, reddit
from dotenv import load_dotenv
import os
import psycopg2
import psycopg2.extras
import json

app = Flask(__name__)

DB_CONFIG = {
    'host': os.getenv('DB_HOST'),
    'port': os.getenv('DB_PORT'),
    'database': os.getenv('DB_NAME'),
    'user': os.getenv('DB_USER'),
    'password': os.getenv('DB_PASSWORD')
}

def get_db_connection():
    print("Creating DB connection...")
    return psycopg2.connect(**DB_CONFIG)

@app.route("/reddit_posts", methods=["POST"])
def reddit_posts():
    data = request.get_json()
    limit = data.get("limit", 10)
    
    if not isinstance(limit, int) or limit < 1 or limit > 100:
        return {"status": "error", "message": "Limit must be an integer between 1 and 100."}, 400

    print(f"Fetching {limit} Reddit posts...")
    posts = fetch_reddit_posts(COINS, QUESTIONS, limit=limit)
    print(f"Fetched {len(posts)} posts")
    return jsonify(posts)

# Twitter endpoint
@app.route("/twitter_posts", methods=["POST"])
def twitter_posts():
    data = request.get_json()
    query = data.get("query", "Bitcoin")
    limit = data.get("limit", 10)

    if not isinstance(limit, int) or limit < 1 or limit > 100:
        return {"status": "error", "message": "Limit must be an integer between 1 and 100."}, 400

    try:
        tweets = fetch_tweets(query=query, limit=limit)
        return jsonify(tweets)
    except Exception as e:
        return {"status": "error", "message": str(e)}, 500

@app.route("/reddit_status", methods=["GET"])
def reddit_status():
    try:
        me = reddit.user.me()
        return {"status": "success", "message": f"Authenticated. User: {me}"}
    except Exception as e:
        return {"status": "error", "message": str(e)}, 401

@app.route('/news', methods=['POST'])
def get_filtered_news():
    data = request.get_json(force=True) or {}
    start_iso      = data.get('start_date',  '2017-09-01')
    end_iso        = data.get('end_date',    '2025-01-31')
    currency_codes = data.get('currency_codes')

    try:
        start_dt = datetime.fromisoformat(start_iso.replace("Z", "+00:00"))
        end_dt   = datetime.fromisoformat(end_iso.replace("Z",   "+00:00"))
    except ValueError:
        return jsonify({"error": "Invalid ISO format for start_date or end_date"}), 400

    try:
        conn = get_db_connection()
        cur  = conn.cursor(cursor_factory=psycopg2.extras.DictCursor)

        query = """
            SELECT
              cn.id,
              cn.title,
              cn.url,
              cn.score,
              cn.newsdatetime,
              c.code AS currency_code
            FROM crypto_news cn
            JOIN news_currency nc ON cn.id = nc.newsid
            JOIN currency c        ON nc.currencyid = c.id
            WHERE cn.newsdatetime >= %s
              AND cn.newsdatetime <= %s
        """

        params = [start_dt, end_dt]

        if currency_codes:
            query += " AND c.code = ANY(%s)"
            params.append(currency_codes)

        query += " ORDER BY cn.newsdatetime DESC"

        cur.execute(query, params)
        rows = cur.fetchall()

        articles = [
            {
                "id":            row["id"],
                "title":         row["title"],
                "url":           row["url"],
                "score":         row["score"],
                "currency_code": row["currency_code"],
                "newsdatetime":  row["newsdatetime"].strftime("%Y-%m-%dT%H:%M:%SZ")
            }
            for row in rows
        ]

        cur.close()
        conn.close()
        return jsonify(articles)

    except Exception as e:
        return jsonify({"error": str(e)}), 500

COINS = [
    {"id": 1, "symbol": "BTC", "subreddit": "Bitcoin"},
    {"id": 2, "symbol": "ETH", "subreddit": "ethereum"},
    {"id": 3, "symbol": "USDT", "subreddit": "Tether+CryptoCurrency"},
    {"id": 4, "symbol": "XRP", "subreddit": "Ripple"},
    {"id": 5, "symbol": "BNB", "subreddit": "binance"},
    {"id": 6, "symbol": "SOL", "subreddit": "solana"},
    {"id": 7, "symbol": "USDC", "subreddit": "CryptoCurrency"},
    {"id": 8, "symbol": "TRX", "subreddit": "Tronix"},
    {"id": 9, "symbol": "DOGE", "subreddit": "dogecoin"},
    {"id": 10, "symbol": "ADA", "subreddit": "cardano"}
]

QUESTIONS = [
    {"id": 1, "label": "features",      "query": '"new features" OR "use cases"'},
    {"id": 2, "label": "leadership",    "query": '"founder" OR "CEO" OR "leadership"'},
    {"id": 3, "label": "security",      "query": '"hack" OR "security breach" OR "exploit"'},
    {"id": 4, "label": "market",        "query": '"price prediction" OR "market trend"'},
    {"id": 5, "label": "regulations",   "query": '"regulation" OR "government policy"'},
    {"id": 6, "label": "community",     "query": '"adoption" OR "community sentiment"'},
    {"id": 7, "label": "partnerships",  "query": '"partnership" OR "integration"'},
    {"id": 8, "label": "staking",       "query": '"mining" OR "staking" OR "validator"'}
]

# New reddit_db_dump endpoint
@app.route("/reddit_db_dump", methods=["POST"])
def reddit_db_dump():
    data = request.get_json()
    limit = data.get("limit", 2)
    time_filter = data.get("time_filter", "all")

    if not isinstance(limit, int) or limit < 1 or limit > 100:
        print("Invalid limit provided.")
        return {"status": "error", "message": "Limit must be an integer between 1 and 100."}, 400

    try:
        print(f"Fetching Reddit posts with limit={limit} and time_filter='{time_filter}'")
        posts = fetch_reddit_posts(COINS, QUESTIONS, limit=limit, time_filter=time_filter)
        print(f"Fetched {len(posts)} posts from Reddit")

        if not posts:
            print("No posts returned from fetch_reddit_posts.")
            return {"status": "success", "message": "No posts to insert."}

        conn = get_db_connection()
        print("Database connection established.")
        cur = conn.cursor()

        insert_query = """
            INSERT INTO raw_messages
            (source, external_id, title, message, metadata, created_utc)
            VALUES (%s, %s, %s, %s, %s, %s)
            ON CONFLICT DO NOTHING;
        """

        for post in posts:
            try:
                print(f"Inserting post: {post['title']} (score: {post['score']})")
                cur.execute(insert_query, (
                    "reddit",
                    post["url"],
                    post["title"],
                    post["text"],
                    json.dumps({
                        "coin": post["coin"],
                        "category": post["category"],
                        "author": post["author"],
                        "score": post["score"],
                        "num_comments": post["num_comments"]
                    }),
                    post["timestamp"]
                ))
                print(f"Inserted row count: {cur.rowcount}")
            except Exception as inner_e:
                print(f"Failed to insert post: {post}")
                print(f"Insert error: {inner_e}")

        conn.commit()
        print("Commit completed. Check the DB to confirm records.")
        print(f"Successfully committed {len(posts)} posts.")
        cur.close()
        conn.close()

        return {"status": "success", "message": f"{len(posts)} Reddit posts inserted into DB."}

    except Exception as e:
        print(f"Error during DB dump: {e}")
        return {"status": "error", "message": str(e)}, 500

if __name__ == "__main__":
    port = int(os.environ.get("PORT", 8080))
    app.run(debug=False, host='0.0.0.0', port=port)
