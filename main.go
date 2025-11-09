package main

import (
	_ "image/png"
	"log"
	"math/rand"

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
	player       Player
	playerBullet *PlayerBullet
	enemies      Enemies
	enemyBullets EnemyBullets
}

type Player struct {
	playerX  int
	vPlayerX int
	vPlayerY int
}

type PlayerBullet struct {
	bulletX  int
	bulletY  int
	vBulletY int
}

type EnemyBullet struct {
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

type EnemyBullets []EnemyBullet

func (g *Game) keyPressed() {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		moveRight(g)
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		moveLeft(g)
	}
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		playerShoot(g)
	}
}

func moveRight(g *Game) {
	g.player.vPlayerX += 4
}

func moveLeft(g *Game) {
	g.player.vPlayerX -= 4
}

func playerShoot(g *Game) {
	if g.playerBullet == nil {
		g.playerBullet = &PlayerBullet{bulletX: g.player.playerX, bulletY: playerY + 4, vBulletY: -4}
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

	// update positions using velocities and possibly shoot.
	for i := range *e {
		(*e)[i].enemyShoot(g)
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

func (e Enemy) enemyShoot(g *Game) {
	if randomNumber := rand.Float64(); randomNumber < 0.0005 {
		// Enemy shoots a bullet (not implemented)
		print("Enemy at (", e.enemyX, ",", e.enemyY, ") shoots!\n")
		g.enemyBullets = append(g.enemyBullets, EnemyBullet{bulletX: e.enemyX, bulletY: e.enemyY + 16, vBulletY: 4})
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

	if g.playerBullet == nil || len(g.enemies) == 0 {
		return false
	}

	for i, enemy := range g.enemies {

		// Check if bullet is within vertical range of enemy
		bulletInVerticalRange := g.playerBullet.bulletY == enemy.enemyY

		// Check if bullet hits enemy horizontally
		bulletHitsHorizontally := g.playerBullet.bulletX >= enemy.enemyX-enemySize/2 &&
			g.playerBullet.bulletX <= enemy.enemyX+enemySize/2

		if bulletInVerticalRange && bulletHitsHorizontally {
			print(len(g.enemies), " Enemies left\n")
			g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
			g.playerBullet = nil
			return true
		}
	}
	return false
}

func (g *Game) playerHit() bool {
	for _, bullet := range g.enemyBullets {
		// Check if bullet is within vertical range of player
		bulletInVerticalRange := bullet.bulletY == playerY

		// Check if bullet hits player horizontally
		bulletHitsHorizontally := bullet.bulletX >= g.player.playerX-16 &&
			bullet.bulletX <= g.player.playerX+16
		if bulletInVerticalRange && bulletHitsHorizontally {
			return true
		}
	}
	return false
}

func (p *Player) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.playerX), float64(playerY))
	screen.DrawImage(playerImage, op)
}

func (b *PlayerBullet) update(g *Game) {
	if b == nil {
		return
	}
	b.bulletY += b.vBulletY
	if b.bulletY < 0 {
		g.playerBullet = nil
	}
}

func (b *PlayerBullet) draw(screen *ebiten.Image) {
	if b == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.bulletX), float64(b.bulletY))
	screen.DrawImage(bulletImage, op)
}

func (b *EnemyBullets) update() {
	if b == nil {
		return
	}
	// Iterate backwards when removing elements to avoid index errors that occur if we remove items while iterating forward with a range
	for i := len(*b) - 1; i >= 0; i-- {
		(*b)[i].bulletY += (*b)[i].vBulletY
		if (*b)[i].bulletY > frameHeight {
			*b = append((*b)[:i], (*b)[i+1:]...)
		}
	}
}

func (b *EnemyBullets) draw(screen *ebiten.Image) {
	if b == nil {
		return
	}
	for _, bullet := range *b {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(bullet.bulletX), float64(bullet.bulletY))
		screen.DrawImage(bulletImage, op)
	}
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
	g.playerBullet.draw(screen)
	g.enemies.draw(screen)
	g.enemyBullets.draw(screen)
}

func (g *Game) Update() error {
	g.keyPressed()
	g.player.update()
	g.playerBullet.update(g)
	g.enemies.update(g)
	g.enemyBullets.update()
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
