package main

// Partner is one rental company listed on the directory
// (e.g. "Toyota Rental", "Store VIP" in the reference UI).
type Partner struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	Name          string `json:"name"`
	LogoURL       string `json:"logo_url"`
	MediaImageURL string `json:"media_image_url"`
	Link          string `json:"link"`
	Detail        string `json:"detail"`
	SortOrder     int    `json:"sort_order"`
}

// SiteConfig holds the single-row site-wide settings:
// logo, site name, hero banner, and social links.
type SiteConfig struct {
	ID                 uint   `gorm:"primaryKey" json:"id"`
	SiteName           string `json:"site_name"`
	LogoURL            string `json:"logo_url"`
	TelegramLabel      string `json:"telegram_label"`
	TelegramURL        string `json:"telegram_url"`
	TelegramIconURL    string `json:"telegram_icon_url"`
	TelegramColor      string `json:"telegram_color"`
	FacebookLabel      string `json:"facebook_label"`
	FacebookURL        string `json:"facebook_url"`
	FacebookIconURL    string `json:"facebook_icon_url"`
	FacebookColor      string `json:"facebook_color"`
	TiktokURL          string `json:"tiktok_url"`
	InstagramURL       string `json:"instagram_url"`
	BannerImageURL     string `json:"banner_image_url"`
	BannerLink         string `json:"banner_link"`
	BackgroundColor    string `json:"background_color"`
	BackgroundImageURL string `json:"background_image_url"`
	FooterNote         string `json:"footer_note"`
}

// User is an admin account that can log in and manage site content.
type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Username     string `gorm:"uniqueIndex" json:"username"`
	PasswordHash string `json:"-"`
}

// PaymentIcon is a small logo shown in the footer —
// either a bank/e-wallet icon or a social icon.
// Category is "bank" or "social".
type PaymentIcon struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	IconURL   string `json:"icon_url"`
	Link      string `json:"link"`
	Category  string `json:"category"`
	SortOrder int    `json:"sort_order"`
}

// Banner is one slide in the homepage hero swiper. Multiple banners
// rotate automatically on the public site; Active lets an admin hide
// a banner without deleting it.
type Banner struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ImageURL  string `json:"image_url"`
	Link      string `json:"link"`
	SortOrder int    `json:"sort_order"`
	Active    bool   `gorm:"default:true" json:"active"`
}
