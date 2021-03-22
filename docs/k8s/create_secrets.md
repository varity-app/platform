# Create Kubernetes Secrets

```
k create secret generic reddit-creds \
    --from-literal=client_id='<reddit_client_id>' \
    --from-literal=client_secret='<reddit_client_secret>' \
    --from-literal=user_agent='<reddit_user_agent>' \
    --from-literal=username='<reddit_username>' \
    --from-literal=password='<reddit_password>'

k create secret generic confluent-cloud-creds \
    --from-literal=client_id='<reddit_client_id>' \
    --from-literal=client_secret='<reddit_client_secret>'

```