# AWS setup
If you have AWS account you might need to install the DevBot there.

## Prerequisites
- install the aws cli into your system. So you will able to use the AWS cli commands.
- create `~/.aws/credentials` file and put there your AWS credentials OR export AWS credentials in memory
- install the docker to your system
- after aws cli were installed and credentials are exported into memory, please run ecr login command
- install terraform to your system
```
aws  ecr get-login-password --region us-east-1
```

## 1 Step - Push your image to AWS ECR repository
### 1. Create the repository
First you need to create ECR repository where you will store your image
```
aws ecr create-repository --repository-name devbot
```
To make sure the repository exists, run next command:
```
aws ecr describe-repositories
```
You should see your repository there. Example:
``` 
{
    "repositories": [
        {
            "repositoryArn": "arn:aws:ecr:us-east-1:aws_account_id:repository/devbot",
            "registryId": "aws_account_id",
            "repositoryName": "devbot",
            "repositoryUri": "aws_account_id.dkr.ecr.us-east-1.amazonaws.com/devbot",
            "createdAt": 1590751353.0
        }
    ]
}
```

### 2. Build and tag the image
First you need build the **devbot** image using `docker`
```
docker build -t devbot:latest .
```
After image build was finished, please tag your image
``` 
docker tag TAG_ID aws_account_id.dkr.ecr.us-east-1.amazonaws.com/devbot:latest
```
`TAG_ID` you can find using this command
```
docker images | grep "devbot"
```
As the result you will see something like this
``` 
devbot    latest    c4dac21ebc43(this is TAG_ID)    20 minutes ago    740MB
```

### 2. Push your image into ECR repository
Run next command to push your image into ECR repository:
``` 
docker push aws_account_id.dkr.ecr.us-east-1.amazonaws.com/devbot
```

## 2 Step - Create terraform ECS task container definition file
You for doing this you can manually run the next command `make tf-container-definition`. This step will create the `local.container_definition.tf` file if it doesn't exists inside of `terraform` directory. 
In that file you can add or update existing variables for your devbot ECS instance.

## 3 Step - Apply changes
Before start you need first to initialize the terraform inside of the terraform folder.

So, please go to the `terraform` folder and run `terraform init` command there. If everything is good, then you can run next command `terraform apply`, which will applies the changes in AWS. Please, make sure you check changes before applying!
If something goes wrong, you always can run `terraform destroy` to destroy all created instances.