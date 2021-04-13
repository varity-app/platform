"""
Declare the Memory helper class
"""

from datetime import datetime

from google.cloud import firestore
from google.cloud.firestore_v1.document import DocumentReference

from util.constants.firestore import DISABLE_FIRESTORE, PROJECT


class Memory:
    """Helper class that helps a reddit scraper determine if a post has been already scraped"""

    def __init__(
        self, memory_size: int, enable_firestore=True, collection=None
    ) -> None:
        self.memory_size = memory_size
        self.store = []
        self.cursor = 0
        self.enable_firestore = enable_firestore and not DISABLE_FIRESTORE
        self.collection = collection

        if self.enable_firestore:
            assert collection is not None

        if self.enable_firestore:
            self.firestore = firestore.Client(project=PROJECT)

    def get_doc_ref(self, data: str) -> DocumentReference:
        """Retrieve a firestore document reference"""
        return self.firestore.collection(self.collection).document(data)

    def add(self, data: str) -> None:
        """Add an item to memory"""
        if len(self.store) >= self.memory_size:
            self.store.pop()
            self.store.insert(self.cursor, data)
            self.cursor += 1
        else:
            self.store.append(data)

        # Save to firestore
        if self.enable_firestore:
            doc_ref = self.get_doc_ref(data)
            doc_ref.set(dict(tiemstamp=datetime.now().isoformat()))

    def reset_cursor(self) -> None:
        """Reset the memory cursor to the beginning"""
        self.cursor = 0

    def contains(self, data: str) -> bool:
        """Return whether an item is in memory"""
        if data in self.store:
            return True

        # Check for key in firestore
        if self.enable_firestore:
            doc_ref = self.get_doc_ref(data)
            return doc_ref.get().exists

        return False
