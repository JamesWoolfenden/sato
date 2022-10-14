resource "aws_sns_topic" "ExampleTopic" {
  name = "example"
  kms_master_key_id = ""
}
