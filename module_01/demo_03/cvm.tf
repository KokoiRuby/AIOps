terraform {
  required_providers {
    tencentcloud = {
      source = "tencentcloudstack/tencentcloud"
    }
  }
}

terraform {
  backend "cos" {
    region = "ap-guangzhou"
    bucket = "terraform-bucket-1329112499"
    prefix = "terraform/state"
    encrypt = true
  }
}


# Configure the TencentCloud Provider
provider "tencentcloud" {
  secret_id  = "BLOCKED BY GITHUB"
  secret_key = "BLOCKED BY GITHUB"
  region     = "ap-guangzhou"
}

# Get availability zones
data "tencentcloud_availability_zones" "default" {}

# Get availability images
data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

# Get availability instance types
data "tencentcloud_instance_types" "default" {
  cpu_core_count = 2
  # memory_size    = 1
}

# Create a web server
resource "tencentcloud_instance" "web" {
  instance_name              = "web server"
  availability_zone          = data.tencentcloud_availability_zones.default.zones.0.name
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  allocate_public_ip         = true
  internet_max_bandwidth_out = 20
  orderly_security_groups    = [tencentcloud_security_group.default.id]
}

# Create security group
resource "tencentcloud_security_group" "default" {
  name        = "web accessibility"
  description = "make it accessible for both production and stage ports"
}

# Create security group rule allow web request
resource "tencentcloud_security_group_rule" "web" {
  security_group_id = tencentcloud_security_group.default.id
  type              = "ingress"
  cidr_ip           = "0.0.0.0/0"
  ip_protocol       = "tcp"
  port_range        = "80,8080"
  policy            = "accept"
}

# Create security group rule allow ssh request
resource "tencentcloud_security_group_rule" "ssh" {
  security_group_id = tencentcloud_security_group.default.id
  type              = "ingress"
  cidr_ip           = "0.0.0.0/0"
  ip_protocol       = "tcp"
  port_range        = "22"
  policy            = "accept"
}

# Tencent COS
# data "tencentcloud_user_info" "info" {}

# locals {
#   app_id = data.tencentcloud_user_info.info.app_id
# }

# resource "tencentcloud_cos_bucket" "private_sbucket" {
#   bucket = "terraform-bucket-${local.app_id}"
#   acl    = "private"
#   encryption_algorithm = "AES256"
# }