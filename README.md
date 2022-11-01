# Sato

[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://GitHub.com/jameswoolfenden/sato/graphs/commit-activity)
[![Build Status](https://github.com/JamesWoolfenden/sato/workflows/CI/badge.svg?branch=master)](https://github.com/JamesWoolfenden/sato)
[![Latest Release](https://img.shields.io/github/release/JamesWoolfenden/sato.svg)](https://github.com/JamesWoolfenden/sato/releases/latest)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/JamesWoolfenden/sato.svg?label=latest)](https://github.com/JamesWoolfenden/sato/releases/latest)
![Terraform Version](https://img.shields.io/badge/tf-%3E%3D0.14.0-blue.svg)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)
[![checkov](https://img.shields.io/badge/checkov-verified-brightgreen)](https://www.checkov.io/)
[![Github All Releases](https://img.shields.io/github/downloads/jameswoolfenden/sato/total.svg)](https://github.com/JamesWoolfenden/sato/releases)

Converts CloudFormation into Terraform. In Go, quickerly.

## Install

Download from latest releases <https://github.com/JamesWoolfenden/sato/releases/tag/v0.0.9> or:

Compile locally:

```golang
git clone https://github.com/JamesWoolfenden/sato
cd sato
go install
```

## Usage

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

That's it (I should make it obvious that its worked). So by default (overidable) the parsed CloudFormation (now Terraform) will be in a .sato sub directory.
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

What's missing?

- Lots of CloudFormation support - this is very new project
- variables and resource relationships need work (Current limitation of a library)

Extra credit?

If you use my other tool, Pike you can now apply that and get the policy requirements:

>pike scan -d .sato -o json

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

Ditch it all,Ok ok.. but some older samples can play fast and lose with the CloudFormation schema and data types.
The Go-formation parser is less accommodating, you may need to be stricter on your typing.

- Booleans are true or false and not "false"
- Ints are 1,2,3 not "1", "2", "3"
