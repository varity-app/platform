"""
Declare the method for generating a K8S job spec for historical reddit scraping
"""

from kubernetes import client

from util.constants.k8s import Images, JobResources
from util.constants.reddit import Misc
from util.constants import RELEASE


def create_scraper_job_object(
    scraping_mode: str,
    subreddit: str,
    scraping_year=2020,
    scraping_month=1,
    scraping_day=-1,
    deployment="dev",
    namespace="default",
    container_name="scraper-container",
):
    """
    Create a k8s Job Object
    Minimum definition of a job object:
    {'api_version': None, - Str
    'kind': None,     - Str
    'metadata': None, - Metada Object
    'spec': None,     -V1JobSpec
    'status': None}   - V1Job Status
    Docs: https://github.com/kubernetes-client/python/blob/master/kubernetes/docs/V1Job.md
    """

    assert scraping_mode in [Misc.COMMENTS, Misc.SUBMISSIONS]

    # Body is the object Body
    body = client.V1Job(api_version="batch/v1", kind="Job")

    # Body needs Metadata
    # Attention: Each JOB must have a different name!
    name = f"scraper-{scraping_mode}-{subreddit}-{scraping_year}-{scraping_month}-{scraping_day}"
    body.metadata = client.V1ObjectMeta(namespace=namespace, name=name)

    # And a Status
    body.status = client.V1JobStatus()

    # Now we start with the Template...
    template = client.V1PodTemplate()
    template.template = client.V1PodTemplateSpec()

    # Passing Arguments in Env:
    env_list = [
        client.V1EnvVar(name="YEAR", value=str(scraping_year)),
        client.V1EnvVar(name="MONTH", value=str(scraping_month)),
        client.V1EnvVar(name="DAY", value=str(scraping_day)),
        client.V1EnvVar(name="MODE", value=scraping_mode),
        client.V1EnvVar(name="SUBREDDIT", value=subreddit),
        client.V1EnvVar(name="DEPLOYMENT", value=deployment),
    ]

    # Secrets
    secret_list = []

    # Resources
    limits = {
        "memory": JobResources.HISTORICAL_SCRAPER.memory_limit,
        "cpu": JobResources.HISTORICAL_SCRAPER.cpu_limit,
    }
    requests = {
        "memory": JobResources.HISTORICAL_SCRAPER.memory_req,
        "cpu": JobResources.HISTORICAL_SCRAPER.cpu_req,
    }

    resource_requirements = client.V1ResourceRequirements(
        limits=limits, requests=requests
    )

    # Create container
    container = client.V1Container(
        name=container_name,
        image=f"{Images.REPO}/{deployment}/historical-reddit-scraper:{RELEASE}",
        env=env_list,
        volume_mounts=[],
        image_pull_policy="Always",
        resources=resource_requirements,
    )

    # Preemtible node taint tolerations
    tolerations = [
        client.V1Toleration(
            key="cloud.google.com/gke-preemptible",
            value="true",
            effect="NoSchedule",
            operator="Equal",
        )
    ]

    # Set affinity to only schedule on preemtible nodes
    requirement = client.V1NodeSelectorRequirement(
        key="cloud.google.com/gke-preemptible",
        values=["true"],
        operator="In",
    )

    term = client.V1NodeSelectorTerm(match_expressions=[requirement])

    node_selector = client.V1NodeSelector(node_selector_terms=[term])

    affinity = client.V1Affinity(
        node_affinity=client.V1NodeAffinity(
            required_during_scheduling_ignored_during_execution=node_selector
        )
    )

    template.template.spec = client.V1PodSpec(
        containers=[container],
        restart_policy="Never",
        image_pull_secrets=secret_list,
        tolerations=tolerations,
        affinity=affinity,
    )

    # And finaly we can create our V1JobSpec!
    body.spec = client.V1JobSpec(
        ttl_seconds_after_finished=600, template=template.template
    )
    return body
