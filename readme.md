Simple service for AWS tutorial 

### version 0.0.4-SNAPSHOT

Sample with:
1) XRay  
2) DynamoDB  
3) API GW + lambda and SNS (extra material)   

### how to run?
> go run cmd/sever/main.go

--- 

TODO list:
1) refactoring main service and remove db layer
2) tests

--- 

### AWS XRAY
0) I must add your credential for AWS XRAY [intstuction](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials) and I higly recommend using new IAM role and user


## 1.1) you can start app with Docker:  
 1.1.2 sudo docker build -t aws-go .  
 1.1.2 sudo docker run -p 8082:8082 aws-go  
 1.1.2 run makefile with xray
   
## 1.2) start with go   
 1.2.1 cd /project_directory 
 1.2.2 go run .  
 1.2.3 download xray: [link](https://docs.aws.amazon.com/xray/latest/devguide/xray-daemon.html#xray-daemon-permissions)  
 1.2.4 run xray (just make xray command)

### useful:
1. [code syntax build spec](https://docs.aws.amazon.com/codebuild/latest/userguide/build-spec-ref.html#build-spec-ref-syntax)
2. [build go application codebuild](https://dev.classmethod.jp/articles/building-go-project-in-codebuild/)i
3. [tutorial rest api with Go](https://tutorialedge.net/golang/creating-restful-api-with-golang/)
4. [DynamoDb sdk](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-scan-table-item.html)  
5. [DynamoDb create table sdk](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-load-table-items-from-json.html)  