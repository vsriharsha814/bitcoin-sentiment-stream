from flask import Flask, jsonify, request
from fetch_posts import fetch_reddit_posts, reddit
from datetime import timezone

app = Flask(__name__)

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

if __name__ == "__main__":
    app.run(debug=True)