package game

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rwxsu/goot/parser"
)

// SectorPosition contains the position equal to the .sec filename
//
// To get the real in-game position, calculate:
//		x = SectorPosition.X * 32 + x offset
//		y = SectorPosition.Y * 32 + y offset
//		z = SectorPosition.Z
//
//		E.g.: 0999-0998-07.sec (02-04: Content={102})
//	 		x = 999 * 32 + 2
//	 		y = 998 * 32 + 4
//			z = 7
type SectorPosition Position

// Column 32 tiles (stacked below each other visually)
type Column [32]*Tile

// Sector 32*32=256 tiles (32 columns)
type Sector [32]*Column

// Map is a collection of sectors loaded from data/map/sectors/*.sec
type Map map[SectorPosition]*Sector

func (m *Map) SetTile(tile *Tile) {
	t := m.GetTile(tile.Position)
	if t != nil {
		t = tile
	}
}

// GetTile gets the SectorPosition by dividing the real x and y position by 32.
// Mod is used to get the remainder, which is the x and y offset.
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

// LoadSector loads a single .sec file and stores the sector at the correct
// SectorPosition extracted from the filename in the map.
func (m *Map) LoadSector(filename string) {
	var p parser.Parser
	p.Filename = filename
	dirOffset := len(filepath.Dir(p.Filename)) + 1
	x, _ := strconv.Atoi(p.Filename[dirOffset+0 : dirOffset+4])
	y, _ := strconv.Atoi(p.Filename[dirOffset+5 : dirOffset+9])
	z, _ := strconv.Atoi(p.Filename[dirOffset+10 : dirOffset+12])
	fmt.Printf(":: Loading %04d-%04d-%02d.sec ", x, y, z)
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
				panic("in LoadSector: could not get item ids")
			}
			if len(tile.Items) > 0 {
				(*m)[spos][offsetX][offsetY] = &tile
			} else {
				(*m)[spos][offsetX][offsetY] = nil
			}
		}
	}

	fmt.Printf("[%v]\n", time.Since(begin))
}

// MoveCreature on map
func (m *Map) MoveCreature(c *Creature, pos Position, direction DirectionType) bool {
	from := m.GetTile(c.Position)
	to := m.GetTile(pos)
	if from == nil || to == nil {
		return false
	}
	if !from.RemoveCreature(c) {
		return false
	}
	to.AddCreature(c)
	c.Position = pos
	c.Direction = direction
	return true
}

// func (m *Map) GetSpectators(pos Position) []*Creature {
// 	pos.X = pos.X - 8
// 	pos.Y = pos.Y - 6
// 	var spectators []*Creature
// 	for x := (uint16)(0); x < 18; x++ {
// 		for y := (uint16)(0); y < 14; y++ {
// 			tile := m.GetTile(pos)
// 			if tile != nil {
// 				spectators = append(spectators, tile.Creatures...)
// 			}
// 		}
// 	}
// 	return spectators
// }
