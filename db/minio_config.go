package db

import (
    "errors"
)

// MinIOConfig stores configuration information for MinIO
type MinIOConfig struct {
    ID        int    `db:"id" json:"id"`
    ProjectID int    `db:"project_id" json:"project_id"`
    Name      string `db:"name" json:"name" binding:"required"`
    URL       string `db:"url" json:"url" binding:"required"`
    AccessKey string `db:"access_key" json:"access_key" binding:"required"`
    SecretKey string `db:"secret_key" json:"secret_key" binding:"required"`
    Bucket    string `db:"bucket" json:"bucket" binding:"required"` // Added Bucket field
}

// Validate checks if the MinIOConfig fields are valid
func (m MinIOConfig) Validate() error {
    if m.Name == "" {
        return errors.New("MinIO configuration name can't be empty")
    }
    if m.URL == "" {
        return errors.New("MinIO URL can't be empty")
    }
    if m.AccessKey == "" {
        return errors.New("MinIO Access Key can't be empty")
    }
    if m.SecretKey == "" {
        return errors.New("MinIO Secret Key can't be empty")
    }
    if m.Bucket == "" {
        return errors.New("MinIO Bucket can't be empty")
    }
    return nil
}
