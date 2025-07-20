package services

import (
	"portfolio-be/internal/models"
	"portfolio-be/internal/repository"
)

type TestimonialService interface {
	CreateTestimonial(request *models.TestimonialRequest) (*models.TestimonialResponse, error)
	GetAllTestimonials() ([]models.TestimonialResponse, error)
	GetTestimonialByID(id uint) (*models.TestimonialResponse, error)
	UpdateTestimonial(id uint, request *models.TestimonialRequest) (*models.TestimonialResponse, error)
	DeleteTestimonial(id uint) error
	GetActiveTestimonials() ([]models.TestimonialResponse, error)
	GetTestimonialsCount() (int64, error)
	UpdateTestimonialsOrder(items []OrderItem) error
}

type testimonialService struct {
	testimonialRepo repository.TestimonialRepository
}

func NewTestimonialService(testimonialRepo repository.TestimonialRepository) TestimonialService {
	return &testimonialService{testimonialRepo: testimonialRepo}
}

func (s *testimonialService) CreateTestimonial(request *models.TestimonialRequest) (*models.TestimonialResponse, error) {
	testimonial := &models.Testimonial{
		Testimonial: request.Testimonial,
		Name:        request.Name,
		Designation: request.Designation,
		Company:     request.Company,
		Image:       request.Image,
		Order:       request.Order,
		IsActive:    request.IsActive,
	}

	err := s.testimonialRepo.Create(testimonial)
	if err != nil {
		return nil, err
	}

	response := testimonial.ToResponse()
	return &response, nil
}

func (s *testimonialService) GetAllTestimonials() ([]models.TestimonialResponse, error) {
	testimonials, err := s.testimonialRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []models.TestimonialResponse
	for _, testimonial := range testimonials {
		responses = append(responses, testimonial.ToResponse())
	}

	return responses, nil
}

func (s *testimonialService) GetTestimonialByID(id uint) (*models.TestimonialResponse, error) {
	testimonial, err := s.testimonialRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := testimonial.ToResponse()
	return &response, nil
}

func (s *testimonialService) UpdateTestimonial(id uint, request *models.TestimonialRequest) (*models.TestimonialResponse, error) {
	testimonial, err := s.testimonialRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	testimonial.Testimonial = request.Testimonial
	testimonial.Name = request.Name
	testimonial.Designation = request.Designation
	testimonial.Company = request.Company
	testimonial.Image = request.Image
	testimonial.Order = request.Order
	testimonial.IsActive = request.IsActive

	err = s.testimonialRepo.Update(testimonial)
	if err != nil {
		return nil, err
	}

	response := testimonial.ToResponse()
	return &response, nil
}

func (s *testimonialService) DeleteTestimonial(id uint) error {
	return s.testimonialRepo.Delete(id)
}

func (s *testimonialService) GetActiveTestimonials() ([]models.TestimonialResponse, error) {
	testimonials, err := s.testimonialRepo.GetActive()
	if err != nil {
		return nil, err
	}

	var responses []models.TestimonialResponse
	for _, testimonial := range testimonials {
		responses = append(responses, testimonial.ToResponse())
	}

	return responses, nil
}

func (s *testimonialService) GetTestimonialsCount() (int64, error) {
	return s.testimonialRepo.GetCount()
}

func (s *testimonialService) UpdateTestimonialsOrder(items []OrderItem) error {
	for _, item := range items {
		testimonial, err := s.testimonialRepo.GetByID(item.ID)
		if err != nil {
			return err
		}

		testimonial.Order = item.Order
		err = s.testimonialRepo.Update(testimonial)
		if err != nil {
			return err
		}
	}
	return nil
}
