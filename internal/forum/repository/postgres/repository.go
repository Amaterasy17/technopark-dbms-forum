package postgres

import (
	"fmt"
	"github.com/jackc/pgx"
	"strings"
	domain "technopark-dbms-forum/internal/forum"
	models "technopark-dbms-forum/models"
	"time"
)

type postgresForumRepository struct {
	Conn *pgx.ConnPool
}

func NewPostgresForumRepository(Conn *pgx.ConnPool) domain.ForumRepository {
	return &postgresForumRepository{Conn: Conn}
}

func (p *postgresForumRepository) InsertForum(forum models.Forum) error {
	_, err := p.Conn.Exec(	`Insert INTO forum(Slug, "user", Title) VALUES ($1, $2, $3);`,
		forum.Slug, forum.UserId, forum.Title)
	if err != nil {
		return err
	}
	return nil
}

func (p *postgresForumRepository) SelectForum(forumName string) (models.Forum, error) {
	var forum models.Forum
	row := p.Conn.QueryRow(`Select slug, "user", title, posts, threads From forum
				Where slug=$1 LIMIT 1`, forumName)
	err := row.Scan(&forum.Slug, &forum.UserId, &forum.Title, &forum.Posts, &forum.Threads)
	if err != nil {
		return models.Forum{}, models.ErrNotFound
	}
	forum.User = p.SelectNicknameForum(forum.UserId)
	return forum, nil
}

func (p *postgresForumRepository) SelectNicknameForum(user_id int) string {
	var result string
	row := p.Conn.QueryRow(`Select nickname from users where id=$1 LIMIT 1`, user_id)
	err := row.Scan(&result)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func (p *postgresForumRepository) CheckForum(forum models.Forum) (models.Forum, bool) {
	resultForum := models.Forum{
		Posts: -1,
	}
	row := p.Conn.QueryRow(`Select slug, user, title, posts, threads From forum
				Where slug=$1`, forum.Slug)
	_ = row.Scan(&resultForum.Slug, &resultForum.User, &resultForum.Title, &resultForum.Posts, &resultForum.Threads)
	if resultForum.Posts == -1 {
		return models.Forum{},false
	}
	return resultForum, true
}

func (p *postgresForumRepository) SelectUsers(user models.User) ([]models.User, error) {
	var users []models.User
	rows, err := p.Conn.Query(`Select Nickname, FullName, About, Email From users Where Nickname=$1 or Email=$2 LIMIT 2;`,
														user.Nickname, user.Email)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return users, err
	}
	for rows.Next() {
		var userModel models.User
		err := rows.Scan(&userModel.Nickname, &userModel.FullName, &userModel.About, &userModel.Email)
		if err != nil {
			return users, err
		}
		users = append(users, userModel)
	}

	return users, nil
}

func (p *postgresForumRepository) InsertUser(user models.User) error {
	_, err := p.Conn.Exec(	`Insert INTO users(Nickname, FullName, About, Email) VALUES ($1, $2, $3, $4);`,
		user.Nickname, user.FullName, user.About, user.Email)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (p *postgresForumRepository) SelectUser(user string) (models.User, error) {
	var userModel models.User
	row := p.Conn.QueryRow(`Select id, Nickname, FullName, About, Email From users Where nickname=$1 LIMIT 1;`, user)
	err := row.Scan(&userModel.ID, &userModel.Nickname, &userModel.FullName, &userModel.About, &userModel.Email)
	if err != nil {
		return models.User{}, models.ErrNotFound
	}
	return userModel, nil
}

func (p *postgresForumRepository) SelectUserByEmail(user models.User) (models.User, error) {
	var userModel models.User
	row := p.Conn.QueryRow(`Select nickname, email from users Where email=$1 LIMIT 1;`, user.Email)
	err := row.Scan(&userModel.Nickname, &userModel.Email)
	if err != nil {
		return models.User{}, nil
	}
	if userModel.Nickname == user.Nickname {
		return models.User{}, nil
	}
	return userModel, models.ErrConflict
}

func (p *postgresForumRepository) UpdateUserInfo(user models.User) (models.User, error) {
	var err error
	var newUser models.User

	err = p.Conn.QueryRow(
		`UPDATE users SET email=COALESCE(NULLIF($1, ''), email), 
							  about=COALESCE(NULLIF($2, ''), about), 
							  fullname=COALESCE(NULLIF($3, ''), fullname) WHERE nickname=$4 RETURNING *`,
		user.Email,
		user.About,
		user.FullName,
		user.Nickname,
	).Scan(&newUser.ID,&newUser.Nickname, &newUser.FullName, &newUser.About, &newUser.Email)
	//if user.FullName != "" {
	//	_, err = p.Conn.Exec(`UPDATE users SET fullname=$1 WHERE nickname=$2;`, user.FullName, user.Nickname)
	//	if err != nil {
	//		return models.User{}, err
	//	}
	//}
	//
	//if user.About != "" {
	//	_, err = p.Conn.Exec(`UPDATE users SET about=$1 WHERE nickname=$2;`, user.About, user.Nickname)
	//	if err != nil {
	//		return models.User{}, err
	//	}
	//}
	//
	//
	//if user.Email != "" {
	//	_, err = p.Conn.Exec(`UPDATE users SET email=$1 WHERE nickname=$2;`, user.Email, user.Nickname)
	//	if err != nil {
	//		return models.User{}, err
	//	}
	//}

	return newUser, err
}

func (p *postgresForumRepository) SelectThreadBySlug(slug string) (models.Thread, error) {
	var thread models.Thread
	row := p.Conn.QueryRow(`Select id, title, author, forum, message, votes, slug, created from thread
							Where slug=$1 LIMIT 1;`, slug)
	err := row.Scan(&thread.Id, &thread.Title, &thread.AuthorId, &thread.Forum, &thread.Message, &thread.Votes,
					&thread.Slug, &thread.Created)
	if err != nil {
		return models.Thread{}, models.ErrNotFound
	}
	thread.Author = p.SelectNicknameForum(thread.AuthorId)
	return thread, nil
}

func (p *postgresForumRepository) InsertThread(thread models.Thread) (models.Thread,error) {
	var newThread models.Thread
	var row *pgx.Row

	row = p.Conn.QueryRow(	`Insert INTO thread(Title, Author, Created, Forum, Message, slug, Votes)
							VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *`, thread.Title, thread.AuthorId, thread.Created,
							thread.Forum,
			thread.Message, thread.Slug, thread.Votes)

	err := row.Scan(&newThread.Id,&newThread.Title, &newThread.AuthorId, &newThread.Created,
		&newThread.Forum, &newThread.Message, &newThread.Slug, &newThread.Votes)
	if err != nil {
		return models.Thread{},err
	}
	newThread.Author = p.SelectNicknameForum(newThread.AuthorId)
	return newThread, nil
}

func (p *postgresForumRepository) SelectThreadById(id int) (models.Thread, error) {
	var thread models.Thread
	row := p.Conn.QueryRow(`Select id, title, author, forum, message, votes, slug, created from thread
							Where id=$1 LIMIT 1;`, id)

	err := row.Scan(&thread.Id, &thread.Title, &thread.AuthorId, &thread.Forum, &thread.Message, &thread.Votes,
		&thread.Slug, &thread.Created)
	if err != nil {
		return models.Thread{}, models.ErrNotFound
	}
	thread.Author = p.SelectNicknameForum(thread.AuthorId)
	return thread, nil
}

func (p *postgresForumRepository) CheckParent(post models.Post) bool {
	fmt.Printf("menya vizvali")
	fmt.Println(post.Parent)
	var id string
	row := p.Conn.QueryRow(`Select author from post where id=$1;`, post.Parent.Int64)

	err := row.Scan(&id)

	if err == pgx.ErrNoRows {
		fmt.Printf("FALSE FLASE FALSE")
		return false
	}
	fmt.Println(id)
	if err != nil {
		fmt.Printf("FALSE FLASE FALSE")
		return false
	}
	return true
}
func (p *postgresForumRepository) InsertPost(post models.Post) (models.Post, error) {
	var row *pgx.Row

	row = p.Conn.QueryRow(`INSERT INTO post(author, created, forum, message, parent, thread) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;`,
			post.Author, post.Created, post.Forum, post.Message, post.Parent, post.Thread)

	var postModel models.Post
	err := row.Scan(&postModel.ID, &postModel.Author, &postModel.Created, &postModel.Forum,  &postModel.IsEdited,
		&postModel.Message, &postModel.Parent, &postModel.Thread, &postModel.Path)



	return postModel, err
}

func (p *postgresForumRepository) StatusOfForum() models.Status {
	var status models.Status
	err := p.Conn.QueryRow(`SELECT COUNT(*) FROM users;`).Scan(&status.User)
	if err != nil {
		status.User = 0
	}
	err = p.Conn.QueryRow(`SELECT COUNT(*) FROM forum;`).Scan(&status.Forum)
	if err != nil {
		status.Forum = 0
	}
	err = p.Conn.QueryRow(`SELECT COUNT(*) FROM thread;`).Scan(&status.Thread)
	if err != nil {
		status.Thread = 0
	}
	err = p.Conn.QueryRow(`SELECT COUNT(*) FROM post;`).Scan(&status.Post)
	if err != nil {
		status.Post = 0
	}
	return status
}

func (p *postgresForumRepository) ClearDB() error {
	var err error
	_, err = p.Conn.Exec(`TRUNCATE users CASCADE;`)
	_, err = p.Conn.Exec(`TRUNCATE forum CASCADE;`)
	_, err = p.Conn.Exec(`TRUNCATE thread CASCADE;`)
	_, err = p.Conn.Exec(`TRUNCATE post CASCADE;`)
	_, err = p.Conn.Exec(`TRUNCATE votes CASCADE;`)
	return err
}

func (p *postgresForumRepository) SelectVote(vote models.Vote) (models.Vote, error) {
	var voteResult models.Vote
	row := p.Conn.QueryRow(`Select author, voice, thread from votes Where author=$1 and thread=$2;`, vote.Nickname, vote.Thread)
	err := row.Scan(&voteResult.Nickname, &voteResult.Voice, &voteResult.Thread)
	if err != nil {
		return models.Vote{}, models.ErrNotFound
	}
	return voteResult, nil
}


func (p *postgresForumRepository) UpdateVote(vote models.Vote) (models.Vote, error) {
	_, err := p.Conn.Exec(`UPDATE votes SET voice=$1 WHERE author=$2 and thread=$3;`, vote.Voice, vote.Nickname, vote.Thread)
	if err != nil {
		return models.Vote{}, err
	}
	return vote, nil
}

func (p *postgresForumRepository) InsertVote(vote models.Vote)  error {
	_, err := p.Conn.Exec(`INSERT INTO votes(author, voice, thread) VALUES ($1, $2, $3);`, vote.Nickname,
							vote.Voice, vote.Thread)
	if err != nil {
		return err
	}
	return nil
}

func (p *postgresForumRepository) SumVotesInThread(id int) int {
	var sum int
	row := p.Conn.QueryRow(`Select SUM(voice) from votes WHERE thread=$1;`, id)
	err := row.Scan(&sum)
	if err != nil {
		return 0
	}
	return sum
}

func (p *postgresForumRepository) UpdatePost(post models.Post, postUpdate models.PostUpdate) (models.Post, error) {

		row := p.Conn.QueryRow(`UPDATE post SET message=COALESCE(NULLIF($1, ''), message),
                             isEdited = CASE WHEN $1 = '' OR message = $1 THEN isEdited ELSE true END
                             WHERE id=$2 RETURNING *`, postUpdate.Message, post.ID)
		err := row.Scan(&post.ID, &post.Author, &post.Created, &post.Forum,  &post.IsEdited,
			&post.Message, &post.Parent, &post.Thread, &post.Path)
		if err != nil {
			return post, err
		}
	return post, nil
}


func (p *postgresForumRepository) SelectPost(id int) (models.Post, error) {
	var postModel models.Post
	row := p.Conn.QueryRow(`Select id, author, created, forum, isEdited, message, parent, thread from post Where id=$1 LIMIT 1;`, id)
	err := row.Scan(&postModel.ID, &postModel.Author, &postModel.Created, &postModel.Forum,  &postModel.IsEdited,
		&postModel.Message, &postModel.Parent, &postModel.Thread)
	if err != nil {
		return models.Post{}, models.ErrNotFound
	}
	//postModel.Author = p.SelectNickById(postModel.AuthorId)
	return postModel, nil
}

func (p *postgresForumRepository) SelectThreads(slug string, params models.Parameters) ([]models.Thread, error) {
	var threads []models.Thread
	var err error
	var rows *pgx.Rows

	if params.Since != "" {
		if params.Desc {
			rows, err = p.Conn.Query(`SELECT id, author, created, forum, message, slug, title, votes FROM thread
		WHERE forum=$1 AND created <= $2 ORDER BY created DESC LIMIT $3;`, slug, params.Since, params.Limit)
		} else {
			rows, err = p.Conn.Query(`SELECT id, author, created, forum, message, slug, title, votes FROM thread
		WHERE forum=$1 AND created >= $2 ORDER BY created ASC LIMIT $3;`, slug, params.Since, params.Limit)
		}
	} else {
		if params.Desc {
			rows, err = p.Conn.Query(`SELECT id, author, created, forum, message, slug, title, votes FROM thread
		WHERE forum=$1 ORDER BY created DESC LIMIT $2;`, slug, params.Limit)
		} else {
			rows, err = p.Conn.Query(`SELECT id, author, created, forum, message, slug, title, votes FROM thread
		WHERE forum=$1 ORDER BY created ASC LIMIT $2;`, slug, params.Limit)
		}
	}

	if err != nil {
		return threads, err
	}
	defer rows.Close()


	for rows.Next() {
		var thread models.Thread
		err = rows.Scan(&thread.Id, &thread.AuthorId, &thread.Created, &thread.Forum, &thread.Message,
			&thread.Slug, &thread.Title, &thread.Votes)
		if err != nil {
			continue
		}
		thread.Author = p.SelectNicknameForum(thread.AuthorId)
		threads = append(threads, thread)
	}
	return threads, nil
}

func (p *postgresForumRepository) SelectUsersByForum(slug string, params models.Parameters) ([]models.User, error) {
	var query string
	if params.Desc {
		if params.Since != "" {
		//	query = fmt.Sprintf(`SELECT users.about, users.Email, users.FullName, users.Nickname FROM users
    	//inner join users_forum uf on users.Nickname = uf.nickname
        //WHERE uf.slug =$1 AND uf.nickname < '%s'
        //ORDER BY users.Nickname DESC LIMIT NULLIF($2, 0)`, params.Since)
			query = fmt.Sprintf(
				`SELECT about, email, fullname, nickname 
				FROM users_forum WHERE slug=$1 AND nickname < '%s' 
				ORDER BY nickname DESC LIMIT NULLIF($2, 0)`, params.Since)
		} else {
		//	query = `SELECT users.about, users.Email, users.FullName, users.Nickname FROM users
    	//inner join users_forum uf on users.Nickname = uf.nickname
        //WHERE uf.slug =$1
        //ORDER BY users.Nickname DESC LIMIT NULLIF($2, 0)`
			query = `SELECT about, email, fullname, nickname 
				FROM users_forum WHERE slug=$1 
				ORDER BY nickname DESC LIMIT NULLIF($2, 0)`
		}
	} else {
		//query = fmt.Sprintf(`SELECT users.about, users.Email, users.FullName, users.Nickname FROM users
    	//inner join users_forum uf on users.Nickname = uf.nickname
        //WHERE uf.slug =$1 AND uf.nickname > '%s'
        //ORDER BY users.Nickname LIMIT NULLIF($2, 0)`, params.Since)
		query = fmt.Sprintf(`SELECT about, email, fullname, nickname
			FROM users_forum WHERE slug=$1 AND nickname > '%s'
			ORDER BY nickname LIMIT NULLIF($2, 0)`, params.Since)
	}
	var data []models.User
	row, err := p.Conn.Query(query, slug, params.Limit)

	if err != nil {
		return data, nil
	}

	defer func() {
		if row != nil {
			row.Close()
		}
	}()

	for row.Next() {

		var u models.User

		err = row.Scan(&u.About, &u.Email, &u.FullName, &u.Nickname)

		if err != nil {
			return data, err
		}

		data = append(data, u)
	}

	return data, err
}


func (p *postgresForumRepository) PostFlatSort(id int, parameters models.Parameters) ([]models.Post, error) {
	var err error
	var rows *pgx.Rows
	var posts []models.Post

	if parameters.Since == "" {
		if parameters.Desc {
			rows, err = p.Conn.Query(`SELECT id, author, created, forum, isEdited, message, parent, thread FROM post
		WHERE thread=$1 ORDER BY id DESC LIMIT $2;`, id, parameters.Limit)
		} else {
			rows, err = p.Conn.Query(`SELECT id, author, created, forum, isEdited, message, parent, thread FROM post
		WHERE thread=$1 ORDER BY id LIMIT $2;`, id, parameters.Limit)
		}
	} else {
		if parameters.Desc {
			rows, err = p.Conn.Query(`SELECT id, author, created, forum, isEdited, message, parent, thread FROM post
		WHERE thread=$1 AND id < $2 ORDER BY id DESC LIMIT $3;`, id, parameters.Since, parameters.Limit)
		} else {
			rows, err = p.Conn.Query(`SELECT id, author, created, forum, isEdited, message, parent, thread FROM post
		WHERE thread=$1 AND id > $2 ORDER BY id LIMIT $3;`, id, parameters.Since, parameters.Limit)
		}
	}

	if err != nil {
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.ID, &post.Author, &post.Created, &post.Forum, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
		if err != nil {
			return posts, err
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func (p *postgresForumRepository) PostTreeSort(threadId int, parameters models.Parameters) ([]models.Post, error) {
	var err error
	var rows *pgx.Rows
	var posts []models.Post

	if parameters.Since == "" {
		if parameters.Desc {
			rows, err = p.Conn.Query(`SELECT id, author, created, forum, isEdited, message, parent, thread FROM post
		WHERE thread=$1 ORDER BY path DESC, id DESC LIMIT $2;`, threadId, parameters.Limit)
		} else {
			rows, err = p.Conn.Query(`SELECT id, author, created, forum, isEdited, message, parent, thread FROM post
		WHERE thread=$1 ORDER BY path ASC, id  ASC LIMIT $2;`, threadId, parameters.Limit)
		}
	} else {
		if parameters.Desc {
			rows, err = p.Conn.Query(`SELECT id, author, created, forum, isEdited, message, parent, thread FROM post
		WHERE thread=$1 AND PATH < (SELECT path FROM post WHERE id = $2)
		ORDER BY path DESC, id  DESC LIMIT $3;`, threadId, parameters.Since, parameters.Limit)
		} else {
			rows, err = p.Conn.Query(`SELECT id, author, created, forum, isEdited, message, parent, thread FROM post
		WHERE thread=$1 AND PATH > (SELECT path FROM post WHERE id = $2)
		ORDER BY path ASC, id  ASC LIMIT $3;`, threadId, parameters.Since, parameters.Limit)
		}
	}

	if err != nil {
		return posts, err
	}
	defer rows.Close()


	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.ID, &post.Author, &post.Created, &post.Forum, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
		if err != nil {
			return posts, err
		}

		//post.Author = p.SelectNickById(post.AuthorId)
		posts = append(posts, post)
	}
	return posts, nil
}

func (p *postgresForumRepository) PostParentTreeSort(threadId int, parameters models.Parameters) ([]models.Post, error) {
	var err error
	var rows *pgx.Rows
	var posts []models.Post

	if parameters.Since == "" {
		if parameters.Desc {
			rows, err = p.Conn.Query(`SELECT id, author, created, forum, isEdited, message, parent, thread FROM post
			WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent IS NULL ORDER BY id DESC LIMIT $2)
			ORDER BY path[1] DESC, path, id;`, threadId, parameters.Limit)
		} else {
			rows, err = p.Conn.Query(`SELECT id, author, created, forum, isEdited, message, parent, thread FROM post
			WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent IS NULL ORDER BY id LIMIT $2)
			ORDER BY path, id;`, threadId, parameters.Limit)
		}
	} else {
		if parameters.Desc {
			rows, err = p.Conn.Query(`SELECT id, author, created, forum, isEdited, message, parent, thread FROM post
				WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent IS NULL AND PATH[1] <
				(SELECT path[1] FROM post WHERE id = $2) ORDER BY id DESC LIMIT $3) ORDER BY path[1] DESC, path, id;`,
				threadId, parameters.Since, parameters.Limit)
		} else {
			rows, err = p.Conn.Query(`SELECT id, author, created, forum, isEdited, message, parent, thread FROM post
				WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent IS NULL AND PATH[1] >
				(SELECT path[1] FROM post WHERE id = $2) ORDER BY id ASC LIMIT $3) ORDER BY path, id;`,
				threadId, parameters.Since, parameters.Limit)
		}
	}

	if err != nil {
		return posts, err
	}
	defer rows.Close()


	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.ID, &post.Author, &post.Created, &post.Forum, &post.IsEdited, &post.Message,
			&post.Parent, &post.Thread)
		if err != nil {
			return posts, err
		}

		//post.Author = p.SelectNickById(post.AuthorId)
		posts = append(posts, post)
	}
	return posts, nil
}

func (p *postgresForumRepository) UpdateThread(thread models.Thread) (models.Thread, error) {
	var row *pgx.Row
	query := `UPDATE thread SET title=COALESCE(NULLIF($1, ''), title), message=COALESCE(NULLIF($2, ''), message) WHERE %s RETURNING *`

	if thread.Slug == "" {
		query = fmt.Sprintf(query, `id=$3`)
		row = p.Conn.QueryRow(query, thread.Title, thread.Message, thread.Id)
		//row = p.Conn.QueryRow(`UPDATE thread SET title=$1, message=$2 WHERE id=$3 RETURNING *`, thread.Title, thread.Message, thread.Id)
	} else {
		query = fmt.Sprintf(query, `slug=$3`)
		row = p.Conn.QueryRow(query, thread.Title, thread.Message, thread.Slug)
		//row = p.Conn.QueryRow(`UPDATE thread SET title=$1, message=$2 WHERE LOWER(slug)=LOWER($3) RETURNING *`, thread.Title, thread.Message, thread.Slug)
	}

	var newThread models.Thread

	err := row.Scan(
		&newThread.Id,
		&newThread.Title,
		&newThread.AuthorId,
		&newThread.Created,
		&newThread.Forum,
		&newThread.Message,
		&newThread.Slug,
		&newThread.Votes,
	)

	newThread.Author = p.SelectNickById(newThread.AuthorId)

	if err != nil {
		fmt.Println(err)
		return models.Thread{}, models.ErrNotFound
	}

	return newThread, nil
}

func (p *postgresForumRepository) NewTransaction() (*pgx.Tx, error) {
	return p.Conn.Begin()
}

func (p *postgresForumRepository) Rollback(tx *pgx.Tx) {
	p.Rollback(tx)
}

func (p *postgresForumRepository) InsertPosts(posts *[]models.Post, thread models.Thread) (*[]models.Post, error) {
	query := `INSERT INTO post(author, created, forum, message, parent, thread) VALUES`

	var values []interface{}
	created := time.Now()
	for i, post := range *posts {
		value := fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d),",
			i * 6 + 1, i * 6 + 2, i * 6 + 3, i * 6 + 4, i * 6 + 5, i * 6 + 6,
		)

		//userId := p.SelectIdByNickname(post.Author)


		query += value

		values = append(values, post.Author)
		values = append(values, created)
		values = append(values, thread.Forum)
		values = append(values, post.Message)
		values = append(values, post.Parent)
		values = append(values, thread.Id)
	}

	query = strings.TrimSuffix(query, ",")
	query += ` RETURNING id, created, forum, isEdited, thread;`

	rows, err := p.Conn.Query(query, values...)
	if err != nil {
		fmt.Println("error of insert")
		return nil, err
	}
	defer rows.Close()
	//var postsResult []models.Post

	for i, _ := range *posts {
		if rows.Next() {
			err := rows.Scan(&(*posts)[i].ID, &(*posts)[i].Created, &(*posts)[i].Forum, &(*posts)[i].IsEdited, &(*posts)[i].Thread)
			if err != nil {
				fmt.Println(err)
				return nil, models.ErrConflict
			}
		}
	}
	if rows.Err() != nil {

		return nil, rows.Err()
		//switch rows.Err().(pgx.PgError).Code {
		//case "23503":
		//	return nil, models.ErrNotFound
		//default:
		//	return nil, models.ErrConflict
		//}
	}

	//fmt.Println(*posts)
	//for rows.Next() {
	//	var postModel models.Post
	//	err := rows.Scan(&postModel.ID, &postModel.Author, &postModel.Created, &postModel.Forum,  &postModel.IsEdited,
	//		&postModel.Message, &postModel.Parent, &postModel.Thread, &postModel.Path)
	//	if err != nil {
	//		fmt.Println("error of SCAN")
	//		return nil, err
	//	}
	//
	//	if !postModel.Parent.Valid {
	//		postModel.Parent.Int64 = 0
	//		postModel.Parent.Valid = true
	//	}
	//	//postModel.Author = p.SelectNicknameForum(postModel.AuthorId)
	//	postsResult = append(postsResult, postModel)
	//}

	return posts, err
	//query := `INSERT INTO post(
    //             author,
    //             created,
    //             message,
    //             parent,
	//			 thread,
	//			 forum) VALUES `
	//data := make([]models.Post, 0, 0)
	//if len(*posts) == 0 {
	//	return &data, nil
	//}
	//
	//slug := thread.Forum
	//
	//
	//timeCreated := time.Now()
	//var valuesNames []string
	//var values []interface{}
	//i := 1
	//for _, element := range *posts {
	//	valuesNames = append(valuesNames, fmt.Sprintf(
	//		"($%d, $%d, $%d, nullif($%d, 0), $%d, $%d)",
	//		i, i+1, i+2, i+3, i+4, i+5))
	//	i += 6
	//	values = append(values, element.Author, timeCreated, element.Message, element.Parent, thread.Id, slug)
	//}
	//
	//query += strings.Join(valuesNames[:], ",")
	//query += " RETURNING *"
	//row, err := p.Conn.Query(query, values...)
	//
	//if err != nil {
	//	return &data, err
	//}
	//defer func() {
	//	if row != nil {
	//		row.Close()
	//	}
	//}()
	//
	//for row.Next() {
	//
	//	var post models.Post
	//
	//
	//	err = row.Scan(&post.ID, &post.Author, &post.Created, &post.Forum,  &post.IsEdited,
	//		&post.Message, &post.Parent, &post.Thread, &post.Path)
	//
	//	if err != nil {
	//		return &data, err
	//	}
	//
	//	data = append(data, post)
	//
	//}
	//
	//if row.Err() != nil {
	//	//
	//		return nil, row.Err()
	//	//	//switch rows.Err().(pgx.PgError).Code {
	//	//	//case "23503":
	//	//	//	return nil, models.ErrNotFound
	//	//	//default:
	//	//	//	return nil, models.ErrConflict
	//	//	//}
	//	}
	//
	//return &data, err
}

func (p *postgresForumRepository) SelectNickById(userId int) string {
	var result string
	row := p.Conn.QueryRow(`Select nickname from users where id=$1 LIMIT 1`, userId)
	err := row.Scan(&result)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func (p *postgresForumRepository) SelectIdByNickname(nick string) int {
	var result int
	row := p.Conn.QueryRow(`Select id from users where nickname=$1 LIMIT 1`, nick)
	err := row.Scan(&result)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func (p *postgresForumRepository) AutocommitOff() {
	p.Conn.QueryRow(`SET autocommit TO 'off';`)
}

func (p *postgresForumRepository) AutocommitOn() {
	p.Conn.QueryRow(`SET autocommit TO 'on';`)
}