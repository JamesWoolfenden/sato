# CloudFormation Intrinsic Functions Support

Sato now supports CloudFormation intrinsic functions with automatic conversion to Terraform equivalents.

## Supported Functions

### !Sub (String Substitution)

**CloudFormation:**

```yaml
!Sub "arn:aws:s3:::${BucketName}"
!Sub "Hello ${AWS::Region}"
```

**Terraform Output:**

```hcl
"arn:aws:s3:::${var.bucketname}"
"Hello ${data.aws_region.current.name}"
```

**Usage in Code:**

```go
ParseSub(input, parameters)
```

**Supported Pseudo-Parameters:**

- `${AWS::Region}` → `${data.aws_region.current.name}`
- `${AWS::AccountId}` → `${data.aws_caller_identity.current.account_id}`
- `${AWS::StackName}` → `${var.stack_name}`
- `${AWS::StackId}` → `${var.stack_id}`
- `${AWS::Partition}` → `${data.aws_partition.current.partition}`
- `${AWS::URLSuffix}` → `${data.aws_partition.current.dns_suffix}`

---

### !Join (Array Joining)

**CloudFormation:**

```yaml
!Join ["-", ["web", "server", "01"]]
```

**Terraform Output:**

```hcl
join("-", ["web", "server", "01"])
```

**Template Usage:**

```go
{{TFJoin "-" (list "web" "server" "01")}}
```

---

### !Split (String Splitting)

**CloudFormation:**

```yaml
!Split ["|", "a|b|c"]
```

**Terraform Output:**

```hcl
split("|", "a|b|c")
```

**Template Usage:**

```go
{{TFSplit "|" "a|b|c"}}
```

---

## Coverage & Testing

All intrinsic function parsers are thoroughly tested:

| Function     | Test Coverage | Test Cases |
|--------------|---------------|------------|
| `ParseSub`   | 92.9%         | 8 cases    |
| `ParseJoin`  | 100%          | 6 cases    |
| `ParseSplit` | 100%          | 6 cases    |

---

## Using in Templates

The intrinsic function parsers are available in template rendering via `funcMap`:

```go
// In your resource template (e.g., aws_s3_bucket.template):
resource "aws_s3_bucket" "{{.item}}" {
  bucket = "{{TFJoin "-" (list "myapp" .resource.Environment)}}"

  tags = {
    Name = "{{TFJoin "-" (list "bucket" .resource.Environment)}}"
  }
}
```

---

## Examples

### Example 1: S3 Bucket with Sub

**CloudFormation:**

```yaml
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Sub "${Environment}-${AWS::Region}-data"
```

**Sato Output:**

```hcl
resource "aws_s3_bucket" "mybucket" {
  bucket = "${var.environment}-${data.aws_region.current.name}-data"
}
```

### Example 2: Security Group with Join

**CloudFormation:**

```yaml
Resources:
  MySecurityGroup:
    Type: AWS::EC2::SecurityGroup
    Properties:
      GroupDescription: !Join [" ", ["Security", "group", "for", !Ref Environment]]
```

**Sato Output:**

```hcl
resource "aws_security_group" "mysecuritygroup" {
  description = join(" ", ["Security", "group", "for", var.environment])
}
```

---

## Future Enhancements

Planned intrinsic function support:

- `!Select` - Array indexing
- `!GetAZs` - Get availability zones
- `!If` - Conditional values
- `!FindInMap` - Map lookups
- `!Base64` - Base64 encoding
- `!Cidr` - CIDR calculations

---

## Implementation Details

### Architecture

1. **Parsing Phase**: CloudFormation template → goformation parser
2. **Variable Conversion**: `ParseVariables()` extracts parameters
3. **Resource Conversion**: Templates render with intrinsic function helpers
4. **Post-processing**: `ReplaceVariables()` and `ReplaceDependant()` clean up references

### Code Location

- Implementation: `src/cf/parse.go` (lines 399-461)
- Tests: `src/cf/intrinsic_test.go`
- Template Functions: Added to `funcMap` in `parse.go` (lines 70-71)

---

## Testing

Run intrinsic function tests:

```bash
# All intrinsic function tests
go test -v ./src/cf -run "TestParse(Sub|Join|Split)"

# Specific function
go test -v ./src/cf -run "TestParseSub"

# With coverage
go test ./src/cf -cover
```

Expected output:

```bash
ok      sato/src/cf    2.193s
```
