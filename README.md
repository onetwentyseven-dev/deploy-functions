# AWS Lambda Function Async Deployer

This is a very simple Golang script that utilizes V2 of the AWS Golang SDK and Golang's Go Routines to Asynchrounsly upload zip files to AWS Lambda Functions.

This script takes a total of two (3) args.

1. [Required] Is the directory where the zip archives reside 
2. [Optional] Is the format string that will be pass the name of the zip file with the `.zip` extension trimmed off. use this if the name of the function being updated has additional identifiers associated with it (i.e. The zip is called animals and the lambda function name is animals-handler then this args would be %s-handler)

### Executing without a format string

```
$ deploy-functions ./dist
```

### Executing with a format string

```
$ deploy-functions ./dist %s-handler
```

