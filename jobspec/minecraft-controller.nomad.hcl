job "minecraft" {

  group "controller" {
    network {
      port "http" {
        to = 8080
      }
    }

    task "api" {
      driver = "docker"

      config {
        image          = "ghcr.io/hashi-at-home/minecraft-operator:latest"
        ports          = ["http"]
        auth_soft_fail = true
      }

      identity {
        env  = true
        file = true
        aud  = ["vault.io"]
      }

      resources {
        cpu    = 500
        memory = 256
      }

      vault {
        change_mode = "noop"
        env = true
        # policies      = ["nomad-read"]
      }

      template {
        data = <<-EOH
{{ with secret "hashiatho.me-v2/data/digitalocean" }}
DIGITALOCEAN_TOKEN={{ .Data.data.minecraft_controller }}
{{ end }}
        EOH
        destination = "/.env"
        env = true
      }


      service {
        name = "minecraft-controller"
        port = "http"
        tags = ["urlprefix-/mc/ strip=/mc", "urlprefix-/mc/healthz strip=/mc"]

        check {
          name = "minecraft-controller-ready"
          type     = "http"
          path     = "/"
          interval = "10s"
          timeout  = "2s"
        }
        check {
          name = "minecraft-controller-healthy"
          type     = "http"
          path     = "/healthz"
          interval = "10s"
          timeout  = "2s"
        }
      }
    }
  }
}
