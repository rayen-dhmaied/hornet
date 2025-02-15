from flask import Flask, jsonify, request
import requests
import os
from datetime import datetime, timezone

app = Flask(__name__)

# Read environment variables for service URLs and Flask port
FOLLOWERS_SERVICE_URL= os.getenv("FOLLOWERS_SERVICE_URL", "http://localhost:5001")
POSTS_SERVICE_URL = os.getenv("POSTS_SERVICE_URL", "http://localhost:5002")
FLASK_PORT = int(os.getenv("FLASK_PORT", 5000))

def calculate_post_score(post):
    """Calculate a score for the post based on engagement and recency."""
    engagement_score = post["replies_count"] + post["shares_count"]
    created_at = datetime.fromisoformat(post["created_at"].replace("Z", "+00:00"))

    # Convert utcnow() to an offset-aware datetime
    now = datetime.utcnow().replace(tzinfo=timezone.utc)
    
    # Boost score for recent posts (within the last 7 days)
    recency_boost = max(0, 7 - (now - created_at).days)

    return engagement_score + recency_boost * 2

@app.route('/feed', methods=['GET'])
def get_recommended_posts():
    current_user_id = request.headers.get('X-User-ID')
    if not current_user_id:
        return jsonify({"error": "X-User-ID header is required"}), 400

    # Step 1: Fetch the user's following list
    followers_response = requests.get(f"{FOLLOWERS_SERVICE_URL}/followers/user/{current_user_id}/following")
    if followers_response.status_code != 200:
        return jsonify({"error": "Failed to fetch following list"}), 500

    following_list = followers_response.json()
    following_ids = [f["receiver_id"] for f in following_list]

    # Step 2: Fetch posts from each followed user
    posts = []
    for user_id in following_ids:
        posts_response = requests.get(f"{POSTS_SERVICE_URL}/posts/author/{user_id}")
        if posts_response.status_code == 200:
            user_posts = posts_response.json()
            posts.extend(user_posts)

    # Step 3: Filter out replies (posts with parent_post_id)
    filtered_posts = [post for post in posts if not post.get("parent_post_id")]

    # Step 4: Calculate a recommendation score for each post
    scored_posts = []
    for post in filtered_posts:
        score = calculate_post_score(post)
        scored_posts.append({**post, "score": score})

    # Step 5: Sort posts by score in descending order
    recommended_posts = sorted(scored_posts, key=lambda x: x["score"], reverse=True)

    # Return top 20 recommended posts
    return jsonify(recommended_posts[:20]), 200

if __name__ == '__main__':
    debug_mode = os.getenv("FLASK_DEBUG", "0") == "1"
    app.run(host="0.0.0.0", port=FLASK_PORT, debug=debug_mode)