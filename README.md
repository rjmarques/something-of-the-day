# Something of the day

[![CircleCI](https://circleci.com/gh/rjmarques/something-of-the-day/tree/master.svg?style=svg)](https://circleci.com/gh/rjmarques/something-of-the-day/tree/master)

https://something.ricardomarq.com

This document explains how to build and deploy the **something-of-the-day App**. The app is a completely over-engineered React Single Page Web App that has a Golang backend and connects to a PG database. The app runs locally on Docker, is deployed to AWS ECS and links to an Heroku database. There's also a CI workflow set up on CircleCI, that automatically deploys the app when code is merged to master and passes all the tests.

The goal was to learn a bit about AWS ECS, Terraform, Docker, as well as to set up a relatively reaslistic CI/CD.

##### DISCLAIMER: It goes without saying that a trivial static webpage does not need: React, Docker, Postgresql, etc. The complexity was introduced to allow me to build realistic CI workflow. I.e., have good frontend build & tests, backend build & test, and DB integration tests. Plus, I wanted to push everything to containers as much as possible. Oh and Heroku provides a tiny database for free so I used that - win! 💰

## Minimal requirements to build and run

You'll need the following accounts

- AWS
- Twitter

Please install the following tools (if you haven't already):

- Docker
- AWS CLI
- Terraform

_NOTE: this documentation was written on April 2020. Commands and Versions might have changed since. For the purpose of this document I used:_

- _Docker 19.03.8_
- _AWS CLI 2.0.6_
- _Terraform v0.12.24_

This document assumes you're able to place env vars in your `~/.bash_profile`, or equivalent. From now on I'll only refer to the `~/.bash_profile`. I'm also going to assume you reload your profile after adding new env vars. For example by running:

```
. ~/.bash_profile
```

## Terraform module

This application exports a Terraform module from [./terraform/aws](https://github.com/rjmarques/something-of-the-day/blob/master/terraform/aws). The module follows a module API defined in my [ECS cluster repo](https://github.com/rjmarques/my-ecs-cluster). Thus, the provisioning is done from that repo and this most initializes some required resources and outputs this app's ECS container definition.

Additionally, When provisioning this module we get the following ecr repository as output:

- _sotd_ecr_repository_url_

We'll need this at a later point in this tutorial, so be sure to save it somewhere.

## Twitter

This app gets its _somethings_ entirely from Twitter, specifically this account: https://twitter.com/lgst_something. You'll need to go to the Twitter [developer section](https://developer.twitter.com/en/apps) and create a new app and copy the API keys. Once you have them add them to your `~/.bash_profile` as:

```
export TF_VAR_twitter_client_id=<your_API_key>
export TF_VAR_twitter_client_secret=<your_API_secret_key>
```

## Postgres DB

This appliction requires a Postgres Database. It is outside of the scope of this documentation to explain how or where to set one up. I got mine from for free from [Heroku](https://elements.heroku.com/addons/heroku-postgresql).

Assuming you already have one running database, we then create the required schema. The script is available at:

[./backend/datastore/schema.sql](https://github.com/rjmarques/something-of-the-day/blob/master/backend/datastore/schema.sql)

## Docker

Ensure you have docker running:

```
docker version
```

If no errors are visible on both client and server/deamon, we can continue!

We need to link docker with the AWS ECR that was created by Terraform. The _sotd_ecr_repository_url_ follows the following format:

`aws_account_id.dkr.ecr.region_name.amazonaws.com/repository_name`

The domain part of the URL identifies your account and region, on AWS's ECR service. While the repository name identifies which repository to where we'll push our docker images.

During deployment we'll need the full URL. As such, we add the following to the `~/.bash_profile`:

```
export ECR_SOTD_REPO=$ECR/something-of-the-day
```

Note: this project assumes some variables set from the [ECS cluster repo](https://github.com/rjmarques/my-ecs-cluster) have already been defined (e.g., _\$ECR_).

## Application building and deployment

TL;DR:

- To build run: `make`
- To deploy run: `make deploy`

If you inspect the _Makefile_ you can see the steps that each of the previous commands is running. I'll quickly go over each one.

### make

Builds the whole app and runs unit + integration tests. The resulting docker image. representing the production ready app, will be available as _rjmarques/something-of-the-day_.

The build steps are run in docker containers that encapsultate the frontend (nodejs) and backend (golang) environments. Additionally, the integration tests will run against a local Postgres database (also running in a container). The build/test cluster is created using _docker-compose_.

The first time `make` is run it can take a while to finish, as new images are pulled from the internet.

### make deploy

Deployment assumes there's already a cluster ready to run this application. Deploying said cluster falls outside the scope of this documentation. If you want more info on how to deploy your own ECS cluster you may checkout [my-ecs-cluster](https://github.com/rjmarques/my-ecs-cluster).

The Makefile assumes the presence of a few env vars used to push the image and restart the ECS service where this app runs. These must be added to `~/.bash_profile` as:

```
export AWS_REGION=eu-west-2            # the region where my cluster is running
export ECS_CLUSTER=hobby-cluster       # the name of my cluster
export ECS_SERVICE=hobby-ecs-service   # the name of the ECS service running on my cluster
```

Once `make` has been run, and _rjmarques/something-of-the-day_ is created, we can push the image to the ECR repository and consequently update the ECS service that's using it, thus fully deploying the app.

To push to ECR we first tag the image with the ECR URL, represented by _ECR_SOTD_REPO_. Afterwards we force update the _ECS_SERVICE_ running on the _ECS_CLUSTER_. After ECS updates the service (which should take a few seconds), the new image should be in use and the app updated.

## Running the application locally

Before we can run the app locally we have to build it: `make`. Additionally we have to define 3 environment variables that the app is expecting: _CLIENT_SECRET_, _CLIENT_ID_ and, _POSTGRES_URL_. The first two are very easy and we can just add the following to the `~/.bash_profile`:

```
export CLIENT_ID=$TF_VAR_twitter_client_id
export CLIENT_SECRET=$TF_VAR_twitter_client_secret
```

Since the app needs a database to run against we can _borrow_ the one that the integration tests use. To do this, simply add the following to the `~/.bash_profile`:

```
export POSTGRES_URL=postgres://postgres:mysecretpassword@postgres:5432/somethingoftheday?sslmode=disable
```

Now run in another terminal window run:

```
docker-compose up db
```

The DB should now be running, and accessible to other containers via the _something-of-the-day_integration-tests_ network. To start the app inside this network we run:

```
make run_with_localdb
```

On the other hand, if you wish to run the app in the default bridge network you can simply run:

```
make run
```

Once you're done with the DB and nothing is using it you can bring it all down by running:

```
make clean
```

## Checking the logs

While the _something-of-the-day_ container is running it will show application logs.

To get the production logs you must first SSH into your EC2 instance. To find it please log in to you AWS console and look for an instance with the name _something-of-the-day_.

Once inside we need to find the name of the container running the app:

```
docker ps
```

To read the logs we then simply run:

```
docker logs <the_container's_name> | less
docker logs <the_container's_name> | grep "whatever I'm trying to find"
```

## Possible Improvements

- The way I'm overriding terraform variables is perhaps not the cleanest one and caused duplications in the `~/.bash_profile`.
- Unprovisioning the env with terraform leaves behind AWS secrets, as AWS does not delete them immediately. This can cause problems if we want later to re-provision everything.
- Non master branches do not save images to the ECR repo but could do so.
