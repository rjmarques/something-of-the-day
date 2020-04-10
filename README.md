## Setting up the env with Terraform

### Required tools

Please install the following tools:

- Terraform
- Heroku CLI

_NOTE: this documentation was written on March 2020. Commands and Versions might have changed since. For the purpose of this document I used Heroku CLI heroku/7.39.1 and Terraform v0.12.24._

### Linking Terraform to Heroku

Login to Heroku

```
heroku login -i
```

Create an API token for Terraform

```
heroku authorizations:create --description terraform-sotd
```

This last command will generate an API token for you. Terraform will use this token to authenticate with Heroku.

We need to save the Heroku creds in a way that Terraform can access them. Thus, add the following lines to your _.bash_profile_, or equivalent, replacing the params with the real values.

```
export HEROKU_EMAIL=<EMAIL>
export HEROKU_API_KEY=<TOKEN>
```

Intialize Terraform and fetch the Heroku provider

```
terraform init
```

Check to see if there are no errors in the Terraform plan, before executing it

```
terraform plan
```

Provision the app's resources

```
terraform apply
```

Retreive the Postgres URL to be used in by the APP. The command will return a Postgres Database URL in the standard format: _postgres://username:password@host:5432/myDB_.

```
heroku config:get DATABASE_URL --app something-of-the-day
```

Apply the schema to the DB. The schema is availabe at:

`./backend/datastore/schema.sql`

You can decide how to run this schema against the DB, for example via a CLI like `psql` or a GUI tool. To aid in running the App locally you can also add the _DATABASE_URL_ to your _.bash_profile_, or equivalent, as:

```
export POSTGRES_URL="postgres://username:password@host:5432/myDB"
```

## Twitter

* Export twitter client id as `export TF_VAR_twitter_client_id=<your_client_id>`
* Export twitter client secret as `export TF_VAR_twitter_client_secret=<your_client_secret>`


### AWS

* create a new AMI User with admin credentials and get the access key id and secret access key
* add the keys into ~/.aws/credentials
* find your AWS account ID `aws sts get-caller-identity`
* Export aws account id as `export TF_VAR_account_id=<the id obtained from the last command>`
* Export a region where you want to run `export TF_VAR_region=eu-west-2`
* Export the availability zone associated to the default VPC `export TF_VAR_availability_zone=eu-west-2a`

Run terraform

* Export the domain part of *ecr_repository_arn* as `export ERC=<ecr_repository_arn>`
* Export the full erc repo url as `export ERC_REPO=$ERC/something-of-the-day`

Link docker to the AWS ERC for your account and region

Link docker to your amazon ECR 

```
aws ecr get-login-password --region $TF_VAR_region | docker login --username AWS --password-stdin $ERC
```

Tag your image with the destination aws ecr url

```
docker tag rjmarques/something-of-the-day:latest $ERC_URL:latest
```

Push the image

```
docker push $ERC_URL:latest
```

Force ECS to update the container

```
aws ecs update-service --cluster sotd-cluster --service sotd-ecs-service --region $TF_VAR_region --force-new-deployment
```