package main

import (
	"log"
	"portfolio-be/internal/config"
	"portfolio-be/internal/database"
	"portfolio-be/internal/models"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.InitSQLite(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run migrations to ensure all tables exist
	err = db.AutoMigrate(
		&models.Content{},
		&models.Upload{},
		&models.Experience{},
		&models.Service{},
		&models.Technology{},
		&models.Project{},
		&models.Testimonial{},
		&models.User{},
	)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Starting database seeding...")

	// Seed Admin User
	adminUser := models.User{
		Username: "admin",
		Email:    "admin@moclaw.dev",
		Password: "admin123", // This will be hashed automatically by BeforeCreate hook
		Role:     "admin",
		IsActive: true,
	}

	// Check if admin user already exists
	var existingUser models.User
	result := db.Where("username = ?", adminUser.Username).First(&existingUser)
	if result.Error != nil {
		// Admin user doesn't exist, create it
		if err := db.Create(&adminUser).Error; err != nil {
			log.Printf("Failed to create admin user: %v", err)
		} else {
			log.Println("✓ Admin user created successfully")
			log.Printf("  Username: %s", adminUser.Username)
			log.Printf("  Email: %s", adminUser.Email)
			log.Printf("  Password: admin123")
		}
	} else {
		log.Println("✓ Admin user already exists")
	}

	// Seed Services
	services := []models.Service{
		{
			Title:    "Full Stack Developer",
			Icon:     "web.png",
			Order:    1,
			IsActive: true,
		},
		{
			Title:    ".NET Developer",
			Icon:     "backend.png",
			Order:    2,
			IsActive: true,
		},
		{
			Title:    "Software Engineer",
			Icon:     "mobile.png",
			Order:    3,
			IsActive: true,
		},
	}

	for _, service := range services {
		if err := db.FirstOrCreate(&service, models.Service{Title: service.Title}).Error; err != nil {
			log.Printf("Failed to create service %s: %v", service.Title, err)
		} else {
			log.Printf("Created/Updated service: %s", service.Title)
		}
	}

	// Seed Technologies
	technologies := []models.Technology{
		{Name: "C#", Icon: "c-sharp.png", Category: "programming", Order: 1, IsActive: true},
		{Name: ".NET", Icon: "NETcore.png", Category: "framework", Order: 2, IsActive: true},
		{Name: "Ubuntu", Icon: "Ubuntu.png", Category: "os", Order: 3, IsActive: true},
		{Name: "RabbitMQ", Icon: "RabbitMQ.png", Category: "messaging", Order: 4, IsActive: true},
		{Name: "Redis", Icon: "Redis.png", Category: "database", Order: 5, IsActive: true},
		{Name: "Google Cloud", Icon: "GoogleCloud.png", Category: "cloud", Order: 6, IsActive: true},
		{Name: "Jenkins", Icon: "Jenkins.png", Category: "ci/cd", Order: 7, IsActive: true},
		{Name: "Jira", Icon: "Jira.png", Category: "project-management", Order: 8, IsActive: true},
		{Name: "Kubernetes", Icon: "Kubernetes.png", Category: "container", Order: 9, IsActive: true},
		{Name: "React", Icon: "React.png", Category: "frontend", Order: 10, IsActive: true},
		{Name: "Github Actions", Icon: "gitaction.png", Category: "ci/cd", Order: 11, IsActive: true},
		{Name: "Docker", Icon: "docker.png", Category: "container", Order: 12, IsActive: true},
		{Name: "AWS", Icon: "AWS.png", Category: "cloud", Order: 13, IsActive: true},
		{Name: "Kafka", Icon: "ApacheKafka.png", Category: "messaging", Order: 14, IsActive: true},
		{Name: "MongoDB", Icon: "mongodb.png", Category: "database", Order: 15, IsActive: true},
	}

	for _, tech := range technologies {
		if err := db.FirstOrCreate(&tech, models.Technology{Name: tech.Name}).Error; err != nil {
			log.Printf("Failed to create technology %s: %v", tech.Name, err)
		} else {
			log.Printf("Created/Updated technology: %s", tech.Name)
		}
	}

	// Seed Experiences
	experiences := []models.Experience{
		{
			Title:       "Full Stack Developer",
			CompanyName: "MaicoGroup",
			Icon:        "maico.png",
			IconBg:      "#ffffff",
			Date:        "Apr 2021 - Jan 2022",
			Points:      `["Engaged in ongoing communication with end users to gather feedback and requirements, ensuring project updates were tailored to user needs and preferences.", "Played a pivotal role in bug reporting and user support, actively addressing issues and ensuring a seamless user experience throughout software usage.", "Took charge of training new employees on software development processes and best practices, contributing to the growth and efficiency of the team.", "Participating in code reviews and providing constructive feedback to other developers."]`,
			Order:       1,
			IsActive:    true,
		},
		{
			Title:       "Full Stack Developer",
			CompanyName: "Freelance",
			Icon:        "freelance.jpg",
			IconBg:      "#ffffff",
			Date:        "Jan 2022 - Aug 2022",
			Points:      `["Communicated with clients to understand requirements for custom software projects", "Developed and implemented custom software solutions for two clients, resulting in tangible benefits and enhanced performance", "Ensured client satisfaction through effective follow-up and support."]`,
			Order:       2,
			IsActive:    true,
		},
		{
			Title:       "Software Engineer",
			CompanyName: "Levinci Co., Ltd",
			Icon:        "levinci.svg",
			IconBg:      "#ffffff",
			Date:        "Jan 2022 - Aug 2023",
			Points:      `["Managed two critical ERP projects for finance and employee management.", "Implemented comprehensive project management plans, ensuring successful goal achievement for numerous clients.", "Utilized modern ERP technologies for increased efficiency and accuracy.", "Delivered tailored solutions with a customer-focused approach."]`,
			Order:       3,
			IsActive:    true,
		},
		{
			Title:       "Software Engineer",
			CompanyName: "Terralogic",
			Icon:        "terralogic.jfif",
			IconBg:      "#ffffff",
			Date:        "Aug 2023 - Present",
			Points:      `["Facilitated effective communication with international colleagues to ensure project alignment.", "Developed and structured base code and module components, contributing to the project's architecture for optimal system performance.", "Conducted code reviews, providing constructive feedback and innovative solutions.", "Actively contributed to continuous improvement initiatives and fostered a collaborative team environment.", "Embraced learning opportunities to stay updated with the latest technologies and best practices."]`,
			Order:       4,
			IsActive:    true,
		},
	}

	for _, exp := range experiences {
		if err := db.FirstOrCreate(&exp, models.Experience{Title: exp.Title, CompanyName: exp.CompanyName}).Error; err != nil {
			log.Printf("Failed to create experience %s at %s: %v", exp.Title, exp.CompanyName, err)
		} else {
			log.Printf("Created/Updated experience: %s at %s", exp.Title, exp.CompanyName)
		}
	}

	// Seed Testimonials
	testimonials := []models.Testimonial{
		{
			Testimonial: "I thought it was impossible to make a website as beautiful as our product, but Moclaw proved me wrong.",
			Name:        "Sara Lee",
			Designation: "CFO",
			Company:     "Acme Co",
			Image:       "https://randomuser.me/api/portraits/women/4.jpg",
			Order:       1,
			IsActive:    true,
		},
		{
			Testimonial: "I've never met a web developer who truly cares about their clients' success like Moclaw does.",
			Name:        "Chris Brown",
			Designation: "COO",
			Company:     "DEF Corp",
			Image:       "https://randomuser.me/api/portraits/men/5.jpg",
			Order:       2,
			IsActive:    true,
		},
		{
			Testimonial: "After Moclaw optimized our website, our traffic increased by 50%. We can't thank them enough!",
			Name:        "Lisa Wang",
			Designation: "CTO",
			Company:     "456 Enterprises",
			Image:       "https://randomuser.me/api/portraits/women/6.jpg",
			Order:       3,
			IsActive:    true,
		},
	}

	for _, testimonial := range testimonials {
		if err := db.FirstOrCreate(&testimonial, models.Testimonial{Name: testimonial.Name}).Error; err != nil {
			log.Printf("Failed to create testimonial from %s: %v", testimonial.Name, err)
		} else {
			log.Printf("Created/Updated testimonial from: %s", testimonial.Name)
		}
	}

	// Seed Projects
	projects := []models.Project{
		{
			Name:           "Car Rent",
			Description:    "Web-based platform that allows users to search, book, and manage car rentals from various providers, providing a convenient and efficient solution for transportation needs.",
			Tags:           `[{"name":"react","color":"blue-text-gradient"},{"name":"mongodb","color":"green-text-gradient"},{"name":"tailwind","color":"pink-text-gradient"}]`,
			Image:          "carrent.png",
			SourceCodeLink: "https://github.com/",
			Order:          1,
			IsActive:       true,
		},
		{
			Name:           "Job IT",
			Description:    "Web application that enables users to search for job openings, view estimated salary ranges for positions, and locate available jobs based on their current location.",
			Tags:           `[{"name":"react","color":"blue-text-gradient"},{"name":"restapi","color":"green-text-gradient"},{"name":"scss","color":"pink-text-gradient"}]`,
			Image:          "jobit.png",
			SourceCodeLink: "https://github.com/",
			Order:          2,
			IsActive:       true,
		},
		{
			Name:           "Trip Guide",
			Description:    "A comprehensive travel booking platform that allows users to book flights, hotels, and rental cars, and offers curated recommendations for popular destinations.",
			Tags:           `[{"name":"nextjs","color":"blue-text-gradient"},{"name":"supabase","color":"green-text-gradient"},{"name":"css","color":"pink-text-gradient"}]`,
			Image:          "tripguide.png",
			SourceCodeLink: "https://github.com/",
			Order:          3,
			IsActive:       true,
		},
	}

	for _, project := range projects {
		if err := db.FirstOrCreate(&project, models.Project{Name: project.Name}).Error; err != nil {
			log.Printf("Failed to create project %s: %v", project.Name, err)
		} else {
			log.Printf("Created/Updated project: %s", project.Name)
		}
	}

	log.Println("Database seeding completed successfully!")
}
