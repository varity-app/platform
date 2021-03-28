"""
Declare Kafka message class for reddit submissions
"""

from util.constants.reddit import SubmissionConstants as SC

from .. import Message


class SubmissionMessage(Message):
    """Kafka message repesenting a reddit submission"""

    def __init__(
        self,
        submission_id: str,
        subreddit: str,
        title: str,
        created_utc: str,
        name: str,
        selftext: str,
        author: str,
        is_original_content: bool,
        is_text: bool,
        nsfw: bool,
        num_comments: int,
        permalink: str,
        upvotes: int,
        upvote_ratio: float,
        url: str,
    ) -> None:
        self.submission_id = submission_id
        self.subreddit = subreddit
        self.title = self.get(title)
        self.created_utc = created_utc
        self.name = self.get(name)
        self.selftext = self.get(selftext)
        self.author = author
        self.is_original_content = self.get(is_original_content)
        self.is_text = self.get(is_text)
        self.nsfw = self.get(nsfw)
        self.num_comments = self.get(num_comments)
        self.permalink = self.get(permalink)
        self.upvotes = self.get(upvotes)
        self.upvote_ratio = self.get(upvote_ratio)
        self.url = self.get(url)

    def to_obj(self) -> dict:
        """Serialize the class to a dict"""
        obj = dict(
            submission_id=self.submission_id,
            subreddit=self.subreddit,
            title=self.title,
            created_utc=self.created_utc,
            name=self.name,
            selftext=self.selftext,
            author=self.author,
            is_original_content=self.is_original_content,
            is_text=self.is_text,
            nsfw=self.nsfw,
            num_comments=self.num_comments,
            permalink=self.permalink,
            upvotes=self.upvotes,
            upvote_ratio=self.upvote_ratio,
            url=self.url,
        )

        return obj

    @classmethod
    def from_obj(cls, data: dict):
        """Create a new instance of the class from a dict"""
        cls.assert_has_fields(data, [SC.ID, SC.SUBREDDIT, SC.CREATED_UTC, SC.AUTHOR])

        return cls(
            data.get(SC.ID),
            data.get(SC.SUBREDDIT),
            data.get(SC.TITLE),
            data.get(SC.CREATED_UTC),
            data.get(SC.NAME),
            data.get(SC.SELFTEXT),
            data.get(SC.AUTHOR),
            data.get(SC.IS_ORIGINAL_CONTENT),
            data.get(SC.IS_TEXT),
            data.get(SC.NSFW),
            data.get(SC.NUM_COMMENTS),
            data.get(SC.PERMALINK),
            data.get(SC.UPVOTES),
            data.get(SC.UPVOTE_RATIO),
            data.get(SC.URL),
        )
