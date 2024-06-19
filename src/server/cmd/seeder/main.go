package main

import (
	"database/sql"
	"log"
	"math/rand"
	"strconv"

	"github.com/cebuh/simpleHolidayPlaner/config"
	"github.com/cebuh/simpleHolidayPlaner/db"
	"github.com/cebuh/simpleHolidayPlaner/service/auth"
	"github.com/cebuh/simpleHolidayPlaner/service/team"
	"github.com/cebuh/simpleHolidayPlaner/service/user"
	"github.com/cebuh/simpleHolidayPlaner/types"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func main() {
	db, err := db.NewMySqlStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)
	log.Println("Start seeding...")
	seed(db)
}

func seed(db *sql.DB) {
	createUsersAndTeams(db)
	log.Println("seeding completed!")
}

func initStorage(db *sql.DB) {
	log.Println("Check Database connection...")
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database successfully connected!")
}

// use prepared statements to execute faster
func createUsersAndTeams(db *sql.DB) {
	userStore := user.NewStore(db)
	log.Println("seed users...")

	generatedUsers := make([]types.User, 0)
	generatedTeams := make([]types.Team, 0)

	for i := 0; i < 20; i++ {
		testpw, err := auth.HashPassword("password" + strconv.Itoa(i))
		if err != nil {
			panic(err)
		}

		user := types.User{
			Id:       uuid.NewString(),
			Name:     "user" + strconv.Itoa(i),
			Email:    "user" + strconv.Itoa(i) + "@seed.com",
			Password: testpw,
		}
		if err := userStore.CreateUser(user); err != nil {
			panic(err)
		}

		generatedUsers = append(generatedUsers, user)
	}
	log.Println("seed teams...")
	teamStore := team.NewStore(db)
	teamNames := [10]string{
		"The Masterminds",
		"Lightning Minds",
		"The Innovation Hunters",
		"Team Titan",
		"The Visionaries",
		"The Trailblazers",
		"Team Dynamo",
		"The Peak Performers",
		"The Power Pioneers",
		"The Strategy Specialists",
	}

	for _, name := range teamNames {
		team := types.Team{
			Id:   uuid.NewString(),
			Name: name,
		}
		if err := teamStore.CreateTeam(team); err != nil {
			panic(err)
		}

		generatedTeams = append(generatedTeams, team)
	}

	log.Println("add users to teams...")
	removeUser := func(s []types.User, index int) []types.User {
		return append(s[:index], s[index+1:]...)
	}
	for range 5 {
		randomNumber := rand.Intn(11)

		index1 := rand.Intn(len(generatedUsers))
		index2 := rand.Intn(len(generatedUsers) - 1)
		if index2 >= index1 {
			index2++
		}

		user1 := generatedUsers[index1]
		user2 := generatedUsers[index2]

		randomTeam := generatedTeams[randomNumber]
		teamStore.AddUserToTeam(user1.Id, randomTeam.Id, types.Member)
		teamStore.AddUserToTeam(user2.Id, randomTeam.Id, types.Member)

		generatedUsers = removeUser(generatedUsers, index1)
		generatedUsers = removeUser(generatedUsers, index2)

	}
}
