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
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/digitalocean/godo"
	"github.com/gin-gonic/gin"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
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
	doEnvCheck := checks.NewEnvCheck("DIGITALOCEAN_TOKEN")
	healthcheck.New(router, config.DefaultConfig(), []checks.Check{doEnvCheck}) //#nosec G104 -- This is a false positive

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

func setupDigitalOcean(router *gin.Engine) *gin.Engine {
	// Set up routes for interacting with digital ocean
	// We need to set up the digital ocean client
	doClient := godo.NewFromToken(os.Getenv("DIGITALOCEAN_TOKEN"))
	if doClient == nil {
		log.Fatal("Failed to create DigitalOcean client")
	}

	// First, we will implement the GET route for getting all of the
	// droplets
	context := context.TODO()
	droplets, dropletsResp, err := doClient.Droplets.ListByTag(context, "minecraft", &godo.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.Info(dropletsResp.StatusCode)

	for _, droplet := range droplets {
		log.Info(droplet.Name)
	}
	router.GET("/droplets", func(c *gin.Context) {
		c.JSON(http.StatusOK, droplets)
	})

	router.GET("/droplets/:id", func(c *gin.Context) {
		id := c.Param("id")
		nid, err := strconv.Atoi(id)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		}

		droplet, _, err := doClient.Droplets.Get(context, nid)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, droplet)
	})

	router.DELETE("/droplets/:id", func(c *gin.Context) {
		id := c.Param("id")
		nid, err := strconv.Atoi(id)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		}

		_, err = doClient.Droplets.Delete(context, nid)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"message": "Droplet deleted"})
	})

	router.DELETE("/droplets", func(c *gin.Context) {
		// Implementation of DELETE /droplets endpoint
		// this should delete all droplets with tag minecraft

		resp, err := doClient.Droplets.DeleteByTag(context, "minecraft")
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete droplets"})
		}
		log.Infof("Deleted droplets: %s", resp.Body)
		// If the response code is not ok, return an error
		if resp.StatusCode != http.StatusNoContent {
			log.Infof("Failed to delete droplets code %v", resp.StatusCode)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete droplets - response was not OK"})
		}

		c.JSON(http.StatusOK, gin.H{"message": resp.Body})
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
	router = setupDigitalOcean(router)
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
