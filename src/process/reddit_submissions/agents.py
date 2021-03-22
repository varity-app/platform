import faust

from process.app import app
from .views import submissions_topic


@app.agent(submissions_topic)
async def process_submission(submissions):
    async for submission in submissions:
        print(submission)
