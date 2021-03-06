from .. import Message

class CommentMessage(Message):
    def __init__(
        self,
        comment_id: str,
        submission_id: str,
        subreddit: str,
        author: str,
        created_utc: str,
        body: str,
        upvotes: str,
        
    ):
        self.comment_id = comment_id
        self.submission_id = submission_id
        self.subreddit = subreddit
        self.author = author
        self.created_utc = created_utc
        self.body = self.get(body)
        self.upvotes = upvotes

    def get(self, value, default=""):
        if value is None:
            return default

        return value

    def serialize(self, as_str=True) -> dict:
        response = dict(
            comment_id=self.comment_id,
            submission_id=self.submission_id,
            subreddit=self.subreddit,
            author=self.author,
            created_utc=self.created_utc,
            body=self.body,
            upvotes=self.upvotes,
        )

        if as_str:
            response = self.to_str(response)

        return response

    @classmethod
    def from_obj(cls, data: dict):
        for field in [CC.ID, CC.SUBMISSION_ID, CC.SUBREDDIT, CC.AUTHOR, CC.CREATED_UTC]:
            assert field in data.keys()

        comment_id = data.get(CC.ID)
        submission_id = data.get(CC.SUBMISSION_ID)
        subreddit = data.get(CC.SUBREDDIT)
        author = data.get(CC.AUTHOR)
        created_utc = data.get(CC.CREATED_UTC)
        body = data.get(CC.BODY)
        upvotes = data.get(CC.UPVOTES)

        return cls(
            comment_id,
            submission_id,
            subreddit,
            author,
            created_utc,
            body,
            upvotes,
        )