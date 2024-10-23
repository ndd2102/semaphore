package projects

import (
    "net/http"
	"os"
    "os/exec"
    "strings"
    "github.com/ansible-semaphore/semaphore/api/helpers"
)

// GetResticData lấy dữ liệu từ repository của Restic
func GetSnapshotData(w http.ResponseWriter, r *http.Request) {

	os.Setenv("RESTIC_REPOSITORY", "s3:https://restic.com")
    os.Setenv("RESTIC_PASSWORD", "123")
    os.Setenv("AWS_ACCESS_KEY_ID", "133")
    os.Setenv("AWS_SECRET_ACCESS_KEY", "123")
    // Ví dụ lệnh để lấy snapshots từ repository của Restic
    cmd := exec.Command("restic", "snapshots")

    // Thực thi lệnh
    output, err := cmd.CombinedOutput()
    if err != nil {
        helpers.WriteJSON(w, http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
        return
    }

    // Trả lại dữ liệu lệnh đã thực thi dưới dạng JSON
    helpers.WriteJSON(w, http.StatusOK, map[string]string{
        "data": strings.TrimSpace(string(output)),
    })
}
