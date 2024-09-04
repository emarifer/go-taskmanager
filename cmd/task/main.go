package main

import (
	"taskmanager/cmd"
	"taskmanager/internal/tasks"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	// dsn := "taskmanager_user:taskmanager_user_password@tcp(127.0.0.1:13306)/taskmanager?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(sqlite.Open("app_data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&tasks.Task{})

	cmd.Execute(db)
}

/* REFERENCES:
https://github.com/goreleaser/goreleaser/releases
https://goreleaser.com/deprecations/?h=name_template#after_23
https://greeeg.com/en/issues/how-to-create-release-distribute-cli-golang
https://github.com/go-gorm/gorm/discussions/6095#discussioncomment-7967488
https://forum.golangbridge.org/t/getting-errors-when-running-go-project-that-uses-github-com-mattn-go-sqlite3-library/31800/2
https://goreleaser.com/deprecations/#-rm-dist
https://goreleaser.com/blog/goreleaser-v2/?h=v2#upgrading
https://github.com/goreleaser/goreleaser-action
https://github.com/mattn/go-sqlite3/issues/303

https://www.google.com/search?q=failed+to+initialize+database%2C+got+error+Binary+was+compiled+with+%27CGO_ENABLED%3D0%27%2C+go-sqlite3+requires+cgo+to+work.+This+is+a+stub+panic%3A+failed+to+connect+database&oq=failed+to+initialize+database%2C+got+error+Binary+was+compiled+with+%27CGO_ENABLED%3D0%27%2C+go-sqlite3+requires+cgo+to+work.+This+is+a+stub+panic%3A+failed+to+connect+database&aqs=chrome..69i57j69i60.1900688j0j7&sourceid=chrome&ie=UTF-8
*/

/* COMMANDS:
goreleaser init
goreleaser release --snapshot --clean
goreleaser check

git tag 1.0.0 && git push --tags
*/
