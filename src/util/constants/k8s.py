"""
Kubernetes API related constants
"""

from . import RELEASE


class Images:
    """Docker images used for k8s jobs"""

    REPO = "cgundlach13"

    HISTORICAL_SCRAPER = f"{REPO}/historical-reddit-scraper:{RELEASE}"


class Resource:
    """Helper class for defining kubernetes resources"""

    cpu_req: str
    cpu_limit: str
    memory_req: str
    memory_limit: str

    def __init__(
        self, cpu_req: str, cpu_limit: str, memory_req: str, memory_limit: str
    ) -> None:
        self.cpu_req = cpu_req
        self.cpu_limit = cpu_limit
        self.memory_req = memory_req
        self.memory_limit = memory_limit


class JobResources:
    """Job resource definitions"""

    HISTORICAL_SCRAPER = Resource("250m", "500m", ".25Gi", ".5Gi")
