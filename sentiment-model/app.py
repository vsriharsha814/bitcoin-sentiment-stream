from flask import Flask, jsonify, request

from views import get_para_sentiments

app = Flask(__name__)

@app.route('/', methods=['GET'])
def index():
    data = "Hello World"
    return jsonify(data)

@app.route('/para-sentiment-analyze', methods=['POST'])
def para_sentiment_analyze():
    try:
        data = request.get_json()
        if not isinstance(data, list) or not all(isinstance(item, str) for item in data):
            return jsonify({"error": "Input must be a list of strings"}), 400

        sentiments = get_para_sentiments(data)
        return jsonify(sentiments)

    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == '__main__':
    app.run()