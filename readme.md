Simple service for AWS tutorial 

### how to run?
> go run cmd/main.go

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
1. [code syntax](https://docs.aws.amazon.com/codebuild/latest/userguide/build-spec-ref.html#build-spec-ref-syntax)
2. [build go application](https://dev.classmethod.jp/articles/building-go-project-in-codebuild/)i
3. [tutorial](https://tutorialedge.net/golang/creating-restful-api-with-golang/)