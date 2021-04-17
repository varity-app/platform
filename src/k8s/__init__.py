"""
Helper module for utilizing the python kubernetes client
"""

import random
import string
import logging
from kubernetes import client, config
from kubernetes.client.rest import ApiException

logger = logging.getLogger(__name__)


def init_k8s():
    """Configure and initialize kubernetes.  Return the API instance"""

    config.load_kube_config()
    api_instance = client.BatchV1Api()

    return api_instance


def create_job(api_instance, body, namespace="default"):
    """Using an existing body, send it to the API and create the job"""

    try:
        api_instance.create_namespaced_job(namespace, body, pretty=True)
        logger.info("Created job...")
    except ApiException as ex:
        logger.exception(
            f"Exception when calling BatchV1Api -> create_namespaced_job: {ex}\n"
        )


def get_secret_env_var(secret_name, key, env_name):
    """Create the k8s environment variable reference and objects for a given secret"""

    secret_ref = client.V1SecretKeySelector(key=key, name=secret_name)
    secret_source = client.V1EnvVarSource(secret_key_ref=secret_ref)
    secret_env = client.V1EnvVar(name=env_name, value_from=secret_source)

    return secret_ref, secret_env


def id_generator(size=4, chars=string.ascii_lowercase + string.digits):
    """Generate a unique ID string"""

    return "".join(random.choice(chars) for _ in range(size))


def delete_pods(api_pods, namespace="default", phase="Succeeded"):
    """
    Pods are never empty, just completed the lifecycle.
    As such they can be deleted.
    Pods can be without any running container in 2 states:
    Succeeded and Failed. This call doesn't terminate Failed pods by default.

    Possible phases: ['Failed', 'Succeeded', 'Pending']
    """

    # We need the api entry point for pods
    api_pods = client.CoreV1Api()

    # List the pods
    try:
        pods = api_pods.list_namespaced_pod(namespace, pretty=True, timeout_seconds=60)
    except ApiException as ex:
        logger.error(f"Exception when calling CoreV1Api -> list_namespaced_pod: {ex}\n")

    for pod in pods.items:
        logger.debug(pod)
        pod_name = pod.metadata.name
        try:
            if pod.status.phase == phase:
                api_response = api_pods.delete_namespaced_pod(pod_name, namespace)
                logger.info(f"Pod: {pod_name} deleted!")
                logger.debug(api_response)
            else:
                logger.info(
                    f"Pod: {pod_name} still not done... Phase: {pod.status.phase}"
                )
        except ApiException as ex:
            logger.exception(
                f"Exception when calling CoreV1Api -> delete_namespaced_pod: {ex}\n"
            )


def delete_jobs(api_instance, namespace="default"):
    """
    Since the TTL flag (ttl_seconds_after_finished) is still in alpha (Kubernetes 1.12)
    jobs need to be cleanup manually.

    As such this method checks for existing Finished Jobs and deletes them.
    By default it only cleans Finished jobs. Failed jobs require manual intervention
    or a second call to this function.
    Docs: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion
    For deletion you need a new object type! V1DeleteOptions! But you can have it empty!
    CAUTION: Pods are not deleted at the moment. They are set to not running, but will
            count for your autoscaling limit, so if
             pods are not deleted, the cluster can hit the autoscaling limit even with
             free, idling pods.
             To delete pods, at this moment the best choice is to use the kubectl tool
             ex: kubectl delete jobs/JOBNAME.
             But! If you already deleted the job via this API call, you now need to delete
             the Pod using Kubectl:
             ex: kubectl delete pods/PODNAME

    Possible phases: ['Failed', 'Finished', 'Pending']
    """
    delete_options = client.V1DeleteOptions(
        grace_period_seconds=0, propagation_policy="Background"
    )
    try:
        jobs = api_instance.list_namespaced_job(
            namespace,
            pretty=True,
            timeout_seconds=60,
        )
    except ApiException as ex:
        logger.exception(
            f"Exception when calling BatchV1Api -> list_namespaced_job: {ex}\n"
        )

    # Now we have all the jobs, lets clean up
    # We are also logging the jobs we didn't clean up because they either failed
    # or are still running
    for job in jobs.items:
        logging.debug(job)
        job_name = job.metadata.name
        job_status = job.status.conditions

        if job.status.succeeded == 1:

            # Clean up Job
            logging.info(
                f"Cleaning up Job: {job_name}. Finished at: {job.status.completion_time}"
            )

            try:
                # What is at work here. Setting Grace Period to 0 means delete ASAP.
                # Otherwise it defaults to
                # some value I can't find anywhere. Propagation policy makes the
                # Garbage cleaning Async
                api_response = api_instance.delete_namespaced_job(
                    name=job_name,
                    namespace=namespace,
                    body=delete_options,
                )
                logging.debug(api_response)

            except ApiException as ex:
                logger.exception(
                    f"Exception when calling BatchV1Api -> delete_namespaced_job: {ex}\n"
                )
        else:
            if job_status is None and job.status.active == 1:
                job_status = "active"
            logging.info(
                f"Job: {job_name} not cleaned up. Current status: {job_status}"
            )

    # Now that we have the jobs cleaned, let's clean the pods
    delete_pods(namespace)
