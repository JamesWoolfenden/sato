package sato

const dataAvailabilityZone = "data \"aws_availability_zone\" \"example\" {\n  name = \"\"\n}\n"

const dataSubnet = "data \"aws_subnet\" \"selected\" {\n  id = \"\"\n}\n"

const dataKeyPair = "data \"aws_key_pair\" \"example\" {\n  key_name           = \"\"\n  include_public_key = true\n\n  filter {\n    name   = \"tag:Component\"\n    values = [\"web\"]\n  }\n}\n"
