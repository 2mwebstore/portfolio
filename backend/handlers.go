package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- Site config ---

func getConfig(c *gin.Context) {
	var cfg SiteConfig
	db.First(&cfg)
	c.JSON(http.StatusOK, cfg)
}

func updateConfig(c *gin.Context) {
	var cfg SiteConfig
	db.First(&cfg)
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Save(&cfg)
	c.JSON(http.StatusOK, cfg)
}

// --- Partners ---

func listPartners(c *gin.Context) {
	var partners []Partner
	db.Order("sort_order asc").Find(&partners)
	c.JSON(http.StatusOK, partners)
}

func createPartner(c *gin.Context) {
	var p Partner
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&p)
	c.JSON(http.StatusCreated, p)
}

func updatePartner(c *gin.Context) {
	id := c.Param("id")
	var p Partner
	if err := db.First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "partner not found"})
		return
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Save(&p)
	c.JSON(http.StatusOK, p)
}

func deletePartner(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Partner{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

// --- Payment / social icons ---

func listPaymentIcons(c *gin.Context) {
	var icons []PaymentIcon
	db.Order("category asc, sort_order asc").Find(&icons)
	c.JSON(http.StatusOK, icons)
}

func createPaymentIcon(c *gin.Context) {
	var icon PaymentIcon
	if err := c.ShouldBindJSON(&icon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&icon)
	c.JSON(http.StatusCreated, icon)
}

func updatePaymentIcon(c *gin.Context) {
	id := c.Param("id")
	var icon PaymentIcon
	if err := db.First(&icon, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "icon not found"})
		return
	}
	if err := c.ShouldBindJSON(&icon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Save(&icon)
	c.JSON(http.StatusOK, icon)
}

func deletePaymentIcon(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&PaymentIcon{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

// --- Banners (homepage hero swiper) ---

// listBanners is used by the public site, so it only returns active
// banners in display order.
func listBanners(c *gin.Context) {
	var banners []Banner
	db.Where("active = ?", true).Order("sort_order asc").Find(&banners)
	c.JSON(http.StatusOK, banners)
}

// listAllBanners is used by the admin panel — includes inactive banners
// so they can be re-enabled later instead of only ever being deleted.
func listAllBanners(c *gin.Context) {
	var banners []Banner
	db.Order("sort_order asc").Find(&banners)
	c.JSON(http.StatusOK, banners)
}

func createBanner(c *gin.Context) {
	var b Banner
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&b)
	c.JSON(http.StatusCreated, b)
}

func updateBanner(c *gin.Context) {
	id := c.Param("id")
	var b Banner
	if err := db.First(&b, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "banner not found"})
		return
	}
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Save(&b)
	c.JSON(http.StatusOK, b)
}

func deleteBanner(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Banner{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}
