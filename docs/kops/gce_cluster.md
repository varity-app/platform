# Launch an autoscaling kubernetes cluster on GCE with kops

### Install KOPS cli
```
curl -LO https://github.com/kubernetes/kops/releases/download/$(curl -s https://api.github.com/repos/kubernetes/kops/releases/latest | grep tag_name | cut -d '"' -f 4)/kops-linux-amd64
chmod +x kops-linux-amd64
sudo mv kops-linux-amd64 /usr/local/bin/kops
```

## Setup cluster on GCE

### Create store
```
export BUCKET_NAME=varity-kops-clusters
gsutil mb gs://${BUCKET_NAME}/
export KOPS_STATE_STORE=gs://${BUCKET_NAME}/
```

### Create cluster
```
export NAME=varity-dev.k8s.local
PROJECT=`gcloud config get-value project`
export KOPS_FEATURE_FLAGS=AlphaAllowGCE # to unlock the GCE features
kops create cluster ${NAME} --zones us-central1-a --state ${KOPS_STATE_STORE}/ --project=${PROJECT} --master-size="e2-small" --node-size="e2-standard-4" --node-count=1 --master-count=1
```