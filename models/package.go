// models/package.go
package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Package struct {
    gorm.Model
    Name       string         `json:"name"`
    Data       string         `json:"data"`
    Duration   string         `json:"duration"`
    Price      float64        `json:"price"`
    Details    datatypes.JSON `json:"details"`
    Categories string         `json:"categories"`
}
