package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	frameWidth  = 640
	frameHeight = 480
	playerY     = frameHeight/2 + frameHeight/4 + frameHeight/8
	enemyY      = frameHeight/2 - frameHeight/4 - frameHeight/6

	_ GameMode = iota
	ModeStart
	ModePlaying
	ModePause
	ModeGameOver

	Squid = iota
	Octopus
	Crab
	Ufo
)

var (
	playerImage   *ebiten.Image
	bulletImage   *ebiten.Image
	squidImage    *ebiten.Image
	crabImage     *ebiten.Image
	octopusImage  *ebiten.Image
	squidImage2   *ebiten.Image
	crabImage2    *ebiten.Image
	octopusImage2 *ebiten.Image
	ufoImage      *ebiten.Image
	deathImage    *ebiten.Image
)

type Game struct {
	player       Player
	playerBullet *PlayerBullet
	enemies      Enemies
	enemyBullets EnemyBullets
	deathAnims   []DeathAnimation
	mode         GameMode
	ufo          *Enemy
	score        int
	lives        int
	level        int
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
	dead      bool
	frame     int
}

type GameMode int

type Enemies []Enemy

type EnemyBullets []EnemyBullet

type DeathAnimation struct {
	enemyX int
	enemyY int
	frame  int
}

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
	if g.player.playerX+32 >= frameWidth {
		return
	}
	g.player.vPlayerX += 4
}

func moveLeft(g *Game) {
	if g.player.playerX <= 0 {
		return
	}
	g.player.vPlayerX -= 4
}

func playerShoot(g *Game) {
	if g.playerBullet == nil {
		g.playerBullet = &PlayerBullet{bulletX: g.player.playerX, bulletY: playerY + 4, vBulletY: -4}
	}
}

func spawnWave(g *Game) {
	for i := range 11 {
		g.enemies = append(g.enemies, Enemy{enemyX: i * 32, enemyY: enemyY, vEnemyY: 0, vEnemyX: 0, enemyType: Squid})
	}
	for i := range 11 {
		g.enemies = append(g.enemies, Enemy{enemyX: i * 32, enemyY: enemyY + 32, vEnemyY: 0, vEnemyX: 0, enemyType: Crab})
	}
	for i := range 11 {
		g.enemies = append(g.enemies, Enemy{enemyX: i * 32, enemyY: enemyY + 64, vEnemyY: 0, vEnemyX: 0, enemyType: Crab})
	}
	for i := range 11 {
		g.enemies = append(g.enemies, Enemy{enemyX: i * 32, enemyY: enemyY + 96, vEnemyY: 0, vEnemyX: 0, enemyType: Octopus})
	}
	for i := range 11 {
		g.enemies = append(g.enemies, Enemy{enemyX: i * 32, enemyY: enemyY + 128, vEnemyY: 0, vEnemyX: 0, enemyType: Octopus})
	}
	g.level++
}

func (e *Enemies) update(g *Game) {
	if len(*e) == 0 {
		spawnWave(g)
		print("Spawned Enemy on level ", g.level, "\n")
	}

	if randomNumber := rand.Float64(); randomNumber < 0.001 && g.ufo == nil {
		g.ufo = &Enemy{enemyX: 10, enemyY: 16, vEnemyX: 2, enemyType: Ufo}
		print("UFO Spawned!\n")
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

	if g.ufo != nil {
		if g.ufo.enemyX > frameWidth {
			g.ufo = nil
		} else {
			g.ufo.enemyX += g.ufo.vEnemyX
		}
	}

	// Apply direction change once per frame
	if hitRight {
		for j := range *e {
			(*e)[j].vEnemyX = -1
			(*e)[j].enemyY += 16
		}
	} else if hitLeft {
		for j := range *e {
			(*e)[j].vEnemyX = 1
			(*e)[j].enemyY += 16
		}
	}

	// update positions using velocities and possibly shoot.
	for i := range *e {
		enemy := &(*e)[i]
		enemy.enemyShoot(g)
		enemy.frame++
		if enemy.frame > 60 {
			enemy.frame = 0
		}
		enemy.enemyX += enemy.vEnemyX
		enemy.enemyY += enemy.vEnemyY
		if enemy.enemyY >= playerY {
			g.mode = ModeGameOver
		}
	}
}

func (e *Enemies) draw(screen *ebiten.Image) {
	for _, enemy := range *e {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(enemy.enemyX), float64(enemy.enemyY))
		switch enemy.enemyType {
		case Squid:
			if enemy.frame < 30 {
				screen.DrawImage(squidImage, op)
			} else {
				screen.DrawImage(squidImage2, op)
			}
		case Crab:
			if enemy.frame < 30 {
				screen.DrawImage(crabImage, op)
			} else {
				screen.DrawImage(crabImage2, op)
			}
		case Octopus:
			if enemy.frame < 30 {
				screen.DrawImage(octopusImage, op)
			} else {
				screen.DrawImage(octopusImage2, op)
			}
		}
	}
}

func (g *Game) updateDeathAnims() {
	for i := len(g.deathAnims) - 1; i >= 0; i-- {
		g.deathAnims[i].frame++
		if g.deathAnims[i].frame > 5 {
			g.deathAnims = append(g.deathAnims[:i], g.deathAnims[i+1:]...)
		}
	}
}

func (g *Game) drawDeathAnims(screen *ebiten.Image) {
	for _, anim := range g.deathAnims {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(anim.enemyX), float64(anim.enemyY))
		screen.DrawImage(deathImage, op)
	}
}

func (e *Enemy) draw(screen *ebiten.Image) {
	if e == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(e.enemyX), float64(e.enemyY))
	screen.DrawImage(ufoImage, op)
}

func (e Enemy) enemyShoot(g *Game) {
	if randomNumber := rand.Float64(); randomNumber < 0.0005+float64(g.level)*0.0001 {
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

	if g.playerBullet == nil || len(g.enemies) == 0 {
		return false
	}

	for i, enemy := range g.enemies {

		// Check if bullet is within vertical range of enemy
		if enemyHitCheck(&enemy, g) {
			g.deathAnims = append(g.deathAnims, DeathAnimation{
				enemyX: enemy.enemyX,
				enemyY: enemy.enemyY,
				frame:  0,
			})
			g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
			g.playerBullet = nil
			if enemy.enemyType == Squid {
				g.score += 10
			}
			if enemy.enemyType == Crab {
				g.score += 20
			}
			if enemy.enemyType == Octopus {
				g.score += 30
			}
			print("Score: ", g.score, "\n")
			return true
		}
	}
	return false
}

func ufoHit(ufo *Enemy, g *Game) bool {
	if ufo == nil || g.playerBullet == nil {
		return false
	}
	if enemyHitCheck(ufo, g) {
		g.deathAnims = append(g.deathAnims, DeathAnimation{
			enemyX: ufo.enemyX,
			enemyY: ufo.enemyY,
			frame:  0,
		})
		g.ufo = nil
		g.score += 100
		return true
	}
	return false
}

func enemyHitCheck(enemy *Enemy, g *Game) bool {
	const (
		enemySize   = 32
		enemyHeight = 32
	)
	// Check if bullet is within vertical range of enemy
	bulletInVerticalRange := g.playerBullet.bulletY == enemy.enemyY

	// Check if bullet hits enemy horizontally
	bulletHitsHorizontally := g.playerBullet.bulletX >= enemy.enemyX-enemySize/2 &&
		g.playerBullet.bulletX <= enemy.enemyX+enemySize/2

	if bulletInVerticalRange && bulletHitsHorizontally {
		return true
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

func (g *Game) scoreDisplay(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Score: "+fmt.Sprint(g.score), 10, 10)
}

func (g *Game) livesDisplay(screen *ebiten.Image) {
	for i := 0; i < g.lives; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(16+i*32), frameHeight-30)
		screen.DrawImage(playerImage, op)
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
	img, _, err = ebitenutil.NewImageFromFile("assets/squid2.png")
	if err != nil {
		log.Fatal(err)
	}
	squidImage2 = ebiten.NewImageFromImage(img)
	img, _, err = ebitenutil.NewImageFromFile("assets/crab2.png")
	if err != nil {
		log.Fatal(err)
	}
	crabImage2 = ebiten.NewImageFromImage(img)
	img, _, err = ebitenutil.NewImageFromFile("assets/octopus2.png")
	if err != nil {
		log.Fatal(err)
	}
	octopusImage2 = ebiten.NewImageFromImage(img)
	img, _, err = ebitenutil.NewImageFromFile("assets/ufo.png")
	if err != nil {
		log.Fatal(err)
	}
	ufoImage = ebiten.NewImageFromImage(img)
	img, _, err = ebitenutil.NewImageFromFile("assets/death.png")
	if err != nil {
		log.Fatal(err)
	}
	deathImage = ebiten.NewImageFromImage(img)
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.mode == ModeStart {
		ebitenutil.DebugPrintAt(screen, "Press SPACE to Start", frameWidth/2-80, frameHeight/2)
		return
	}
	if g.mode == ModePause {
		ebitenutil.DebugPrintAt(screen, "Game Paused. Press P to Resume", frameWidth/2-120, frameHeight/2)
		return
	}
	if g.mode == ModeGameOver {
		ebitenutil.DebugPrintAt(screen, "Game Over! Press R to Restart", frameWidth/2-120, frameHeight/2)
		return
	}
	g.player.draw(screen)
	g.playerBullet.draw(screen)
	g.enemies.draw(screen)
	g.drawDeathAnims(screen)
	g.ufo.draw(screen)
	g.enemyBullets.draw(screen)
	g.scoreDisplay(screen)
	g.livesDisplay(screen)
}

func (g *Game) Update() error {
	if g.mode == ModeGameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			*g = Game{mode: ModePlaying}
		}
		return nil
	}
	if g.mode == ModeStart {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.mode = ModePlaying
		}
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if g.mode == ModePause {
			g.mode = ModePlaying
		} else {
			g.mode = ModePause
		}
	}
	if g.mode == ModePlaying {
		g.keyPressed()
		g.player.update()
		g.playerBullet.update(g)
		g.enemies.update(g)
		g.enemyBullets.update()
		g.updateDeathAnims()
		if g.enemyHit() {
			print("Enemy Hit!\n")
		}
		if g.playerHit() {
			g.lives--
			print("Player hit! Lives left: ", g.lives, "\n")
			if g.lives == 0 {
				g.mode = ModeGameOver
			}
		}
		if ufoHit(g.ufo, g) {
			print("UFO Hit!\n")
		}
	}
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return frameWidth, frameHeight
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Space Invaders")
	if err := ebiten.RunGame(&Game{mode: ModeStart, lives: 3, level: 0}); err != nil {
		log.Fatal(err)
	}
}
