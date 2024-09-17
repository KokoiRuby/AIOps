实践 Terraform，开通腾讯云虚拟机，并安装 Docker。

`module_01/assignment_01/modules/cvm/main.tf`

```json
resource "tencentcloud_instance" "web" {
  # ...

  # install docker
  connection {
    type = "ssh"
    user = var.user
    password = var.password
    host = self.public_ip
  }

  provisioner "remote-exec" {
    inline = [ 
      # https://docs.docker.com/engine/install/ubuntu/
      # Set up Docker's apt repository.
      "sudo apt-get update",
      "sudo apt-get install -y ca-certificates curl",
      "sudo install -m 0755 -d /etc/apt/keyrings",
      "sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc",
      "sudo chmod a+r /etc/apt/keyrings/docker.asc",
      "echo \"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo \"$VERSION_CODENAME\") stable\" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null",
      "sudo apt-get update",
      # Install the Docker packages.
      "sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin",
    ]
  }
}

# docker version
resource "null_resource" "get_docker_ver" {
  provisioner "local-exec" {
    command = "sshpass -p '${var.password}' ssh -o StrictHostKeyChecking=no ${var.user}@${tencentcloud_instance.web[0].public_ip} 'sudo docker version' > /tmp/docker_ver.txt"
  }

  depends_on = [ tencentcloud_instance.web ]
}

data "local_file" "docker_ver" {
  filename = "/tmp/docker_ver.txt"
  depends_on = [ null_resource.get_docker_ver ]
}
```

`module_01/assignment_01/modules/cvm/output.tf`

```json
# ...

output "docker_ver" {
  description = "docker version"
  value = data.local_file.docker_ver
}
```

`module_01/assignment_01/dev/output.tf`

```json
output "docker_ver" {
  description = "docker version"
  value = module.cvm.docker_ver
}
```

