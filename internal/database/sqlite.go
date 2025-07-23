package database

import (
	"portfolio-be/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

func InitSQLite(databaseURL string) (*gorm.DB, error) {
	// Use the pure Go SQLite driver by specifying the driver name
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        databaseURL,
	}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Content{},
		&models.Upload{},
		&models.Resource{},
		&models.Experience{},
		&models.Service{},
		&models.Technology{},
		&models.Project{},
		&models.Testimonial{},
		&models.Contact{},
	)
}

// IsEmpty checks if the database has been seeded with initial data
func IsEmpty(db *gorm.DB) bool {
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	return userCount == 0
}

// Seed populates the database with initial data
func Seed(db *gorm.DB) error {
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
			return err
		}
	}

	// Seed Services with image links
	services := []models.Service{
		{
			Title:    "Full Stack Developer",
			Icon:     "https://images.unsplash.com/photo-1461749280684-dccba630e2f6?w=100&h=100&fit=crop&crop=center",
			Order:    1,
			IsActive: true,
		},
		{
			Title:    ".NET Developer",
			Icon:     "https://images.unsplash.com/photo-1555949963-aa79dcee981c?w=100&h=100&fit=crop&crop=center",
			Order:    2,
			IsActive: true,
		},
		{
			Title:    "Software Engineer",
			Icon:     "https://images.unsplash.com/photo-1517077304055-6e89abbf09b0?w=100&h=100&fit=crop&crop=center",
			Order:    3,
			IsActive: true,
		},
	}

	for _, service := range services {
		if err := db.FirstOrCreate(&service, models.Service{Title: service.Title}).Error; err != nil {
			return err
		}
	}

	// Seed Technologies with image links
	technologies := []models.Technology{
		{Name: "C#", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/csharp/csharp-original.svg", Category: "programming", Order: 1, IsActive: true},
		{Name: ".NET", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/dotnetcore/dotnetcore-original.svg", Category: "framework", Order: 2, IsActive: true},
		{Name: "Ubuntu", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/ubuntu/ubuntu-plain.svg", Category: "os", Order: 3, IsActive: true},
		{Name: "RabbitMQ", Icon: "https://www.vectorlogo.zone/logos/rabbitmq/rabbitmq-icon.svg", Category: "messaging", Order: 4, IsActive: true},
		{Name: "Redis", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/redis/redis-original.svg", Category: "database", Order: 5, IsActive: true},
		{Name: "Google Cloud", Icon: "https://www.vectorlogo.zone/logos/google_cloud/google_cloud-icon.svg", Category: "cloud", Order: 6, IsActive: true},
		{Name: "Jenkins", Icon: "https://www.vectorlogo.zone/logos/jenkins/jenkins-icon.svg", Category: "ci/cd", Order: 7, IsActive: true},
		{Name: "Jira", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/jira/jira-original.svg", Category: "project-management", Order: 8, IsActive: true},
		{Name: "Kubernetes", Icon: "https://www.vectorlogo.zone/logos/kubernetes/kubernetes-icon.svg", Category: "container", Order: 9, IsActive: true},
		{Name: "React", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/react/react-original.svg", Category: "frontend", Order: 10, IsActive: true},
		{Name: "Github Actions", Icon: "https://www.vectorlogo.zone/logos/github/github-icon.svg", Category: "ci/cd", Order: 11, IsActive: true},
		{Name: "Docker", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/docker/docker-original.svg", Category: "container", Order: 12, IsActive: true},
		{Name: "AWS", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/amazonwebservices/amazonwebservices-original.svg", Category: "cloud", Order: 13, IsActive: true},
		{Name: "Kafka", Icon: "https://www.vectorlogo.zone/logos/apache_kafka/apache_kafka-icon.svg", Category: "messaging", Order: 14, IsActive: true},
		{Name: "MongoDB", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mongodb/mongodb-original.svg", Category: "database", Order: 15, IsActive: true},
	}

	for _, tech := range technologies {
		if err := db.FirstOrCreate(&tech, models.Technology{Name: tech.Name}).Error; err != nil {
			return err
		}
	}

	// Seed Experiences with image links
	experiences := []models.Experience{
		{
			Title:       "Full Stack Developer",
			CompanyName: "MaicoGroup",
			Icon:        "https://images.unsplash.com/photo-1560472354-b33ff0c44a43?w=100&h=100&fit=crop&crop=center",
			IconBg:      "#ffffff",
			Date:        "Apr 2021 - Jan 2022",
			Points:      `["Engaged in ongoing communication with end users to gather feedback and requirements, ensuring project updates were tailored to user needs and preferences.", "Played a pivotal role in bug reporting and user support, actively addressing issues and ensuring a seamless user experience throughout software usage.", "Took charge of training new employees on software development processes and best practices, contributing to the growth and efficiency of the team.", "Participating in code reviews and providing constructive feedback to other developers."]`,
			Order:       1,
			IsActive:    true,
		},
		{
			Title:       "Full Stack Developer",
			CompanyName: "Freelance",
			Icon:        "https://images.unsplash.com/photo-1553877522-43269d4ea984?w=100&h=100&fit=crop&crop=center",
			IconBg:      "#ffffff",
			Date:        "Jan 2022 - Aug 2022",
			Points:      `["Communicated with clients to understand requirements for custom software projects", "Developed and implemented custom software solutions for two clients, resulting in tangible benefits and enhanced performance", "Ensured client satisfaction through effective follow-up and support."]`,
			Order:       2,
			IsActive:    true,
		},
		{
			Title:       "Software Engineer",
			CompanyName: "Levinci Co., Ltd",
			Icon:        "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&h=100&fit=crop&crop=center",
			IconBg:      "#ffffff",
			Date:        "Jan 2022 - Aug 2023",
			Points:      `["Managed two critical ERP projects for finance and employee management.", "Implemented comprehensive project management plans, ensuring successful goal achievement for numerous clients.", "Utilized modern ERP technologies for increased efficiency and accuracy.", "Delivered tailored solutions with a customer-focused approach."]`,
			Order:       3,
			IsActive:    true,
		},
		{
			Title:       "Software Engineer",
			CompanyName: "Terralogic",
			Icon:        "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&h=100&fit=crop&crop=center",
			IconBg:      "#ffffff",
			Date:        "Aug 2023 - Present",
			Points:      `["Facilitated effective communication with international colleagues to ensure project alignment.", "Developed and structured base code and module components, contributing to the project's architecture for optimal system performance.", "Conducted code reviews, providing constructive feedback and innovative solutions.", "Actively contributed to continuous improvement initiatives and fostered a collaborative team environment.", "Embraced learning opportunities to stay updated with the latest technologies and best practices."]`,
			Order:       4,
			IsActive:    true,
		},
	}

	for _, exp := range experiences {
		if err := db.FirstOrCreate(&exp, models.Experience{Title: exp.Title, CompanyName: exp.CompanyName}).Error; err != nil {
			return err
		}
	}

	// Seed Testimonials (already using image links)
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
			return err
		}
	}

	// Seed Projects with image links
	projects := []models.Project{
		{
			Name:           "Car Rent",
			Description:    "Web-based platform that allows users to search, book, and manage car rentals from various providers, providing a convenient and efficient solution for transportation needs.",
			Tags:           `[{"name":"react","color":"blue-text-gradient"},{"name":"mongodb","color":"green-text-gradient"},{"name":"tailwind","color":"pink-text-gradient"}]`,
			Image:          "https://images.unsplash.com/photo-1449824913935-59a10b8d2000?w=800&h=600&fit=crop&crop=center",
			SourceCodeLink: "https://github.com/",
			Order:          1,
			IsActive:       true,
		},
		{
			Name:           "Job IT",
			Description:    "Web application that enables users to search for job openings, view estimated salary ranges for positions, and locate available jobs based on their current location.",
			Tags:           `[{"name":"react","color":"blue-text-gradient"},{"name":"restapi","color":"green-text-gradient"},{"name":"scss","color":"pink-text-gradient"}]`,
			Image:          "https://images.unsplash.com/photo-1486312338219-ce68d2c6f44d?w=800&h=600&fit=crop&crop=center",
			SourceCodeLink: "https://github.com/",
			Order:          2,
			IsActive:       true,
		},
		{
			Name:           "Trip Guide",
			Description:    "A comprehensive travel booking platform that allows users to book flights, hotels, and rental cars, and offers curated recommendations for popular destinations.",
			Tags:           `[{"name":"nextjs","color":"blue-text-gradient"},{"name":"supabase","color":"green-text-gradient"},{"name":"css","color":"pink-text-gradient"}]`,
			Image:          "https://images.unsplash.com/photo-1488646953014-85cb44e25828?w=800&h=600&fit=crop&crop=center",
			SourceCodeLink: "https://github.com/",
			Order:          3,
			IsActive:       true,
		},
	}

	for _, project := range projects {
		if err := db.FirstOrCreate(&project, models.Project{Name: project.Name}).Error; err != nil {
			return err
		}
	}

	return nil
}
