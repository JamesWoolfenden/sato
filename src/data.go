package sato

const dataAvailabilityZone = "data \"aws_availability_zone\" \"example\" {\n  name = \"\"\n}\n"

const dataSubnet = "data \"aws_subnet\" \"selected\" {\n  id = \"\"\n}\n"

const dataKeyPair = "data \"aws_key_pair\" \"example\" {\n  key_name           = \"\"\n  include_public_key = true\n }\n"

const dataVpc = "data \"aws_vpc\" \"selected\" {\n  id = \"\"\n}\n"

const dataRegion = "data \"aws_region\" \"current\" {}\n"

const dataSecurityGroup = "data \"aws_security_group\" \"selected\" {\n  id = \"\"\n}\n"

const provider = "provider \"aws\" {\n  region=\"eu-west-2\"\n}\n"
