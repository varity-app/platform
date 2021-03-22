# Launch a kubernetes cluster on AWS with kops

### Install KOPS cli
```
curl -LO https://github.com/kubernetes/kops/releases/download/$(curl -s https://api.github.com/repos/kubernetes/kops/releases/latest | grep tag_name | cut -d '"' -f 4)/kops-linux-amd64
chmod +x kops-linux-amd64
sudo mv kops-linux-amd64 /usr/local/bin/kops
```

## Setup cluster on AWS


### Set kops AWS user policies:
```
AmazonEC2FullAccess
AmazonRoute53FullAccess
AmazonS3FullAccess
IAMFullAccess
AmazonVPCFullAccess
```

### Export env vars for kops
```
export AWS_ACCESS_KEY_ID=$(aws configure get aws_access_key_id)
export AWS_SECRET_ACCESS_KEY=$(aws configure get aws_secret_access_key)
```

### Create s3 bucket to store cluster state
```
aws s3api create-bucket \
    --bucket varity-kops-state \
    --region us-east-2 \
    --create-bucket-configuration LocationConstraint=us-east-2
```

### Create Cluster
```	
export NAME=varity-dev.k8s.local
export KOPS_STATE_STORE=s3://varity-kops-state
kops create cluster --zones=us-east-2b \
    --master-size="t2.micro" \
    --node-size="t3a.small" \
    --master-count=1 \
    --node-count=3 \
    --ssh-public-key ~/.ssh/id_rsa.pub \
    ${NAME}

kops update cluster --name ${NAME} --yes
```

### Verify Cluster
Note: may take a few minutes to create
```
kops validate cluster
```

## Enable Autoscaling
Help comes from https://varlogdiego.com/kubernetes-cluster-with-autoscaling-on-aws-and-kops

### Edit node instances group
Run `kops edit instancegroups nodes`
Edit config to look like:
```
kind: InstanceGroup
...
spec:
  cloudLabels:
    service: k8s_node
    k8s.io/cluster-autoscaler/enabled: ""
    k8s.io/cluster-autoscaler/<YOUR CLUSTER NAME>: ""
  maxSize: 10
  minSize: 2
...
```

### Edit Cluster
Run `kops edit cluster`
Add following policies to nodes:
```
kind: Cluster
...
spec:
  additionalPolicies:
    node: |
      [
        {
          "Effect": "Allow",
          "Action": [
            "autoscaling:DescribeAutoScalingGroups",
            "autoscaling:DescribeAutoScalingInstances",
            "autoscaling:DescribeLaunchConfigurations",
            "autoscaling:SetDesiredCapacity",
            "autoscaling:TerminateInstanceInAutoScalingGroup",
            "autoscaling:DescribeTags"
          ],
          "Resource": "*"
        }
      ]
...
```

### Apply changes
```
kops update cluster
kops rolling-update cluster
```

### Deploy Kubernetes Autoscaler
Apply the yaml descriptor:
```
kubectl apply -f ${REPO_BASE}/res/k8s/cluster-autoscaler-autodiscover.yaml
````

### Test deployment
```
kubectl create deployment webserver-demo --image=nginx:latest
kubectl scale deployment webserver-demo --replicas=20

# Check the pods if are in "Pending" status
kubectl get pods

# Check the amount of nodes, will take a few minutes until the new node join to the cluster
kubectl get nodes
```