package mocks

const dummyPublicMarkdownText1 = `+++ 
title = "First Post
date = "2024-08-26" 
description = "Wrap up creating a continuous deployment pipeline by integrating Google Cloud's Artifact Registry and Cloud Run." 
private = false 
headerImage = "https://storage.googleapis.com/timengledev-blog-bucket/static/dist/img/icon_sm.png" 
+++

In the [previous post](https://timengle.dev/posts/view/Continuous++Deployment%3A+Creating+a+GitHub+Workflow%2C+Build+Script%2C+and+Dockerfile), we created our GitHub Actions workflow, build script, and Dockerfile. In this post, we will focus on setting up GCP and enabling the APIs that allow our workflow to automate each step of the deployment process.

Here’s what we’ll cover:

- Setting up a GCP account and project
- Creating and configuring a service account for GitHub workflow access
- Building and pushing Docker images to Artifact Registry
- Preparing the Cloud Run API for automatic container deployment

By the end of this post we will have a fully automated pipeline which will deploy new containers to GCP whenever updates are pushed to our main branch of our repository.

---
## Creating a GCP Account & Project
### Creating an account
Creating a GCP account is pretty straight forward. You just need to follow the [sign up process](https://cloud.google.com/?hl=en). Starting account activates a free year long trial which you should be able to use to finish this series of posts. 
### Creating a project
Once your account is created you will want to [start your first project](https://developers.google.com/workspace/guides/create-project). Give your project a name that will be easy to remember, I will name mine timengledev, and then we will be ready to move on to create a service account.

Before continuing, we will also need to [set up a billing account](https://cloud.google.com/billing/docs/how-to/create-billing-account) to verify ourself with GCP. Don’t worry—you should stay within the free tier if you follow the instructions in this series.

Now that our GCP account is set up we are ready to start enabling the services we need in order to store, host, and deploy Docker images and containers.

---

## Creating an Artifact Registry
We will be using Google Artifact Registry to store our Docker images. Google Artifact Registry is similar to Docker Hub, but it is private and hosted by GCP. Our goal is to create and push a new Docker image to the Artifact Registry whenever a new version of our Go application is pushed to our main branch.

### Enabling the Artifact Registry
1. Head to the Artifact Registry using the search input in the GCP console and enable the API. 
2. Create a new repository. I used the following settings
	- Name: app
	- Format: Docker
	- Mode: Standard
	- Location Type: Region
	- Region For Deployment: us-west1
	- Left 'Google-managed encryption key' checked
	
Great now our Artifact Registry is set up for our project and we can continue on to create a Service Account so we can automate this build process in our GitHub Action workflow.

---
## Creating a Service Account
To gain access to our GCP project we will want to set up and create a key for a Service Account. We can follow these instructions to create a JSON key which we can use as a GitHub Secret to gain access to our GCP account in our workflow.

1. Head to the [IAM & Service Accounts](https://console.cloud.google.com/iam-admin/serviceaccounts) and select the current project you are working with.
2. At the top of the console click 'Create Service Account'.
3. Create a Service Account (I named mine timengleblog_cloud) and enable these permissions.
	- Cloud Build Editor
	- Cloud Build Service Account
	- Cloud Run Admin
	- Service Account User
	- Viewer
4. Click the actions ellipsis next to the account you just created
5. Click Manage keys
6.  You will be taken to a new page where you can click "ADD KEY" and select "Create new key" and select JSON
7. The key will download

Great we now have the credentials to begin automating our CD pipeline within our workflow.

---
## Accessing GCP Within our Workflow
We can begin to access our GCP account in our workflow by adding a new step that uses this [google-github-actions](https://github.com/google-github-actions/setup-gcloud) which will set up and configure, gcloud the cli tool used to access and make changes to our project via our service worker.
### Adding our JSON key
We will head to the GitHub repository for our project and navigate to 
Settings -> Secrets and Variables -> Actions
Then copy the contents of the JSON file which contains the GCP service account key.
We will now create a new secret called GCP_CREDENTIALS by clicking 'New repository secret'.

Great now that we have added our credentials of our service worker we are ready to access GCP within our workflow.

### Adding a step to authenticate GCP
To allow our workflow to interact with the Google Cloud we need to authenticate the workflow using the secret we added to our repository in the previous step. This will be done using setup-cloud action which will configure the gcloud CLI in our workflow. We will add this step following our step that runs our script responsible for building our application binary.
`

const dummyPrivateMarkdownText1 = `+++ 
title = "Continuous Deployment: Automating GCP Deployments with Artifact Registry and Cloud Run" 
date = "2024-08-26" 
description = "Wrap up creating a continuous deployment pipeline by integrating Google Cloud's Artifact Registry and Cloud Run." 
private = true 
headerImage = "https://storage.googleapis.com/timengledev-blog-bucket/static/dist/img/icon_sm.png" 
+++

In the [previous post](https://timengle.dev/posts/view/Continuous++Deployment%3A+Creating+a+GitHub+Workflow%2C+Build+Script%2C+and+Dockerfile), we created our GitHub Actions workflow, build script, and Dockerfile. In this post, we will focus on setting up GCP and enabling the APIs that allow our workflow to automate each step of the deployment process.

Here’s what we’ll cover:

- Setting up a GCP account and project
- Creating and configuring a service account for GitHub workflow access
- Building and pushing Docker images to Artifact Registry
- Preparing the Cloud Run API for automatic container deployment

By the end of this post we will have a fully automated pipeline which will deploy new containers to GCP whenever updates are pushed to our main branch of our repository.

---
## Creating a GCP Account & Project
### Creating an account
Creating a GCP account is pretty straight forward. You just need to follow the [sign up process](https://cloud.google.com/?hl=en). Starting account activates a free year long trial which you should be able to use to finish this series of posts. 
### Creating a project
Once your account is created you will want to [start your first project](https://developers.google.com/workspace/guides/create-project). Give your project a name that will be easy to remember, I will name mine timengledev, and then we will be ready to move on to create a service account.

Before continuing, we will also need to [set up a billing account](https://cloud.google.com/billing/docs/how-to/create-billing-account) to verify ourself with GCP. Don’t worry—you should stay within the free tier if you follow the instructions in this series.

Now that our GCP account is set up we are ready to start enabling the services we need in order to store, host, and deploy Docker images and containers.

---

## Creating an Artifact Registry
We will be using Google Artifact Registry to store our Docker images. Google Artifact Registry is similar to Docker Hub, but it is private and hosted by GCP. Our goal is to create and push a new Docker image to the Artifact Registry whenever a new version of our Go application is pushed to our main branch.

### Enabling the Artifact Registry
1. Head to the Artifact Registry using the search input in the GCP console and enable the API. 
2. Create a new repository. I used the following settings
	- Name: app
	- Format: Docker
	- Mode: Standard
	- Location Type: Region
	- Region For Deployment: us-west1
	- Left 'Google-managed encryption key' checked
	
Great now our Artifact Registry is set up for our project and we can continue on to create a Service Account so we can automate this build process in our GitHub Action workflow.

---
## Creating a Service Account
To gain access to our GCP project we will want to set up and create a key for a Service Account. We can follow these instructions to create a JSON key which we can use as a GitHub Secret to gain access to our GCP account in our workflow.

1. Head to the [IAM & Service Accounts](https://console.cloud.google.com/iam-admin/serviceaccounts) and select the current project you are working with.
2. At the top of the console click 'Create Service Account'.
3. Create a Service Account (I named mine timengleblog_cloud) and enable these permissions.
	- Cloud Build Editor
	- Cloud Build Service Account
	- Cloud Run Admin
	- Service Account User
	- Viewer
4. Click the actions ellipsis next to the account you just created
5. Click Manage keys
6.  You will be taken to a new page where you can click "ADD KEY" and select "Create new key" and select JSON
7. The key will download

Great we now have the credentials to begin automating our CD pipeline within our workflow.

---
## Accessing GCP Within our Workflow
We can begin to access our GCP account in our workflow by adding a new step that uses this [google-github-actions](https://github.com/google-github-actions/setup-gcloud) which will set up and configure, gcloud the cli tool used to access and make changes to our project via our service worker.
### Adding our JSON key
We will head to the GitHub repository for our project and navigate to 
Settings -> Secrets and Variables -> Actions
Then copy the contents of the JSON file which contains the GCP service account key.
We will now create a new secret called GCP_CREDENTIALS by clicking 'New repository secret'.

Great now that we have added our credentials of our service worker we are ready to access GCP within our workflow.

### Adding a step to authenticate GCP
To allow our workflow to interact with the Google Cloud we need to authenticate the workflow using the secret we added to our repository in the previous step. This will be done using setup-cloud action which will configure the gcloud CLI in our workflow. We will add this step following our step that runs our script responsible for building our application binary.
`

const dummyContentText1 = `In the [previous post](https://timengle.dev/posts/view/Continuous++Deployment%3A+Creating+a+GitHub+Workflow%2C+Build+Script%2C+and+Dockerfile), we created our GitHub Actions workflow, build script, and Dockerfile. In this post, we will focus on setting up GCP and enabling the APIs that allow our workflow to automate each step of the deployment process.

Here’s what we’ll cover:

- Setting up a GCP account and project
- Creating and configuring a service account for GitHub workflow access
- Building and pushing Docker images to Artifact Registry
- Preparing the Cloud Run API for automatic container deployment

By the end of this post we will have a fully automated pipeline which will deploy new containers to GCP whenever updates are pushed to our main branch of our repository.

---
## Creating a GCP Account & Project
### Creating an account
Creating a GCP account is pretty straight forward. You just need to follow the [sign up process](https://cloud.google.com/?hl=en). Starting account activates a free year long trial which you should be able to use to finish this series of posts. 
### Creating a project
Once your account is created you will want to [start your first project](https://developers.google.com/workspace/guides/create-project). Give your project a name that will be easy to remember, I will name mine timengledev, and then we will be ready to move on to create a service account.

Before continuing, we will also need to [set up a billing account](https://cloud.google.com/billing/docs/how-to/create-billing-account) to verify ourself with GCP. Don’t worry—you should stay within the free tier if you follow the instructions in this series.

Now that our GCP account is set up we are ready to start enabling the services we need in order to store, host, and deploy Docker images and containers.

---

## Creating an Artifact Registry
We will be using Google Artifact Registry to store our Docker images. Google Artifact Registry is similar to Docker Hub, but it is private and hosted by GCP. Our goal is to create and push a new Docker image to the Artifact Registry whenever a new version of our Go application is pushed to our main branch.

### Enabling the Artifact Registry
1. Head to the Artifact Registry using the search input in the GCP console and enable the API. 
2. Create a new repository. I used the following settings
	- Name: app
	- Format: Docker
	- Mode: Standard
	- Location Type: Region
	- Region For Deployment: us-west1
	- Left 'Google-managed encryption key' checked
	
Great now our Artifact Registry is set up for our project and we can continue on to create a Service Account so we can automate this build process in our GitHub Action workflow.

---
## Creating a Service Account
To gain access to our GCP project we will want to set up and create a key for a Service Account. We can follow these instructions to create a JSON key which we can use as a GitHub Secret to gain access to our GCP account in our workflow.

1. Head to the [IAM & Service Accounts](https://console.cloud.google.com/iam-admin/serviceaccounts) and select the current project you are working with.
2. At the top of the console click 'Create Service Account'.
3. Create a Service Account (I named mine timengleblog_cloud) and enable these permissions.
	- Cloud Build Editor
	- Cloud Build Service Account
	- Cloud Run Admin
	- Service Account User
	- Viewer
4. Click the actions ellipsis next to the account you just created
5. Click Manage keys
6.  You will be taken to a new page where you can click "ADD KEY" and select "Create new key" and select JSON
7. The key will download

Great we now have the credentials to begin automating our CD pipeline within our workflow.

---
## Accessing GCP Within our Workflow
We can begin to access our GCP account in our workflow by adding a new step that uses this [google-github-actions](https://github.com/google-github-actions/setup-gcloud) which will set up and configure, gcloud the cli tool used to access and make changes to our project via our service worker.
### Adding our JSON key
We will head to the GitHub repository for our project and navigate to 
Settings -> Secrets and Variables -> Actions
Then copy the contents of the JSON file which contains the GCP service account key.
We will now create a new secret called GCP_CREDENTIALS by clicking 'New repository secret'.

Great now that we have added our credentials of our service worker we are ready to access GCP within our workflow.

### Adding a step to authenticate GCP
To allow our workflow to interact with the Google Cloud we need to authenticate the workflow using the secret we added to our repository in the previous step. This will be done using setup-cloud action which will configure the gcloud CLI in our workflow. We will add this step following our step that runs our script responsible for building our application binary.
`

const dummyPublicMarkdownText2 = `+++
title = "Second Post" 
date = "2024-08-27" 
description = "Another Post" 
private = false 
headerImage = "https://storage.googleapis.com/timengledev-blog-bucket/static/dist/img/icon_sm.png" 
+++

# Header 1

Some Content
`

const dummyPrivateMarkdownText2 = `+++
title = "Second Post" 
date = "2024-08-27" 
description = "Another Post" 
private = true 
headerImage = "https://storage.googleapis.com/timengledev-blog-bucket/static/dist/img/icon_sm.png" 
+++

# Header 1

Some Content
`

const dummyContentText2 = `# Header 1

Some Content
`
