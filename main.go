package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"go-game/entities"
)

type Game struct {
	// the image and position variables for our player
	player               *entities.Player
	enemies              []*entities.Enemy
	potions              []*entities.Potion
	tilemapJSON          *TilemapJSON
	tilesets             []Tileset
	tilemapImg           *ebiten.Image
	cam                  *Camera
	hardColliders        []image.Rectangle
	softColliders        []entities.Collider
	Tick                 int
	Points               int
	spawnEnemies         bool
	killEnemies          bool
	showColliders        bool
	enemiesFollowsPlayer bool
	audioContext         *audio.Context
	musicPlayer          *audio.Player
}

func (g *Game) Update() error {
	g.HandleControls()

	g.UpdatePlayer()

	g.spawnEnemy()

	g.updateEnemies()

	g.UpdateCamera()

	g.UpdateHitbox()

	g.Tick++

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}

	g.DrawTilemap(screen, opts)

	g.DrawPlayer(screen, opts)

	g.DrawEnemies(screen, opts)

	g.DrawColliders(screen)

	g.DrawHitbox(screen)

	// draw fps counter
	msg := fmt.Sprintf(
		"TPS: %0.2f\nEnemies: %d\nPoints: %d",
		ebiten.ActualTPS(),
		len(g.enemies),
		g.Points,
	)
	ebitenutil.DebugPrintAt(screen, msg, 0, 0)
	ebitenutil.DebugPrintAt(screen, "Controls: [W/A/S/D] [LButton] [Q/E/G] [F] [R]", 0, 225)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Ebiten Zelda Clone")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// load the image from file
	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/player.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	// potionImg, _, err := ebitenutil.NewImageFromFile("assets/images/potion.png")
	// if err != nil {
	// 	// handle error
	// 	log.Fatal(err)
	// }

	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/images/TilesetFloor.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	tilemapJSON, err := NewTilemapJSON("assets/maps/spawn.json")
	if err != nil {
		log.Fatal(err)
	}

	tilesets, _, err := tilemapJSON.GenTilesets()
	if err != nil {
		log.Fatal(err)
	}

	// collider for buildings
	hardColliders := make([]image.Rectangle, 0)
	for layerIndex, layer := range tilemapJSON.Layers {
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

			img := tilesets[layerIndex].Img(id)

			// check if the tile is a building
			if layer.Name == "buildings" {
				// create a collider based on the position and size of the image
				collider := image.Rect(
					x,
					y-16,
					x+img.Bounds().Dx(),
					y-img.Bounds().Dy()-16,
				)
				hardColliders = append(hardColliders, collider)
			}
		}
	}

	game := Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   150.0,
				Y:   150.0,
			},
			Health: 3,
			Damage: 10,
			Collider: entities.Collider{
				Rect: &entities.FloatRect{
					MinX: 150.0,
					MinY: 150.0,
					MaxX: 150.0 + 16,
					MaxY: 150.0 + 16,
				},
				Weight: 40,
			},
			Hitbox: &entities.Hitbox{
				Vertices: [4][2]float64{
					{0, 0},
					{0, 0},
					{0, 0},
					{0, 0},
				},
			},
		},
		enemies:              []*entities.Enemy{},
		potions:              []*entities.Potion{},
		tilemapJSON:          tilemapJSON,
		tilemapImg:           tilemapImg,
		tilesets:             tilesets,
		cam:                  NewCamera(0.0, 0.0),
		hardColliders:        hardColliders,
		spawnEnemies:         false,
		killEnemies:          false,
		enemiesFollowsPlayer: true,
	}

	game.softColliders = append(game.softColliders, game.player.Collider)

	if err := game.PlayOGGSound("assets/sounds/music.ogg"); err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
