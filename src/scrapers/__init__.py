from .memory import Memory

class Scraper:
    def __init__(self, limit):
        self.memory = Memory(limit)

    def memory_add(self, data):
        self.memory.add(data)

    def memory_contains(self, data):
        return self.memory.contains(data)

    def memory_reset_cursor(self):
        self.memory.reset_cursor()

    def run(self):
        pass