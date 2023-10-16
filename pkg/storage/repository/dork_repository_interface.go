package repository

import "magma/pkg/dto"

type DorkRepository interface {
	SaveDork(dork dto.DorkDTO) error
	GetDorks() ([]dto.DorkDTO, error)
	UpdateDork(dork dto.DorkDTO) error
}
