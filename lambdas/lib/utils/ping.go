package utils

import (
	"fmt"
	"io"
	"net/http"

	"gorm.io/gorm"

	"github.com/kmesiab/equilibria/lambdas/lib/log"
)

func PingDatabase(globalDB *gorm.DB) error {

	sqlDB, err := globalDB.DB()

	if err != nil {
		return err
	}

	err = sqlDB.Ping()

	if err != nil {
		return err
	}

	return nil
}

func PingGoogle() error {

	log.New("Checking internet connectivity.....").Log()
	response, err := http.Get("https://www.google.com/")
	log.New("Connected to the internet...").AddError(err).Log()

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.New("Error closing response while checking internet connectivity").
				AddError(err).Log()
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return nil

}
