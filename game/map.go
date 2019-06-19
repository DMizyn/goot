package game

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/rwxsu/goot/parser"
)

const (
	// EOT End Of Tile
	EOT uint8 = 0xff
)

type Tile struct {
	Position
	Items     []*Item
	Creatures []*Creature
}

func (t *Tile) SetPosition(x, y uint16, z uint8) {
	t.X = x
	t.Y = y
	t.Z = z
}

func (t *Tile) AddItem(i *Item) {
	t.Items = append(t.Items, i)
}

func (t *Tile) AddCreature(c *Creature) {
	t.Creatures = append(t.Creatures, c)
}

type SectorPosition Position

type Column [32]*Tile
type Sector [32]*Column
type Map map[SectorPosition]*Sector

func (m *Map) SetTile(tile *Tile) {
	var spos SectorPosition
	spos.X = tile.X / 32
	spos.Y = tile.Y / 32
	spos.Z = tile.Z
	if (*m)[spos] == nil {
		return
	}
	(*m)[spos][tile.X%32][tile.Y%32] = tile
}

func (m *Map) GetTile(pos Position) *Tile {
	var spos SectorPosition
	spos.X = pos.X / 32
	spos.Y = pos.Y / 32
	spos.Z = pos.Z
	if (*m)[spos] == nil {
		return nil
	}
	return (*m)[spos][pos.X%32][pos.Y%32]
}

func (m *Map) LoadSector(filename string) {
	var p parser.Parser
	p.Filename = filename
	x, _ := strconv.Atoi(p.Filename[0:4])
	y, _ := strconv.Atoi(p.Filename[5:9])
	z, _ := strconv.Atoi(p.Filename[10:12])
	fmt.Printf("Loading %04d-%04d-%02d.sec ", x, y, z)
	begin := time.Now()

	if fileBytes, err := ioutil.ReadFile(p.Filename); err == nil {
		p.Buffer = bytes.NewBuffer(fileBytes)
	} else {
		panic(err)
	}

	spos := SectorPosition{X: (uint16)(x), Y: (uint16)(y), Z: (uint8)(z)}
	(*m)[spos] = new(Sector)
	for offsetX := (uint16)(0); offsetX < 32; offsetX++ {
		(*m)[spos][offsetX] = new(Column)
		for offsetY := (uint16)(0); offsetY < 32; offsetY++ {
			var tile Tile
			tile.SetPosition(spos.X*32+offsetX, spos.Y*32+offsetY, spos.Z)
			p.NextToken() // skip offsetX
			p.NextToken() // skip offsetY
			itemids := p.NextToken()
			switch itemids := itemids.(type) {
			case []int:
				for _, id := range itemids {
					tile.AddItem(&Item{ID: (uint16)(id)})
				}
				break
			default:
				panic("OOOOPS")
			}
			(*m)[spos][offsetX][offsetY] = &tile
		}
	}
	fmt.Printf("[%v]\n", time.Since(begin))
}