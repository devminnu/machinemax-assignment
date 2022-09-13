package register

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/devminnu/assignment/common/utils"
	"github.com/devminnu/assignment/internal/app/model"
)

func (r *registration) Register(ctx context.Context) (registeredDevEUIDs []string, err error) {
	wait := make(chan bool, 10)
	wg := new(sync.WaitGroup)
	sm := new(sync.Map)
	// generate 100 unique ids
	for i := 0; i < 100; i++ {
		generateAndStoreUniqueId(sm)
	}
	sm.Range(func(key, value any) bool {
		select {
		case <-ctx.Done():
			return false
		default:
			wait <- true
			wg.Add(1)
			go func() {
				defer func() {
					<-wait
				}()
				defer wg.Done()
				r.doRegistration(ctx, key.(string), sm)
			}()
			return true
		}
	})
	wg.Wait()
	sm.Range(func(key, value any) bool {
		if value.(bool) {
			registeredDevEUIDs = append(registeredDevEUIDs, key.(string))
		}
		return true
	})

	return
}

func (r *registration) doRegistration(ctx context.Context, devEUI string, sm *sync.Map) {
	resp, err := r.httpDo(ctx, devEUI, sm)
	if err != nil {
		return
	}
	if resp.StatusCode == http.StatusOK {
		sm.Store(devEUI, true)
		log.Print("success")
		return
	}
	if resp.StatusCode == http.StatusUnprocessableEntity {
		r.doRegistration(ctx, generateAndStoreUniqueId(sm), sm)
		return
	}
	log.Print("unknown error")
}

func (r *registration) httpDo(ctx context.Context, devEUI string, sm *sync.Map) (resp *http.Response, err error) {
	registrationRequest := model.RegistrationRequest{DevEUI: devEUI}
	body, err := json.Marshal(registrationRequest)
	if err != nil {
		log.Print(err)
		return
	}
	request, err := http.NewRequest(http.MethodPost, utils.GetRegistrationURL(), bytes.NewReader(body))
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Print("registration error", err)
		return
	}
	return http.DefaultClient.Do(request)
}

func generateAndStoreUniqueId(sm *sync.Map) string {
	for {
		devEUID, err := utils.GenerateHexString(16)
		if err != nil {
			log.Print("error generating dev euid")
			continue
		}
		if _, ok := sm.Load(devEUID); ok {
			log.Print("deveui already exists")
			continue
		}
		sm.Store(devEUID, false)
		return devEUID
	}
}
