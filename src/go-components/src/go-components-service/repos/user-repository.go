// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0

package repos

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"go-component-service/models"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var Users models.Users
var UsersById map[string]int
var UsersByUsername map[string]int
var UsersByIdentityId map[string]int
var UsersByPrimaryPersona map[string][]int
var UsersByAgeRange map[string][]int
var UsersClaimedByIdentityId map[int]bool

// Init
func init() {
	loadedUsers, err := loadUsers("/bin/data/users.json.gz")
	if err != nil {
		log.Panic("Unable to load users file: ", err)
	}
	Users = loadedUsers
}

func loadUsers(filename string) (models.Users, error) {

	log.Println("Attempting to load users file: ", filename)

	var r models.Users
	UsersById = make(map[string]int)
	UsersByUsername = make(map[string]int)
	UsersByIdentityId = make(map[string]int)
	UsersByPrimaryPersona = make(map[string][]int)
	UsersByAgeRange = make(map[string][]int)
	UsersClaimedByIdentityId = make(map[int]bool)

	file, err := os.Open(filename)
	if err != nil {
		return r, err
	}

	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return r, err
	}

	defer gz.Close()

	dec := json.NewDecoder(gz)

	err = dec.Decode(&r)
	if err != nil {
		return r, err
	}

	// Load maps with user array index
	for i, u := range r {
		UsersById[u.ID] = i
		UsersByUsername[u.Username] = i
		UsersByPrimaryPersona[strings.Split(u.Persona, "_")[0]] = append(UsersByPrimaryPersona[strings.Split(u.Persona, "_")[0]], i)
		UsersByAgeRange[getAgeRange(u.Age)] = append(UsersByAgeRange[getAgeRange(u.Age)], i)
	}

	log.Println("Users successfully loaded into memory structures")

	return r, nil
}

func getAgeRange(age int) string {
	if age < 18 {
		return ""
	} else if age < 25 {
		return "18-24"
	} else if age < 35 {
		return "25-34"
	} else if age < 45 {
		return "35-44"
	} else if age < 55 {
		return "45-54"
	} else if age < 70 {
		return "54-70"
	} else {
		return "70-and-above"
	}
}

// containsInt returns a bool indicating whether the given []int contained the given int
func containsInt(slice []int, value int) bool {
	for _, v := range slice {
		if value == v {
			return true
		}
	}
	return false
}

// RepoFindUsersIdByAgeRange Function
func RepoFindUserIdsByAgeRange(ageRange string) []int {
	return UsersByAgeRange[ageRange]
}

// RepoFindUsersIdByPrimaryPersona Function
func RepoFindUsersIdByPrimaryPersona(persona string) []int {
	return UsersByPrimaryPersona[persona]
}

// RepoFindRandomUsersByPrimaryPersonaAndAgeRage Function
func RepoFindRandomUsersByPrimaryPersonaAndAgeRange(primaryPersona string, ageRange string, count int) models.Users {
	var unclaimedUsers models.Users
	var primaryPersonaFilteredUserIds = RepoFindUsersIdByPrimaryPersona(primaryPersona)
	var ageRangeFilteredUserIds = RepoFindUserIdsByAgeRange(ageRange)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(ageRangeFilteredUserIds), func(i, j int) {
		ageRangeFilteredUserIds[i], ageRangeFilteredUserIds[j] = ageRangeFilteredUserIds[j], ageRangeFilteredUserIds[i]
	})
	for _, idx := range ageRangeFilteredUserIds {
		if len(unclaimedUsers) >= count {
			break
		}
		if containsInt(primaryPersonaFilteredUserIds, idx) && !(UsersClaimedByIdentityId[idx]) {
			if Users[idx].SelectableUser {
				log.Println("User found matching filter criteria:", idx)
				unclaimedUsers = append(unclaimedUsers, Users[idx])
			}
		}
	}
	return unclaimedUsers
}

// RepoClaimUser Function
// Function used to map which shopper user ids have been claimed by the user Id.
func RepoClaimUser(userId int) bool {
	log.Println("An identity has claimed the user id:", userId)
	UsersClaimedByIdentityId[userId] = true
	return true
}

func RepoFindRandomUser(count int) models.Users {
	rand.Seed(time.Now().UnixNano())
	var randomUserId int
	var randomUsers models.Users
	if len(Users) > 0 {
		for len(randomUsers) < count {
			randomUserId = rand.Intn(len(Users))
			log.Println("Random number Selected:", randomUserId)
			if randomUserId != 0 {
				if !(UsersClaimedByIdentityId[randomUserId]) {
					if Users[randomUserId].SelectableUser {
						log.Println("Random user id selected:", randomUserId)
						randomUsers = append(randomUsers, RepoFindUserByID(strconv.Itoa(randomUserId)))
						log.Println("Random users :", randomUsers)
					}
				}
			}
		}
	}
	return randomUsers
}

// RepoFindUserByID Function
func RepoFindUserByID(id string) models.User {
	if idx, ok := UsersById[id]; ok {
		return Users[idx]
	} else {
		return models.User{}
	}
}

// RepoFindUserByUsername Function
func RepoFindUserByUsername(username string) models.User {
	if idx, ok := UsersByUsername[username]; ok {
		return Users[idx]
	} else {
		return models.User{}
	}
}

// RepoFindUserByIdentityID Function
func RepoFindUserByIdentityID(identityID string) models.User {
	if idx, ok := UsersByIdentityId[identityID]; ok {
		return Users[idx]
	} else {
		return models.User{}
	}
}

// RepoUpdateUser Function
func RepoUpdateUser(t models.User) models.User {
	if idx, ok := UsersById[t.ID]; ok {
		u := &Users[idx]
		u.FirstName = t.FirstName
		u.LastName = t.LastName
		u.Email = t.Email
		u.Addresses = t.Addresses
		u.SignUpDate = t.SignUpDate
		u.LastSignInDate = t.LastSignInDate
		u.PhoneNumber = t.PhoneNumber

		if len(u.IdentityId) > 0 && u.IdentityId != t.IdentityId {
			delete(UsersByIdentityId, u.IdentityId)
		}

		u.IdentityId = t.IdentityId

		if len(t.IdentityId) > 0 {
			UsersByIdentityId[t.IdentityId] = idx
		}

		return RepoFindUserByID(t.ID)
	}

	// return empty User if not found
	return models.User{}
}

// RepoCreateUser Function
func RepoCreateUser(t models.User) (models.User, error) {
	if _, ok := UsersByUsername[t.Username]; ok {
		return models.User{}, errors.New("User with this username already exists")
	}

	idx := len(Users)

	if len(t.ID) > 0 {
		// ID provided by caller (provisionally created on storefront) so make
		// sure it's not already taken.
		if _, ok := UsersById[t.ID]; ok {
			return models.User{}, errors.New("User with this ID already exists")
		}
	} else {
		t.ID = strconv.Itoa(idx)
	}

	Users = append(Users, t)
	UsersById[t.ID] = idx
	UsersByUsername[t.Username] = idx
	if len(t.IdentityId) > 0 {
		UsersByIdentityId[t.IdentityId] = idx
	}

	return t, nil
}
