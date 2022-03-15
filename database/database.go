package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rwxsu/goot/game"
	"log"
	"time"
)

// $ docker run --name gootsMariaDB -p 3306:3306 -e MARIADB_RANDOM_ROOT_PASSWORD=true -e MYSQL_USER=goots -e MYSQL_PASSWORD=goots -e MYSQL_DATABASE=goots -d mariadb:10.4

var DB *sql.DB

func init() {
	DB, _ = newDatabase()
}

func newDatabase() (*sql.DB, error) {
	log.Println("Connecting to database...")
	db, err := sql.Open("mysql", "goots:goots@tcp(127.0.0.1:3306)/goots?parseTime=true")

	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to database")
	return db, nil
}

func GetAccountById(id uint32) (*game.Account, error) {
	account := game.Account{}

	row := DB.QueryRow("SELECT `id`, `password`, `type` FROM `account` WHERE `id` = ?", id)

	if err := row.Scan(&account.Id, &account.Password, &account.AccountType); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("getAccountById %d: no such account", id)
		}
		return nil, fmt.Errorf("getAccountById %d: %v", id, err)
	}

	players, _ := GetCharactersByAccountId(id)
	account.Characters = players

	return &account, nil
}

func GetCharactersByAccountId(id uint32) ([]game.Player, error) {
	query := "SELECT `id`, `name`, `account_id`, `sex`, `vocation`, `experience`, `level`, `maglevel`," +
		" `health`, `healthmax`, `mana`, `manamax`, `manaspent`, `lookbody`, `lookfeet`, `lookhead`," +
		" `looklegs`, `looktype`, `posx`, `posy`, `posz`, `cap`, `lastlogin`, `lastlogout`, `lastip`," +
		" `skulltime`, `skull`, `town_id`, `skill_fist`, `skill_fist_tries`, `skill_club`, `skill_club_tries`, `skill_sword`," +
		" `skill_sword_tries`, `skill_axe`, `skill_axe_tries`, `skill_dist`, `skill_dist_tries`, `skill_shielding`," +
		" `skill_shielding_tries`, `skill_fishing`, `skill_fishing_tries`, `direction` FROM `players` WHERE `account_id` = ?"
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	stmt, err := DB.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return []game.Player{}, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return []game.Player{}, err
	}
	defer rows.Close()
	var players = []game.Player{}
	for rows.Next() {
		var player game.Player
		if err := rows.Scan(&player.ID, &player.Name, &player.AccountNumber, &player.Sex, &player.Vocation,
			&player.Combat.Experience, &player.Combat.Level, &player.Magic.Level, &player.HealthNow, &player.HealthMax,
			&player.ManaNow, &player.ManaMax, &player.Magic.Experience, &player.Outfit.Body,
			&player.Outfit.Feet, &player.Outfit.Head, &player.Outfit.Legs, &player.Outfit.Type, &player.Position.X,
			&player.Position.Y, &player.Position.Z, &player.Cap, &player.LastLogin, &player.LastLogout, &player.LastIp,
			&player.SkullTime, &player.Skull, &player.TownId, &player.Fist.Level, &player.Fist.Experience, &player.Club.Level,
			&player.Club.Experience, &player.Sword.Level, &player.Sword.Experience, &player.Axe.Level, &player.Axe.Experience,
			&player.Distance.Level, &player.Distance.Experience, &player.Shielding.Level, &player.Shielding.Experience,
			&player.Fishing.Level, &player.Fishing.Experience, &player.Direction); err != nil {
			return []game.Player{}, err
		}
		players = append(players, player)
	}
	if err := rows.Err(); err != nil {
		return []game.Player{}, err
	}
	return players, nil
}

func getAccountIdByPlayerName(playerName string) (uint32, error) {
	var accountID uint32
	row := DB.QueryRow("SELECT `account_id` FROM `players` WHERE `name` = ?", playerName)
	if err := row.Scan(&accountID); err != nil {
		if err == sql.ErrNoRows {
			return accountID, fmt.Errorf("getAccountIdByPlayerName %s: no such player", playerName)
		}
		return accountID, fmt.Errorf("getAccountIdByPlayerName %s: %v", playerName, err)
	}
	return accountID, nil
}

func getAccountIdByPlayerId(playerId uint32) (uint32, error) {
	var accountID uint32
	row := DB.QueryRow("SELECT `account_id` FROM `players` WHERE `id` = ?", playerId)
	if err := row.Scan(&accountID); err != nil {
		if err == sql.ErrNoRows {
			return accountID, fmt.Errorf("getAccountIdByPlayerId %d: no such player", playerId)
		}
		return accountID, fmt.Errorf("getAccountIdByPlayerId %d: %v", playerId, err)
	}
	return accountID, nil
}

func getAccountType(accountId uint32) (uint8, error) {
	var accountType uint8
	row := DB.QueryRow("SELECT `type` FROM `accounts` WHERE `id` = ?", accountType)
	if err := row.Scan(&accountType); err != nil {
		if err == sql.ErrNoRows {
			return accountType, fmt.Errorf("getAccountType %d: no such account", accountId)
		}
		return accountType, fmt.Errorf("getAccountType %d: %v", accountId, err)
	}
	return accountType, nil
}

func updateOnlineStatus(guid uint32, login bool) (bool, error) {
	if login {
		_, err := DB.Query("INSERT INTO `players_online` VALUES (?)", guid)
		if err != nil {
			return false, fmt.Errorf("updateOnlineStatus guid: %d, %v ", guid, login)
		}
	} else {
		_, err := DB.Query("DELETE FROM `players_online` WHERE `player_id` = ?", guid)
		if err != nil {
			return false, fmt.Errorf("updateOnlineStatus guid: %d, %v ", guid, login)
		}
	}
	return true, nil
}

func preloadPlayer(player *game.Player, name string) (bool, error) {
	row := DB.QueryRow("SELECT `p`.`id`, `p`.`account_id`, `a`.`type`, `a`.`premium_ends_at` FROM `players` as `p` JOIN `accounts` as `a` ON `a`.`id` = `p`.`account_id` WHERE `p`.`name` = ? AND `p`.`deletion` = 0", name)
	if err := row.Scan(&player.ID, &player.AccountNumber, &player.AccountType, &player.PremiumEndsAt); err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("preloadPlayer %s: no such player", name)
		}
		return false, fmt.Errorf("preloadPlayer %s: %v", name, err)
	}
	return true, nil
}

func LoadPlayerById(player *game.Player, id uint32) error {
	row := DB.QueryRow("SELECT `id`, `name`, `account_id`, `sex`, `vocation`, `experience`, `level`, `maglevel`,"+
		" `health`, `healthmax`, `mana`, `manamax`, `manaspent`, `lookbody`, `lookfeet`, `lookhead`,"+
		" `looklegs`, `looktype`, `posx`, `posy`, `posz`, `cap`, `lastlogin`, `lastlogout`, `lastip`,"+
		" `skulltime`, `skull`, `town_id`, `skill_fist`, `skill_fist_tries`, `skill_club`, `skill_club_tries`, `skill_sword`,"+
		" `skill_sword_tries`, `skill_axe`, `skill_axe_tries`, `skill_dist`, `skill_dist_tries`, `skill_shielding`,"+
		" `skill_shielding_tries`, `skill_fishing`, `skill_fishing_tries`, `direction` FROM `players` WHERE `id` = ?", id)
	return loadPlayer(player, row)
}

func LoadPlayerByName(player *game.Player, name string) error {
	log.Println("SCID_4087765557 LoadPlayerByName: " + name + ".")
	row := DB.QueryRow("SELECT id, name, account_id, sex, vocation, experience, level, maglevel, health, healthmax, mana, manamax, manaspent, lookbody, lookfeet, lookhead, looklegs, looktype, posx, posy, posz, cap, lastlogin, lastlogout, lastip, skulltime, skull, town_id, skill_fist, skill_fist_tries, skill_club, skill_club_tries, skill_sword, skill_sword_tries, skill_axe, skill_axe_tries, skill_dist, skill_dist_tries, skill_shielding, skill_shielding_tries, skill_fishing, skill_fishing_tries, direction FROM players WHERE name = ?", name)
	return loadPlayer(player, row)
}

func loadPlayer(player *game.Player, row *sql.Row) error {
	switch err := row.Scan(&player.ID, &player.Name, &player.AccountNumber, &player.Sex, &player.Vocation,
		&player.Combat.Experience, &player.Combat.Level, &player.Magic.Level, &player.HealthNow, &player.HealthMax,
		&player.ManaNow, &player.ManaMax, &player.Magic.Experience, &player.Outfit.Body,
		&player.Outfit.Feet, &player.Outfit.Head, &player.Outfit.Legs, &player.Outfit.Type, &player.Position.X,
		&player.Position.Y, &player.Position.Z, &player.Cap, &player.LastLogin, &player.LastLogout, &player.LastIp,
		&player.SkullTime, &player.Skull, &player.TownId, &player.Fist.Level, &player.Fist.Experience, &player.Club.Level,
		&player.Club.Experience, &player.Sword.Level, &player.Sword.Experience, &player.Axe.Level, &player.Axe.Experience,
		&player.Distance.Level, &player.Distance.Experience, &player.Shielding.Level, &player.Shielding.Experience,
		&player.Fishing.Level, &player.Fishing.Experience, &player.Direction); err {
	case sql.ErrNoRows:
		{
			return fmt.Errorf("SCID_4134685069 loadPlayer: no such player")
		}
	case sql.ErrConnDone:
		{
			return fmt.Errorf("SCID_8100774773 loadPlayer: ErrConnDone")
		}
	}

	player.Access = game.Tutor
	player.Combat.Percent = 20
	player.Magic.Percent = 50
	player.Fist.Percent = 50
	player.Club.Percent = 50
	player.Sword.Percent = 50
	player.Axe.Percent = 50
	player.Distance.Percent = 50
	player.Shielding.Percent = 50
	player.Fishing.Percent = 50
	player.Icons = 1
	player.Light = game.Light{Level: 0x7, Color: 0xd7}
	player.World = game.World{Light: game.Light{Level: 0x00, Color: 0xd7}}
	player.Speed = 200

	return nil
}

func GetPlayerByName(name string) (game.Player, error) {
	var player game.Player
	row := DB.QueryRow("SELECT * FROM player WHERE name = ?", name)

	if err := row.Scan(&player.ID, &player.Name); err != nil {
		if err == sql.ErrNoRows {
			return player, fmt.Errorf("PlayerByName %s: no such player", name)
		}
		return player, fmt.Errorf("PlayerByName %s: %v", name, err)
	}
	return player, nil
}
