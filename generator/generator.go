package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/watchedsky-social/core/internal/config"
	"github.com/watchedsky-social/core/internal/database/models"
	"github.com/watchedsky-social/core/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
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

	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require TimeZone=UTC", args.Host, args.User, args.Password, args.Name)))
	if err != nil {
		log.Fatal(err)
	}

	dataTypeMap := map[string]func(columnType gorm.ColumnType) (dataType string){
		"geometry": func(columnType gorm.ColumnType) (dataType string) {
			ct, _ := columnType.ColumnType()
			if strings.Contains(strings.ToLower(ct), "geometry(") {
				return "Geometry"
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
		g.GenerateModel("zones", gen.FieldJSONTagWithNS(func(columnName string) (tagContent string) {
			tagContent = fmt.Sprintf("%s,omitempty", columnName)
			return
		})),
	)
	g.ApplyInterface(
		func(CustomMapSearchQueries) {},
		g.GenerateModel("mapsearch", gen.FieldIgnore("display_name"), gen.FieldJSONTagWithNS(func(columnName string) (tagContent string) {
			tagContent = fmt.Sprintf("%s,omitempty", columnName)
			return
		})),
	)
	g.ApplyBasic(
		g.GenerateModel("saved_areas"),
		g.GenerateModel("alerts", gen.FieldJSONTagWithNS(func(columnName string) (tagContent string) {
			tagContent = fmt.Sprintf("%s,omitempty", columnName)
			return
		}), gen.FieldType("skeet_info", "*SkeetInfo")),
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
