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
	bucket := cfg.S3Config.Bucket
	if bucket == "" {
		bucket = "portfolio-bucket"
	}

	// Initialize database based on configuration
	db, err := database.InitDatabase(cfg.DatabaseType, cfg.DatabaseURL)
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
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
	)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Starting database seeding...")

	// First, seed Permissions
	permissions := []models.Permission{
		{Name: "read", Resource: "users", Description: "View users"},
		{Name: "create", Resource: "users", Description: "Create users"},
		{Name: "update", Resource: "users", Description: "Update users"},
		{Name: "delete", Resource: "users", Description: "Delete users"},
		{Name: "read", Resource: "roles", Description: "View roles"},
		{Name: "create", Resource: "roles", Description: "Create roles"},
		{Name: "update", Resource: "roles", Description: "Update roles"},
		{Name: "delete", Resource: "roles", Description: "Delete roles"},
		{Name: "read", Resource: "permissions", Description: "View permissions"},
		{Name: "create", Resource: "permissions", Description: "Create permissions"},
		{Name: "update", Resource: "permissions", Description: "Update permissions"},
		{Name: "delete", Resource: "permissions", Description: "Delete permissions"},
		{Name: "read", Resource: "projects", Description: "View projects"},
		{Name: "create", Resource: "projects", Description: "Create projects"},
		{Name: "update", Resource: "projects", Description: "Update projects"},
		{Name: "delete", Resource: "projects", Description: "Delete projects"},
		{Name: "read", Resource: "experiences", Description: "View experiences"},
		{Name: "create", Resource: "experiences", Description: "Create experiences"},
		{Name: "update", Resource: "experiences", Description: "Update experiences"},
		{Name: "delete", Resource: "experiences", Description: "Delete experiences"},
		{Name: "read", Resource: "technologies", Description: "View technologies"},
		{Name: "create", Resource: "technologies", Description: "Create technologies"},
		{Name: "update", Resource: "technologies", Description: "Update technologies"},
		{Name: "delete", Resource: "technologies", Description: "Delete technologies"},
		{Name: "read", Resource: "services", Description: "View services"},
		{Name: "create", Resource: "services", Description: "Create services"},
		{Name: "update", Resource: "services", Description: "Update services"},
		{Name: "delete", Resource: "services", Description: "Delete services"},
		{Name: "read", Resource: "testimonials", Description: "View testimonials"},
		{Name: "create", Resource: "testimonials", Description: "Create testimonials"},
		{Name: "update", Resource: "testimonials", Description: "Update testimonials"},
		{Name: "delete", Resource: "testimonials", Description: "Delete testimonials"},
		{Name: "read", Resource: "contacts", Description: "View contacts"},
		{Name: "create", Resource: "contacts", Description: "Create contacts"},
		{Name: "update", Resource: "contacts", Description: "Update contacts"},
		{Name: "delete", Resource: "contacts", Description: "Delete contacts"},
	}

	for _, permission := range permissions {
		var existingPermission models.Permission
		result := db.Where("name = ? AND resource = ?", permission.Name, permission.Resource).First(&existingPermission)
		if result.Error != nil {
			if err := db.Create(&permission).Error; err != nil {
				log.Printf("Failed to create permission %s:%s: %v", permission.Resource, permission.Name, err)
			} else {
				log.Printf("✓ Created permission: %s:%s", permission.Resource, permission.Name)
			}
		}
	}

	// Then, seed Roles
	roles := []models.Role{
		{Name: "admin", Description: "Administrator with full access"},
		{Name: "user", Description: "Regular user with read-only access"},
	}

	var adminRole, userRole models.Role
	for _, role := range roles {
		var existingRole models.Role
		result := db.Where("name = ?", role.Name).First(&existingRole)
		if result.Error != nil {
			if err := db.Create(&role).Error; err != nil {
				log.Printf("Failed to create role %s: %v", role.Name, err)
			} else {
				log.Printf("✓ Created role: %s", role.Name)
			}
		}

		// Store references for permission assignment
		if role.Name == "admin" {
			db.Where("name = ?", "admin").First(&adminRole)
		} else if role.Name == "user" {
			db.Where("name = ?", "user").First(&userRole)
		}
	}

	// Assign all permissions to admin role
	if adminRole.ID != 0 {
		var allPermissions []models.Permission
		db.Find(&allPermissions)
		for _, permission := range allPermissions {
			var existingRolePermission models.RolePermission
			result := db.Where("role_id = ? AND permission_id = ?", adminRole.ID, permission.ID).First(&existingRolePermission)
			if result.Error != nil {
				rolePermission := models.RolePermission{
					RoleID:       adminRole.ID,
					PermissionID: permission.ID,
				}
				if err := db.Create(&rolePermission).Error; err != nil {
					log.Printf("Failed to assign permission %s:%s to admin role: %v", permission.Resource, permission.Name, err)
				}
			}
		}
		log.Println("✓ Assigned all permissions to admin role")
	}

	// Assign only read permissions to user role (for all resources)
	if userRole.ID != 0 {
		var readPermissions []models.Permission
		db.Where("name = ?", "read").Find(&readPermissions)
		for _, permission := range readPermissions {
			var existingRolePermission models.RolePermission
			result := db.Where("role_id = ? AND permission_id = ?", userRole.ID, permission.ID).First(&existingRolePermission)
			if result.Error != nil {
				rolePermission := models.RolePermission{
					RoleID:       userRole.ID,
					PermissionID: permission.ID,
				}
				if err := db.Create(&rolePermission).Error; err != nil {
					log.Printf("Failed to assign permission %s:%s to user role: %v", permission.Resource, permission.Name, err)
				}
			}
		}
		log.Println("✓ Assigned read permissions for all resources to user role")
	}

	// Now seed Admin User with proper role_id
	adminUser := models.User{
		Username: "admin",
		Email:    "admin@moclaw.dev",
		Password: "admin123", // This will be hashed automatically by BeforeCreate hook
		Role:     "admin",
		RoleID:   &adminRole.ID,
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
		// Update admin user to have admin role if not set
		if existingUser.RoleID == nil || *existingUser.RoleID != adminRole.ID {
			db.Model(&existingUser).Update("role_id", adminRole.ID)
			log.Println("✓ Updated admin user with admin role")
		}
	}

	// Seed Services
	services := []models.Service{
		{
			Title:    ".NET Backend Developer",
			Icon:     "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/dotnetcore/dotnetcore-original.svg",
			Order:    1,
			IsActive: true,
		},
		{
			Title:    "Microservices Architect",
			Icon:     "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/kubernetes/kubernetes-plain.svg",
			Order:    2,
			IsActive: true,
		},
		{
			Title:    "Cloud Engineer",
			Icon:     "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/amazonwebservices/amazonwebservices-original-wordmark.svg",
			Order:    3,
			IsActive: true,
		},
		{
			Title:    "Database Specialist",
			Icon:     "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/postgresql/postgresql-original.svg",
			Order:    4,
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

	// Seed Technologies - Based on CV skills
	technologies := []models.Technology{
		{Name: "C#", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/csharp/csharp-original.svg", Category: "programming", Order: 1, IsActive: true},
		{Name: ".NET Core", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/dotnetcore/dotnetcore-original.svg", Category: "framework", Order: 2, IsActive: true},
		{Name: "SQL Server", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/microsoftsqlserver/microsoftsqlserver-plain.svg", Category: "database", Order: 3, IsActive: true},
		{Name: "PostgreSQL", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/postgresql/postgresql-original.svg", Category: "database", Order: 4, IsActive: true},
		{Name: "MongoDB", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mongodb/mongodb-original.svg", Category: "database", Order: 5, IsActive: true},
		{Name: "Redis", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/redis/redis-original.svg", Category: "database", Order: 6, IsActive: true},
		{Name: "AWS", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/amazonwebservices/amazonwebservices-original-wordmark.svg", Category: "cloud", Order: 7, IsActive: true},
		{Name: "Google Cloud Platform", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/googlecloud/googlecloud-original.svg", Category: "cloud", Order: 8, IsActive: true},
		{Name: "Azure", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/azure/azure-original.svg", Category: "cloud", Order: 9, IsActive: true},
		{Name: "Docker", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/docker/docker-original.svg", Category: "container", Order: 10, IsActive: true},
		{Name: "Kubernetes", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/kubernetes/kubernetes-plain.svg", Category: "container", Order: 11, IsActive: true},
		{Name: "RabbitMQ", Icon: "https://www.vectorlogo.zone/logos/rabbitmq/rabbitmq-icon.svg", Category: "messaging", Order: 12, IsActive: true},
		{Name: "Kafka", Icon: "https://www.vectorlogo.zone/logos/apache_kafka/apache_kafka-icon.svg", Category: "messaging", Order: 13, IsActive: true},
		{Name: "React", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/react/react-original.svg", Category: "frontend", Order: 14, IsActive: true},
		{Name: "TypeScript", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/typescript/typescript-original.svg", Category: "programming", Order: 15, IsActive: true},
		{Name: "Go", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original-wordmark.svg", Category: "programming", Order: 16, IsActive: true},
		{Name: "Git", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/git/git-original.svg", Category: "tools", Order: 17, IsActive: true},
		{Name: "GitHub Actions", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/githubactions/githubactions-original.svg", Category: "ci/cd", Order: 18, IsActive: true},
		{Name: "Jenkins", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/jenkins/jenkins-original.svg", Category: "ci/cd", Order: 19, IsActive: true},
		{Name: "Jira", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/jira/jira-original.svg", Category: "tools", Order: 20, IsActive: true},
		{Name: "Linux", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/linux/linux-original.svg", Category: "os", Order: 21, IsActive: true},
		{Name: "Nginx", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nginx/nginx-original.svg", Category: "server", Order: 22, IsActive: true},
		{Name: "GraphQL", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/graphql/graphql-plain.svg", Category: "api", Order: 23, IsActive: true},
	}

	for _, tech := range technologies {
		if err := db.FirstOrCreate(&tech, models.Technology{Name: tech.Name}).Error; err != nil {
			log.Printf("Failed to create technology %s: %v", tech.Name, err)
		} else {
			log.Printf("Created/Updated technology: %s", tech.Name)
		}
	}

	// Seed Experiences - Based on CV profile.csv
	experiences := []models.Experience{
		{
			Title:       "Senior Software Engineer",
			CompanyName: "Keyloop",
			Icon:        "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/dotnetcore/dotnetcore-original.svg",
			IconBg:      "#1a1a2e",
			Date:        "October 2025 - Present",
			Points:      `["Currently working as Senior Software Engineer at Keyloop", "Contributing to automotive software solutions", "Working with cutting-edge technologies in the automotive industry"]`,
			Order:       1,
			IsActive:    true,
		},
		{
			Title:       "Software Engineer",
			CompanyName: "Terralogic",
			Icon:        "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/azure/azure-original.svg",
			IconBg:      "#0078d4",
			Date:        "August 2023 - September 2025",
			Points:      `["Spearheaded a digital transformation project for an education platform spanning 11 countries", "Developed and refined three core modules, achieving 99.8% uptime while ensuring seamless integration with existing infrastructure", "Supported a platform with 100,000 users, facilitating over 1 million transactions per month while guaranteeing response times under 200ms", "Contributed from the early development phase, shaping solution architecture and optimizing system performance, resulting in a 25% reduction in infrastructure costs", "Collaborated with a team of over 100 professionals, including developers, designers, and product managers, to align business goals with cutting-edge technology", "Improved project execution efficiency by reducing development cycles by 15% through structured planning and streamlined workflows"]`,
			Order:       2,
			IsActive:    true,
		},
		{
			Title:       "Software Engineer",
			CompanyName: "Levinci CO Ltd.",
			Icon:        "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/csharp/csharp-original.svg",
			IconBg:      "#68217a",
			Date:        "August 2022 - August 2023",
			Points:      `["Championed the development of two key projects: an ERP system and a financial management platform, streamlining company operations and reducing costs by 20%", "Enhanced workflow efficiency by automating 50% of manual tasks, thereby increasing process accuracy", "Managed a system actively utilized by over 100 internal and external users, maintaining 99.5% uptime to ensure seamless daily operations", "Implemented strategic optimizations that reduced financial reporting errors by 30%, leading to more reliable data-driven decisions"]`,
			Order:       3,
			IsActive:    true,
		},
		{
			Title:       "Freelance Web Developer",
			CompanyName: "Freelancer",
			Icon:        "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nodejs/nodejs-original.svg",
			IconBg:      "#339933",
			Date:        "January 2022 - August 2022",
			Points:      `["Successfully developed and delivered three custom software projects for individual clients, overseeing end-to-end development from requirement analysis to deployment", "Engineered and implemented tailored solutions, ensuring high performance, security, and scalability", "Collaborated closely with clients to refine requirements, provide technical consultation, and integrate feedback for optimal user experience", "Facilitated seamless project handovers with comprehensive documentation, training sessions, and post-launch support", "Achieved 100% client satisfaction, with all projects successfully deployed and operating smoothly"]`,
			Order:       4,
			IsActive:    true,
		},
		{
			Title:       "Software Engineer",
			CompanyName: "MAICO GROUP",
			Icon:        "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/docker/docker-original.svg",
			IconBg:      "#2496ed",
			Date:        "April 2021 - January 2022",
			Points:      `["Developed and delivered two major projects: CRM and Call Management, effectively meeting business requirements", "Managed and supported a system servicing over 10,000 customers, along with nearly 50 internal users, primarily company employees", "Engaged closely with end users to collect and implement over 100 feature requests, enhancing user satisfaction by 30%", "Reduced bug occurrence by 40%, ensuring system stability while receiving positive user feedback", "Utilized JIRA for 80% of project management and issue tracking, ensuring smooth collaboration and progress monitoring"]`,
			Order:       5,
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
			Testimonial: "Moclaw is an exceptional developer who consistently delivers high-quality solutions. His expertise in .NET and microservices architecture has been invaluable to our education platform serving 100,000+ users across 11 countries.",
			Name:        "Tech Lead",
			Designation: "Senior Technical Lead",
			Company:     "Terralogic",
			Image:       "https://randomuser.me/api/portraits/men/32.jpg",
			Order:       1,
			IsActive:    true,
		},
		{
			Testimonial: "Working with Moclaw was a great experience. He delivered our ERP system on time, automating 50% of manual tasks and reducing costs by 20%. The quality exceeded our expectations. Highly recommended!",
			Name:        "Project Manager",
			Designation: "PM",
			Company:     "Levinci",
			Image:       "https://randomuser.me/api/portraits/women/44.jpg",
			Order:       2,
			IsActive:    true,
		},
		{
			Testimonial: "Moclaw's technical skills and problem-solving abilities are outstanding. He helped us achieve 99.8% uptime on our platform and reduced infrastructure costs by 25%. A true professional!",
			Name:        "CTO",
			Designation: "Chief Technology Officer",
			Company:     "Education Platform",
			Image:       "https://randomuser.me/api/portraits/men/67.jpg",
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

	// Seed Projects - Based on actual work experience
	projects := []models.Project{
		{
			Name:           "Education Platform",
			Description:    "A digital transformation project for an education platform spanning 11 countries. Supports 100,000+ users with over 1 million transactions per month, achieving 99.8% uptime and response times under 200ms. Reduced infrastructure costs by 25%.",
			Tags:           `[{"name":".NET Core","color":"blue-text-gradient"},{"name":"Microservices","color":"green-text-gradient"},{"name":"AWS","color":"pink-text-gradient"}]`,
			Image:          "https://images.unsplash.com/photo-1503676260728-1c00da094a0b?w=800&h=600&fit=crop&crop=center",
			SourceCodeLink: "https://github.com/Moclaw",
			Order:          1,
			IsActive:       true,
		},
		{
			Name:           "ERP & Financial Management System",
			Description:    "Enterprise Resource Planning system that streamlined company operations and reduced costs by 20%. Automated 50% of manual tasks, reduced financial reporting errors by 30%, and maintained 99.5% uptime for 100+ users.",
			Tags:           `[{"name":"C#","color":"blue-text-gradient"},{"name":"SQL Server","color":"green-text-gradient"},{"name":".NET","color":"pink-text-gradient"}]`,
			Image:          "https://images.unsplash.com/photo-1454165804606-c3d57bc86b40?w=800&h=600&fit=crop&crop=center",
			SourceCodeLink: "https://github.com/Moclaw",
			Order:          2,
			IsActive:       true,
		},
		{
			Name:           "CRM & Call Management System",
			Description:    "Customer Relationship Management and Call Management system servicing over 10,000 customers and 50 internal users. Implemented 100+ feature requests, enhanced user satisfaction by 30%, and reduced bug occurrence by 40%.",
			Tags:           `[{"name":".NET","color":"blue-text-gradient"},{"name":"SQL Server","color":"green-text-gradient"},{"name":"JIRA","color":"pink-text-gradient"}]`,
			Image:          "https://images.unsplash.com/photo-1553877522-43269d4ea984?w=800&h=600&fit=crop&crop=center",
			SourceCodeLink: "https://github.com/Moclaw",
			Order:          3,
			IsActive:       true,
		},
		{
			Name:           "NodeTL - Visual ETL Platform",
			Description:    "A powerful, visual data mapping and transformation platform for building ETL pipelines and automating data workflows. Features drag-and-drop workflow designer, visual schema mapping, AI-assisted mapping, RBAC authentication, and Docker/Kubernetes deployment support.",
			Tags:           `[{"name":"Go","color":"blue-text-gradient"},{"name":"React","color":"green-text-gradient"},{"name":"MongoDB","color":"pink-text-gradient"}]`,
			Image:          "https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=800&h=600&fit=crop&crop=center",
			SourceCodeLink: "https://github.com/nodetl/nodetl",
			Order:          4,
			IsActive:       true,
		},
		{
			Name:           "Portfolio Website",
			Description:    "Personal portfolio website built with React, Three.js for 3D animations, and Go backend with PostgreSQL. Features admin panel, AWS S3 integration for media storage, Redis caching, and Docker deployment.",
			Tags:           `[{"name":"React","color":"blue-text-gradient"},{"name":"Go","color":"green-text-gradient"},{"name":"Three.js","color":"pink-text-gradient"}]`,
			Image:          "https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=800&h=600&fit=crop&crop=center",
			SourceCodeLink: "https://github.com/Moclaw/portfolio-be",
			Order:          5,
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

	// Seed Content (blog posts / case studies)
	contents := []models.Content{
		{
			Title:       "Designing Resilient Backend APIs",
			Description: "Lessons from building portfolio-scale APIs with Go and SQLite",
			Body: `Building resilient backend services starts with predictable contracts and well-defined fallbacks.
In this article I walk through the guardrails I rely on—structured logging, idempotent handlers, and
observability hooks—to keep personal projects running even when traffic spikes unexpectedly.`,
			Category: "engineering",
			Tags:     "go,backend,resilience",
			Status:   "published",
			ImageURL: "https://images.unsplash.com/photo-1503023345310-bd7c1de61c7d?w=1200&auto=format&fit=crop",
		},
		{
			Title:       "Improving Developer Experience with Automation",
			Description: "Automating testing, builds, and deployments for a solo portfolio project",
			Body: `Automation gives me headroom to focus on storytelling instead of plumbing.
Here is how I wired Air for hot reloads, swag for documentation, and Makefile targets to keep chores one command away.`,
			Category: "productivity",
			Tags:     "automation,devx,make",
			Status:   "published",
			ImageURL: "https://images.unsplash.com/photo-1483058712412-4245e9b90334?w=1200&auto=format&fit=crop",
		},
		{
			Title:       "From Sketch to Launch: Crafting Portfolio Case Studies",
			Description: "A framework for turning raw project notes into engaging portfolio stories",
			Body: `Great case studies share just enough context, highlight constraints, and celebrate outcomes.
This write-up details how I gather metrics, design visuals, and structure the narrative for the projects listed on my site.`,
			Category: "storytelling",
			Tags:     "portfolio,writing,case-study",
			Status:   "draft",
			ImageURL: "https://images.unsplash.com/photo-1529333166437-7750a6dd5a70?w=1200&auto=format&fit=crop",
		},
	}

	for _, content := range contents {
		seed := content
		if err := db.Where("title = ?", seed.Title).FirstOrCreate(&seed).Error; err != nil {
			log.Printf("Failed to create content %s: %v", seed.Title, err)
		} else {
			log.Printf("Created/Found content: %s", seed.Title)
		}
	}

	// Seed Upload metadata so resources can reference stable assets
	uploads := []models.Upload{
		{
			FileName:     "hero-banner.webp",
			OriginalName: "hero-banner.webp",
			FileSize:     245760,
			ContentType:  "image/webp",
			S3Key:        "resources/hero-banner.webp",
			S3Bucket:     bucket,
			URL:          "https://" + bucket + ".s3.amazonaws.com/resources/hero-banner.webp",
			IsActive:     true,
		},
		{
			FileName:     "profile-avatar.png",
			OriginalName: "profile-avatar.png",
			FileSize:     98304,
			ContentType:  "image/png",
			S3Key:        "resources/profile-avatar.png",
			S3Bucket:     bucket,
			URL:          "https://" + bucket + ".s3.amazonaws.com/resources/profile-avatar.png",
			IsActive:     true,
		},
		{
			FileName:     "moclaw-resume.pdf",
			OriginalName: "moclaw-resume.pdf",
			FileSize:     524288,
			ContentType:  "application/pdf",
			S3Key:        "documents/moclaw-resume.pdf",
			S3Bucket:     bucket,
			URL:          "https://" + bucket + ".s3.amazonaws.com/documents/moclaw-resume.pdf",
			IsActive:     true,
		},
	}

	uploadIDs := make(map[string]uint)
	for _, upload := range uploads {
		seed := upload
		if err := db.Where("s3_key = ?", seed.S3Key).FirstOrCreate(&seed).Error; err != nil {
			log.Printf("Failed to create upload %s: %v", seed.S3Key, err)
		} else {
			uploadIDs[seed.S3Key] = seed.ID
			log.Printf("Created/Found upload: %s", seed.S3Key)
		}
	}

	// Seed Resources referencing the uploads above
	type resourceSeed struct {
		Record    models.Resource
		UploadKey string
	}

	resourceSeeds := []resourceSeed{
		{
			Record: models.Resource{
				Name:        "Hero Banner",
				Description: "Primary hero background used on the landing page",
				Type:        models.ResourceTypeImage,
				Category:    "hero",
				Tags:        "portfolio,hero,landing",
				Alt:         "Code editor screenshot on gradient background",
				IsPublic:    true,
				IsActive:    true,
			},
			UploadKey: "resources/hero-banner.webp",
		},
		{
			Record: models.Resource{
				Name:        "Profile Avatar",
				Description: "Square avatar used for testimonials and about section",
				Type:        models.ResourceTypeImage,
				Category:    "profile",
				Tags:        "avatar,brand",
				Alt:         "Moclaw avatar illustration",
				IsPublic:    true,
				IsActive:    true,
			},
			UploadKey: "resources/profile-avatar.png",
		},
		{
			Record: models.Resource{
				Name:        "Resume PDF",
				Description: "Latest résumé shared via the portfolio",
				Type:        models.ResourceTypeDocument,
				Category:    "documents",
				Tags:        "resume,cv",
				Alt:         "Downloadable Moclaw resume",
				IsPublic:    true,
				IsActive:    true,
			},
			UploadKey: "documents/moclaw-resume.pdf",
		},
	}

	for _, seed := range resourceSeeds {
		uploadID, ok := uploadIDs[seed.UploadKey]
		if !ok {
			log.Printf("Skipping resource %s because upload %s was not created", seed.Record.Name, seed.UploadKey)
			continue
		}

		record := seed.Record
		record.UploadID = uploadID
		if err := db.Where("name = ?", record.Name).FirstOrCreate(&record).Error; err != nil {
			log.Printf("Failed to create resource %s: %v", record.Name, err)
		} else {
			log.Printf("Created/Found resource: %s", record.Name)
		}
	}

	// Seed sample contacts for dashboard/testing
	contacts := []models.Contact{
		{
			Name:    "Jane Product",
			Email:   "jane.product@example.com",
			Subject: "Interested in a consultation",
			Message: "Hi! I'm exploring a freelance engagement and would love to chat about availability.",
			Status:  "unread",
		},
		{
			Name:    "Startup Founder",
			Email:   "founder@launchpad.dev",
			Subject: "Great work on Car Rent",
			Message: "Your Car Rent case study resonated with our roadmap. Could we schedule a call next week?",
			Status:  "read",
		},
	}

	for _, contact := range contacts {
		seed := contact
		if err := db.Where("email = ? AND subject = ?", seed.Email, seed.Subject).FirstOrCreate(&seed).Error; err != nil {
			log.Printf("Failed to create contact from %s: %v", seed.Email, err)
		} else {
			log.Printf("Created/Found contact lead: %s", seed.Email)
		}
	}

	log.Println("Database seeding completed successfully!")
}
