# Slack Timer

Simple periodically timer on Slack, sending a message to DM Channel.

> ⚠️ This project is [Kata](https://en.wikipedia.org/wiki/Kata_(programming)) with [Golang](https://golang.org/) and [Clean Architecuture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).  Please use [Reminder](https://slack.com/help/articles/208423427-Set-a-reminder) if the timer feature on Slack is needed.

![Timer](website/peter-yost-I9jJXmNkXR4-unsplash.jpg?raw=true)

<span>Photo by <a href="https://unsplash.com/@odysseus_?utm_source=unsplash&amp;utm_medium=referral&amp;utm_content=creditCopyText">Peter Yost</a> on <a href="https://unsplash.com/s/photos/stopwatch?utm_source=unsplash&amp;utm_medium=referral&amp;utm_content=creditCopyText">Unsplash</a></span>


<!-- ![GitHub Actions](https://img.shields.io/github/workflow/status/takakd/sandbox/echo)
&nbsp;&nbsp;![GitHub Actions](https://img.shields.io/github/workflow/status/takakd/sandbox/echo)
&nbsp;&nbsp;![Badge-Rhs](https://img.shields.io/badge/Badge-Rhs-orange?style=flat)
&nbsp;&nbsp;![Badge-Rhs](https://img.shields.io/badge/Badge-Rhs-orange?style=flat)
&nbsp;&nbsp; -->
![CircleCI](https://circleci.com/gh/takakd/slack-timer.svg?style=shield&circle-token=195739304092ae914a95802605704f56171b0627)
&nbsp;
![License-MIT](https://img.shields.io/badge/License-MIT-informational?style=flat)

## Table of Contents

- [Features](#features)
- [Setup](#setup)
- [Usage](#usage)
- [Development](#development)
- [License](#license)

## Features

* Notify message by interval minutes we set. 
* Set On/Off notification.

## Setup

### Requirements

* Slack account: necessary to create Slack App.
* AWS account: necessary to control Lambda, DynamoDB, SQS, CloudWatch, and CloudFormation.
* [AWS SAM CLI, version 1.7.0](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

### 1. Deploy to AWS

Deploy codes to AWS by SAM CLI.

```
$ cd ./deployments/sam/slack-timer

# Build
$ sam build
Building codeuri: ../../../cmd/lambda/enqueue runtime: go1.x metadata: {} functions: ['EnqueueFunction']
...
Build Succeeded
...

# Deploy to AWS.
$ sam deploy --guided
Configuring SAM deploy
======================
...
        Stack Name [sam-app]: slacktimer                 
...
... enter attributes as you like
...
CloudFormation outputs from deployed stack
...
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
Outputs                                                                                                                                                                                                                                                          
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
Key                 SetTimerAPI                                                                                                                                                                                                                                  
Value               https://<API_ID>.execute-api.ap-northeast-1.amazonaws.com/Prod/callback/slack                                                                                                                                                                      
...
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

Successfully created/updated stack - slacktimer in ap-northeast-1
```

### 2. Creating Slack App

Create Slack App on your workspace which has Event Subscriptions, Bots and Permissions features.
Ref: https://api.slack.com/start/overview#what

#### Event Subscriptions

* Set enable events.
* Enter the API URL deployed at "1. Deploy to AWS" to Request URL Field.  It will show "Verified✅" if it had succeeded deployment.
* Add "message.im" event to "Subscribe to events on behalf of users" section.

#### Bots

* Enter "App Display Name" as you like.
* Enable "Message Tab" on "Show Tabs" section.

#### Persmissions

Add scopes as follow.

Features and functionality | permission
--- | ----
Bot Token Scopes | channels:read
Bot Token Scopes | chat:write
Bot Token Scopes | im:history
Bot Token Scopes | im:read
Bot Token Scopes | im:write
Bot Token Scopes | mpim:read
User Token Scopes | im.history

## Usage

Provide three commands to control this app. Add the app channel on Slack workspace and enter commands in the app channel message window.

Command | Format | Action
--- | --- | ---
Set | set `minutes` `text` | Notify `text` by `minutes`
On | on | Start to notify.
Off | off | Suspend to notify.

**e.g.**

Receive `I'm active` by 15 minutes.
```
set 15 I'm active!
```

Suspend the notification.
```
off
```

Start the notification.
```
on
```

## Development

### Tech Stacks

* AWS Serverless: Lambda, DynamoDB, SQS, CloudWatch, API Gateway and CloudFormation.
* Golang version go1.14.4 darwin/amd64
* Clean Architecture

### Setup

1. Install Golang by following [Download and install](https://golang.org/doc/install)
2. Run `go mod vendor` for getting moduels.
2. Install AWS SAM CLI by following [Installing the AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

### Command

**Testing**
```
# With details; "-v" and "-cover".
$ make test

# No details; no flags.
$ make test_light
```

**Code format**
```
# Run "go fmt", in addition to "goimports" and "go lint".
$ make fmt
```

### Structure

* Directory structure is followed by [golang-standards/project-layout](https://github.com/golang-standards/project-layout).
* Design is followed by Clean architecture.
https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

#### Design

#### Source

```shell
.
|-- .circleci
|   `-- config.yml      <-- Circle CI config
|-- .gitignore
|-- LICENSE.md
|-- Makefile            <-- Defines make command targets
|-- README.md           <-- This instruction file
|-- cmd
|   `-- lambda
|       |-- enqueue
|       |   `-- main.go     <-- Entrypoint of Lambda; Enqueue
|       |-- notify
|       |   `-- main.go     <-- Entrypoint of Lambda; Notification
|       `-- settime
|           `-- main.go     <-- Entrypoint of Lambda; Set command
|
|-- deployments
|   `-- sam     <-- AWS SAM project directory
|       `-- slack-timer
|           |-- .gitignore
|           |-- Makefile
|           |-- README.md
|           `-- template.yaml     <-- CloudFormation template
|
|-- docs
|   `-- note.md     <-- Notes
|-- go.mod      <-- go modules list
|-- go.sum      <-- go modules hash list
|-- internal
|   |-- app     <-- This app directory
|   |   |-- adapter         <-- Interface Adapters layer
|   |   |   |-- enqueue     <-- Codes related to SQS Enqueueing
|   |   |   |   |-- cloudwatchlogspresenter.go
|   |   |   |   |-- ...
|   |   |   |-- notify      <-- Codes related to Notifiy
|   |   |   |   |-- cloudwatchlogspresenter.go
|   |   |   |   |-- ...
|   |   |   |-- settime     <-- Codes related to set command
|   |   |   |   |-- controller.go
|   |   |   |   |-- ...
|   |   |   |-- slackhandler      <-- Codes related to Slack API handling
|   |   |   |   |-- slackhandler.go
|   |   |   `-- validator         <-- Validation methods.
|   |   |       |-- validate_error_bag.go
|   |   |       `-- validate_error_bag_test.go
|   |   |
|   |   |-- driver                <-- Frameworks & Drivers layer
|   |   |   |-- lambdahandler     <-- Codes related to Lambdahandler
|   |   |   |   |-- enqueue
|   |   |   |   |   |-- lambdafunctor.go
|   |   |   |   |   |-- ...
|   |   |   |   |-- notify
|   |   |   |   |   |-- lambdafunctor.go
|   |   |   |   |   |-- ...
|   |   |   |   `-- settime
|   |   |   |       |-- lambdafunctor.go
|   |   |   |       `-- ...
|   |   |   |-- queue           <-- Codes related to handle SQS
|   |   |   |   |-- sqs.go
|   |   |   |   |-- ...
|   |   |   |-- repository      <-- Codes related to handle DynamoDB
|   |   |   |   |-- dynamodb.go
|   |   |   |   |-- ...
|   |   |   `-- slack           <-- Codes related to handle Slack API
|   |   |       |-- api.go
|   |   |       `-- ...
|   |   |
|   |   |-- enterpriserule      <-- Enterprise Business Rules layer
|   |   |   |-- timerevent.go
|   |   |   `-- ...
|   |   |
|   |   |-- usecase                 <-- Application Business Rules layer
|   |   |   |-- enqueueevent        <-- Enqueueing usecase
|   |   |   |   |-- inputport.go
|   |   |   |   `-- ...
|   |   |   |-- notifyevent         <-- Notification usecase
|   |   |   |   |-- inputport.go
|   |   |   |   `-- ...
|   |   |   |-- timeroffevent       <-- Set off notification usecase
|   |   |   |   |-- inputport.go
|   |   |   |   `-- ...
|   |   |   |-- timeronevent        <-- Set on notification usecase
|   |   |   |   |-- inputport.go
|   |   |   |   `-- ...
|   |   |   `-- updatetimerevent    <-- Set minutes and notification text usecase
|   |   |       |-- inputport.go
|   |   |       `-- ...
|   |   |
|   |   `-- util      <-- Codes shared throughout the app
|   |       |-- appcontext          <-- Context includes Lambda handler context
|   |       |   |-- appcontext.go
|   |       |   `-- ...
|   |       |-- appinitializer      <-- Initialize the app function
|   |       |   `-- appinitializer.go
|   |       |-- config              <-- Config
|   |       |   |-- config.go
|   |       |   |-- driver          <-- Concrete implementation of Config methods
|   |       |   |   |-- envconfig.go
|   |       |   |   `-- ...
|   |       |   `-- ...
|   |       |-- di                <-- Dependency Injection methods
|   |       |   |-- container     <-- Concrete implementation of DI methods
|   |       |   |   |-- dev
|   |       |   |   `-- ...
|   |       |   |-- di.go
|   |       |   `-- ...
|   |       `-- log             <-- Logging
|   |           |-- driver      <-- Concrete implementation of Logging methods
|   |           |   |-- cloudwatchlogger.go
|   |           |   `-- ...
|   |           |-- logger.go
|   |           `-- ...
|   |
|   `-- pkg               <-- Codes shared, which are not dependent on the app
|       |-- collection    <-- Collection structure
|       |   |-- set.go
|       |   `-- ...
|       `-- helper        <-- Helper functions
|           |-- file.go
|           |-- http.go
|           |-- time.go
|           |-- type.go
|   
|-- scripts     <-- Scripts used by Makefile
|   `-- local.sh
|-- website     <-- GitHub readme assets
|   `-- peter-yost-I9jJXmNkXR4-unsplash.jpg
```

## Support

Get in touch is here.

- [Dev.to](https://dev.to/takakd)
- [Twitter](https://twitter.com/takakdkd)

## Contributing

This is just [Kata](https://en.wikipedia.org/wiki/Kata_(programming)) project, but welcome to issues and reviews, don't hesitate to create issue and PR.

## License

- **[MIT license](http://opensource.org/licenses/mit-license.php)**
- Copyright 2020 © takakd</a>.
