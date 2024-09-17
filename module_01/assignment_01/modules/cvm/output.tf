output "public_ip" {
  description = "cloud vm public ip address"
  value       = tencentcloud_instance.web[0].public_ip
}

output "user" {
    description = "cloud vm ssh user"
    value = var.user
}

output "password" {
    description = "cloud vm ssh password"
    value = var.password
}

output "docker_ver" {
  description = "docker version"
  value = data.local_file.docker_ver
}