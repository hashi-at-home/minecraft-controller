package main

// This application exposes a simple REST API for controlling Minecraft servers in DigitalOcean droplets.
// The API exposes only one endpoint: /instance which accepts GETs, and POSTs and DELETEs
// GET requests are used to retrieve the status of the Minecraft server
// POST requests are used to start the Minecraft server
// DELETE requests are used to stop the Minecraft server

// The main function starts the server, authenticating to Vault via a token
// provided by the vault agent.
//
// The API is provided by gin-gonic
import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"

	vault "github.com/hashicorp/vault/api"
)

func main() {
	// Setup logging
	logger := log.New(os.Stdout)
	logger.SetFormatter(log.JSONFormatter)

	// configure Vault client
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = "http://active.vault.service.consul:8200"

	router := gin.Default()
	log.Info("Starting server")
	router.GET("/health", func(c *gin.Context) {
		healthy := getSystemHealth()
		log.Info("Health check")
		c.JSON(200, gin.H{"status": healthy})
	})

	// Add a readiness check endpoint, which will return a 200 and a message indicating the status of access to secrets
	// If the vault token is valid, vault_initialized will be true.
	// If the digitalocean token is valid, digitalocean_initialized will be true.
	router.GET("/readiness", func(c *gin.Context) {
		// Check Vault token
		vaultTokenReady := checkVaultToken()
		// Check DigitalOcean token
		doTokenReady := checkDigitalOceanToken()

		c.JSON(200, gin.H{"vault_initialized": vaultTokenReady, "digitalocean_initialized": doTokenReady})
	})

	router.Run(":8080")
}

func checkVaultToken() bool {
	// Implementation of checkVaultToken function
	// this should run a look up of the
	return false
}

func checkDigitalOceanToken() bool {
	// Implementation of checkDigitalOceanToken function
	// this should run a look up of the
	return false
}

func getSystemHealth() string {
	// Implementation of getMinecraftServerStatus function
	// this should run a look up of the
	return "unknown"
}
