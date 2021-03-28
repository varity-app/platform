"""
Declare the base scraper class
"""

from .memory import Memory


class Scraper:
    """Base class that will be inherited by every other scraper class"""

    def __init__(self, limit):
        self.memory = Memory(limit)

    def memory_add(self, data):
        """Add item to memory"""
        self.memory.add(data)

    def memory_contains(self, data):
        """Return whether an item is in memory"""
        return self.memory.contains(data)

    def memory_reset_cursor(self):
        """Reset the cursor to the beginning"""
        self.memory.reset_cursor()

    async def run(self):
        """Placeholder that should be implemented by every child class"""
        return
