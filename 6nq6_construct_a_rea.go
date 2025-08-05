package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// GamePrototype represents a real-time game prototype simulator
type GamePrototype struct {
	window  *pixelgl.Window
	player   *pixel.Sprite
	enemies  []*pixel.Sprite
	bullets  []*pixel.Sprite
	score    int
	level    int
	fps      int
	lastTime time.Time
}

// NewGamePrototype creates a new game prototype simulator
func NewGamePrototype(winTitle string) (*GamePrototype, error) {
	cfg := pixelgl.WindowConfig{
		Title:  winTitle,
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return nil, err
	}

	player, err := loadSprite("player.png")
	if err != nil {
		return nil, err
	}

	return &GamePrototype{
		window: win,
		player: player,
		enemies: make([]*pixel.Sprite, 0),
		bullets: make([]*pixel.Sprite, 0),
		score:  0,
		level:  1,
		fps:    60,
	}, nil
}

func (g *GamePrototype) run() {
	last := time.Now()
	for !g.window.Closed() {
		g.window.Clear(pixel.RGB(0.2, 0.2, 0.2))
		g.update(time.Since(last).Seconds())
		g.draw()
		g.window.Update()
		last = time.Now()
	}
}

func (g *GamePrototype) update(dt float64) {
	// Update player position
	g.player.Pos = g.window.MousePosition()

	// Update enemy positions
	for _, enemy := range g.enemies {
		enemy.Pos = enemy.Pos.Add(pixel.V(math.Sin(g.lastTime.Sub(time.Now()).Seconds()), 0))
	}

	// Update bullet positions
	for _, bullet := range g.bullets {
		bullet.Pos = bullet.Pos.Add(pixel.V(0, 10))
	}

	// Check collisions
	for _, enemy := range g.enemies {
		if g.player.Frame().Intersect(enemy.Frame()) {
			g.score++
			g.enemies = removeEnemy(g.enemies, enemy)
		}
	}

	// Spawn new enemies
	if len(g.enemies) < g.level {
		g.spawnEnemy()
	}

	g.lastTime = time.Now()
}

func (g *GamePrototype) draw() {
	g.player.Draw(g.window, pixel.IM)
	for _, enemy := range g.enemies {
		enemy.Draw(g.window, pixel.IM)
	}
	for _, bullet := range g.bullets {
		bullet.Draw(g.window, pixel.IM)
	}

	g.window.SetTitle(fmt.Sprintf("FPS: %d | Score: %d | Level: %d", g.fps, g.score, g.level))
}

func (g *GamePrototype) spawnEnemy() {
	enemy, err := loadSprite("enemy.png")
	if err != nil {
		return
	}
	enemy.Pos = pixel.V(rand.Intn(1024), rand.Intn(768))
	g.enemies = append(g.enemies, enemy)
}

func loadSprite(filename string) (*pixel.Sprite, error) {
	pic, err := pixelgl.LoadPicture(filename)
	if err != nil {
		return nil, err
	}
	sprite := pixel.NewSprite(pic, pic.Bounds())
	return sprite, nil
}

func removeEnemy(enemies []*pixel.Sprite, enemy *pixel.Sprite) []*pixel.Sprite {
	for i, e := range enemies {
		if e == enemy {
			enemies = append(enemies[:i], enemies[i+1:]...)
			return enemies
		}
	}
	return enemies
}

func main() {
	game, err := NewGamePrototype("Real-Time Game Prototype Simulator")
	if err != nil {
		panic(err)
	}
	game.run()
}