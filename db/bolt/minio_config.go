package bolt

import (
    "github.com/ansible-semaphore/semaphore/db"
)

func (d *BoltDb) GetMinIOConfig(projectID int, minioconfigID int) (minioconfig db.MinIOConfig, err error) {
    err = d.getObject(projectID, db.MinIOConfigProps, intObjectID(minioconfigID), &minioconfig)
	if err != nil {
		return
	}
    return
}

func (d *BoltDb) GetMinIOConfigs(projectID int, params db.RetrieveQueryParams) (minioconfigs []db.MinIOConfig, err error) {
    err = d.getObjects(projectID, db.MinIOConfigProps, params, nil, &minioconfigs)
    return
}

func (d *BoltDb) CreateMinIOConfig(minioconfig db.MinIOConfig) (db.MinIOConfig, error) {
    err := minioconfig.Validate()
    if err != nil {
        return db.MinIOConfig{}, err
    }

    newMinio, err := d.createObject(minioconfig.ProjectID, db.MinIOConfigProps, minioconfig)
    return newMinio.(db.MinIOConfig), nil
}

func (d *BoltDb) UpdateMinIOConfig(minioconfig db.MinIOConfig) error {
    err := minioconfig.Validate()
    if err != nil {
        return err
    }

    return d.updateObject(minioconfig.ProjectID, db.MinIOConfigProps, minioconfig)
}

func (d *BoltDb) DeleteMinIOConfig(projectID int, minioconfigID int) error {
    return d.deleteObject(projectID, db.MinIOConfigProps, intObjectID(minioconfigID), nil)
}
