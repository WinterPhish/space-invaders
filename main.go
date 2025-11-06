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
	bulletImage *ebiten.Image
)

type Game struct {
	player   Player
	bullet   *Bullet
	ennemies []string
}

type Player struct {
	playerX  int
	vPlayerX int
	vPlayerY int
}

type Bullet struct {
	bulletX  int
	bulletY  int
	vBulletY int
}

func (g *Game) Update() error {
	g.keyPressed()
	g.player.update()
	g.bullet.update(g)
	return nil
}

func (g *Game) keyPressed() {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		moveRight(g)
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		moveLeft(g)
	}
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		shoot(g)
	}
}

func moveRight(g *Game) {
	g.player.vPlayerX += 4
}

func moveLeft(g *Game) {
	g.player.vPlayerX -= 4
}

func shoot(g *Game) {
	if g.bullet == nil {
		g.bullet = &Bullet{bulletX: g.player.playerX, bulletY: playerY, vBulletY: -4}
	}
}

func (p *Player) update() {
	p.playerX += p.vPlayerX
	if p.vPlayerX != 0 {
		p.vPlayerX = 0
	}
}

func (p *Player) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.playerX), float64(playerY))
	screen.DrawImage(playerImage, op)
}

func (b *Bullet) update(g *Game) {
	if b == nil {
		return
	}
	b.bulletY += b.vBulletY
	if b.bulletY < 0 {
		g.bullet = nil
	}
}

func (b *Bullet) draw(screen *ebiten.Image) {
	if b == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.bulletX), float64(b.bulletY))
	screen.DrawImage(bulletImage, op)
}

func init() {
	img, _, err := ebitenutil.NewImageFromFile("assets/player.png")
	if err != nil {
		log.Fatal(err)
	}
	playerImage = ebiten.NewImageFromImage(img)
	img, _, err = ebitenutil.NewImageFromFile("assets/bullet.png")
	if err != nil {
		log.Fatal(err)
	}
	bulletImage = ebiten.NewImageFromImage(img)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.player.draw(screen)
	g.bullet.draw(screen)
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
