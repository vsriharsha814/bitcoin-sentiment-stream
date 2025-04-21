from datetime import datetime

from flask import Flask, jsonify, request
from fetch_posts import fetch_reddit_posts, reddit
from dotenv import load_dotenv
import os
import psycopg2
import psycopg2.extras

app = Flask(__name__)

DB_CONFIG = {
    'host': os.getenv('DB_HOST'),
    'port': os.getenv('DB_PORT'),
    'database': os.getenv('DB_NAME'),
    'user': os.getenv('DB_USER'),
    'password': os.getenv('DB_PASSWORD')
}

def get_db_connection():
    return psycopg2.connect(**DB_CONFIG)

@app.route("/reddit_posts", methods=["POST"])
def reddit_posts():
    data = request.get_json()
    limit = data.get("limit", 10)
    
    if not isinstance(limit, int) or limit < 1 or limit > 100:
        return {"status": "error", "message": "Limit must be an integer between 1 and 100."}, 400

    posts = fetch_reddit_posts(limit=limit)
    return jsonify(posts)

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

if __name__ == "__main__":
    port = int(os.environ.get("PORT", 8080))
    app.run(debug=False, host='0.0.0.0', port=port)