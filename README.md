# Something of the day 

This document explains how to build and deploy the **something-of-the-day App**. The app is a completely over-engineered React Single Page Web App that has a Golang backend and connects to a PG database. The app runs locally on Docker, is deployed to AWS ECS and links to an Heroku database. There's also a CI workflow set up on CircleCI, that automatically deploys the app when code is merged to master and passes all the tests.

The goal was to learn a bit about AWS ECS, Terraform, Docker, as well as to set up a relatively reaslistic CI/CD.

##### DISCLAIMER: It goes without saying that a trivial static webpage does not need: React, Docker, Postgresql, etc. The complexity was introduced to allow me to build realistic CI workflow. I.e., have good frontend build & tests, backend build & test, and DB integration tests. Plus, I wanted to push everything to containers as much as possible. Oh and Heroku provides a tiny database for free - win!

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

This app was build to be able to provision itself on AWS from the ground without much faff. However, you still need to specify enough information about where and how that provisioning will take place.

First thing we will need is a new AMI User that has enough permissions to create a whole ECS environment on the AWS account. For the sake of simplicity I ran my user as admin, albeit in a real environment the permissions should be more locked down.

Once the AMI user is created set the access key id and secret access key into `~/.aws/credentials`, as you would to enable access to AWS via the aws cli.

We then need to find your AWS account id:

```
aws sts get-caller-identity
```

And add it to the `~/.bash_profile`:

```
export TF_VAR_account_id=<the id obtained from the last command>
```

While we're at it we also need to specify where in the AWS cloud we will run the app. I picked London, as such I added the following variables to my `~/.bash_profile`:

```
export TF_VAR_region=eu-west-2
export TF_VAR_availability_zone=eu-west-2a
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

This command will show us what changes will be made on the remote providers. It's useful as a safety check before trying to change something. 

To actually provision the environment run:

```
terraform apply
```

If everything went well, your new something-of-the-day environment is now up and provisioned! 🎉

The command should also have outputed a few values:
* _ecr_repository_url_
* _heroku_db_url_

We'll need for the rest of the tutorial, so be sure to save them somewhere. 

_NOTE: To avoid needing an extra service to provide runtime params, this application reads its runtime params from the AWS secrets manager. For example, the Twitter API tokens and Postgres Database URL. However, Terraform saves these values in plain-text in its state files. Be mindful of that, as it can be a potential security breach in a real application!_

## Postgres DB

Before we run our app we need to initialize the DB with the schema and tables the app needs to access. The _heroku_db_url_, that was outputed when we ran `terraform apply`, follows the typical PG URL normally use to connect to a PG datase. That is, `postgres://user:password@host:5432/database_name`.

Using this URL we should now initialize the DB with the app's schema. The latter are available at:

[./backend/datastore/schema.sql](https://github.com/rjmarques/something-of-the-day/blob/master/backend/datastore/schema.sql) 

You can decide how best to run this schema against the DB, for example via a CLI like _psql_ or a GUI tool.

**Optionally**: to run the App locally against the real database, you can also add the _POSTGRES_URL_ to your _.bash_profile_, as:

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

The domain part of the URL identifies your account and region, on AWS's ECR service. While the repository name naturally identifies which repository to where we'll push our docker images.

During deployment we'll need the full URL and also the URL without the repository name. As such we add the following lines to the `~/.bash_profile`: 

```
export ERC=<the domain part of the ecr_repository_url> # without the /something-of-the-day
export ERC_REPO=$ERC/something-of-the-day
```

If you then run `echo $ERC_REPO` the output should be equal to _ecr_repository_url_.

## Application building and deployment

TL;DR:

* To build run: `make`
* To deploy run: `make deploy`

If you inspect the _Makefile_ you can see the steps that each of the previous commands is issuing. I'll shortly go over each one to explain what they do.

### make

### make deploy

## Running the application locally


# TODO terraform import aws_secretsmanager_secret.example arn:aws:secretsmanager:us-east-1:123456789012:secret:example-123456
# TODO aws ecr get-login-password --region $TF_VAR_region