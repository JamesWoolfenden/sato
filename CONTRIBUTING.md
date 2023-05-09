# Contributing

I'll be more than happy to get any contributions.
I'm currently using the pre-commit framework, you should have that installed
and running on the code you want to add.
James

## adding resource support

This is my set-up there are others but this one is mine.

Set-up IDE/Goland to debug target cfn, this config might help:
 
```xml
  <component name="RunManager">
    <configuration name="go build sato" type="GoApplicationRunConfiguration" factoryName="Go Application" nameIsGenerated="true">
      <module name="sato" />
      <working_directory value="$PROJECT_DIR$" />
      <parameters value="parse -f $PROJECT_DIR$/../aws-security-survival-kit/cfn-global.yml " />
      <kind value="PACKAGE" />
      <package value="sato" />
      <directory value="$PROJECT_DIR$" />
      <filePath value="$PROJECT_DIR$" />
      <method v="2" />
    </configuration>
  </component>
```

Add a new resource mapping to **lookup.go**

```golang
"AWS::Logs::MetricFilter":                          awsCloudwatchLogMetricFilter,
```

Add template files to **src/resources/** e.g. aws_cloudwatch_log_metric_filter.template

Import resources into go - **src/resource.go**

```golang

//go:embed resources/aws_cloudwatch_log_metric_filter.template
var awsCloudwatchLogMetricFilter []byte
```

Iterate on creation of tf using a debugger
Set a break point to see schema of new resources in parse resources in lookup.go line 34 :

``` golang
err = Write(ReplaceDependant(ReplaceVariables(output.String())), destination, fmt.Sprint(ToTFName(myType), ".", strings.ToLower(item)))
```
Find the resource and copy the schema to the template file and work out the templating:

```gotemplate
resource "aws_cloudwatch_log_metric_filter" "{{.item}}" {
{{- if .resource.FilterName }}
  name          = {{.resource.FilterName|Quote}}
{{- else}}
  name          = "{{.item}}"
{{- end}}
  pattern       = {{Replace .resource.FilterPattern "\"" "\\\x22" |Quote}}
  log_group_name= {{Nil .resource.LogGroupName|Quote}}

{{- if .resource.MetricTransformations}}
{{- range $a, $i:= .resource.MetricTransformations}}
  metric_transformation     {
{{- if $i.DefaultValue}}
    default_value= {{$i.DefaultValue}}
{{- end}}
{{- if $i.Dimensions }}
    dimensions   {
    {{$i.Dimensions}}
    }
{{- end}}
    name         = {{$i.MetricName |Quote}}
    namespace    = {{$i.MetricNamespace|Quote}}
    value        = {{$i.MetricValue}}
{{- if $i.Unit}}
    unit         = {{$i.Unit |Quote}}
{{- end}}
  }
{{- end}}
{{- end}}
}
```
