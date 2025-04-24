just fetch the posts and fetch sentiment scores for them and send them into the database.

you'll definitely need the .env files to make the reddit calls and database calls.

run the app.py for the app to start

the api call should look something like http://localhost:8080/reddit_db_dump and the data for this something like { "limit" : 2, "time_filter" : "day" }