# Deploying with Terraform

## Prerequisites
Terraform provisions almost everything you need, but there are a few things that must be done first:
1. Create the Google Cloud project, `varity`
2. Create the secrets in Secret Manager necessary for reddit scraping
```
varity-reddit-client-id
varity-reddit-client-secret
varity-reddit-username
varity-reddit-password
varity-reddit-user-agent
```

This can be done via the command line:
```
gcloud secrets create <secret-id> --project varity
echo -n '<secret>' | gcloud secrets versions add <secret-id> --data-file=-
```


## Login to Terraform Cloud
Make sure you logged into Terraform Cloud, where the remote state is stored.  Do this with `terraform login`

## Deploy with Terraform
1. `cd` into corresponding deployment folder (e.g `tf/dev`)
2. `terraform init`
3. `terraform plan`
4. `terraform apply`

Bonus: Format `.tf` files with `terraform fmt <dir>`

## Use kubectl with a GKE cluster
To generate a kubecfg entry for a cluster hosted on GKE, use the following command:
```
gcloud container clusters get-credentials <cluster-name> --zone <zone>
```