package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/shobhit-Creator/IPLocationService/internal/controllers"
	"github.com/shobhit-Creator/IPLocationService/internal/models"
	"github.com/shobhit-Creator/IPLocationService/internal/utils"
	"github.com/shobhit-Creator/IPLocationService/internal/workerpool"
)

func LocationHandler(wp *workerpool.WorkerPool) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		ipAddress := utils.GetIPAddress(r)
		if(ipAddress == ""){
			http.Error(w, "Ip address not found", http.StatusBadRequest)
			return
		}
		resultChan := make(chan models.ApiResult, 1)
		
		var wg sync.WaitGroup
		wg.Add(1)

		err := wp.Submit(func() {
			defer wg.Done()
			controllers.GetLocation(ipAddress, resultChan)
	    })
		if err != nil {
			http.Error(w, "Worker pool is full", http.StatusServiceUnavailable)
			return
		}
		
		go func() {
			wg.Wait()
			close(resultChan)
		}()

		select {
			case result := <-resultChan:
				if result.Error != nil {
					http.Error(w, result.Error.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]string{"result": result.Result}); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			case <-time.After(time.Second * 300):
				http.Error(w, "Request timed out", int(http.StatusGatewayTimeout))
		}
	}
}