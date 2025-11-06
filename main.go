package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	frameWidth  = 320
	frameHeight = 240
	playerY     = frameHeight/2 + frameHeight/4
)

var (
	playerImage *ebiten.Image
)

type Game struct {
	player   Player
	ennemies []string
}

type Player struct {
	playerX int
	vX      int
	vY      int
}

func (g *Game) Update() error {
	g.keyPressed()
	g.player.update()
	return nil
}

func (g *Game) keyPressed() {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		moveRight(g)
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		moveLeft(g)
	}
}

func moveRight(g *Game) {
	g.player.vX += 8
	print("move right")
}

func moveLeft(g *Game) {
	g.player.vX -= 8
	print("move left")
}

func (p *Player) update() {
	p.playerX += p.vX
	if p.vX != 0 {
		p.vX = 0
	}
}

func (p *Player) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.playerX), float64(playerY))
	screen.DrawImage(playerImage, op)
}

func init() {
	img, _, err := ebitenutil.NewImageFromFile("assets/player.png")
	if err != nil {
		log.Fatal(err)
	}
	playerImage = ebiten.NewImageFromImage(img)
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(frameWidth/2-16, frameHeight/2+frameHeight/4)
	g.player.draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return frameWidth, frameHeight
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Space Invaders")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
