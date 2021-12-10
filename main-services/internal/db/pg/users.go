package pg

func AddUsers(username string) error {
	const q = `INSERT INTO users (
		user_name)
	VALUES ($1)`
	_, err := UserDB.Exec(
		q, username)
	if err != nil {
		return err
	}
	return nil
}
