import faust
from nltk.sentiment.vader import SentimentIntensityAnalyzer

from util.constants.sentiment import Estimators

from process.app import app
from .views import scraped_posts_topic
from ..sentiment.views import sentiment_topic
from ..sentiment.models import SentimentEstimate
  

@app.agent(scraped_posts_topic)
async def estimate_sentiment(posts):
    """Estimate sentiment for scraped posts"""
    model = SentimentIntensityAnalyzer()

    async for post in posts:
        # Estimate sentiment
        scores = model.polarity_scores(post.text)

        # Parse object
        sentiment = SentimentEstimate(
            data_source=post.data_source,
            parent_source=post.parent_source,
            parent_id=post.parent_id,
            estimator=Estimators.VADER,
            compound=scores['compound'],
            positive=scores['pos'],
            neutral=scores['neu'],
            negative=scores['neg'],
        )

        # Publish to Kafka
        await sentiment_topic.send(value=sentiment)
