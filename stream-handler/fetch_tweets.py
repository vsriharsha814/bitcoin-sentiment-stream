# import tweepy
# import os
# from dotenv import load_dotenv
# from pathlib import Path
#
# # Load environment variables
# load_dotenv(dotenv_path=Path(__file__).resolve().parent / ".env")
#
# # Set up API credentials
# api_key = os.getenv("TWITTER_API_KEY")
# api_secret = os.getenv("TWITTER_API_SECRET")
# access_token = os.getenv("TWITTER_ACCESS_TOKEN")
# access_secret = os.getenv("TWITTER_ACCESS_SECRET")
#
# # Authenticate
# auth = tweepy.OAuth1UserHandler(api_key, api_secret, access_token, access_secret)
# api = tweepy.API(auth)
#
# # Fetch tweets
# def fetch_tweets(query, limit=10):
#     tweets = []
#     for tweet in tweepy.Cursor(api.search_tweets, q=query, lang="en", tweet_mode="extended").items(limit):
#         tweets.append({
#             "user": tweet.user.screen_name,
#             "text": tweet.full_text,
#             "created_at": tweet.created_at.isoformat(),
#             "likes": tweet.favorite_count,
#             "retweets": tweet.retweet_count
#         })
#     return tweets
#
# # Example test
# if __name__ == "__main__":
#     results = fetch_tweets("Bitcoin price prediction", limit=5)
#     for r in results:
#         print(r)