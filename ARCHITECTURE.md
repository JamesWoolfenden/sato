# Sato Architecture

## Overview

Sato is a CLI tool that converts Infrastructure-as-Code templates from AWS CloudFormation and Azure Resource Manager (ARM) formats into Terraform configuration files.

## Design Principles

1. **Template-Driven Conversion**: Uses Go's text/template engine to generate Terraform HCL
2. **Schema-Based Parsing**: Leverages CloudFormation schemas for accurate type information
3. **Two-Pass Processing**: First parses input, then transforms to Terraform syntax
4. **Resource Mapping**: Maps cloud provider resource types to Terraform providers

## Package Structure

### src/cf

**Purpose**: CloudFormation to Terraform conversion

Key files:

- `parse.go`: Main parsing logic and orchestration
- `helpers.go`: Template functions for CF-specific transformations
- `resources.go`: Resource conversion logic
- `variables.go`: Parameter/variable handling
- `data.go`: Data source generation
- `lookup.go`: CF resource type to Terraform type mapping

### src/arm

**Purpose**: Azure Resource Manager to Terraform conversion

Key files:

- `parse.go`: Main parsing logic and orchestration
- `helpers.go`: Template functions for ARM-specific transformations
- `resources.go`: Resource conversion logic
- `variables.go`: Parameter/variable handling
- `outputs.go`: Output value conversion
- `data.go`: Data source generation

### src/see

**Purpose**: Shared resource mapping and lookup

Key files:

- `lookup.go`: Resource type lookup across providers
- `resource_mapping.go`: Mapping tables for resource types

### src/version

**Purpose**: Application version information

- Version is injected at build time via ldflags

## Data Flow

### CloudFormation Conversion

```map
JSON Template → Parse() → Goformation Library → Extract Components
                                                        ↓
                                                   Preprocess
                                                        ↓
                                         ┌──────────────┴──────────────┐
                                         ↓                             ↓
                                  ParseVariables                 ParseResources
                                         ↓                             ↓
                                  variables.tf              resource templates
                                                                      ↓
                                                              ParseOutputs
                                                                      ↓
                                                                outputs.tf
```

### ARM Conversion

```map
JSON Template → Parse() → json.Unmarshal → Preprocess
                                                ↓
                              ┌─────────────────┴─────────────┐
                              ↓                               ↓
                      ParseVariables                   ParseResources
                              ↓                               ↓
                      variables.tf/locals.tf          resource templates
                                                              ↓
                                                       ParseOutputs
                                                              ↓
                                                          outputs.tf
```

## Template System

### Template Functions

Both cf and arm packages provide custom template functions that are available in Go templates:

Common functions:

- `Quote/Dequote`: Handle string quoting
- `Snake/Kebab`: Convert naming conventions
- `Split/SplitOn`: String manipulation
- `Array/ArrayReplace`: Array formatting
- `Tags/Map`: Complex type formatting
- `Marshal`: JSON encoding

CloudFormation-specific:

- Functions map CF intrinsic functions like `Ref`, `GetAtt`, `Sub` to Terraform syntax

ARM-specific:

- Functions map ARM template functions like `resourceId`, `concat`, `parameters` to Terraform

### Template Flow

1. Parse input template (JSON)
2. Extract resources, variables, outputs
3. For each component:
   - Load corresponding .template file
   - Execute template with parsed data
   - Template functions transform cloud-specific syntax to Terraform
4. Write generated .tf files

## Resource Naming

Generates meaningful resource names based on resource type and name:

```go
// In generateResourceName()
// Example: "Microsoft.Network/virtualNetworks" + "myVNET"
//       -> "virtual_networks_myvnet"
baseName := generateResourceName(resourceType, resourceName, item)
```

The naming logic:

- Extracts resource type from path (e.g., `virtualNetworks` from `Microsoft.Network/virtualNetworks`)
- Converts to snake_case
- Cleans and appends resource name if available
- Handles duplicates with numeric suffixes
- Falls back to sequential naming (`sato0`, `sato1`) if type/name unavailable

## Type Mapping

### CloudFormation → Terraform

Handled in `src/cf/lookup.go`:

- Maps `AWS::*` resource types to `azurerm_*` or `aws_*` providers
- Example: `AWS::EC2::Instance` → `aws_instance`

### ARM → Terraform

Handled in `src/see/lookup.go`:

- Maps `Microsoft.*` resource types to `azurerm_*` providers
- Example: `Microsoft.Compute/virtualMachines` → `azurerm_virtual_machine`

## String Processing

Uses `strings.Builder` for efficient string concatenation in:

- Tag formatting
- Map formatting
- Complex object building

This provides better performance than string concatenation in loops.

## Version Management

Version is managed via Go build ldflags:

```makefile
LDFLAGS=-ldflags "-X sato/src/version.Version=$(VERSION)"
```

For releases, goreleaser injects the version automatically.

## Error Handling

Custom error types for different failure scenarios:

- `parseError`: Template parsing failures
- `filepathError`: File path resolution issues
- `writeFileError`: File writing failures
- `unmarshalError`: JSON unmarshaling failures

All errors are logged using zerolog for structured logging.

## Testing Strategy

### Unit Tests

- Each package has *_test.go files
- Tests focus on individual functions and transformations
- Use testify for assertions

### Integration Tests

- Located in tests/integration/
- Test full conversion workflows
- Compare output against expected .tf files

## Key Design Decisions

### 1. Why Template-Based?

Templates provide flexibility to adjust Terraform output format without changing Go code. The template syntax is closer to HCL, making it easier to maintain.

### 2. Why Separate arm/cf Packages?

CloudFormation and ARM have different:

- JSON structures
- Intrinsic functions
- Resource naming conventions
- Type systems

Separate packages allow specialized handling while sharing common code via src/see.

### 3. Schema Usage

CloudFormation schemas provide:

- Property types and constraints
- Required vs optional fields
- Enumeration values

This enables more accurate conversion than guessing from example templates.

## Extension Points

To add support for a new resource type:

1. Add mapping in `src/see/resource_mapping.go`
2. Create template file if needed
3. Add resource-specific transformation logic if required
4. Add test cases

To add a new template function:

1. Implement function in helpers.go
2. Add to funcMap in parse.go
3. Document in template files
4. Add unit tests

## Performance Considerations

- Template parsing is done once per file type (cached by Go's template engine)
- String building uses strings.Builder for efficiency
- File I/O is minimized by batching operations
- Vendor directory ensures consistent dependency versions

## Security Considerations

- File paths are validated before use
- File permissions set to restrictive values (0600 for files, 0750 for dirs)
- #nosec directives document intentional security exceptions
- gosec linter catches common security issues
- govulncheck ensures dependencies have no known vulnerabilities
