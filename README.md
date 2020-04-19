# Something of the day 

[![CircleCI](https://circleci.com/gh/rjmarques/something-of-the-day/tree/master.svg?style=svg)](https://circleci.com/gh/rjmarques/something-of-the-day/tree/master)

http://something.ricardomarq.com

This document explains how to build and deploy the **something-of-the-day App**. The app is a completely over-engineered React Single Page Web App that has a Golang backend and connects to a PG database. The app runs locally on Docker, is deployed to AWS ECS and links to an Heroku database. There's also a CI workflow set up on CircleCI, that automatically deploys the app when code is merged to master and passes all the tests.

The goal was to learn a bit about AWS ECS, Terraform, Docker, as well as to set up a relatively reaslistic CI/CD.

##### DISCLAIMER: It goes without saying that a trivial static webpage does not need: React, Docker, Postgresql, etc. The complexity was introduced to allow me to build realistic CI workflow. I.e., have good frontend build & tests, backend build & test, and DB integration tests. Plus, I wanted to push everything to containers as much as possible. Oh and Heroku provides a tiny database for free - win! 💰

## Minimal requirements to build and run

You'll need the following accounts
* Heroku
* AWS
* Twitter

Please install the following tools (if you haven't already):

- Docker
- Heroku CLI
- AWS CLI
- Terraform

_NOTE: this documentation was written on April 2020. Commands and Versions might have changed since. For the purpose of this document I used:_
* _Docker 19.03.8_
* _Heroku CLI 7.39.1_
* _AWS CLI 2.0.6_
* _Terraform v0.12.24_

This document assumes you're able to place env vars in your `~/.bash_profile`, or equivalent. From now on I'll only refer to the `~/.bash_profile`. I'm also going to assume you reload your profile after adding new env vars. For example by running:

```
. ~/.bash_profile
```

All env variables starting with _TF_VAR_ are used by Terraform to set up values for some of its module variables, as defined in:

[./provision/variables.tf`](https://github.com/rjmarques/something-of-the-day/blob/master/provision/variables.tf)

## Linking Heroku to Terraform 

Login to Heroku

```
heroku login -i
```

Create an API token for Terraform

```
heroku authorizations:create --description terraform-sotd
```

This last command will generate an API token for you. Terraform will use this token to authenticate with Heroku.

We need to save the Heroku creds in a way that Terraform can access them. Thus, add the following lines to your `~/.bash_profile`, replacing the params with the real values.

```
export HEROKU_EMAIL=<EMAIL>
export HEROKU_API_KEY=<TOKEN>
```

## Twitter

This app gets its _somethings_ entirely from Twitter, specifically this account: https://twitter.com/lgst_something. You'll need to go to the Twitter [developer section](https://developer.twitter.com/en/apps) and create a new app and copy the API keys. Once you have them add them to your `~/.bash_profile` as:

```
export TF_VAR_twitter_client_id=<your_API_key>
export TF_VAR_twitter_client_secret=<your_API_secret_key>
```

## AWS

Although the AWS provisioning is almost fully automatic, you still need to specify enough information about where and how that provisioning will take place.

First thing we will need is a new AMI User that has enough permissions to create a whole ECS environment on the AWS account. For the sake of simplicity I ran my user as admin, albeit in a real environment the permissions should be more locked down.

Once the AMI user is created set the access key id and secret access key into `~/.aws/credentials`, as you would to enable access to AWS via the aws cli.

We then need to find your AWS account id:

```
aws sts get-caller-identity
```

And add it to the `~/.bash_profile` as:

```
export TF_VAR_account_id=<the id obtained from the last command>
```

While we're at it we also need to specify where in the AWS cloud the app will be running. I picked London, as such I added the following variables to my `~/.bash_profile`:

```
export TF_VAR_region=eu-west-2
export TF_VAR_availability_zone=eu-west-2a
```

Lastly, you must also add a keypair to AWS, or use a pre-existing one, that will allow you to SSH into the EC2 instance that belongs to the cluster. Add the name of the keypair to your `~/.bash_profile` as:

```
export TF_VAR_ecs_key_pair_name=<my-aws-ssh-keypair>
```

## Terraform provisioning
Make sure all your changes to the `~/.bash_profile` are available from the terminal you're trying to run terraform. Also, all Terraform commands should be run from within the _provision_ folder, so please:

```
cd ./provision
```

To intialize Terraform we first need to fetch its providers. We're using the Heroku and AWS providers:

```
terraform init
```

We're almost ready to provision our environment! But first let's check to see if there are no errors in the Terraform plan, before executing it:

```
terraform plan
```

This command will show us what changes will be made on the remote providers. It's useful as a safety check before trying to change something. To actually provision the environment run:

```
terraform apply
```

If everything went well, your new something-of-the-day environment is now up and provisioned! 🎉

The command should also have outputted a few values:
* _ecr_repository_url_
* _heroku_db_url_

We'll need these at a later point in this tutorial, so be sure to save them somewhere. 

_NOTE: To avoid needing an extra service to provide runtime params, this application reads its runtime params from the AWS secrets manager. For example, the Twitter API tokens and Postgres Database URL. However, Terraform saves these values in plain-text in its state files. Be mindful of that, as it can be a potential security breach in a real application!_

## Postgres DB

Before we run our app, we need to initialize the DB with the schema and tables the app needs to access. The _heroku_db_url_, that was outputed when we ran `terraform apply`, follows the typical PG URL normally use to connect to a PG datase. That is, `postgres://user:password@host:5432/database_name`.

Using this URL we should now initialize the DB with the app's schema. The latter are available at:

[./backend/datastore/schema.sql](https://github.com/rjmarques/something-of-the-day/blob/master/backend/datastore/schema.sql) 

You can decide how best to run this schema against the DB, for example via a CLI like _psql_ or a GUI tool.

**Optionally**: to run the app locally against the real database, you can also add the _POSTGRES_URL_ to your _.bash_profile_, as:

```
export POSTGRES_URL="the_heroku_db_url"
```

## Docker

Ensure you have docker running:
```
docker version
```
If no errors are visible on both client and server/deamon, we can continue!

We need to link docker with the AWS ECR that was created by Terraform, and whose URL was outputed when we ran `terraform apply`. The _ecr_repository_url_ follows the following format:

`aws_account_id.dkr.ecr.region_name.amazonaws.com/repository_name`

The domain part of the URL identifies your account and region, on AWS's ECR service. While the repository name identifies which repository to where we'll push our docker images.

During deployment we'll need the full URL and also the URL without the repository name. As such we add the following lines to the `~/.bash_profile`: 

```
export ECR=<the domain part of the ecr_repository_url> # without the /something-of-the-day
export ECR_REPO=$ECR/something-of-the-day
```

If you then run `echo $ECR_REPO` the output should be equal to _ecr_repository_url_.

## Application building and deployment

TL;DR:

* To build run: `make`
* To deploy run: `make deploy`

If you inspect the _Makefile_ you can see the steps that each of the previous commands is running. I'll quickly go over each one.

### make

Builds the whole app and runs unit + integration tests. The resulting docker image. representing the production ready app, will be available as _rjmarques/something-of-the-day_. 

The build steps are run in docker containers that encapsultate the frontend (nodejs) and backend (golang) environments. Additionally, the integration tests will run against a local Postgres database (also running in a container). The build/test cluster is created using _docker-compose_. 

The first time `make` is run it can take a while to finish, as new images are pulled from the internet.

### make deploy

Once `make` has been run, and _rjmarques/something-of-the-day_ is created, we can push the image to the ECR repository and consequently update the ECS service that's using it, thus fully deploying the app.

To push to ECR we first tag the image with the ECR URL, represented by _ECR_REPO_. Afterwards we force update the _sotd-ecs-service_ running on the _sotd-cluster_. After ECS updates the service (which should take a few seconds), the new image should be in use and the app updated.

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

* The way I'm overriding terraform variables is perhaps not the cleanest one and caused duplications in the `~/.bash_profile`.
* Unprovisioning the env with terraform leaves behind AWS secrets, as AWS does not delete them immediately. This can cause problems if we want later to re-provision everything.
* Non master branches do not save images to the ECR repo but could do so. 
* The ECS roles created for this app can potentially collide with ones already existing if the terraform scripts are ran against AWS accounts with existing ECS clusters.


