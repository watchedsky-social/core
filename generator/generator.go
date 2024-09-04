package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/watchedsky-social/core/internal/config"
	"github.com/watchedsky-social/core/internal/database/models"
	"github.com/watchedsky-social/core/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	files := []string{".env"}
	if envFile, ok := os.LookupEnv("WATCHEDSKY_ENV_FILE"); ok {
		files = append(files, envFile)
	}

	files = utils.Filter(files, func(file string) bool {
		_, err := os.Stat(file)
		return err == nil
	})

	fmt.Printf("%v", files)

	if err := godotenv.Load(files...); err != nil {
		log.Fatal(err)
	}
	config.Load()

	args := config.Config.DB

	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require TimeZone=UTC", args.Host, args.User, args.Password, args.Name)), &gorm.Config{
		Logger: logger.Default,
	})
	if err != nil {
		log.Fatal(err)
	}

	dataTypeMap := map[string]func(columnType gorm.ColumnType) (dataType string){
		"geometry": func(columnType gorm.ColumnType) (dataType string) {
			ct, _ := columnType.ColumnType()
			if strings.Contains(strings.ToLower(ct), "geometry(") {
				return "*Geometry"
			}

			return "string"
		},
		"text[]": func(columnType gorm.ColumnType) (dataType string) {
			return "StringSlice"
		},
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:           "./query",
		OutFile:           "gen_query.go",
		ModelPkgPath:      "./models",
		WithUnitTest:      true,
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(db)
	g.WithDataTypeMap(dataTypeMap)
	g.WithImportPkgPath("github.com/paulmach/orb")

	g.ApplyInterface(
		func(CustomZoneQueries) {},
		g.GenerateModel("zones", gen.FieldJSONTagWithNS(func(columnName string) string {
			return fmt.Sprintf("%s,omitempty", columnName)
		})),
	)
	g.ApplyInterface(
		func(CustomMapSearchQueries) {},
		g.GenerateModel("mapsearch", gen.FieldIgnore("display_name"),
			gen.FieldJSONTagWithNS(func(columnName string) string {
				return fmt.Sprintf("%s,omitempty", columnName)
			})),
	)
	g.ApplyInterface(
		func(CustomAlertsQueries) {},
		g.GenerateModel("alerts", gen.FieldJSONTagWithNS(func(columnName string) string {
			return fmt.Sprintf("%s,omitempty", columnName)
		}), gen.FieldType("skeet_info", "*SkeetInfo")),
	)
	g.ApplyInterface(
		func(CustomSavedAreaQueries) {},
		g.GenerateModel("saved_areas"),
	)
	g.ApplyBasic(
		g.GenerateModel("didwids"),
	)
	g.Execute()
}

type CustomZoneQueries interface {
	// SELECT count(*) FROM zones WHERE ST_Intersects(border, ST_SetSRID(ST_MakeBox2D(@southEast, @northWest), 4326));
	CountVisibleZones(southEast models.Geometry, northWest models.Geometry) (int64, error)

	// SELECT * FROM zones WHERE type='county' AND ST_Intersects(border, ST_SetSRID(ST_MakeBox2D(@southEast, @northWest), 4326)) ORDER BY concat(name, ' ', type, ' ', state) LIMIT 20;
	ShowVisibleZones(southEast models.Geometry, northWest models.Geometry) ([]*gen.T, error)

	// SELECT id FROM zones;
	ListIDs() ([]string, error)

	// SELECT z2.id FROM zones z1 INNER JOIN zones z2 ON z1.border = z2.border WHERE z1.id IN (@zoneList);
	FindCongruentZones(zoneList []string) ([]string, error)

	// SELECT ST_Union(border) FROM zones WHERE id IN (@affectedZones);
	ResolveGeometry(affectedZones []string) (*models.Geometry, error)
}

type CustomMapSearchQueries interface {
	/* WITH searchResults AS (
	       SELECT * FROM mapsearch
	           WHERE display_name ILIKE @searchText || '%' OR id LIKE @searchText || '%'
	       )
	   SELECT DISTINCT ON (display_name) id, name, state, county, centroid
	       FROM searchResults ORDER by display_name;  */
	PrefixSearch(searchText string) ([]*gen.T, error)
}

type CustomAlertsQueries interface {
	/*
		INSERT INTO alerts (
		    id, area_desc, headline, description, severity, certainty, urgency, event,
		    sent, effective, onset, expires, ends, reference_ids, border, message_type
		    ) VALUES (
		    @id, @areaDesc, @headline, @description, @severity, @certainty,
		    @urgency, @event, @sent, @effective, @onset, @expires, @ends,
		    @referenceIDs, ST_UnaryUnion(@border),
		    @messageType
		    ) ON CONFLICT(id) DO UPDATE SET
		    area_desc = EXCLUDED.area_desc,
		    headline = EXCLUDED.headline,
		    description = EXCLUDED.description,
		    severity = EXCLUDED.severity,
		    certainty = EXCLUDED.certainty,
		    urgency = EXCLUDED.urgency,
		    event = EXCLUDED.event,
		    sent = EXCLUDED.sent,
		    effective = EXCLUDED.effective,
		    onset = EXCLUDED.onset,
		    expires = EXCLUDED.expires,
		    ends = EXCLUDED.ends,
		    reference_ids = EXCLUDED.reference_ids,
		    border = EXCLUDED.border,
		    message_type = EXCLUDED.message_type;
	*/
	InsertOptimizedAlert(
		id, areaDesc, headline, description string,
		severity, certainty, urgency, event *string,
		sent, effective time.Time, onset, expires, ends *time.Time,
		referenceIDs *models.StringSlice, border *models.Geometry, messageType *string,
	) error

	/*
	   WITH target_area AS (SELECT border FROM saved_areas WHERE id = @watchID LIMIT 1)
	     SELECT a.skeet_info AS skeet_info, EXTRACT(EPOCH FROM a.sent) * 1000 as sent
	     FROM alerts a, target_area t
	     WHERE a.skeet_info IS NOT NULL
	     AND (a.border && t.border)
	     AND ST_Intersects(a.border, t.border)
	     LIMIT @limit
	     ORDER BY a.sent DESC;
	*/
	GetCustomAlertURIs(watchID string, limit uint) ([]*gen.T, error)

	/*
	   WITH target_area AS (SELECT border FROM saved_areas WHERE ID = @watchID LIMIT 1)
	     SELECT a.skeet_info AS skeet_info, EXTRACT(EPOCH FROM a.sent) * 1000 as sent
	     FROM alerts a, target_area t
	     WHERE a.skeet_info IS NOT NULL
	     AND a.border && t.border AND ST_Intersects(a.border, t.border)
	     AND sent < @cursor
	     ORDER BY a.sent LIMIT @limit DESC;
	*/
	GetCustomAlertURIsWithCursor(watchID string, limit uint, cursor uint) ([]*gen.T, error)
}

type CustomSavedAreaQueries interface {
	/*
	  INSERT INTO saved_areas (
	    id, passed_zones, calculated_zones, border
	  ) VALUES (
	    @id, @passedZones, @calculatedZones,
	    ST_UnaryUnion(@border)
	  ) ON CONFLICT(id) DO UPDATE SET
	    passed_zones = EXCLUDED.passed_zones,
	    calculated_zones = EXCLUDED.calculated_zones,
	    border = EXCLUDED.border;
	*/
	InsertOptimizedSavedArea(id, passedZones, calculatedZones string, border *models.Geometry) error
}
