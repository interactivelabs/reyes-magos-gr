package store

import (
	"database/sql"
	"reyes-magos-gr/store/dtos"
)

type LibSQLVolunteersOrdersStore struct {
	DB *sql.DB
}

func NewVolunteersOrdersStore(db *sql.DB) *LibSQLVolunteersOrdersStore {
	return &LibSQLVolunteersOrdersStore{DB: db}
}

type VolunteersOrdersStore interface {
	GetActiveVolunteersWithStats() ([]dtos.VolunteerListItem, error)
}

const activeVolunteersWithStatsQuery string = `
	SELECT
		v.volunteer_id,
		v.name,
		v.email,
		COALESCE(v.phone, ''),
		v.address,
		COALESCE(v.address2, ''),
		v.country,
		v.state,
		v.city,
		COALESCE(v.province, ''),
		v.zip_code,
		v.deleted,
		(
			SELECT COUNT(*)
			FROM codes c
			INNER JOIN volunteer_codes vc ON c.code_id = vc.code_id
			WHERE
				vc.volunteer_id = v.volunteer_id
				AND vc.deleted = 0
				AND c.deleted = 0
				AND c.cancelled = 0
				AND c.given = 1
		) AS given_count,
		(
			SELECT COUNT(*)
			FROM codes c
			INNER JOIN volunteer_codes vc ON c.code_id = vc.code_id
			WHERE
				vc.volunteer_id = v.volunteer_id
				AND vc.deleted = 0
				AND c.deleted = 0
				AND c.cancelled = 0
				AND c.used = 0
				AND c.given = 0
				AND date(c.expiration) > date('now')
		) AS held_count,
		(
			SELECT MAX(o.order_date)
			FROM orders o
			WHERE
				o.volunteer_id = v.volunteer_id
				AND o.deleted = 0
				AND o.cancelled = 0
		) AS last_activity_date
	FROM volunteers v
	WHERE v.deleted = 0`

func (r *LibSQLVolunteersOrdersStore) GetActiveVolunteersWithStats() (items []dtos.VolunteerListItem, err error) {
	rows, err := r.DB.Query(activeVolunteersWithStatsQuery)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var item dtos.VolunteerListItem
		var lastActivityDate sql.NullString

		err = rows.Scan(
			&item.VolunteerID,
			&item.Name,
			&item.Email,
			&item.Phone,
			&item.Address,
			&item.Address2,
			&item.Country,
			&item.State,
			&item.City,
			&item.Province,
			&item.ZipCode,
			&item.Deleted,
			&item.GivenCount,
			&item.HeldCount,
			&lastActivityDate,
		)
		if err != nil {
			return nil, err
		}

		if lastActivityDate.Valid {
			item.LastActivityDate = &lastActivityDate.String
		}

		items = append(items, item)
	}

	return items, nil
}
