package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// uploadsDir is where uploaded images are written and served from.
// Set from DATA_DIR so it can point at a mounted persistent volume in production.
var uploadsDir string

func main() {
	// DATA_DIR should point at a persistent volume in production (e.g. the
	// Droplet's ./data folder mounted at /app/data). Defaults to the working
	// directory for local dev.
	dataDir := getEnvOr("DATA_DIR", ".")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatal("failed to prepare data dir:", err)
	}
	uploadsDir = filepath.Join(dataDir, "uploads")

	dbPath := filepath.Join(dataDir, "car_rental.db")

	// If car_rental.db already exists (e.g. in the Droplet's ./data folder),
	// this is a real, already-running site — skip seeding placeholder data
	// entirely so a redeploy never touches existing content. Only a
	// brand-new database gets the starter rows.
	_, statErr := os.Stat(dbPath)
	isNewDB := os.IsNotExist(statErr)

	var err error
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	if err := db.AutoMigrate(&Partner{}, &SiteConfig{}, &PaymentIcon{}, &User{}, &Banner{}); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	if isNewDB {
		log.Println("no existing car_rental.db found — seeding starter data")
		seedData()
	} else {
		log.Println("car_rental.db already exists — skipping seed data")
	}
	seedAdmin()

	r := gin.Default()
	r.Use(cors.Default())

	// Serve uploaded logos/banners (from the persistent data dir)
	r.Static("/uploads", uploadsDir)

	api := r.Group("/api")
	{
		// Public reads — the storefront needs these without logging in.
		api.GET("/config", getConfig)
		api.GET("/partners", listPartners)
		api.GET("/payment-icons", listPaymentIcons)
		api.GET("/banners", listBanners)

		// Login is public by definition.
		api.POST("/auth/login", login)

		// Everything that changes data requires a valid admin token.
		admin := api.Group("/")
		admin.Use(authRequired)
		{
			admin.PUT("/config", updateConfig)

			admin.POST("/partners", createPartner)
			admin.PUT("/partners/:id", updatePartner)
			admin.DELETE("/partners/:id", deletePartner)

			admin.POST("/payment-icons", createPaymentIcon)
			admin.PUT("/payment-icons/:id", updatePaymentIcon)
			admin.DELETE("/payment-icons/:id", deletePaymentIcon)

			admin.GET("/banners/all", listAllBanners)
			admin.POST("/banners", createBanner)
			admin.PUT("/banners/:id", updateBanner)
			admin.DELETE("/banners/:id", deleteBanner)

			admin.PUT("/auth/password", changePassword)

			admin.POST("/upload", uploadImage)
		}
	}

	// Serve the built Vue frontend (frontend/dist, copied into ./web by the
	// Dockerfile) so the whole app is one deployable service. Any path that
	// isn't /api or /uploads falls back to index.html so Vue Router's
	// history-mode routes (e.g. /admin) work on a hard refresh.
	webDir := getEnvOr("WEB_DIR", "./web")
	if _, err := os.Stat(webDir); err == nil {
		r.Static("/assets", filepath.Join(webDir, "assets"))
		r.NoRoute(func(c *gin.Context) {
			p := c.Request.URL.Path
			if strings.HasPrefix(p, "/api") || strings.HasPrefix(p, "/uploads") {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
			c.File(filepath.Join(webDir, "index.html"))
		})
	} else {
		log.Println("no frontend build found at", webDir, "— API only")
	}

	port := getEnvOr("PORT", "8080")
	log.Println("server running on :" + port)
	r.Run(":" + port)
}

// seedAdmin creates a single default admin account on first run.
// CHANGE THIS PASSWORD IMMEDIATELY after first login (via the
// /api/auth/password endpoint, or the admin UI's password field).
func seedAdmin() {
	var count int64
	db.Model(&User{}).Count(&count)
	if count > 0 {
		return
	}

	defaultPassword := getEnvOr("ADMIN_DEFAULT_PASSWORD", "123123123")
	hash, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("failed to hash default admin password:", err)
	}

	db.Create(&User{Username: "admin", PasswordHash: string(hash)})
	log.Println("======================================================")
	log.Println(" Created default admin account:")
	log.Println("   username: admin")
	log.Printf("   password: %s\n", defaultPassword)
	log.Println(" Log in and change this password right away.")
	log.Println("======================================================")
}

// seedData inserts placeholder rows on first run so the frontend
// has something to render immediately. Replace via the API
// (or directly in car_rental.db) once you have real content.
func seedData() {
	var count int64

	db.Model(&SiteConfig{}).Count(&count)
	if count == 0 {
		db.Create(&SiteConfig{
			SiteName:        "VP168 OFFICIAL",
			LogoURL:         "https://imagedelivery.net/NG23G0bfLLwjBy-wmkJ4Aw/52249037-776c-4fa0-d4af-5a635322d200/public",
			TelegramLabel:   "តេលេក្រាមផ្លូវការ",
			TelegramURL:     "https://t.me/yourchannel",
			TelegramIconURL: "https://res.cloudinary.com/dzfrjxhvl/image/upload/v1722498213/telegram_uyayab.jpg",
			TelegramColor:   "#16294f",
			FacebookLabel:   "ហ្វេសប៊ុកផ្លូវការ",
			FacebookURL:     "https://facebook.com/yourpage",
			FacebookIconURL: "https://res.cloudinary.com/dzfrjxhvl/image/upload/v1722498209/facebook_pfu9qj.png",
			FacebookColor:   "#16294f",
			BannerImageURL:  "#",
			BannerLink:      "#",
			BackgroundColor: "#2da5e1",
			FooterNote:      "Sell Car Group",
		})
	}

	db.Model(&Partner{}).Count(&count)
	if count == 0 {
		partners := []Partner{
			{Name: "SB24", LogoURL: "https://res.cloudinary.com/dzfrjxhvl/image/upload/v1722498213/sb24-logo_rmx3fd.png", MediaImageURL: "https://res.cloudinary.com/dzfrjxhvl/image/upload/v1722498215/sb24_hazko8.jpg", Link: "https://aamm24.win/login", Detail: "មាន់ជល់ច្រើនប៉ុស៍ ហាងឆេងខ្ពស់ ភាគរយឈ្នះច្រើន មាននៅ SB24", SortOrder: 1},
			{Name: "Lotto8888", LogoURL: "https://landing-v1.2m-sy.com/images/lotto-logo.png", MediaImageURL: "https://landing-v1.2m-sy.com/images/lotto-bg.png", Link: "https://today8888.net/", Detail: "គ្រាន់តែចុះឈ្មោះបង្កើតអាខោន និងកំសាន្តជាមួយ ឡូតូ8888 ក៏មានឪកាស់ឈ្នះប្រាក់យ៉ាងច្រើនសន្ធឹកសន្ធាប់", SortOrder: 2},
			{Name: "SBC369", LogoURL: "https://landing-v1.2m-sy.com/images/sbc369-ogo.png", MediaImageURL: "#", Link: "https://cat369.com/", Detail: "ជាគេហទំព័រកំសាន្ដអនឡាញ ដែលទទួលបាននូវការជឿទុកចិត្តខ្ពស់។", SortOrder: 3},
		}
		db.Create(&partners)
	}

	db.Model(&PaymentIcon{}).Count(&count)
	if count == 0 {
		icons := []PaymentIcon{
			{Name: "ABA", IconURL: "https://res.cloudinary.com/dzfrjxhvl/image/upload/v1722498210/aba_igvq2r.png", Category: "bank", SortOrder: 1},
			{Name: "Wing", IconURL: "https://res.cloudinary.com/dzfrjxhvl/image/upload/v1722498214/wing_xqtdrh.png", Category: "bank", SortOrder: 2},
			{Name: "TrueMoney", IconURL: "https://res.cloudinary.com/dzfrjxhvl/image/upload/v1722498214/true_slhf9j.png", Category: "bank", SortOrder: 3},
			{Name: "Telegram", IconURL: "https://res.cloudinary.com/dzfrjxhvl/image/upload/v1722498213/telegram_uyayab.jpg", Category: "social", SortOrder: 1},
			{Name: "Facebook", IconURL: "https://res.cloudinary.com/dzfrjxhvl/image/upload/v1722498209/facebook_pfu9qj.png", Category: "social", SortOrder: 2},
			{Name: "Tiktok", IconURL: "https://res.cloudinary.com/dzfrjxhvl/image/upload/v1722498213/tiktok_pbdtgg.webp", Category: "social", SortOrder: 3},
		}
		db.Create(&icons)
	}

	db.Model(&Banner{}).Count(&count)
	if count == 0 {
		banners := []Banner{
			{ImageURL: "https://vp-168.com/uploads/1785423756332313107.jpg", Link: "#", SortOrder: 1, Active: true},
			{ImageURL: "https://imagedelivery.net/NG23G0bfLLwjBy-wmkJ4Aw/25b85533-7fc1-429f-d87f-dc572ea4b400/public", Link: "#", SortOrder: 2, Active: true},
			{ImageURL: "https://vp-168.com/uploads/1785422722850879222.jpg", Link: "#", SortOrder: 3, Active: true},
		}
		db.Create(&banners)
	}
}
