package helpers

import (
	"github.com/InspectorGadget/realtime-polling-system/db"
	"github.com/InspectorGadget/realtime-polling-system/models"
)

func FetchPollFromDB(pollID string) (models.Poll, error) {
	var p models.Poll
	err := db.GetDB().QueryRow("SELECT id, topic FROM polls WHERE id=$1", pollID).Scan(&p.ID, &p.Topic)
	if err != nil {
		return p, err
	}

	rows, err := db.GetDB().Query("SELECT id, text, votes FROM options WHERE poll_id=$1 ORDER BY id", pollID)
	if err != nil {
		return p, err
	}
	defer rows.Close()

	for rows.Next() {
		var o models.Option
		rows.Scan(&o.ID, &o.Text, &o.Votes)
		p.Options = append(p.Options, o)
	}

	return p, nil
}
