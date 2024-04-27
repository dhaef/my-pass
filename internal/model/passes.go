package model

import (
	"database/sql"

	"github.com/dhaef/my-pass/internal/db"
)

type PassItem struct {
	Id     sql.NullString
	Value  sql.NullString
	PassId string
}

type Pass struct {
	Id        string
	UserId    string
	Name      sql.NullString
	Username  sql.NullString
	Password  sql.NullString
	Tags      []PassItem
	Websites  []PassItem
	CreatedAt string         `json"createdAt"`
	UpdatedAt sql.NullString `json"updatedAt"`
}

func CreatePass(pass *Pass) (*Pass, error) {
	createdAt := getNowTimeStamp()
	if err := db.GetDB().QueryRow(
		"INSERT INTO passes(userId, name, username, password, createdAt) VALUES($1, $2, $3, $4, $5) RETURNING id",
		pass.UserId,
		pass.Name.String,
		pass.Username.String,
		pass.Password.String, // TODO: encrypt
		createdAt,
	).Scan(&pass.Id); err != nil {
		return pass, err
	}
	pass.CreatedAt = createdAt

	for i, tag := range pass.Tags {
		if err := db.GetDB().QueryRow(
			"INSERT INTO tags(id, passId, value, createdAt) VALUES($1, $2, $3, $4) RETURNING id",
			tag.Id,
			pass.Id,
			tag.Value,
			createdAt,
		).Scan(&pass.Tags[i].Id); err != nil {
			return pass, err
		}
	}

	for i, website := range pass.Websites {
		if err := db.GetDB().QueryRow(
			"INSERT INTO websites(id, passId, value, createdAt) VALUES($1, $2, $3, $4) RETURNING id",
			website.Id,
			pass.Id,
			website.Value,
			createdAt,
		).Scan(&pass.Websites[i].Id); err != nil {
			return pass, err
		}
	}

	return pass, nil
}

func UpdatePass(pass *Pass) (*Pass, error) {
	updatedAt := getNowTimeStamp()
	_, err := db.GetDB().Exec(
		"UPDATE passes SET userId = $1, name = $2, username = $3, password = $4, updatedAt = $5 WHERE id = $6",
		pass.UserId,
		pass.Name.String,
		pass.Username.String,
		pass.Password.String, // TODO: encrypt
		updatedAt,
		pass.Id,
	)
	if err != nil {
		return pass, err
	}
	pass.UpdatedAt = sql.NullString{
		String: updatedAt,
		Valid:  true,
	}

	currentTags, err := getTagsByPassId(pass.Id)
	if err != nil {
		return pass, err
	}

	tagsToDelete := getPassItemsToDelete(currentTags, pass.Tags)

	for _, tag := range tagsToDelete {
		_, err := db.GetDB().Exec(
			"DELETE FROM tags WHERE id = $1",
			tag,
		)
		if err != nil {
			return pass, err
		}
	}

	for _, tag := range pass.Tags {
		_, err := db.GetDB().Exec(
			`INSERT INTO 
			tags(id, passId, value, createdAt) 
			VALUES($1, $2, $3, $4) 
			ON CONFLICT (id) DO UPDATE
			SET value = $3, updatedAt = $4`,
			tag.Id,
			pass.Id,
			tag.Value,
			updatedAt,
		)
		if err != nil {
			return pass, err
		}
	}

	currentWebsites, err := getWebsitesByPassId(pass.Id)
	if err != nil {
		return pass, err
	}

	websitesToDelete := getPassItemsToDelete(currentWebsites, pass.Websites)

	for _, website := range websitesToDelete {
		_, err := db.GetDB().Exec(
			"DELETE FROM websites WHERE id = $1",
			website,
		)
		if err != nil {
			return pass, err
		}
	}

	for _, website := range pass.Websites {
		_, err := db.GetDB().Exec(
			`INSERT INTO 
			websites(id, passId, value, createdAt) 
			VALUES($1, $2, $3, $4) 
			ON CONFLICT (id) DO UPDATE
			SET value = $3, updatedAt = $4`,
			website.Id,
			pass.Id,
			website.Value,
			updatedAt,
		)
		if err != nil {
			return pass, err
		}
	}

	return pass, nil
}

func getPassItemsToDelete(currentItems []PassItem, newItems []PassItem) []string {
	passItemsToDelete := []string{}
	for _, item := range currentItems {
		exists := false
		for _, i := range newItems {
			if i.Id.String == item.Id.String {
				exists = true
				break
			}
		}

		if !exists {
			passItemsToDelete = append(passItemsToDelete, item.Id.String)
		}
		continue
	}

	return passItemsToDelete
}

func getTagsByPassId(id string) ([]PassItem, error) {
	rows, err := db.GetDB().Query(
		"SELECT id, value FROM tags WHERE passId = $1",
		id,
	)
	if err != nil {
		return []PassItem{}, err
	}
	defer rows.Close()

	tags := []PassItem{}
	for rows.Next() {
		tag := PassItem{}
		if err := rows.Scan(
			&tag.Id,
			&tag.Value,
		); err != nil {
			return []PassItem{}, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func getWebsitesByPassId(id string) ([]PassItem, error) {
	rows, err := db.GetDB().Query(
		"SELECT id, value FROM websites WHERE passId = $1",
		id,
	)
	if err != nil {
		return []PassItem{}, err
	}
	defer rows.Close()

	websites := []PassItem{}
	for rows.Next() {
		tag := PassItem{}
		if err := rows.Scan(
			&tag.Id,
			&tag.Value,
		); err != nil {
			return []PassItem{}, err
		}

		websites = append(websites, tag)
	}

	return websites, nil
}

func GetPass(id string, userId string) (Pass, error) {
	rows, err := db.GetDB().Query(
		`SELECT 
		p.id, 
		p.userId, 
		p.name, 
		p.username, 
		p.password, 
		p.createdAt, 
		p.updatedAt,
		w.id,
		w.value,
		t.id,
		t.value 
		FROM passes p 
		LEFT JOIN websites w ON p.id = w.passId
		LEFT JOIN tags t ON p.id = t.passId
		WHERE p.id = $1 AND p.userId = $2`,
		id,
		userId,
	)
	if err != nil {
		return Pass{}, err
	}
	defer rows.Close()

	fullPass := Pass{}

	for rows.Next() {
		var pass Pass
		tag := PassItem{}
		website := PassItem{}
		if err := rows.Scan(
			&pass.Id,
			&pass.UserId,
			&pass.Name,
			&pass.Username,
			&pass.Password,
			&pass.CreatedAt,
			&pass.UpdatedAt,
			&website.Id,
			&website.Value,
			&tag.Id,
			&tag.Value,
		); err != nil {
			return fullPass, err
		}

		// set base values in the first loop
		// if fullPass.Id != "" {
		fullPass.Id = pass.Id
		fullPass.UserId = pass.UserId
		fullPass.Name = pass.Name
		fullPass.Username = pass.Username
		fullPass.Password = pass.Password
		fullPass.CreatedAt = pass.CreatedAt
		fullPass.UpdatedAt = pass.UpdatedAt
		// }

		if tag.Id.Valid {
			fullPass.Tags = append(fullPass.Tags, tag)
		}

		if website.Id.Valid {
			fullPass.Websites = append(fullPass.Websites, website)
		}

		// reset values
		tag = PassItem{}
		website = PassItem{}
		continue
	}
	if err = rows.Err(); err != nil {
		return fullPass, err
	}
	return fullPass, nil
}

func GetPasses(userId string) ([]Pass, error) {
	rows, err := db.GetDB().Query(
		"SELECT id, userId, name, username, password, createdAt FROM passes WHERE userId = $1",
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	passes := []Pass{}

	for rows.Next() {
		var pass Pass
		if err := rows.Scan(
			&pass.Id,
			&pass.UserId,
			&pass.Name,
			&pass.Username,
			&pass.Password,
			&pass.CreatedAt,
		); err != nil {
			return passes, err
		}
		passes = append(passes, pass)
	}
	if err = rows.Err(); err != nil {
		return passes, err
	}
	return passes, nil
}

func DeletePass(id string) error {
	_, err := db.GetDB().Exec(
		`DELETE FROM passes WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}
