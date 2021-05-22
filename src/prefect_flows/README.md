# ETL with Prefect

Prefect was chosen over Airflow primarily due to the benefits offered from Prefect Cloud.  The hybrid compute model is extremely cost effective, as well as a reliable alternative to hosting the UI in-house.

## Installation
The UI can be accessed at https://cloud.prefect.io, but the agent will need to be installed to the GKE cluster provisioned by Terraform.

### GCS
Prefect flows will be uploaded to a GCS bucket, `varity-prefect-<deployment>`, provisioned by Terraform.

### Python Version
A flow must be built with the same python version as the runner, which is `3.8.X`.  Be sure the local python version used to build and upload the flows is correct.

### Install CLI and get token
1. Install the CLI with `pip`:
```bash
pip install prefect
```
2. In the UI, create an API Key.  Copy the key, you will need this for the next step.
3. Create the file `~/.prefect/config.toml` and write the following contents:
```
[cloud]
use_local_secrets = false
auth_token = "<token>"
```
3. Set backend by running the following command:
```bash
prefect backend cloud
```

### Create the kubernetes agent
Once you can interact with the cluster with `kubectl`, the Prefect agent installation should be simple.
```bash
# Create a namespace specifically for prefect
kubectl create ns prefect

# Generate the yaml file and pipe it to kubectl
prefect agent kubernetes install -t <token> --namespace prefect --rbac | kubectl apply -n prefect -f -
```

You should see the agent running as a deployment in the `prefect` namespace.
```bash
kubectl get deployments -n prefect
```

Check in the UI and see if the agent is running.

### Secrets
Certain tasks require credentials for tasks such as authenticating to a PostgreSQL.  These credentials are stored as secrets in Prefect Cloud.

#### Required Secrets
These are the secrets required by various tasks that will need to be created before tasks can be initiated:
* `IEX_TOKEN` Token for authenticating to the IEXCloud financial data API.
