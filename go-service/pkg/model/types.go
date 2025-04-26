package model

// SentimentRequest represents the incoming JSON payload.
type SentimentRequest struct {
	MessageScores map[string][]float64 `json:"messageScores"`
}

// SentimentResponse is what we send back.
type SentimentResponse struct {
	FinalSentiment float64 `json:"finalSentiment"`
}

// DefaultWeights are your per-question weights.
var DefaultWeights = map[string]float64{
	"New Features or Use Cases of \"coin_name\"":                                  0.15,
	"Founders or Leadership of \"coin_name\"":                                     0.10,
	"Security Concerns or Hacks related to \"coin_name\"":                         0.20,
	"Market Trends and Price Predictions of \"coin_name\"":                        0.20,
	"Regulatory Updates and Government Policies affecting \"coin_name\"":          0.10,
	"Community Sentiment and Adoption for \"coin_name\"":                          0.10,
	"Partnerships and Integrations involving \"coin_name\"":                       0.10,
	"Mining and Staking Discussions around \"coin_name\"":                         0.05,
}

// DefaultMessageScores is used by the WebSocket fallback.
var DefaultMessageScores = map[string][]float64{
	"1": {0.1, 0.2, 0.15, 0.1, 0.3, 0.05, 0.2, 0.15, 0.1, 0.2},
	"2": {0.2, 0.1, 0.15, 0.25, 0.1, 0.2, 0.15, 0.2, 0.1, 0.2},
	"3": {-0.1, -0.2, -0.15, -0.1, -0.3, -0.05, -0.2, -0.15, -0.1, -0.2},
	"4": {0.3, 0.25, 0.2, 0.35, 0.3, 0.25, 0.2, 0.35, 0.3, 0.25},
	"5": {0.0, 0.05, -0.05, 0.0, 0.1, 0.0, -0.1, 0.05, 0.0, 0.1},
	"6": {0.2, 0.3, 0.25, 0.2, 0.15, 0.2, 0.25, 0.3, 0.2, 0.15},
	"7": {0.1, 0.05, 0.1, 0.15, 0.1, 0.05, 0.1, 0.15, 0.1, 0.05},
	"8": {0.05, 0.0, 0.05, 0.1, 0.05, 0.0, 0.05, 0.1, 0.05, 0.0},
}

// CalculateFinalSentiment re-weights only across questions
// that actually have scores in messageScores.
func CalculateFinalSentiment(
    weights map[string]float64,
    messageScores map[string][]float64,
) float64 {
    // 1) Figure out the sum of weights for questions we do have
    totalWeight := 0.0
    for q, scores := range messageScores {
        if len(scores) == 0 {
            continue
        }
        if w, ok := weights[q]; ok {
            totalWeight += w
        }
    }
    // Nothing to do if no questions contributed
    if totalWeight == 0 {
        return 0
    }

    // 2) Build the weighted average across only those questions
    final := 0.0
    for q, scores := range messageScores {
        if len(scores) == 0 {
            continue
        }
        w, ok := weights[q]
        if !ok {
            continue
        }
        // normalize this question’s weight
        norm := w / totalWeight

        // compute the average score for this question
        sum := 0.0
        for _, s := range scores {
            sum += s
        }
        avg := sum / float64(len(scores))

        final += norm * avg
    }
    return final
}
