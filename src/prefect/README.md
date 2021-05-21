# ETL with Prefect

Prefect was chosen over Airflow primarily due to the benefits offered from Prefect Cloud.  The hybrid compute model is extremely cost effective, as well as a reliable alternative to hosting the UI in-house.

## Installation
The UI can be accessed at https://cloud.prefect.io, but the agent will need to be installed to the GKE cluster provisioned by Terraform.

### Install CLI and get token
1. Install the CLI with `pip`:
```bash
pip install prefect
```
2. In the UI, create an API Key.  Copy the key, you will need this for the next step

### Create the agent
Once you can interact with the cluster with `kubectl`, the Prefect agent installation should be simple.
```bash
# Create a namespace specifically for prefect
kubectl create ns prefect

# Generate the yaml file and pipe it to kubectl
prefect agent kubernetes install -t <token> | kubectl apply --namespace=prefect -f -\n
```

You should see the agent running as a deployment in the `prefect` namespace.
```bash
kubectl get deployments -n prefect
```

Check in the UI and see if the agent is running.

