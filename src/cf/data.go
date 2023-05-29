package cf

const dataAvailabilityZone = `data "aws_availability_zone" "example" {
  name = ""
}
`

const dataSubnet = `data "aws_subnet" "selected" {
	id = ""
}
`

const dataKeyPair = `data "aws_key_pair" "example" {
  key_name = ""
  include_public_key = true
}
`

const dataVpc = `data "aws_vpc" "selected" {
  id = ""
}
`

const dataRegion = "data \"aws_region\" \"current\" {}\n"

const dataSecurityGroup = "data \"aws_security_group\" \"selected\" {\n  id = \"\"\n}\n"

const provider = "provider \"aws\" {\n  region=\"eu-west-2\"\n}\n"
