package jobs

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type HTTP struct {
	db *sql.DB
}

func NewHTTP(db *sql.DB) HTTP {
	return HTTP{db: db}
}

// TODO: move to repository layer
type Job struct {
	ID             string     `json:"id"`
	State          string     `json:"state"` // queued | waiting | running | success | error
	ExpirationTime *time.Time `json:"expiration_time"`
	RunnerID       *string    `json:"runner_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// GET /jobs
func (h *HTTP) ListJobs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rows, err := h.db.Query("SELECT * FROM jobs")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer rows.Close()

	response := []Job{}
	for rows.Next() {
		var row Job
		err = rows.Scan(&row.ID, &row.State, &row.ExpirationTime, &row.RunnerID, &row.CreatedAt, &row.UpdatedAt)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		response = append(response, row)
	}

	serialized, err := json.Marshal(response)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(serialized)
}

// POST /jobs/enqueue
func (h *HTTP) EnqueueJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var job Job

	err := h.db.QueryRow(`INSERT INTO jobs (state)
			VALUES ('queued')
			RETURNING *`).Scan(&job.ID, &job.State, &job.ExpirationTime, &job.RunnerID, &job.CreatedAt, &job.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			// if no rows, return nil
			w.Write([]byte("null"))
			return
		} else {
			w.Write([]byte(err.Error()))
			return
		}
	}

	serialized, err := json.Marshal(job)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(serialized)
}

type ClaimJobBody struct {
	RunnerID string `json:"runner_id"`
}

// POST /jobs
// job: queued -> waiting
// This is for a runner to claim a job
func (h *HTTP) ClaimJob(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var b ClaimJobBody

	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var job Job
	err = h.db.QueryRow(`UPDATE jobs
			SET state = 'waiting', runner_id = $1
			WHERE id = (
				SELECT id FROM jobs
				WHERE state = 'queued'
				LIMIT 1
			)
			RETURNING *`, b.RunnerID).Scan(&job.ID, &job.State, &job.ExpirationTime, &job.RunnerID, &job.CreatedAt, &job.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			// if no rows, return nil
			w.Write([]byte("null"))
			return
		} else {
			w.Write([]byte(err.Error()))
			return
		}
	}

	serialized, err := json.Marshal(job)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(serialized)
}

// POST /jobs/:job_id/ack

// POST /jobs/:job_id/nack

// POST /jobs/:job_id/heartbeat
// job: running -> running
// timeout:     -> error

// POST /jobs/:job_id/complete
// job: running -> success

// POST /jobs/:job_id/error
// job: running -> error

// GET /jobs/:job_id/log
