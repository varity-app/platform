class Memory:
    def __init__(self, memory_size):
        self.memory_size = memory_size
        self.store = []

    def add(self, data):
        if len(self.store) >= self.memory_size:
            self.store.pop(0)
        
        self.store.append(data)

    def contains(self, data):
        return data in self.store