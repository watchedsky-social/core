package alerts

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/watchedsky-social/core/internal/config"
	"github.com/watchedsky-social/core/internal/database/models"
	"github.com/watchedsky-social/core/internal/database/query"
	"github.com/watchedsky-social/core/internal/utils"
	"gorm.io/gorm"
)

const alertURL = "https://api.weather.gov/alerts?status=actual&urgency=Immediate,Expected,Future,Unknown"

type featureReference struct {
	ID string `mapstructure:"identifier"`
}

type fcPagination struct {
	NextLink string `mapstructure:"next"`
}

func CollectNewAlerts(ctx context.Context) error {
	log.Println("Getting new alerts from NWS...")

	var cfg config.AlertConfig
	if err := envconfig.Process("collector", &cfg); err != nil {
		return err
	}

	id, err := getLatestAlertID(ctx)
	if err != nil {
		return err
	}

	return loadLatestAlerts(ctx, id, cfg)
}

func getLatestAlertID(ctx context.Context) (string, error) {
	alert := query.Alert
	res, err := alert.WithContext(ctx).Select(alert.ID).Order(alert.Sent.Desc()).First()
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}

		return "", nil
	}

	return res.ID, nil
}

func loadLatestAlerts(ctx context.Context, latestID string, cfg config.AlertConfig) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
				return
			}
			err = nil
		}
	}()
	url := alertURL

	maxAge := cfg.MaxAge
	if latestID != "" {
		log.Println("since we have alerts already, ignore max age and fill everything in")
		maxAge = time.Duration(math.MaxInt64)
	}

	for url != "" {
		dbAlerts := []*models.Alert{}

		var resp *http.Response

		resp, err = http.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		url = ""

		var fc *geojson.FeatureCollection
		var data []byte

		if data, err = io.ReadAll(resp.Body); err != nil {
			return
		}
		if fc, err = geojson.UnmarshalFeatureCollection(data); err != nil {
			return
		}

		foundTheEnd := false
		for _, feature := range fc.Features {
			id := feature.Properties.MustString("id")
			if id == latestID {
				foundTheEnd = true
				break
			}
			refIDs := models.StringSlice{}
			refs, ok := feature.Properties["references"].([]any)

			var border *models.Geometry
			border, err = resolveGeometry(ctx, feature)
			if err != nil {
				return
			}

			if ok {
				for _, ref := range refs {
					var r featureReference
					if err = mapstructure.Decode(ref, &r); err != nil {
						continue
					}

					refIDs = append(refIDs, r.ID)
				}
			}
			dba := &models.Alert{
				ID:           id,
				AreaDesc:     feature.Properties.MustString("areaDesc", ""),
				Headline:     feature.Properties.MustString("headline", ""),
				Description:  feature.Properties.MustString("description", ""),
				Severity:     getOptional(feature, "severity"),
				Certainty:    getOptional(feature, "certainty"),
				Urgency:      getOptional(feature, "urgency"),
				Event:        getOptional(feature, "event"),
				Sent:         mustTime(feature, "sent"),
				Effective:    mustTime(feature, "effective"),
				Onset:        optionalTime(feature, "onset"),
				Expires:      optionalTime(feature, "expires"),
				Ends:         optionalTime(feature, "ends"),
				MessageType:  getOptional(feature, "messageType"),
				ReferenceIds: &refIDs,
				Border:       border,
			}

			if time.Since(dba.Sent) > maxAge {
				log.Printf("alert %s is older than %s, do not fetch", dba.ID, maxAge)
				foundTheEnd = true
				break
			}

			dbAlerts = append(dbAlerts, dba)
		}

		if len(dbAlerts) > 0 {
			log.Printf("adding %d alerts to the database", len(dbAlerts))
			if err = query.Q.Transaction(func(tx *query.Query) error {
				dao := tx.Alert.WithContext(ctx)
				for _, alert := range dbAlerts {
					log.Printf("Inserting alert %q with optimized geography", alert.ID)
					if err := dao.InsertOptimizedAlert(alert.ID, alert.AreaDesc,
						alert.Headline, alert.Description, alert.Severity, alert.Certainty,
						alert.Urgency, alert.Event, alert.Sent, alert.Effective, alert.Onset,
						alert.Expires, alert.Ends, alert.ReferenceIds, alert.Border, alert.MessageType); err != nil {
						log.Println(err)
						return err
					}
				}

				return nil
			}); err != nil {
				log.Println(err)
				continue
			}

			dbAlerts = nil
			if p, ok := fc.ExtraMembers["pagination"]; ok {
				var pagination fcPagination
				if err = mapstructure.Decode(p, &pagination); err == nil {
					url = pagination.NextLink
				}
			}
		}

		if foundTheEnd {
			return
		}

		select {
		case <-time.After(100 * time.Millisecond):
			continue
		case <-ctx.Done():
			err = nil
			return
		}
	}

	err = nil
	return
}

func resolveGeometry(ctx context.Context, f *geojson.Feature) (*models.Geometry, error) {
	if f.Geometry != nil {
		return utils.Ref(models.NewGenericGeometry(f.Geometry)), nil
	}

	zones := query.Zone

	affectedZones := []string{}

	if intf, ok := f.Properties["affectedZones"]; ok {
		if intfSl, ok := intf.([]any); ok {
			affectedZones = utils.FromAnySlice[string](intfSl)
		}
	}

	if len(affectedZones) > 0 {
		affZones, err := zones.WithContext(ctx).Select(zones.Border).Where(zones.ID.In(affectedZones...)).Find()
		if err != nil {
			return nil, err
		}

		geos := utils.Map(affZones, func(z *models.Zone) orb.Geometry {
			return z.Border.ToOrbGeometry()
		})

		merged := utils.MergeGeometries(geos)
		return utils.Ref(models.NewGenericGeometry(merged)), nil
	}

	return nil, errors.New("no geometry available")
}

func getOptional(feature *geojson.Feature, property string) *string {
	v, ok := feature.Properties[property]
	if !ok {
		return nil
	}

	s, ok := v.(string)
	if !ok {
		return nil
	}

	return &s
}

func optionalTime(feature *geojson.Feature, property string) (t *time.Time) {
	defer func() {
		if r := recover(); r != nil {
			t = nil
		}
	}()

	timeStr := feature.Properties.MustString(property)
	ts, err := time.Parse(time.RFC3339, timeStr)

	if err != nil {
		t = nil
		return
	}

	t = &ts
	return
}
func mustTime(feature *geojson.Feature, property string) time.Time {
	t := optionalTime(feature, property)
	if t == nil {
		panic(fmt.Errorf("%s is not a time or not present on the feature", property))
	}

	return *t
}
