package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/sahapranta/kong-check/config"
	"github.com/sahapranta/kong-check/models"
	"github.com/sahapranta/kong-check/utils"
)

func GetConnection(conf *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", conf.DatabaseURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetAllRoutes(conf *config.Config) ([]models.Route, error) {
	db, err := GetConnection(conf)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT 
			r.id, 
			r.name, 
			r.paths, 
			r.methods, 
			r.protocols,
			r.hosts,
			r.headers,
			r.service_id, 
			s.name as service_name
		FROM routes r
		JOIN services s ON r.service_id = s.id
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	routes := []models.Route{}
	for rows.Next() {
		var r models.Route
		var paths, methods, protocols, hosts, headers sql.NullString

		err := rows.Scan(
			&r.ID,
			&r.Name,
			&paths,
			&methods,
			&protocols,
			&hosts,
			&headers,
			&r.ServiceID,
			&r.ServiceName,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		if paths.Valid {
			r.Paths = utils.ParsePostgresArray(paths.String)
		}
		if methods.Valid {
			r.Methods = utils.ParsePostgresArray(methods.String)
		}
		if protocols.Valid {
			r.Protocols = utils.ParsePostgresArray(protocols.String)
		}
		if hosts.Valid {
			r.HostNames = utils.ParsePostgresArray(hosts.String)
		}
		if headers.Valid {
			r.Headers = utils.ParsePostgresJSONB(headers.String)
		}

		routes = append(routes, r)
	}

	return routes, nil
}

func GetRoutesByServiceNames(conf *config.Config, serviceNames []string) ([]models.Route, error) {
	db, err := GetConnection(conf)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Build query with placeholders for services
	query := `
		SELECT 
			r.id, 
			r.name, 
			r.paths, 
			r.methods, 
			r.protocols,
			r.hosts,
			r.headers,
			r.service_id, 
			s.name as service_name
		FROM routes r
		JOIN services s ON r.service_id = s.id
		WHERE s.name = ANY($1)
	`

	// Convert slice to PostgreSQL array syntax
	pgArray := "{" + utils.JoinStrings(serviceNames, ",") + "}"

	rows, err := db.Query(query, pgArray)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	routes := []models.Route{}
	for rows.Next() {
		var r models.Route
		var paths, methods, protocols, hosts, headers sql.NullString

		err := rows.Scan(
			&r.ID,
			&r.Name,
			&paths,
			&methods,
			&protocols,
			&hosts,
			&headers,
			&r.ServiceID,
			&r.ServiceName,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		if paths.Valid {
			r.Paths = utils.ParsePostgresArray(paths.String)
		}
		if methods.Valid {
			r.Methods = utils.ParsePostgresArray(methods.String)
		}
		if protocols.Valid {
			r.Protocols = utils.ParsePostgresArray(protocols.String)
		}
		if hosts.Valid {
			r.HostNames = utils.ParsePostgresArray(hosts.String)
		}
		if headers.Valid {
			r.Headers = utils.ParsePostgresJSONB(headers.String)
		}

		routes = append(routes, r)
	}

	return routes, nil
}
