SELECT a.id, a.user_id, a.title, a.description, a.price, a.image_url, a.created_at, u.login
FROM ads AS a
JOIN users AS u ON a.user_id = u.id
WHERE a.price >= $1 AND a.price <= $2
ORDER BY %s %s
LIMIT $3 OFFSET $4