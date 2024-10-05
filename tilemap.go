package main

import (
	"embed"
	_ "embed"
	"encoding/json"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
)

// data we want for one layer in our list of layers
type TilemapLayerJSON struct {
	Data   []int  `json:"data"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Name   string `json:"name"`
}

//go:embed assets/maps/*.json
var EmbeddedFiles embed.FS

// all layers in a tilemap
type TilemapJSON struct {
	Layers []TilemapLayerJSON `json:"layers"`
	// raw data for each tileset (path, gid)
	Tilesets []map[string]any `json:"tilesets"`
}

// temp function to generate all of our tilesets and return a slice of them
func (t *TilemapJSON) GenTilesets() ([]Tileset, []*Tileset, error) {
	tilesets := make([]Tileset, 0)
	buildings := make([]*Tileset, 0)

	for _, tilesetData := range t.Tilesets {
		// convert map relative path to project relative path
		tilesetPath := path.Join("assets/maps/", tilesetData["source"].(string))
		tileset, isBuilding, err := NewTileset(tilesetPath, int(tilesetData["firstgid"].(float64)))
		if err != nil {
			return nil, nil, err
		}

		if isBuilding {
			buildings = append(buildings, &tileset)
		}

		tilesets = append(tilesets, tileset)
	}

	return tilesets, buildings, nil
}

// opens the file, parses it, and returns the json object + potential error
func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	// contents, err := os.ReadFile(filepath)
	// if err != nil {
	// 	return nil, err
	// }

	// Lê o conteúdo do arquivo embutido
	contents, err := EmbeddedFiles.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tilemapJSON TilemapJSON
	err = json.Unmarshal(contents, &tilemapJSON)
	if err != nil {
		return nil, err
	}

	return &tilemapJSON, nil
}

func (g *Game) DrawTilemap(screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	for layerIndex, layer := range g.tilemapJSON.Layers {
		// loop over the tiles in the layer data
		for index, id := range layer.Data {

			if id == 0 {
				continue
			}

			// get the tile position of the tile
			x := index % layer.Width
			y := index / layer.Width

			// convert the tile position to pixel position
			x *= 16
			y *= 16

			img := g.tilesets[layerIndex].Img(id)

			opts.GeoM.Translate(float64(x), float64(y))

			// opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 16))

			opts.GeoM.Translate(g.cam.X, g.cam.Y)

			screen.DrawImage(img, &opts)

			// reset the opts for the next tile
			opts.GeoM.Reset()
		}
	}
}
