package database

import (
	"log"

	"github.com/brunobarlari/inventorypulse/internal/config"
	"github.com/brunobarlari/inventorypulse/internal/domain/models"
	"gorm.io/gorm"
)

// RunSeeder seeds initial data into the database
func RunSeeder(db *gorm.DB, cfg *config.Config) error {
	log.Println("Running database seeder...")

	// Migrate quantity to stock for existing products
	if err := migrateQuantityToStock(db); err != nil {
		log.Printf("Warning: Could not migrate quantity to stock: %v", err)
	}

	// Seed admin user
	if err := seedAdminUser(db, cfg); err != nil {
		return err
	}

	// Seed sample categories
	if err := seedCategories(db); err != nil {
		return err
	}

	// Seed sample products
	if err := seedProducts(db); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully")
	return nil
}

// migrateQuantityToStock migrates data from quantity column to stock column
func migrateQuantityToStock(db *gorm.DB) error {
	// Check if quantity column exists
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_name = 'products' AND column_name = 'quantity'").Scan(&count)
	if count == 0 {
		return nil // Column doesn't exist, nothing to migrate
	}

	// Update stock from quantity where stock is 0 and quantity > 0
	result := db.Exec("UPDATE products SET stock = quantity WHERE stock = 0 AND quantity > 0")
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		log.Printf("Migrated %d products from quantity to stock", result.RowsAffected)
	}

	return nil
}

func seedAdminUser(db *gorm.DB, cfg *config.Config) error {
	var count int64
	db.Model(&models.User{}).Count(&count)

	if count > 0 {
		log.Println("Users already exist, skipping admin user seeding")
		return nil
	}

	admin := &models.User{
		Email: cfg.Admin.Email,
		Role:  models.RoleAdmin,
	}

	if err := admin.SetPassword(cfg.Admin.Password); err != nil {
		return err
	}

	if err := db.Create(admin).Error; err != nil {
		return err
	}

	log.Printf("Admin user created: %s", admin.Email)
	return nil
}

func seedCategories(db *gorm.DB) error {
	var count int64
	db.Model(&models.Category{}).Count(&count)

	if count > 0 {
		log.Println("Categories already exist, skipping category seeding")
		return nil
	}

	categories := []models.Category{
		{Name: "Electronics", Description: "Electronic devices and gadgets"},
		{Name: "Clothing", Description: "Apparel and fashion items"},
		{Name: "Books", Description: "Books and publications"},
		{Name: "Home & Garden", Description: "Home improvement and garden supplies"},
		{Name: "Sports", Description: "Sports equipment and accessories"},
	}

	for _, category := range categories {
		if err := db.Create(&category).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d categories", len(categories))
	return nil
}

func seedProducts(db *gorm.DB) error {
	var count int64
	db.Model(&models.Product{}).Count(&count)

	if count > 0 {
		log.Println("Products already exist, skipping product seeding")
		return nil
	}

	products := []models.Product{
		{Name: "Laptop Pro 15", Description: "High-performance laptop with 15-inch display", SKU: "ELEC-001", Stock: 50, Price: 1299.99, CategoryID: 1},
		{Name: "Wireless Mouse", Description: "Ergonomic wireless mouse with USB receiver", SKU: "ELEC-002", Stock: 200, Price: 29.99, CategoryID: 1},
		{Name: "USB-C Hub", Description: "7-in-1 USB-C hub with HDMI and SD card reader", SKU: "ELEC-003", Stock: 150, Price: 49.99, CategoryID: 1},
		{Name: "Cotton T-Shirt", Description: "100% cotton casual t-shirt", SKU: "CLTH-001", Stock: 300, Price: 19.99, CategoryID: 2},
		{Name: "Denim Jeans", Description: "Classic fit denim jeans", SKU: "CLTH-002", Stock: 150, Price: 59.99, CategoryID: 2},
		{Name: "Programming in Go", Description: "Complete guide to Go programming language", SKU: "BOOK-001", Stock: 75, Price: 39.99, CategoryID: 3},
		{Name: "Clean Code", Description: "A handbook of agile software craftsmanship", SKU: "BOOK-002", Stock: 100, Price: 34.99, CategoryID: 3},
		{Name: "Garden Tool Set", Description: "5-piece stainless steel garden tool set", SKU: "HOME-001", Stock: 80, Price: 44.99, CategoryID: 4},
		{Name: "LED Desk Lamp", Description: "Adjustable LED desk lamp with USB charging", SKU: "HOME-002", Stock: 120, Price: 35.99, CategoryID: 4},
		{Name: "Yoga Mat", Description: "Non-slip yoga mat with carrying strap", SKU: "SPRT-001", Stock: 200, Price: 24.99, CategoryID: 5},
	}

	for _, product := range products {
		if err := db.Create(&product).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d products", len(products))
	return nil
}

