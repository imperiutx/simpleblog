package db

import (
	"context"
	"fmt"
)

func (q *QueriesDynamic) ListPostsDynamic(ctx context.Context, title string, tags []string, filters Filters) ([]*Post, Metadata, error) {

	query := fmt.Sprintf(`
	SELECT COUNT(*) OVER(), id, username, title, content, tags, status, created_at, edited_at
	FROM "posts"
	WHERE
		(TO_TSVECTOR('simple', title) @@ PLAINTO_TSQUERY('simple', $1) OR $1 = '')
		AND (tags @> $2 OR $2 = '{}')
	ORDER BY %s %s, id ASC
	LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	args := []interface{}{title, tags, filters.limit(), filters.offset()}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0

	posts := []*Post{}
	for rows.Next() {
		var p Post
		if err := rows.Scan(
			&totalRecords,
			&p.ID,
			&p.Username,
			&p.Title,
			&p.Content,
			&p.Tags,
			&p.Status,
			&p.CreatedAt,
			&p.EditedAt,
		); err != nil {
			return nil, Metadata{}, err
		}
		posts = append(posts, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, int(filters.Page), int(filters.PageSize))

	return posts, metadata, nil
}
