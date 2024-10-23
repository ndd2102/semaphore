package projects

import (
    "net/http"
    "errors"

    "github.com/ansible-semaphore/semaphore/api/helpers"
    "github.com/ansible-semaphore/semaphore/db"
    "github.com/gorilla/context"
)

// MinIOConfigMiddleware ensures that the MinIO configuration exists and loads it into the context
func MinIOConfigMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        project := context.Get(r, "project").(db.Project)
        minioconfigID, err := helpers.GetIntParam("minioconfig_id", w, r)
        if err != nil {
            helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
                "error": "Invalid or missing 'minioconfig_id' parameter",
            })
            return
        }

        minioconfig, err := helpers.Store(r).GetMinIOConfig(project.ID, minioconfigID)
        if err != nil {
            helpers.WriteError(w, err)
            return
        }

        context.Set(r, "minioconfig", minioconfig)
        next.ServeHTTP(w, r)
    })
}

// GetMinIOConfigs returns a list of MinIO configurations in the project
func GetMinIOConfigs(w http.ResponseWriter, r *http.Request) {
    if minio := context.Get(r, "minioconfig"); minio != nil {
        helpers.WriteJSON(w, http.StatusOK, minio.(db.MinIOConfig))
        return
    }

    project := context.Get(r, "project").(db.Project)

	minios, err := helpers.Store(r).GetMinIOConfigs(project.ID, helpers.QueryParams(r.URL))

    if err != nil {
		helpers.WriteError(w, err)
		return
	}

    helpers.WriteJSON(w, http.StatusOK, minios)
}

// AddMinIOConfig creates a new MinIO configuration
func AddMinIOConfig(w http.ResponseWriter, r *http.Request) {
    project := context.Get(r, "project").(db.Project)

    var minioconfig db.MinIOConfig

    if !helpers.Bind(w, r, &minioconfig) {
        return
    }

    if minioconfig.ProjectID != project.ID {
        helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
            "error": "Project ID in body and URL must be the same",
        })
    }

    newMinio, err := helpers.Store(r).CreateMinIOConfig(minioconfig)
    if err != nil {
        helpers.WriteError(w, err)
        return
    }

    helpers.EventLog(r, helpers.EventLogCreate, helpers.EventLogItem{
        UserID:      helpers.UserFromContext(r).ID,
        ProjectID:   newMinio.ProjectID,
        ObjectType:  db.EventMinIOConfig,
        ObjectID:    newMinio.ID,
        Description: "MinIO configuration created",
    })

    // Remove secret_key before returning
    newMinio.SecretKey = ""

    helpers.WriteJSON(w, http.StatusCreated, newMinio)
}

// UpdateMinIOConfig updates an existing MinIO configuration
func UpdateMinIOConfig(w http.ResponseWriter, r *http.Request) {
    oldMinio := context.Get(r, "minioconfig").(db.MinIOConfig)
    var minioconfig db.MinIOConfig

    if !helpers.Bind(w, r, &minioconfig) {
        return
    }

    if minioconfig.ID != oldMinio.ID {
        helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
            "error": "Minio ID in body and URL must be the same",
        })
        return
    }

    if minioconfig.ProjectID != oldMinio.ProjectID {
        helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
            "error": "Project ID in body and URL must be the same",
        })
        return
    }


    if err := helpers.Store(r).UpdateMinIOConfig(minioconfig); err != nil {
        helpers.WriteError(w, err)
        return
    }

    helpers.EventLog(r, helpers.EventLogUpdate, helpers.EventLogItem{
        UserID:      helpers.UserFromContext(r).ID,
        ProjectID:   oldMinio.ProjectID,
        ObjectType:  db.EventMinIOConfig,
        ObjectID:    oldMinio.ID,
        Description: "MinIO configuration updated",
    })

    w.WriteHeader(http.StatusNoContent)
}

// RemoveMinIOConfig deletes a MinIO configuration
func RemoveMinIOConfig(w http.ResponseWriter, r *http.Request) {
    minioconfig := context.Get(r, "minioconfig").(db.MinIOConfig)

    var err error

    err = helpers.Store(r).DeleteMinIOConfig(minioconfig.ProjectID, minioconfig.ID)
	if errors.Is(err, db.ErrInvalidOperation) {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": "Minio is in use by one or more templates",
			"inUse": true,
		})
		return
	}

    if err != nil {
		helpers.WriteError(w, err)
		return
	}
    
    helpers.EventLog(r, helpers.EventLogDelete, helpers.EventLogItem{
        UserID:      helpers.UserFromContext(r).ID,
        ProjectID:   minioconfig.ProjectID,
        ObjectType:  db.EventMinIOConfig,
        ObjectID:    minioconfig.ID,
        Description: "MinIO configuration deleted",
    })

    w.WriteHeader(http.StatusNoContent)
}
