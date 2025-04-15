<!--suppress HtmlDeprecatedAttribute -->
<p align="right">
    <a href="https://github.com/tommed" title="See Project Ducto">
        <img src="../assets/ducto-logo-small.png" alt="A part of Project Ducto"/>
    </a>
</p>

# Google Cloud Platform: Cloud Run

<p align="center">
    <img src="../assets/ducto-orchestrator-gcp-logo-small.jpg" 
         alt="Ducto Orchestrator for Google Cloud Platform" />
</p>

You can run Ducto Orchestrator as a Cloud Run service using our public Docker repo.

Your config and program are stored in a GCS private bucket, the service can access.
Given Cloud Run is based on a web endpoint, you need to configure your Event Source 
a certain way for it to work properly.

```yaml
program_file: gs://ducto-public/02-gcs_program.json
# Or, we can define the program inline for convenience
#program:
#  version: 1
#  instructions:
#    - op: set
#      ...

source:
  type: http
  config:
    meta_field: _http  # If you want it
    use_env: true     # Awareness of Cloud Run's `PORT` env for addr binding

#output:
#  ...
```

## Security

We strongly recommend you publish your service with IAM authentication left enabled 
(i.e., **do not** use the `--allow-unauthenticated` parameter when 
calling `gcloud run deploy`), otherwise anyone with the URL of your service 
can call it and send events.

---
## 🛣️ Roadmap

- [x] Use the Cloud Run runtime
- [ ] Host our **Feature Flags** service on Cloud Run and have orchestrator call it 
- [ ] Support static Token-based authentication
- [ ] Support OAuth2 & OIDC authentication

---
## Repo

Our Ducto repo for Ducto Orchestrator for Cloud Run is:

```
europe-west2-docker.pkg.dev/testsandbox-3cb3f/ducto-docker/ducto-orchestrator-cloudrun:[version|latest]
```

---
## Permissions Needed
Your Cloud Run service should have its own identity, which needs the following permissions:

- Service Account User
- Legacy Storage Object Reader (on the bucket hosting your config and program)

---
## Deploy

```bash
# Your GCP Project
YOUR_PROJECT=your-project-id

# Service Account to Run As
SERVICE_ACCOUNT_EMAIL=ducto-orchestrator@${YOUR_PROJECT}.iam.gserviceaccount.com

# Note: the service account must have read access to this file, and your program json if stored on gcs
YOUR_CONFIG_URI=gs://ducto-public/02-gcs_config.yaml

gcloud run deploy ducto-cloudrun \
  --image=europe-west2-docker.pkg.dev/testsandbox-3cb3f/ducto-docker/ducto-orchestrator-cloudrun:latest \
  --region=europe-west2 \
  --set-env-vars CONFIG_URI=${YOUR_CONFIG_URI} \
  --service-account=${SERVICE_ACCOUNT_EMAIL} \
  --platform=managed \
  --project=${YOUR_PROJECT}
```
