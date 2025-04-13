from flask import Flask, jsonify
from fetch_posts import fetch_reddit_posts

app = Flask(__name__)

# In-memory post cache
cached_posts = []

@app.route("/fetch-posts", methods=["GET"])
def fetch_and_cache():
    global cached_posts
    cached_posts = fetch_reddit_posts()
    return {"message": "Posts fetched and cached."}

@app.route("/reddit-posts", methods=["GET"])
def get_reddit_posts():
    return jsonify(cached_posts)

if __name__ == "__main__":
    app.run(debug=True)