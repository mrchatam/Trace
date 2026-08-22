import os
from pathlib import Path
from collections import defaultdict as dd

def helper():
    return 1

class Worker:
    def run(self):
        return helper()

def main():
    w = Worker()
    w.run()
