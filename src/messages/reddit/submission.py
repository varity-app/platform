from .. import Message

class SubmissionMessage(Message):
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
    ):
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

    def get(self, value, default=""):
        if value is None:
            return default

        return value

    def serialize(self, as_str=True) -> dict:
        response = dict(
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

        if as_str:
            response = self.to_str(response)

        return response

    @classmethod
    def from_obj(cls, data: dict):
        for field in [SC.ID, SC.SUBREDDIT, SC.CREATED_UTC, SC.AUTHOR]:
            assert field in data.keys()

        submission_id = data.get(SC.ID)
        subreddit = data.get(SC.SUBREDDIT)
        title = data.get(SC.TITLE)
        created_utc = data.get(SC.CREATED_UTC)
        name = data.get(SC.NAME)
        selftext = data.get(SC.SELFTEXT)
        author = data.get(SC.AUTHOR)
        is_original_content = data.get(SC.IS_ORIGINAL_CONTENT)
        is_text = data.get(SC.IS_TEXT)
        nsfw = data.get(SC.NSFW)
        num_comments = data.get(SC.NUM_COMMENTS)
        permalink = data.get(SC.PERMALINK)
        upvotes = data.get(SC.UPVOTES)
        upvote_ratio = data.get(SC.UPVOTE_RATIO)
        url = data.get(SC.URL)

        return cls(
            submission_id,
            subreddit,
            title,
            created_utc,
            name,
            selftext,
            author,
            is_original_content,
            is_text,
            nsfw,
            num_comments,
            permalink,
            upvotes,
            upvote_ratio,
            url,
        )