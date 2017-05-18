package store

import (
	"encoding/json"
	"fmt"
	"github.com/mergermarket/gotools"
	"github.com/mergermarket/notifications-scheduler-service/models"
	"net/http"
	"time"
)

type queryable interface {
	GetByIntelIDAndProfileID(intelID, profileID string) (*models.OutboundMatch, error)
	GetByProfileID(profileID string) ([]*models.OutboundMatch, error)
	UnsentDigestsCount(before time.Time) (count int, err error)
	GetDeletedProfiles() (DeletedProfiles, error)
	GetWithdrawnContents() (WithdrawnContents, error)
}

type server struct {
	store  queryable
	logger tools.Logger
	statsd tools.StatsD
}

func newServer(logger tools.Logger, statsd tools.StatsD, store queryable) *server {
	server := &server{store: store, logger: logger, statsd: statsd}
	return server
}

type unsentResponse struct {
	NumberUnsent int
	Since        string
}

//ServeHTTP Serves HTTP server
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path == "/unsent" {
		s.unsent(w, r)

		return
	}

	if r.URL.Path == "/withdrawn-contents" {
		s.withdrawnContent(w, r)

		return
	}

	if r.URL.Path == "/deleted-profiles" {
		s.deleteProfiles(w, r)

		return
	}

	s.root(w, r)
}

func (s *server) root(w http.ResponseWriter, r *http.Request) {
	profileID := r.URL.Query().Get("profileId")
	intelID := r.URL.Query().Get("intelId")

	if profileID == "" {
		http.Error(w, "Please supply a profileId and intelId", http.StatusBadRequest)
		return
	}

	if intelID != "" {
		result, err := s.store.GetByIntelIDAndProfileID(intelID, profileID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Fatal error with database request %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, result.JSON())
	} else {
		result, err := s.store.GetByProfileID(profileID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Fatal error with database request %v", err), http.StatusInternalServerError)
			return
		}

		if len(result) == 0 {
			http.Error(w, fmt.Sprintf("No results found for profile %s", profileID), http.StatusNotFound)
			return
		}

		resultsAsJSON, _ := json.Marshal(result)

		fmt.Fprint(w, string(resultsAsJSON))
	}
}

func (s *server) unsent(w http.ResponseWriter, r *http.Request) {

	cutoff := time.Now().UTC()
	count, err := s.store.UnsentDigestsCount(cutoff)

	if err != nil {
		http.Error(w, fmt.Sprintf("Problem getting unsent count %s", err.Error()), http.StatusInternalServerError)
		return
	}

	res := unsentResponse{
		NumberUnsent: count,
		Since:        cutoff.Format(time.RFC822),
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *server) withdrawnContent(w http.ResponseWriter, r *http.Request) {
	result, err := s.store.GetWithdrawnContents()
	if err != nil {
		http.Error(w, fmt.Sprintf("Fatal error with database request %v", err), http.StatusInternalServerError)
		return
	}

	if result == nil {
		fmt.Fprint(w, "No withdrawn contents found")
		return
	}

	fmt.Fprint(w, result.JSON())
}

func (s *server) deleteProfiles(w http.ResponseWriter, r *http.Request) {
	result, err := s.store.GetDeletedProfiles()
	if err != nil {
		http.Error(w, fmt.Sprintf("Fatal error with database request %v", err), http.StatusInternalServerError)
		return
	}

	if result == nil {
		fmt.Fprint(w, "No deleted profiles found")
		return
	}

	fmt.Fprint(w, result.JSON())
}
