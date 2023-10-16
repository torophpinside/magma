package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"magma/pkg/common"
	"magma/pkg/dto"
	"magma/pkg/service"
	repository_sqlite "magma/pkg/storage/repository/sqlite"
	"net/http"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func NewDorksListRunCmd(rootCmd *cobra.Command) {
	newCmd := &cobra.Command{
		Use:   "dorks:run",
		Short: "Run dorks against target",
		Run: func(cmd *cobra.Command, args []string) {
			target, err := cmd.Flags().GetString("target")
			if err != nil {
				log.Fatalln("Error getting config target", err)
			}

			if target == "" {
				log.Fatalln("Target cannot be empty")
			}

			regenerate, err := cmd.Flags().GetBool("regenerate")
			if err != nil {
				regenerate = false
			}

			db := service.NewSqliteDatabase(regenerate)
			err = db.Ping()
			if err != nil {
				log.Fatalln("Error connecting to database", err)
			}
			service.RunSqliteMigration(db)

			err = runDorksList(target, regenerate, common.NewHttpClient(), db)
			if err != nil {
				log.Fatalln("Error running dorks list", err)
			}

			log.Println("Successfully runned dorks list")
		},
	}

	newCmd.Flags().StringP("target", "t", "", "Target url to search vulnerabilities or regex. Eg: domain.com, *domain.com")
	newCmd.Flags().BoolP("regenerate", "r", false, "Regenerate dorks list")
	rootCmd.AddCommand(newCmd)
}

func runDorksList(target string, regenerate bool, client *http.Client, db *sql.DB) error {
	dorkRepository := repository_sqlite.NewSqliteDorkRepository(db)
	if regenerate {
		dorksData, err := common.FileReadLines("./dorklist")
		if err != nil {
			return err
		}

		for _, dork := range dorksData {
			err = dorkRepository.SaveDork(dto.DorkDTO{
				Dork:  dork,
				Score: 0,
			})
			if err != nil {
				return err
			}
		}
		log.Println("Successfully generated dorks")
	}

	dorks, err := dorkRepository.GetDorks()
	if err != nil {
		return err
	}

	var results []dto.ParsedLinkDTO
	for _, dork := range dorks {

		links, err := parseLinks(target, client, dork.Dork)
		if err != nil {
			return err
		}
		if len(links) == 0 {
			fmt.Print(" no items found\n")
			continue

		}
		fmt.Print(" total found:", len(links), "\n")

		dork.Score = dork.Score + 1
		err = dorkRepository.UpdateDork(dork)
		if err != nil {
			return err
		}

		results = append(results, dto.ParsedLinkDTO{
			Links:  links,
			Target: target,
			Dork:   dork.Dork,
		})
	}

	log.Println("saving result to csv")
	parseAndSaveToCSV(results)

	return nil
}

func parseLinks(target string, client *http.Client, dork string) ([]string, error) {
	query := "site:" + target + " " + dork

	fmt.Print("Scanning > ", query)

	cmd := exec.Command("./go-dork", "-q", query, "-e", "google", "-p", "2")
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	results := strings.Split(string(stdout), "\n")

	return results[0 : len(results)-1], nil
}

func parseAndSaveToCSV(data []dto.ParsedLinkDTO) error {
	var csvData [][]string
	for _, row := range data {
		for _, link := range row.Links {
			csvData = append(csvData, []string{
				row.Target,
				row.Dork,
				link,
			})
		}
	}

	common.SaveToCSV(csvData, "./dorks_result.csv")
	return nil
}
