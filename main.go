package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"go-game/entities"
)

type Game struct {
	GameState            string
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
	attackCounter        int
	Points               int
	spawnEnemies         bool
	killEnemies          bool
	showColliders        bool
	enemiesFollowsPlayer bool
	globalVolume         float64
}

func (g *Game) Update() error {

	if g.GameState == "RUNNING" {

		g.HandleControls()

		g.UpdatePlayer()

		g.spawnEnemy()

		g.updateEnemies()

		g.UpdateCamera()

		g.UpdateHitbox()

		MusicLoop()

		g.Tick++
		g.attackCounter++
	} else if g.GameState == "GAMEOVER" {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.RestartGame()
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.GameState == "RUNNING" || g.GameState == "PAUSED" {
		screen.Fill(color.RGBA{120, 180, 255, 255})

		opts := ebiten.DrawImageOptions{}

		g.DrawTilemap(screen, opts)

		g.DrawPlayer(screen, opts)

		g.DrawEnemies(screen, opts)

		g.DrawColliders(screen)

		g.DrawHitbox(screen)

		g.DrawAttack(screen, opts)

		g.DrawPotions(screen, opts)

		// draw fps counter
		msg := fmt.Sprintf(
			"TPS: %0.2f\nEnemies: %d\nScore: %d\n\nDamage: %d",
			ebiten.ActualTPS(),
			len(g.enemies),
			g.Points,
			g.player.Attack.Damage,
		)
		ebitenutil.DebugPrintAt(screen, msg, 0, 0)
		ebitenutil.DebugPrintAt(screen, "Controls: [W/A/S/D] [LButton] [Q/E/G] [F] [R]", 0, 225)

		g.DrawHUD(screen)

	} else if g.GameState == "GAMEOVER" {
		MusicPlayer.Pause()
		GameoverSoundPlayer.Play()
		ebitenutil.DebugPrintAt(screen, "GAME OVER", 130, 110)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.Points), 130, 125)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Press R to try again..."), 90, 170)

	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func initializeGame() *Game {
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

	playerHealth := uint(10)

	game := Game{
		GameState:    "RUNNING",
		globalVolume: 0.1,
		player: &entities.Player{
			Experience: 0,
			Level:      1,
			Invencible: false,
			Status:     "IDLE",
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   150.0,
				Y:   150.0,
			},
			MaxHealth: playerHealth,
			Health:    &playerHealth,
			Attack: entities.Attack{
				Damage: 5,
				Img:    attackImg,
			},
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
		spawnEnemies:         true,
		killEnemies:          false,
		enemiesFollowsPlayer: true,
	}

	game.softColliders = append(game.softColliders, game.player.Collider)

	return &game
}

func (g *Game) RestartGame() {
	*g = *initializeGame()
	MusicPlayer.Rewind()
	MusicPlayer.Play()
}

var (
	playerImg   *ebiten.Image
	attackImg   *ebiten.Image
	tilemapImg  *ebiten.Image
	tilemapJSON *TilemapJSON
	tilesets    []Tileset
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Ebiten Zelda Clone")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	var err error
	// load the image from file
	playerImg, _, err = ebitenutil.NewImageFromFile("assets/images/player.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	attackImg, _, err = ebitenutil.NewImageFromFile("assets/images/attack.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	tilemapImg, _, err = ebitenutil.NewImageFromFile("assets/images/TilesetFloor.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	tilemapJSON, err = NewTilemapJSON("assets/maps/spawn.json")
	if err != nil {
		log.Fatal(err)
	}

	tilesets, _, err = tilemapJSON.GenTilesets()
	if err != nil {
		log.Fatal(err)
	}

	if err := CreateMusicSound("assets/sounds/music.ogg"); err != nil {
		log.Fatal(err)
	}

	if err := CreateHitSound("assets/sounds/hit.wav"); err != nil {
		log.Fatal(err)
	}

	if err := CreateKillSound("assets/sounds/kill.wav"); err != nil {
		log.Fatal(err)
	}

	if err := CreatePlayerHitSound("assets/sounds/playerhit.wav"); err != nil {
		log.Fatal(err)
	}

	if err := CreateGameoverSound("assets/sounds/gameover.wav"); err != nil {
		log.Fatal(err)
	}

	if err := CreateHealSound("assets/sounds/heal.wav"); err != nil {
		log.Fatal(err)
	}

	if err := CreateLevelUpSoundPlayer("assets/sounds/levelup.wav"); err != nil {
		log.Fatal(err)
	}

	game := initializeGame()

	MusicPlayer.Play()
	SetVolumeValue(game.globalVolume)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
