# Sato

[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://GitHub.com/jameswoolfenden/sato/graphs/commit-activity)
[![Build Status](https://github.com/JamesWoolfenden/sato/workflows/CI/badge.svg?branch=master)](https://github.com/JamesWoolfenden/sato)
[![Latest Release](https://img.shields.io/github/release/JamesWoolfenden/sato.svg)](https://github.com/JamesWoolfenden/sato/releases/latest)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/JamesWoolfenden/sato.svg?label=latest)](https://github.com/JamesWoolfenden/sato/releases/latest)
![Terraform Version](https://img.shields.io/badge/tf-%3E%3D0.14.0-blue.svg)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)
[![checkov](https://img.shields.io/badge/checkov-verified-brightgreen)](https://www.checkov.io/)
[![Github All Releases](https://img.shields.io/github/downloads/jameswoolfenden/sato/total.svg)](https://github.com/JamesWoolfenden/sato/releases)
[![codecov](https://codecov.io/gh/JamesWoolfenden/sato/graph/badge.svg?token=AT1DREJQPR)](https://codecov.io/gh/JamesWoolfenden/sato)

Converts CloudFormation (and now also ARM) into Terraform. In Go, but quickerly.

## Table of Contents

<!--toc:start-->
- [Sato](#sato)
  - [Table of Contents](#table-of-contents)
  - [Install](#install)
    - [Compile](#compile)
    - [MacOS](#macos)
    - [Windows](#windows)
    - [Docker](#docker)
  - [Usage](#usage)
    - [Bisect](#bisect)
    - [Parse](#parse)
    - [See](#see)
    - [Version](#version)
  - [Help](#help)

<!--toc:end-->

## Install

### Compile

Download the latest releases <https://github.com/JamesWoolfenden/sato/releases/tag/v0.1.19> or:

Compile locally:

```bash
git clone https://github.com/JamesWoolfenden/sato
cd sato
go install
```

### MacOS

```shell
brew tap jameswoolfenden/homebrew-tap
brew install jameswoolfenden/tap/pike
```

### Windows

I'm now using Scoop to distribute releases, it's much quicker to update and easier to manage than previous methods,
you can install scoop from <https://scoop.sh/>.

Add my scoop bucket:

```shell
scoop bucket add iac https://github.com/JamesWoolfenden/scoop.git
```

Then you can install a tool:

```pwsh
scoop install sato
```

### Docker

```shell
docker pull jameswoolfenden/sato
docker run --tty --volume /local/path/to/tf:/tf jameswoolfenden/sato scan -d /tf
```

<https://hub.docker.com/repository/docker/jameswoolfenden/sato>

## Usage

### Parse

Get yourself some valid CloudFormation*

```bash
 git clone https://github.com/JamesWoolfenden/aws-cloudformation-templates
 >cd aws-cloudformation-templates/community/codestar/custom-ci-cd-pipeline
 ❯ ls
 README.md    template.yml
 >sato parse -f template.yml
 9:17PM INF Created .sato\variables.tf
 9:17PM INF Created .sato\data.tf
 9:17PM INF Created .sato\aws_codebuild_project.productionbuild.tf
 9:17PM INF Created .sato\aws_codebuild_project.productiondeploy.tf
 9:17PM INF Created .sato\aws_codebuild_project.stagingbuild.tf
 9:17PM INF Created .sato\aws_codebuild_project.stagingdeploy.tf
 9:17PM INF Created .sato\aws_iam_role.codebuildrole.tf
 9:17PM INF Created .sato\aws_codepipeline_pipeline.pipeline.tf
 9:17PM INF Created .sato\aws_iam_role.pipelinerole.tf
 9:17PM INF Created .sato\aws_s3_bucket.pipelines3bucket.tf
```

That's it. So by default (overridable) the parsed CloudFormation (now Terraform) will be in a .sato subdirectory.
So let's have a look see:

```bash
> ls .sato
aws_codebuild_project.productionbuild.tf  aws_codebuild_project.stagingbuild.tf     aws_codepipeline_pipeline.pipeline.tf     aws_iam_role.pipelinerole.tf              variables.tf
aws_codebuild_project.productiondeploy.tf aws_codebuild_project.stagingdeploy.tf    aws_iam_role.codebuildrole.tf             aws_s3_bucket.pipelines3bucket.tf
```

So there are some files that could be Terraform.

### The Cats Pyjamas

Testing...

```bash
>terraform init
...
Terraform has been successfully initialized!
....
>terraform plan
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create
...
Plan: 12 to add, 0 to change, 0 to destroy.
...
```

### See

Shows the Terraform resource equivalent to a CloudFormation resource. Or vice-versa.

This tells you the equivalent resource required, given a CF ..... or an ARM resource;

```bash
$ sato see -r Microsoft.Storage/storageAccounts
azurerm_storage_account
```

or

```bash
$sato see -r AWS::EC2::Instance
aws_instance%
```

### Bisect

ARM to Terraform conversion.

What? You've got these legacy ARM templates, and you'd dearly love to drop them, but you really don't fancy Bicep
and the rework.
I got you covered. Sato now bisects ARM into Terraform - Take one of the Azure quickstart examples from here
<https://github.com/Azure/azure-quickstart-templates/tree/master/quickstarts/microsoft.compute/vm-simple-windows>:

Clone it:

```bash
git clone https://github.com/Azure/azure-quickstart-templates.git
```

Then bisect it!

```bash
$ sato bisect -f /Users/jwoolfenden/code/azure-quickstart-templates/quickstarts/microsoft.compute/vm-simple-windows/azuredeploy.json
1:56PM INF Created /Users/jwoolfenden/code/sato/.sato/variables.tf
1:56PM INF Created /Users/jwoolfenden/code/sato/.sato/locals.tf
1:56PM INF Created /Users/jwoolfenden/code/sato/.sato/azurerm_storage_account.sato0.tf
1:56PM INF Created /Users/jwoolfenden/code/sato/.sato/azurerm_public_ip.sato1.tf
1:56PM INF Created /Users/jwoolfenden/code/sato/.sato/azurerm_network_security_group.sato2.tf
1:56PM INF Created /Users/jwoolfenden/code/sato/.sato/azurerm_virtual_network.sato3.tf
1:56PM INF Created /Users/jwoolfenden/code/sato/.sato/azurerm_network_interface.sato4.tf
1:56PM INF Created /Users/jwoolfenden/code/sato/.sato/azurerm_virtual_machine.sato5.tf
1:56PM INF Created /Users/jwoolfenden/code/sato/.sato/azurerm_virtual_machine_extension.sato6.tf
1:56PM INF Created /Users/jwoolfenden/code/sato/.sato/outputs.tf
1:56PM INF Created /Users/jwoolfenden/code/sato/.sato/data.tf
```

I make an opinionated translation, in Terraform there are no parameters, resources and dependencies are very different,
there's no one for one - ARM to Terraform, so the aim is to get you close to 100%.

There needs to be a lot of work supporting resources and built-in functions/template as yet.
If you want to use this, let me know so, then I'll know to do so, or even better send me a PR.

### Version

```bash
$sato version
9.9.9
```

### Help

```bash
$ sato
NAME:
sato - Translate Cloudformation to Terraform

USAGE:
sato [global options] command [command options] [arguments...]

VERSION:
9.9.9

AUTHOR:
James Woolfenden <jim.wolf@duck.com>

COMMANDS:
bisect      translate ARM to Terraform
parse       translate CFN to Terraform
see         shows equivalent Terraform resource
version, v  Outputs the application version
help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
--help, -h     show help
--version, -v  print the version
```

## Extra credit - <small>Pike</small>

If you use my other tool, Pike you can now apply that and get the policy requirements:

> pike scan -d .sato -o json

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "VisualEditor0",
      "Effect": "Allow",
      "Action": [
        "codebuild:BatchGetProjects",
        "codebuild:CreateProject",
        "codebuild:DeleteProject",
        "codebuild:UpdateProject"
      ],
      "Resource": [
        "*"
      ]
    },
    {
      "Sid": "VisualEditor1",
      "Effect": "Allow",
      "Action": [
        "codepipeline:CreatePipeline",
        "codepipeline:DeletePipeline",
        "codepipeline:GetPipeline",
        "codepipeline:ListTagsForResource"
      ],
      "Resource": [
        "*"
      ]
    },
    {
      "Sid": "VisualEditor2",
      "Effect": "Allow",
      "Action": [
        "iam:CreateRole",
        "iam:DeleteRole",
        "iam:DeleteRolePolicy",
        "iam:GetRole",
        "iam:GetRolePolicy",
        "iam:ListAttachedRolePolicies",
        "iam:ListInstanceProfilesForRole",
        "iam:ListRolePolicies",
        "iam:PassRole",
        "iam:PutRolePolicy"
      ],
      "Resource": [
        "*"
      ]
    },
    {
      "Sid": "VisualEditor3",
      "Effect": "Allow",
      "Action": [
        "s3:CreateBucket",
        "s3:DeleteBucket",
        "s3:GetAccelerateConfiguration",
        "s3:GetBucketAcl",
        "s3:GetBucketCORS",
        "s3:GetBucketLogging",
        "s3:GetBucketObjectLockConfiguration",
        "s3:GetBucketPolicy",
        "s3:GetBucketRequestPayment",
        "s3:GetBucketTagging",
        "s3:GetBucketVersioning",
        "s3:GetBucketWebsite",
        "s3:GetEncryptionConfiguration",
        "s3:GetLifecycleConfiguration",
        "s3:GetObject",
        "s3:GetObjectAcl",
        "s3:GetReplicationConfiguration",
        "s3:ListBucket"
      ],
      "Resource": [
        "*"
      ]
    }
  ]
}

```

## Valid CloudFormation

Ditch it all, OK, OK but some older samples can play fast and lose with the CloudFormation schema and data types.
The Go-formation parser is less accommodating, you may need to be stricter on your typing.

- Booleans are true or false and not "false"
- Ints are 1,2,3 not "1", "2", "3"
