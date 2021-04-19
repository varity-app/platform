# Setting up GitHub Actions

## Create the GCP Service Account
First, we'll need to setup a service account within Google Cloud.  The service account should have the following roles:
```
roles/storage.admin
roles/secretmanager.admin
roles/pubsub.admin
roles/container.admin
roles/bigquery.admin
```

## Create the GitHub Secrets
Create the following GitHub secrets:

### Google Cloud
*`GOOGLE_APPLICATION_CREDENTIALS`: base64 encoded service account key created in the previous step. You can encode the json file via:
```
base64 <google_key>.json
```
### Reddit API
* `REDDIT_USERNAME`
* `REDDIT_PASSWORD`
* `REDDIT_CLIENT_ID`
* `REDDIT_CLIENT_SECRET`
* `REDDIT_USER_AGENT`

## Ready!
GitHub Actions should be ready to go now.
