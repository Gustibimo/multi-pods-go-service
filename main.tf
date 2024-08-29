terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.0"
    }
  }
}

provider "kubernetes" {
  config_context_cluster = "minikube"
  config_path            = "~/.kube/config"
}


provider "helm" {
  kubernetes {
    config_context_cluster = "minikube"
    config_path            = "~/.kube/config"
  }
}

resource "null_resource" "docker_build_push" {
  triggers = {
    always_run = "${timestamp()}"
  }

  provisioner "local-exec" {
    command = "docker build -t gstbimo/bom-import-backend:latest ."
  }

  provisioner "local-exec" {
    command = "docker push gstbimo/bom-import-backend:latest"
  }
}

resource "helm_release" "bom_import_backend" {
  depends_on = [null_resource.docker_build_push]  # Wait for Docker build and push to complete
  name       = "bom-import-backend"
  chart      = "./deploy"  # Path to your Helm chart directory (current directory in this case)
  repository = ""  # Leave empty if your chart is in the local directory
  recreate_pods = true
  set {
    name  = "image.tag"
    value = "latest"
  }

  # Add more configuration values here as needed
}

output "deployed_release" {
  value = helm_release.bom_import_backend.name
}
 
