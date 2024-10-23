package sql

import (
    "github.com/ansible-semaphore/semaphore/db"
)

func (d *SqlDb) GetMinIOConfig(projectID int, configID int) (db.MinIOConfig, error) {
    var minioconfig db.MinIOConfig
    err := d.getObject(projectID, db.MinIOConfigProps, configID, &minioconfig)
    return minioconfig, err
}

func (d *SqlDb) GetMinIOConfigs(projectID int, params db.RetrieveQueryParams) (minioconfigs []db.MinIOConfig, err error) {
    err = d.getObjects(projectID, db.MinIOConfigProps, params, nil, &minioconfigs)
    return
}

func (d *SqlDb) CreateMinIOConfig(minioconfig db.MinIOConfig) (newMinio db.MinIOConfig, err error) {
    err = minioconfig.Validate()
    
    if err != nil {
        return db.MinIOConfig{}, err
    }

    insertID, err := d.insert(
        "id",
        "INSERT INTO project__minio_config (project_id, name, url, access_key, secret_key, bucket) VALUES (?, ?, ?, ?, ?, ?)",
        minioconfig.ProjectID,
        minioconfig.Name,
        minioconfig.URL,
        minioconfig.AccessKey,
        minioconfig.SecretKey,
        minioconfig.Bucket,
    )

    if err != nil {
        return db.MinIOConfig{}, err
    }

    newMinio = minioconfig
    minioconfig.ID = insertID
    return
}

func (d *SqlDb) UpdateMinIOConfig(minioconfig db.MinIOConfig) error {
    err := minioconfig.Validate()
    if err != nil {
        return err
    }

    _, err = d.exec(
        "UPDATE project__minio_config SET name=?, url=?, access_key=?, secret_key=?, bucket=? WHERE id=? AND project_id=?",
        minioconfig.Name,
        minioconfig.URL,
        minioconfig.AccessKey,
        minioconfig.SecretKey,
        minioconfig.Bucket,
        minioconfig.ID,
        minioconfig.ProjectID,
    )

    return err
}

func (d *SqlDb) DeleteMinIOConfig(projectID int, minioconfigID int) error {
    return d.deleteObject(projectID, db.MinIOConfigProps, minioconfigID)
}
