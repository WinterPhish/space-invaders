package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	frameWidth  = 640
	frameHeight = 480
	playerY     = frameHeight/2 + frameHeight/4 + frameHeight/8
	enemyY      = frameHeight/2 - frameHeight/4 - frameHeight/8

	_ = iota
	Squid
	Octopus
	Crab
)

var (
	playerImage  *ebiten.Image
	bulletImage  *ebiten.Image
	squidImage   *ebiten.Image
	crabImage    *ebiten.Image
	octopusImage *ebiten.Image
)

type Game struct {
	player  Player
	bullet  *Bullet
	enemies Enemies
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

type Enemy struct {
	enemyX    int
	enemyY    int
	vEnemyY   int
	vEnemyX   int
	enemyType int
}

type Enemies []Enemy

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

func spawnWave(g *Game, squid int, octopus int, crab int) {
	// create 15 enemies spaced horizontally
	for i := range squid {
		g.enemies = append(g.enemies, Enemy{enemyX: i * 32, enemyY: enemyY, vEnemyY: 0, vEnemyX: 0, enemyType: Squid})
	}
	for i := range crab {
		g.enemies = append(g.enemies, Enemy{enemyX: i * 32, enemyY: enemyY + 32, vEnemyY: 0, vEnemyX: 0, enemyType: Crab})
	}
	for i := range octopus {
		g.enemies = append(g.enemies, Enemy{enemyX: i * 32, enemyY: enemyY + 64, vEnemyY: 0, vEnemyX: 0, enemyType: Octopus})
	}
}

func (e *Enemies) update(g *Game) {
	if len(*e) == 0 {
		spawnWave(g, 10, 10, 10)
		print("Spawned Enemy\n")
	}

	// Detect if any enemy hit a boundary. Do not modify velocities
	hitRight := false
	hitLeft := false
	for i := range *e {
		if (*e)[i].enemyX > frameWidth-32 {
			hitRight = true
		}
		if (*e)[i].enemyX <= 0 {
			hitLeft = true
		}
	}

	// Apply direction change once per frame
	if hitRight {
		for j := range *e {
			(*e)[j].vEnemyX = -2
			(*e)[j].enemyY += 16
		}
	} else if hitLeft {
		for j := range *e {
			(*e)[j].vEnemyX = 2
		}
	}

	// update positions using velocities.
	for i := range *e {
		(*e)[i].enemyX += (*e)[i].vEnemyX
		(*e)[i].enemyY += (*e)[i].vEnemyY
	}
}

func (e *Enemies) draw(screen *ebiten.Image) {
	for _, enemy := range *e {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(enemy.enemyX), float64(enemy.enemyY))
		switch enemy.enemyType {
		case Squid:
			screen.DrawImage(squidImage, op)
		case Crab:
			screen.DrawImage(crabImage, op)
		case Octopus:
			screen.DrawImage(octopusImage, op)
		}

	}
}

func (p *Player) update() {
	p.playerX += p.vPlayerX
	if p.vPlayerX != 0 {
		p.vPlayerX = 0
	}
}

func (g *Game) enemyHit() bool {
	const (
		enemySize   = 32
		enemyHeight = 32
	)

	if g.bullet == nil || len(g.enemies) == 0 {
		return false
	}

	for i, enemy := range g.enemies {

		// Check if bullet is within vertical range of enemy
		bulletInVerticalRange := g.bullet.bulletY == enemy.enemyY

		// Check if bullet hits enemy horizontally
		bulletHitsHorizontally := g.bullet.bulletX >= enemy.enemyX-enemySize/2 &&
			g.bullet.bulletX <= enemy.enemyX+enemySize/2

		if bulletInVerticalRange && bulletHitsHorizontally {
			print(len(g.enemies), " Enemies left\n")
			g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
			g.bullet = nil
			return true
		}
	}
	return false
}

func (g *Game) playerHit() bool {
	// for when enemies will shoot bullets
	return false
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
	img, _, err = ebitenutil.NewImageFromFile("assets/squid.png")
	if err != nil {
		log.Fatal(err)
	}
	squidImage = ebiten.NewImageFromImage(img)
	img, _, err = ebitenutil.NewImageFromFile("assets/crab.png")
	if err != nil {
		log.Fatal(err)
	}
	crabImage = ebiten.NewImageFromImage(img)
	img, _, err = ebitenutil.NewImageFromFile("assets/octopus.png")
	if err != nil {
		log.Fatal(err)
	}
	octopusImage = ebiten.NewImageFromImage(img)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.player.draw(screen)
	g.bullet.draw(screen)
	g.enemies.draw(screen)
}

func (g *Game) Update() error {
	g.keyPressed()
	g.player.update()
	g.bullet.update(g)
	g.enemies.update(g)
	if g.enemyHit() {
		print("Enemy Hit!\n")
	}
	if g.playerHit() {
		log.Fatal("Game Over")
	}
	return nil
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
