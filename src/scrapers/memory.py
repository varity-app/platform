class Memory:
    def __init__(self, memory_size):
        self.memory_size = memory_size
        self.store = []
        self.cursor = 0

    def add(self, data):
        if len(self.store) >= self.memory_size:
            self.store.pop()
            self.store.insert(self.cursor, data)
            self.cursor += 1
        else:
            self.store.append(data)

    def reset_cursor(self):
        self.cursor = 0

    def contains(self, data):
        return data in self.store