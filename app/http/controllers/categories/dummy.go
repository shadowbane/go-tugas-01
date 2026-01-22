package categories

import (
	"github.com/shadowbane/go-tugas-01/app/models"
	"sync"
)

var (
	categoryData []models.Category
	mutex        sync.RWMutex
)
