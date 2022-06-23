package bot

import (
	"discord-fight/config"
	"discord-fight/data"
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

//global bot id
var BotId string
var CurrentPhase Phase
var Round int

//game phases
type Phase string

const (
	Idle       Phase = "Idle"
	Collecting Phase = "Collecting"
	Playing    Phase = "Playing"
)

//command regex
var (
	collectRegex *regexp.Regexp
	startRegex   *regexp.Regexp
	stopRegex    *regexp.Regexp
	clearRegex   *regexp.Regexp
	nextRegex    *regexp.Regexp
	fighterRegex *regexp.Regexp
	powerRegex   *regexp.Regexp
	nerfRegex    *regexp.Regexp
)

func Start() {
	//creating new bot session
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//make bot a user
	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Store bot id
	BotId = u.ID

	//compile regex
	collectRegex, _ = regexp.Compile("^/collect($|\\s.*)")
	startRegex, _ = regexp.Compile("^/start($|\\s.*)")
	stopRegex, _ = regexp.Compile("^/stop($|\\s.*)")
	clearRegex, _ = regexp.Compile("^/clear($|\\s.*)")
	nextRegex, _ = regexp.Compile("^/next($|\\s.*)")
	fighterRegex, _ = regexp.Compile("^/fighter($|\\s.*)")
	powerRegex, _ = regexp.Compile("^/power($|\\s.*)")
	nerfRegex, _ = regexp.Compile("^/nerf($|\\s.*)")

	//handlers
	goBot.AddHandler(collectHandler)

	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	CurrentPhase = Idle
	fmt.Println("Discord fight bot is running!")
}

func collectHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != BotId {
		if collectRegex.MatchString(m.Content) {
			fmt.Println("got call to /collect")
			CurrentPhase = Collecting
			if _, err := s.ChannelMessageSend(m.ChannelID, data.CollectMsg); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func startHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != BotId {
		if startRegex.MatchString(m.Content) {
			fmt.Println("got call to /start")

			//check for data
			if !data.Ready() {
				if _, err := s.ChannelMessageSend(m.ChannelID, data.FalseStart); err != nil {
					fmt.Println(err.Error())
				}
			} else {
				CurrentPhase = Playing
				Round = 0
				if _, err := s.ChannelMessageSend(m.ChannelID, data.TrueStart); err != nil {
					fmt.Println(err.Error())
				}

				//begin fight
				if newFight, err := GetFightString(Round); err != nil {
					fmt.Println(err.Error())
				} else {
					if _, err := s.ChannelMessageSend(m.ChannelID, newFight); err != nil {
						fmt.Println(err.Error())
					}
				}
			}
		}
	}
}

func stopHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != BotId {
		if stopRegex.MatchString(m.Content) && CurrentPhase != Idle {
			fmt.Println("got call to /stop")
			CurrentPhase = Idle
			if _, err := s.ChannelMessageSend(m.ChannelID, "Thanks for playing! Resume at any time with \"/collect\"!"); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func clearHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != BotId {
		if clearRegex.MatchString(m.Content) {
			fmt.Println("got call to /clear")
			if CurrentPhase == Playing {
				if _, err := s.ChannelMessageSend(m.ChannelID, "Cannot clear data in the middle of a round. Send a command to \"/stop\" first."); err != nil {
					fmt.Println(err.Error())
				}
			} else {
				data.Clear()
				if _, err := s.ChannelMessageSend(m.ChannelID, "Data cleared."); err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}
}

func nextHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != BotId {
		if nextRegex.MatchString(m.Content) && CurrentPhase == Playing {
			fmt.Println("got call to /next")
			Round++

			//next fight
			if newFight, err := GetFightString(Round); err != nil {
				fmt.Println(err.Error())
			} else {
				if _, err := s.ChannelMessageSend(m.ChannelID, newFight); err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}
}

func fighterHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != BotId {
		if fighterRegex.MatchString(m.Content) && CurrentPhase == Collecting {
			fmt.Println("got call to /fighter")

			//parse fighter
			if len(m.Content) > len("/fighter ") {
				fighter := m.Content[len("/fighter "):]

				//add fighter
				data.AddFighter(fighter)
			}

			//delete message so no one sees it
			if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func powerHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != BotId {
		if powerRegex.MatchString(m.Content) && CurrentPhase == Collecting {
			fmt.Println("got call to /power")

			//parse power
			if len(m.Content) > len("/power ") {
				power := m.Content[len("/power "):]

				//add power
				data.AddPower(power)
			}

			//delete message so no one sees it
			if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func nerfHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != BotId {
		if nerfRegex.MatchString(m.Content) && CurrentPhase == Collecting {
			fmt.Println("got call to /nerf")

			//parse nerf
			if len(m.Content) > len("/nerf ") {
				nerf := m.Content[len("/nerf "):]

				//add nerf
				data.AddNerf(nerf)
			}

			//delete message so no one sees it
			if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

//game-related functions

func GetFightString(round int) (string, error) {
	fighter1, err := data.GetFighter()
	if err != nil {
		return "", err
	}

	power1, err := data.GetPower()
	if err != nil {
		return "", err
	}

	nerf1, err := data.GetNerf()
	if err != nil {
		return "", err
	}

	fighter2, err := data.GetFighter()
	if err != nil {
		return "", err
	}

	power2, err := data.GetPower()
	if err != nil {
		return "", err
	}

	nerf2, err := data.GetNerf()
	if err != nil {
		return "", err
	}

	return "Round " + fmt.Sprint(round) + "...\n\nWho would win?\n" +
		fighter1 + " with " + power1 + " but " + nerf1 + "\nOR\n" +
		fighter2 + " with " + power2 + " but " + nerf2, nil
}
