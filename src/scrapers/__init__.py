from .memory import Memory

class Scraper:
    def __init__(self, limit):
        self.memory = Memory(limit)

    def memory_add(self, data):
        self.memory.add(data)

    def memory_contains(self, data):
        return self.memory.contains(data)

    def run(self):
        pass