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
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

// Ready is a type that represents the readiness of the system
type Ready struct {
	VaultInitialized        bool `json:"vault_initialized"`
	DigitalOceanInitialized bool `json:"digitalocean_initialized"`
	Ready                   bool `json:"ready"`
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	log.Info("Starting server")
	router.GET("/health", func(c *gin.Context) {
		healthy := getSystemHealth()
		log.Debug("Health check called")
		c.JSON(http.StatusOK, gin.H{"status": healthy})
	})

	return router
}

func setupReadiness(router *gin.Engine) *gin.Engine {
	router.GET("/readiness", func(c *gin.Context) {
		// Check Vault token
		vaultTokenReady := checkVaultToken()
		// Check DigitalOcean token
		doTokenReady := checkDigitalOceanToken()
		ready := vaultTokenReady && doTokenReady

		c.JSON(http.StatusOK, gin.H{"vault_initialized": vaultTokenReady, "digitalocean_initialized": doTokenReady, "ready": ready})
	})
	return router
}

func setupStaticAssets(router *gin.Engine) *gin.Engine {
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")
	router.LoadHTMLGlob("templates/*.tmpl")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"title": "Minecraft Controller"})
	})
	return router
}

func main() {
	// Setup logging
	logger := log.New(os.Stdout)
	logger.SetFormatter(log.JSONFormatter)

	router := setupRouter()
	if router == nil {
		log.Fatal("Failed to setup router")
	}
	router = setupReadiness(router)
	router = setupStaticAssets(router)
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}

}

func checkVaultToken() bool {
	// Implementation of checkVaultToken function
	// this should run a look up of the
	doToken := os.Getenv("VAULT_TOKEN")
	if doToken == "" {
		log.Debug("VAULT_TOKEN environment variable not set")
		return false
	}
	return true
}

func checkDigitalOceanToken() bool {
	// Implementation of checkDigitalOceanToken function
	// this should run a look up of the
	doToken := os.Getenv("DIGITALOCEAN_TOKEN")
	if doToken == "" {
		log.Debug("DIGITALOCEAN_TOKEN environment variable not set")
		return false
	}
	return true

}

func getSystemHealth() string {
	// Implementation of getMinecraftServerStatus function
	// this should run a look up of the
	return "unknown"
}
