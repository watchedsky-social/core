package zones

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/paulmach/orb/geojson"
	"github.com/watchedsky-social/core/internal/database/models"
	"github.com/watchedsky-social/core/internal/database/query"
	"github.com/watchedsky-social/core/internal/utils"
)

func CollectMissingZones(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)
	failures := map[string]error{}
	for scanner.Scan() {
		id := scanner.Text()
		if err := collectZone(ctx, id); err != nil {
			failures[id] = err
		}
	}

	if scanner.Err() != nil {
		failures["scanner"] = scanner.Err()
	}

	if len(failures) > 0 {
		bigErrs := make([]error, 0, len(failures))
		for id, e := range failures {
			bigErrs = append(bigErrs, fmt.Errorf("id: %s, err: %w", id, e))
		}

		return errors.Join(bigErrs...)
	}

	return nil
}

func collectZone(ctx context.Context, id string) error {
	dao := query.Zone.WithContext(ctx)
	log.Printf("Load zone %s", id)
	resp, err := http.Get(id)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("wanted 200, got %d", resp.StatusCode)
	}

	var f geojson.Feature
	if err = json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return err
	}

	z := &models.Zone{
		ID:     f.Properties.MustString("@id"),
		Name:   f.Properties.MustString("name"),
		Type:   f.Properties.MustString("type"),
		State:  utils.Ref(f.Properties.MustString("state")),
		Border: utils.Ref(models.NewGenericGeometry(f.Geometry)),
	}
	return dao.Create(z)
}
