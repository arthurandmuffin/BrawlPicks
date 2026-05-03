package app

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"BrawlPicks/cli/brawlers"
	"BrawlPicks/cli/client"
	"BrawlPicks/cli/config"
)

type App struct {
	cfg             *config.Config
	client          *client.RecommendationClient
	brawlers        *brawlers.Catalog
	mapName         string
	mode            string
	rank            int
	topK            int
	allyBrawlers    []int
	enemyBrawlers   []int
	bannedBrawlers  []int
	recommendations []*client.Recommendation
	modelID         string
	message         string
	showHelp        bool
}

func New(cfg *config.Config, recommendationClient *client.RecommendationClient, brawlerCatalog *brawlers.Catalog) *App {
	return &App{
		cfg:            cfg,
		client:         recommendationClient,
		brawlers:       brawlerCatalog,
		mapName:        cfg.UI.DefaultMapName,
		mode:           cfg.UI.DefaultMode,
		rank:           cfg.UI.DefaultRank,
		topK:           cfg.UI.DefaultTopK,
		allyBrawlers:   []int{},
		enemyBrawlers:  []int{},
		bannedBrawlers: []int{},
		message:        "type help for commands",
	}
}

func (a *App) Run() error {
	reader := bufio.NewReader(os.Stdin)

	for {
		a.render()
		fmt.Print("\ncmd> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		quit, err := a.handle(strings.TrimSpace(line))
		if err != nil {
			a.message = err.Error()
			continue
		}
		if quit {
			a.clearScreen()
			return nil
		}
	}
}

func (a *App) handle(line string) (bool, error) {
	if line == "" {
		return false, nil
	}

	parts := strings.Fields(line)
	switch parts[0] {
	case "q", "quit", "exit":
		return true, nil
	case "help":
		a.showHelp = true
		a.message = "help open"
		return false, nil
	case "map":
		a.showHelp = false
		if len(parts) < 2 {
			return false, fmt.Errorf("usage: map <name>")
		}
		a.mapName = strings.TrimSpace(line[len("map "):])
		a.message = "map updated"
		return false, nil
	case "mode":
		a.showHelp = false
		if len(parts) != 2 {
			return false, fmt.Errorf("usage: mode <value>")
		}
		a.mode = parts[1]
		a.message = "mode updated"
		return false, nil
	case "rank":
		a.showHelp = false
		value, err := parseIntArg(parts, 1, "usage: rank <number>")
		if err != nil {
			return false, err
		}
		a.rank = value
		a.message = "rank updated"
		return false, nil
	case "topk":
		a.showHelp = false
		value, err := parseIntArg(parts, 1, "usage: topk <number>")
		if err != nil {
			return false, err
		}
		a.topK = value
		a.message = "top_k updated"
		return false, nil
	case "ally":
		a.showHelp = false
		return false, a.editList("ally", parts, &a.allyBrawlers)
	case "enemy":
		a.showHelp = false
		return false, a.editList("enemy", parts, &a.enemyBrawlers)
	case "ban":
		a.showHelp = false
		return false, a.editList("ban", parts, &a.bannedBrawlers)
	case "clear":
		a.showHelp = false
		return false, a.clear(parts)
	case "rec", "recommend":
		a.showHelp = false
		return false, a.recommend()
	default:
		a.showHelp = false
		return false, fmt.Errorf("unknown command: %s", parts[0])
	}
}

func (a *App) editList(kind string, parts []string, target *[]int) error {
	if len(parts) < 3 {
		return fmt.Errorf("usage: %s <add|rm> <brawler_name>", kind)
	}

	value, err := a.brawlers.IDForName(strings.Join(parts[2:], " "))
	if err != nil {
		return err
	}

	switch parts[1] {
	case "add":
		if err := a.ensureAvailable(kind, value); err != nil {
			return err
		}
		*target = addUnique(*target, value)
		a.message = fmt.Sprintf("%s added", kind)
		return nil
	case "rm":
		*target = removeValue(*target, value)
		a.message = fmt.Sprintf("%s removed", kind)
		return nil
	default:
		return fmt.Errorf("usage: %s <add|rm> <brawler_name>", kind)
	}
}

func (a *App) ensureAvailable(kind string, value int) error {
	if contains(a.allyBrawlers, value) && kind != "ally" {
		return fmt.Errorf("brawler %d already in ally list", value)
	}
	if contains(a.enemyBrawlers, value) && kind != "enemy" {
		return fmt.Errorf("brawler %d already in enemy list", value)
	}
	if contains(a.bannedBrawlers, value) && kind != "ban" {
		return fmt.Errorf("brawler %d already in banned list", value)
	}
	return nil
}

func (a *App) clear(parts []string) error {
	if len(parts) != 2 {
		return fmt.Errorf("usage: clear <ally|enemy|ban|all>")
	}

	switch parts[1] {
	case "ally":
		a.allyBrawlers = []int{}
	case "enemy":
		a.enemyBrawlers = []int{}
	case "ban":
		a.bannedBrawlers = []int{}
	case "all":
		a.allyBrawlers = []int{}
		a.enemyBrawlers = []int{}
		a.bannedBrawlers = []int{}
		a.recommendations = nil
		a.modelID = ""
	default:
		return fmt.Errorf("usage: clear <ally|enemy|ban|all>")
	}

	a.message = "state cleared"
	return nil
}

func (a *App) recommend() error {
	topK := a.topK
	resp, err := a.client.Recommend(&client.RecommendRequest{
		MapName:        a.mapName,
		Mode:           a.mode,
		Rank:           a.rank,
		AllyBrawlers:   a.allyBrawlers,
		EnemyBrawlers:  a.enemyBrawlers,
		BannedBrawlers: a.bannedBrawlers,
		TopK:           &topK,
	})
	if err != nil {
		return err
	}

	a.recommendations = resp.Recommendations
	a.modelID = resp.ModelID
	a.message = "recommendations refreshed"
	return nil
}

func (a *App) render() {
	a.clearScreen()

	fmt.Println("BrawlPicks CLI")
	fmt.Println("==============")
	fmt.Printf("Map:    %s\n", a.mapName)
	fmt.Printf("Mode:   %s\n", a.mode)
	fmt.Printf("Rank:   %d\n", a.rank)
	fmt.Printf("Top K:  %d\n", a.topK)
	fmt.Println()
	fmt.Printf("Allies: %s\n", a.formatBrawlers(a.allyBrawlers))
	fmt.Printf("Enemy:  %s\n", a.formatBrawlers(a.enemyBrawlers))
	fmt.Printf("Banned: %s\n", a.formatBrawlers(a.bannedBrawlers))
	fmt.Println()
	fmt.Println("Recommendations")
	fmt.Println("---------------")
	if len(a.recommendations) == 0 {
		fmt.Println("none yet")
	} else {
		for index, item := range a.recommendations {
			fmt.Printf("%d. %s  %.4f\n", index+1, a.brawlers.NameForID(item.BrawlerID), item.Score)
		}
	}
	fmt.Println()
	if a.modelID != "" {
		fmt.Printf("Model:   %s\n", a.modelID)
	}
	fmt.Printf("Status:  %s\n", a.message)
	if a.showHelp {
		fmt.Println()
		fmt.Println("Commands")
		fmt.Println("--------")
		fmt.Println("map <name>")
		fmt.Println("mode <value>")
		fmt.Println("rank <number>")
		fmt.Println("topk <number>")
		fmt.Println("ally add|rm <brawler_name>")
		fmt.Println("enemy add|rm <brawler_name>")
		fmt.Println("ban add|rm <brawler_name>")
		fmt.Println("clear ally|enemy|ban|all")
		fmt.Println("rec")
		fmt.Println("quit")
	}
}

func (a *App) clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func parseIntArg(parts []string, index int, usage string) (int, error) {
	if len(parts) <= index {
		return 0, fmt.Errorf(usage)
	}

	value, err := strconv.Atoi(parts[index])
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", parts[index])
	}
	return value, nil
}

func addUnique(values []int, value int) []int {
	if contains(values, value) {
		return values
	}
	values = append(values, value)
	slices.Sort(values)
	return values
}

func removeValue(values []int, value int) []int {
	next := make([]int, 0, len(values))
	for _, current := range values {
		if current != value {
			next = append(next, current)
		}
	}
	return next
}

func contains(values []int, target int) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func (a *App) formatBrawlers(values []int) string {
	if len(values) == 0 {
		return "-"
	}

	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, a.brawlers.NameForID(value))
	}
	return strings.Join(parts, ", ")
}
