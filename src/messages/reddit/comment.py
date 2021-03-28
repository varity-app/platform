"""
Declare Kafka message class for reddit comments
"""

from util.constants.reddit import CommentConstants as CC

from .. import Message


class CommentMessage(Message):
    """Kafka message repesenting a reddit comment"""

    def __init__(
        self,
        comment_id: str,
        submission_id: str,
        subreddit: str,
        author: str,
        created_utc: str,
        body: str,
        upvotes: str,
    ) -> None:
        self.comment_id = comment_id
        self.submission_id = submission_id
        self.subreddit = subreddit
        self.author = author
        self.created_utc = created_utc
        self.body = self.get(body)
        self.upvotes = upvotes

    def to_obj(self) -> dict:
        """Serialize the class to a dict"""
        obj = dict(
            comment_id=self.comment_id,
            submission_id=self.submission_id,
            subreddit=self.subreddit,
            author=self.author,
            created_utc=self.created_utc,
            body=self.body,
            upvotes=self.upvotes,
        )

        return obj

    @classmethod
    def from_obj(cls: CommentMessage, data: dict) -> CommentMessage:
        """Create a new instance of the class from a dict"""
        cls.assert_has_fields(
            data, [CC.ID, CC.SUBMISSION_ID, CC.SUBREDDIT, CC.AUTHOR, CC.CREATED_UTC]
        )

        return cls(
            data.get(CC.ID),
            data.get(CC.SUBMISSION_ID),
            data.get(CC.SUBREDDIT),
            data.get(CC.AUTHOR),
            data.get(CC.CREATED_UTC),
            data.get(CC.BODY),
            data.get(CC.UPVOTES),
        )
