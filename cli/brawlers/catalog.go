package brawlers

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Catalog struct {
	idToName map[int]string
	nameToID map[string]int
}

type rawCatalog struct {
	Items []*rawBrawler `json:"items"`
}

type rawBrawler struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func Load(path string) (*Catalog, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	payload := new(rawCatalog)
	if err := json.Unmarshal(raw, payload); err != nil {
		return nil, err
	}

	catalog := &Catalog{
		idToName: make(map[int]string, len(payload.Items)),
		nameToID: make(map[string]int, len(payload.Items)),
	}

	for _, item := range payload.Items {
		if item == nil || item.Name == "" {
			continue
		}
		catalog.idToName[item.ID] = item.Name
		catalog.nameToID[normalizeName(item.Name)] = item.ID
	}

	return catalog, nil
}

func (c *Catalog) IDForName(name string) (int, error) {
	id, ok := c.nameToID[normalizeName(name)]
	if !ok {
		return 0, fmt.Errorf("unknown brawler: %s", name)
	}
	return id, nil
}

func (c *Catalog) NameForID(id int) string {
	name, ok := c.idToName[id]
	if !ok {
		return fmt.Sprintf("%d", id)
	}
	return name
}

func normalizeName(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
