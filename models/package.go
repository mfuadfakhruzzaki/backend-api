package models

import (
	"time"

	"gorm.io/datatypes"
)

type Package struct {
    ID          uint           `gorm:"primarykey" json:"id"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   *time.Time     `json:"deleted_at,omitempty"`

    Name        string         `json:"name"`
    Data        string         `json:"data"`
    Duration    string         `json:"duration"`
    Price       float64        `json:"price"`
    Details     datatypes.JSON `json:"details" swaggertype:"string"`  // Override to string
    Categories  string         `json:"categories"`
}
