# Deploying with Terraform

## Prerequisites
Terraform provisions almost everything you need, but there are a few things that are done directly through AWS:
1. Create Confluent Cloud cluster and relevant topics.  See documentation [here](../docs/confluent/create_kafka_topics) for configuring topics.
2. Create the necessary secrets in [AWS System Manager Parameter Store](https://aws.amazon.com/systems-manager/).  Secrets used by ECS Containers need to be provisioned by hand.  Here is a list of necessary secrets:
```
varity-confluent-bootstrap-prod
varity-confluent-key-prod
varity-confluent-secret-prod
varity-confluent-bootstrap-dev
varity-confluent-key-dev
varity-confluent-secret-dev
varity-reddit-client-id
varity-reddit-client-secret
varity-reddit-username
varity-reddit-password
varity-reddit-user-agent
```

3. Create IAM Policy for accessing secrets created in previous step.  The ARN of this policy is recorded in terraform as a variable named `secrets_policy_arn`.  The policy can look like this:
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": "kms:Decrypt",
            "Resource": "arn:aws:kms:us-east-2:178852309825:key/017ef22d-b1b9-4528-ad1c-8c51d161bbc0"
        },
        {
            "Sid": "VisualEditor1",
            "Effect": "Allow",
            "Action": [
                "ssm:DescribeParameters"
            ],
            "Resource": "*"
        },
        {
            "Sid": "VisualEditor2",
            "Effect": "Allow",
            "Action": [
                "ssm:GetParameter",
                "ssm:GetParameters"
            ],
            "Resource": "arn:aws:ssm:us-east-2:178852309825:parameter/*"
        }
    ]
}
```

## Login to Terraform Cloud
Make sure you logged into Terraform Cloud, where the remote state is stored.  Do this with `terraform login`

## Deploy with Terraform
1. `cd` into corresponding deployment folder (e.g `tf/dev`)
2. `terraform init`
3. `terraform plan`
4. `terraform apply`

Bonus: Format `.tf` files with `terraform fmt <dir>`