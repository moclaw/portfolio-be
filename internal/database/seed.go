package database

import (
	"log"
	"portfolio-be/internal/models"

	"gorm.io/gorm"
)

// Seed populates the database with initial data from Moclaw's CV
func Seed(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	// First, seed permissions
	if err := seedPermissions(db); err != nil {
		return err
	}

	// Then seed roles
	if err := seedRoles(db); err != nil {
		return err
	}

	// Seed Admin User
	if err := seedAdminUser(db); err != nil {
		return err
	}

	// Seed Services
	if err := seedServices(db); err != nil {
		return err
	}

	// Seed Technologies from CV
	if err := seedTechnologies(db); err != nil {
		return err
	}

	// Seed Experiences from CV
	if err := seedExperiences(db); err != nil {
		return err
	}

	// Seed Testimonials
	if err := seedTestimonials(db); err != nil {
		return err
	}

	// Seed Projects
	if err := seedProjects(db); err != nil {
		return err
	}

	log.Println("✓ Database seeding completed successfully")
	return nil
}

func seedAdminUser(db *gorm.DB) error {
	adminUser := models.User{
		Username:           "admin",
		Email:              "admin@example.com", // Change this in production
		Password:           "admin123",          // Will be hashed by BeforeCreate hook
		Role:               "admin",
		IsActive:           true,
		MustChangePassword: true, // Force password change on first login
	}

	var existingUser models.User
	result := db.Where("username = ?", adminUser.Username).First(&existingUser)
	if result.Error != nil {
		if err := db.Create(&adminUser).Error; err != nil {
			return err
		}
		log.Println("✓ Admin user created (password change required on first login)")

		var adminRole models.Role
		if err := db.Where("name = ?", "admin").First(&adminRole).Error; err == nil {
			adminUser.RoleID = &adminRole.ID
			db.Save(&adminUser)
		}
	}

	return nil
}

func seedServices(db *gorm.DB) error {
	// Services based on CV skills and experience
	services := []models.Service{
		{
			Title:    ".NET Backend Developer",
			Icon:     "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/dotnetcore/dotnetcore-original.svg",
			Order:    1,
			IsActive: true,
		},
		{
			Title:    "Full Stack Developer",
			Icon:     "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/react/react-original.svg",
			Order:    2,
			IsActive: true,
		},
		{
			Title:    "Cloud & DevOps Engineer",
			Icon:     "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/amazonwebservices/amazonwebservices-original.svg",
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
			return err
		}
	}

	return nil
}

func seedTechnologies(db *gorm.DB) error {
	// Technologies from CV - Top Skills and Certifications
	technologies := []models.Technology{
		// Programming Languages
		{Name: "C#", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/csharp/csharp-original.svg", Category: "programming", Order: 1, IsActive: true},
		{Name: "Go", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original.svg", Category: "programming", Order: 2, IsActive: true},
		{Name: "JavaScript", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/javascript/javascript-original.svg", Category: "programming", Order: 3, IsActive: true},
		{Name: "TypeScript", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/typescript/typescript-original.svg", Category: "programming", Order: 4, IsActive: true},

		// Frameworks
		{Name: ".NET", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/dotnetcore/dotnetcore-original.svg", Category: "framework", Order: 5, IsActive: true},
		{Name: "React", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/react/react-original.svg", Category: "framework", Order: 6, IsActive: true},

		// Databases - SQL/NoSQL from CV
		{Name: "PostgreSQL", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/postgresql/postgresql-original.svg", Category: "database", Order: 7, IsActive: true},
		{Name: "MySQL", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mysql/mysql-original.svg", Category: "database", Order: 8, IsActive: true},
		{Name: "MongoDB", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mongodb/mongodb-original.svg", Category: "database", Order: 9, IsActive: true},
		{Name: "Redis", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/redis/redis-original.svg", Category: "database", Order: 10, IsActive: true},
		{Name: "Cloud SQL", Icon: "https://www.vectorlogo.zone/logos/google_cloud/google_cloud-icon.svg", Category: "database", Order: 11, IsActive: true},

		// Cloud - GCP & AWS from CV
		{Name: "Google Cloud Platform", Icon: "https://www.vectorlogo.zone/logos/google_cloud/google_cloud-icon.svg", Category: "cloud", Order: 12, IsActive: true},
		{Name: "AWS", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/amazonwebservices/amazonwebservices-original.svg", Category: "cloud", Order: 13, IsActive: true},

		// DevOps & Container
		{Name: "Docker", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/docker/docker-original.svg", Category: "devops", Order: 14, IsActive: true},
		{Name: "Kubernetes", Icon: "https://www.vectorlogo.zone/logos/kubernetes/kubernetes-icon.svg", Category: "devops", Order: 15, IsActive: true},
		{Name: "Jenkins", Icon: "https://www.vectorlogo.zone/logos/jenkins/jenkins-icon.svg", Category: "devops", Order: 16, IsActive: true},
		{Name: "GitHub Actions", Icon: "https://www.vectorlogo.zone/logos/github/github-icon.svg", Category: "devops", Order: 17, IsActive: true},

		// Messaging & Architecture - Microservices from CV
		{Name: "RabbitMQ", Icon: "https://www.vectorlogo.zone/logos/rabbitmq/rabbitmq-icon.svg", Category: "messaging", Order: 18, IsActive: true},
		{Name: "Kafka", Icon: "https://www.vectorlogo.zone/logos/apache_kafka/apache_kafka-icon.svg", Category: "messaging", Order: 19, IsActive: true},

		// Tools
		{Name: "Git", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/git/git-original.svg", Category: "tools", Order: 20, IsActive: true},
		{Name: "Jira", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/jira/jira-original.svg", Category: "tools", Order: 21, IsActive: true},

		// OS
		{Name: "Linux", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/linux/linux-original.svg", Category: "os", Order: 22, IsActive: true},
		{Name: "Ubuntu", Icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/ubuntu/ubuntu-plain.svg", Category: "os", Order: 23, IsActive: true},
	}

	for _, tech := range technologies {
		if err := db.FirstOrCreate(&tech, models.Technology{Name: tech.Name}).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedExperiences(db *gorm.DB) error {
	// Experiences from CV - Moclaw Nguyen
	experiences := []models.Experience{
		{
			Title:       "Senior Software Engineer",
			CompanyName: "Keyloop",
			Icon:        "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/dotnetcore/dotnetcore-original.svg",
			IconBg:      "#1a1a2e",
			Date:        "October 2025 - Present",
			Points: `[
				"Currently working as Senior Software Engineer at Keyloop",
				"Contributing to automotive software solutions",
				"Working with cutting-edge technologies in the automotive industry"
			]`,
			Order:    1,
			IsActive: true,
		},
		{
			Title:       "Software Engineer",
			CompanyName: "Terralogic",
			Icon:        "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/azure/azure-original.svg",
			IconBg:      "#0078d4",
			Date:        "August 2023 - September 2025",
			Points: `[
				"Spearheaded a digital transformation project for an education platform spanning 11 countries",
				"Developed and refined three core modules, achieving 99.8% uptime while ensuring seamless integration with existing infrastructure",
				"Supported a platform with 100,000 users, facilitating over 1 million transactions per month while guaranteeing response times under 200ms",
				"Contributed from the early development phase, shaping solution architecture and optimizing system performance, resulting in a 25% reduction in infrastructure costs",
				"Collaborated with a team of over 100 professionals, including developers, designers, and product managers, to align business goals with cutting-edge technology",
				"Improved project execution efficiency by reducing development cycles by 15% through structured planning and streamlined workflows"
			]`,
			Order:    2,
			IsActive: true,
		},
		{
			Title:       "Software Engineer",
			CompanyName: "Levinci CO Ltd.",
			Icon:        "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/csharp/csharp-original.svg",
			IconBg:      "#68217a",
			Date:        "August 2022 - August 2023",
			Points: `[
				"Championed the development of two key projects: an ERP system and a financial management platform, streamlining company operations and reducing costs by 20%",
				"Enhanced workflow efficiency by automating 50% of manual tasks, thereby increasing process accuracy",
				"Managed a system actively utilized by over 100 internal and external users, maintaining 99.5% uptime to ensure seamless daily operations",
				"Implemented strategic optimizations that reduced financial reporting errors by 30%, leading to more reliable data-driven decisions"
			]`,
			Order:    3,
			IsActive: true,
		},
		{
			Title:       "Freelance Web Developer",
			CompanyName: "Freelancer",
			Icon:        "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nodejs/nodejs-original.svg",
			IconBg:      "#339933",
			Date:        "January 2022 - August 2022",
			Points: `[
				"Successfully developed and delivered three custom software projects for individual clients, overseeing end-to-end development from requirement analysis to deployment",
				"Engineered and implemented tailored solutions, ensuring high performance, security, and scalability",
				"Collaborated closely with clients to refine requirements, provide technical consultation, and integrate feedback for optimal user experience",
				"Facilitated seamless project handovers with comprehensive documentation, training sessions, and post-launch support",
				"Achieved 100% client satisfaction, with all projects successfully deployed and operating smoothly"
			]`,
			Order:    4,
			IsActive: true,
		},
		{
			Title:       "Software Engineer",
			CompanyName: "MAICO GROUP",
			Icon:        "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/docker/docker-original.svg",
			IconBg:      "#2496ed",
			Date:        "April 2021 - January 2022",
			Points: `[
				"Developed and delivered two major projects: CRM and Call Management, effectively meeting business requirements",
				"Managed and supported a system servicing over 10,000 customers, along with nearly 50 internal users, primarily company employees",
				"Engaged closely with end users to collect and implement over 100 feature requests, enhancing user satisfaction by 30%",
				"Reduced bug occurrence by 40%, ensuring system stability while receiving positive user feedback",
				"Utilized JIRA for 80% of project management and issue tracking, ensuring smooth collaboration and progress monitoring"
			]`,
			Order:    5,
			IsActive: true,
		},
	}

	for _, exp := range experiences {
		if err := db.FirstOrCreate(&exp, models.Experience{Title: exp.Title, CompanyName: exp.CompanyName}).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedTestimonials(db *gorm.DB) error {
	testimonials := []models.Testimonial{
		{
			Testimonial: "Moclaw is an exceptional developer who consistently delivers high-quality solutions. His expertise in .NET and microservices architecture has been invaluable to our projects.",
			Name:        "Tech Lead",
			Designation: "Senior Technical Lead",
			Company:     "Terralogic",
			Image:       "https://randomuser.me/api/portraits/men/32.jpg",
			Order:       1,
			IsActive:    true,
		},
		{
			Testimonial: "Working with Moclaw was a great experience. He delivered our ERP system on time and the quality exceeded our expectations. Highly recommended!",
			Name:        "Project Manager",
			Designation: "PM",
			Company:     "Levinci",
			Image:       "https://randomuser.me/api/portraits/women/44.jpg",
			Order:       2,
			IsActive:    true,
		},
		{
			Testimonial: "Moclaw's technical skills and problem-solving abilities are outstanding. He helped us achieve 99.8% uptime on our platform serving 100,000+ users.",
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
			return err
		}
	}

	return nil
}

func seedProjects(db *gorm.DB) error {
	projects := []models.Project{
		{
			Name:        "Education Platform",
			Description: "A digital transformation project for an education platform spanning 11 countries. Supports 100,000+ users with over 1 million transactions per month, achieving 99.8% uptime and response times under 200ms.",
			Tags: `[
				{"name": ".NET", "color": "blue-text-gradient"},
				{"name": "microservices", "color": "green-text-gradient"},
				{"name": "AWS", "color": "pink-text-gradient"}
			]`,
			Image:          "https://images.unsplash.com/photo-1503676260728-1c00da094a0b?w=800&h=600&fit=crop&crop=center",
			SourceCodeLink: "https://github.com/Moclaw",
			Order:          1,
			IsActive:       true,
		},
		{
			Name:        "ERP System",
			Description: "Enterprise Resource Planning system that streamlined company operations and reduced costs by 20%. Automated 50% of manual tasks and reduced financial reporting errors by 30%.",
			Tags: `[
				{"name": "C#", "color": "blue-text-gradient"},
				{"name": "SQL Server", "color": "green-text-gradient"},
				{"name": ".NET", "color": "pink-text-gradient"}
			]`,
			Image:          "https://images.unsplash.com/photo-1454165804606-c3d57bc86b40?w=800&h=600&fit=crop&crop=center",
			SourceCodeLink: "https://github.com/Moclaw",
			Order:          2,
			IsActive:       true,
		},
		{
			Name:        "CRM & Call Management",
			Description: "Customer Relationship Management and Call Management system servicing over 10,000 customers. Implemented 100+ feature requests and reduced bug occurrence by 40%.",
			Tags: `[
				{"name": ".NET", "color": "blue-text-gradient"},
				{"name": "SQL", "color": "green-text-gradient"},
				{"name": "React", "color": "pink-text-gradient"}
			]`,
			Image:          "https://images.unsplash.com/photo-1553877522-43269d4ea984?w=800&h=600&fit=crop&crop=center",
			SourceCodeLink: "https://github.com/Moclaw",
			Order:          3,
			IsActive:       true,
		},
		{
			Name:        "Portfolio Website",
			Description: "Personal portfolio website built with React, Three.js, and Go backend. Features 3D animations, admin panel, and AWS S3 integration for media storage.",
			Tags: `[
				{"name": "React", "color": "blue-text-gradient"},
				{"name": "Go", "color": "green-text-gradient"},
				{"name": "Three.js", "color": "pink-text-gradient"}
			]`,
			Image:          "https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=800&h=600&fit=crop&crop=center",
			SourceCodeLink: "https://github.com/Moclaw/portfolio-be",
			Order:          4,
			IsActive:       true,
		},
		{
			Name:        "NodeTL - Visual ETL Platform",
			Description: "A powerful, visual data mapping and transformation platform for building ETL pipelines and automating data workflows. Features drag-and-drop workflow designer, visual schema mapping, AI-assisted mapping, RBAC authentication, and Docker/Kubernetes deployment support.",
			Tags: `[
				{"name": "Go", "color": "blue-text-gradient"},
				{"name": "React", "color": "green-text-gradient"},
				{"name": "MongoDB", "color": "pink-text-gradient"}
			]`,
			Image:          "https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=800&h=600&fit=crop&crop=center",
			SourceCodeLink: "https://github.com/nodetl/nodetl",
			Order:          5,
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

// seedPermissions creates default permissions
func seedPermissions(db *gorm.DB) error {
	resources := []string{"users", "roles", "permissions", "projects", "technologies", "experiences", "testimonials", "contacts", "services", "uploads"}
	actions := []string{"create", "read", "update", "delete"}

	for _, resource := range resources {
		for _, action := range actions {
			permission := models.Permission{
				Name:        resource + ":" + action,
				Description: "Permission to " + action + " " + resource,
				Resource:    resource,
				Action:      action,
				IsActive:    true,
			}

			if err := db.FirstOrCreate(&permission, models.Permission{Name: permission.Name}).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// seedRoles creates default roles with permissions
func seedRoles(db *gorm.DB) error {
	// Create admin role with all permissions
	adminRole := models.Role{
		Name:        "admin",
		Description: "Administrator with full access",
		IsActive:    true,
	}

	if err := db.FirstOrCreate(&adminRole, models.Role{Name: adminRole.Name}).Error; err != nil {
		return err
	}

	var allPermissions []models.Permission
	if err := db.Find(&allPermissions).Error; err != nil {
		return err
	}

	if err := db.Model(&adminRole).Association("Permissions").Replace(allPermissions); err != nil {
		return err
	}

	// Create user role with only read permissions
	userRole := models.Role{
		Name:        "user",
		Description: "Regular user with read-only access",
		IsActive:    true,
	}

	if err := db.FirstOrCreate(&userRole, models.Role{Name: userRole.Name}).Error; err != nil {
		return err
	}

	var readPermissions []models.Permission
	if err := db.Where("action = ?", "read").Find(&readPermissions).Error; err != nil {
		return err
	}

	if err := db.Model(&userRole).Association("Permissions").Replace(readPermissions); err != nil {
		return err
	}

	// Create viewer role
	viewerRole := models.Role{
		Name:        "viewer",
		Description: "Viewer with limited read access",
		IsActive:    true,
	}

	if err := db.FirstOrCreate(&viewerRole, models.Role{Name: viewerRole.Name}).Error; err != nil {
		return err
	}

	viewerResources := []string{"projects", "technologies", "experiences", "testimonials", "services"}
	var viewerPermissions []models.Permission
	for _, resource := range viewerResources {
		var permissions []models.Permission
		if err := db.Where("resource = ? AND action = ?", resource, "read").Find(&permissions).Error; err != nil {
			return err
		}
		viewerPermissions = append(viewerPermissions, permissions...)
	}

	if err := db.Model(&viewerRole).Association("Permissions").Replace(viewerPermissions); err != nil {
		return err
	}

	return nil
}
