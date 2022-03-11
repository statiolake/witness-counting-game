package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/nsf/termbox-go"
	"github.com/statiolake/witness-counting-game/aiplay"
	"github.com/statiolake/witness-counting-game/game"
	"github.com/statiolake/witness-counting-game/geom"
)

// inclusive
type rect struct {
	minX, minY, maxX, maxY int
}

func (r *rect) Width() int {
	return r.maxX - r.minX + 1
}

func (r *rect) Height() int {
	return r.maxY - r.minY + 1
}

func (r *rect) Clamp(x, y int) (int, int) {
	if x < r.minX {
		x = r.minX
	}

	if x > r.maxX {
		x = r.maxX
	}

	if y < r.minY {
		y = r.minY
	}

	if y > r.maxY {
		y = r.maxY
	}

	return x, y
}

type visualizer struct {
	snapshots   []game.Game
	snapshotIdx int

	fieldBorder rect
	fieldArea   rect
	infoArea    rect
}

func (v *visualizer) currentSnapshot() *game.Game {
	return &v.snapshots[v.snapshotIdx]
}

func (v *visualizer) stepSnapshot(step int) {
	v.snapshotIdx += step

	if v.snapshotIdx >= len(v.snapshots) {
		v.snapshotIdx = len(v.snapshots) - 1
	}

	if v.snapshotIdx < 0 {
		v.snapshotIdx = 0
	}
}

// ターミナルの現在のサイズに合わせて、画面の各コンポーネントをレイアウトする。
// ターミナルが小さすぎる場合は false を返す。
func (v *visualizer) recalcDrawAreas(width, height int) bool {
	// FIXME: この数字は適当に選んじゃってる
	if width <= 30 || height <= 15 {
		// 小さすぎる
		return false
	}

	// 座標は縦は 1 マス、横は 2 マスを占めるように作る。
	// フィールドはとりあえず正方形と考える。
	// TODO: 正方形ではなく実際のフィールドの縦横比に合わせる

	// infoArea は右側に置く。
	// ターミナルが十分横長であれば、縦に合わせたフィールド正方形の残りの部分
	// を割り当ててもよいが、もともとターミナルが正方形に近い場合は、無理やり
	// フィールドのサイズを縮小してでも一定の幅は確保する。
	// TODO: 端末が縦長なら下に配置してもよいのでは
	infoAreaWidth := width - (height * 2)
	if infoAreaWidth < 10 {
		// infoArea として小さすぎるので height を縮めて infoArea を強引にもっ
		// ていく
		infoAreaWidth = 10
		height = (width - infoAreaWidth) / 2
	}

	v.fieldBorder = rect{1, 0, width - infoAreaWidth - 2, height - 1}
	v.fieldArea = rect{
		v.fieldBorder.minX + 1,
		v.fieldBorder.minY + 1,
		v.fieldBorder.maxX - 1,
		v.fieldBorder.maxY - 1,
	}

	// 幅が奇数のときは偶数にする
	if v.fieldArea.Width()%2 != 0 {
		v.fieldBorder.maxX--
		v.fieldArea.maxX--
	}

	v.infoArea = rect{v.fieldBorder.maxX + 2, 0, width - 1, height - 1}

	return true
}

// ch1 は Squad ID, ch2 はエージェントの種類 とするべし
func (v *visualizer) drawAtFieldArea(x, y int, ch1, ch2 rune) {
	termbox.SetCell(
		x*2+v.fieldArea.minX, y+v.fieldArea.minY, ch1,
		termbox.ColorDefault, termbox.ColorDefault,
	)
	termbox.SetCell(
		x*2+1+v.fieldArea.minX, y+v.fieldArea.minY, ch2,
		termbox.ColorDefault, termbox.ColorDefault,
	)
}

func (v *visualizer) drawAtInfoArea(y int, msg string) {
	for idx, ch := range msg {
		termbox.SetCell(
			idx+v.infoArea.minX, y+v.infoArea.minY, ch,
			termbox.ColorDefault, termbox.ColorDefault,
		)
	}
}

func (v *visualizer) update(width, height int) {
	if !v.recalcDrawAreas(width, height) {
		// ターミナルが小さすぎるので表示は諦める
		termbox.SetCell(
			width/2, height/2, '!',
			termbox.ColorDefault, termbox.ColorDefault,
		)
		return
	}

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	v.drawFieldBorder()
	v.drawAgents()
	v.drawInfo()

	if err := termbox.Flush(); err != nil {
		panic(err)
	}
}

func (v *visualizer) drawFieldBorder() {
	border := v.fieldBorder
	for x := border.minX; x <= border.maxX; x++ {
		for y := border.minY; y <= border.maxY; y++ {
			ch := ' '
			if (x == border.minX && y == border.minY) ||
				(x == border.minX && y == border.maxY) ||
				(x == border.maxX && y == border.minY) ||
				(x == border.maxX && y == border.maxY) {
				ch = '+'
			} else if x == border.minX || x == border.maxX {
				ch = '|'
			} else if y == border.minY || y == border.maxY {
				ch = '-'
			}

			termbox.SetCell(x, y, ch,
				termbox.ColorDefault, termbox.ColorDefault)
		}
	}
}

func (v *visualizer) drawAgents() {
	snapshot := v.currentSnapshot()
	fieldWidth, fieldHeight := v.getFieldSize()
	minX := snapshot.Field.Rect.LT.X
	minY := snapshot.Field.Rect.LT.Y
	for idx := range snapshot.Agents {
		agent := &snapshot.Agents[idx]
		x := int(((agent.Pos.X - minX) / fieldWidth) * float64(v.fieldArea.Width()/2))
		y := int(((agent.Pos.Y - minY) / fieldHeight) * float64(v.fieldArea.Height()))
		// 数学的な座標と考え、画面座標とは Y 座標を反転する
		y = v.fieldArea.Height() - y
		x, y = v.fieldArea.Clamp(x, y)

		kind := '?'
		if agent.Kind == game.Hunter {
			kind = 'h'
		} else if agent.Kind == game.Runner {
			kind = 'r'
		}

		v.drawAtFieldArea(x, y, []rune(strconv.Itoa(agent.SquadID))[0], kind)
	}
}

func (v *visualizer) drawInfo() {
	v.drawAtInfoArea(
		0,
		fmt.Sprintf(
			"Snapshot #%d (in %d)",
			v.snapshotIdx+1,
			len(v.snapshots),
		),
	)

	// 各 Squad の総得点を表示する
	snapshot := v.currentSnapshot()
	for idx := range snapshot.Squads {
		squad := &snapshot.Squads[idx]
		v.drawAtInfoArea(
			2+idx,
			fmt.Sprintf(
				"- %s: %f (%+f)",
				squad.Name,
				squad.TotalPoint,
				squad.TotalPointGain,
			),
		)
	}
}

func (v *visualizer) getFieldSize() (width, height float64) {
	fieldRect := v.currentSnapshot().Field.Rect
	width = fieldRect.RB.X - fieldRect.LT.X
	height = fieldRect.RB.Y - fieldRect.LT.Y
	return
}

type randomAI struct {
	speed float64
}

func (ai *randomAI) Init(config game.GameConfig) error {
	return nil
}

func (ai *randomAI) Think(
	knowledge game.Knowledge,
	agent game.Agent,
) (*game.ActionMove, error) {
	angle := rand.Float64() * 2 * math.Pi
	return &game.ActionMove{Dir: geom.NewPolarVector(ai.speed, angle)}, nil
}

func createAIPlay() aiplay.AIPlay {
	numSquads := 5
	config := aiplay.DefaultAIPlayConfig()

	for i := 1; i <= numSquads; i++ {
		name := fmt.Sprintf("squad-%02d", i)
		agentBase := fmt.Sprintf("agent-%02d", i)
		config.
			WithSquadAdded(
				aiplay.NewSquadConfig(name).
					WithAgentAdded(
						game.NewAgentConfig(agentBase+"h", game.Hunter),
						&randomAI{speed: 1.0},
					).
					WithAgentAdded(
						game.NewAgentConfig(agentBase+"r", game.Runner),
						&randomAI{speed: 1.0},
					),
			)
	}

	config.GameConfig.Field.
		WithObstructionAdded(
			game.ObstructionConfig{
				Segment: geom.NewSegment(
					geom.NewCoord(0, 2),
					geom.NewCoord(0, -2),
				),
			},
		)

	return config.BuildAIPlay()
}

func main() {
	aiplay := createAIPlay()
	snapshots, err := aiplay.StepAll()
	if err != nil {
		panic(err)
	}

	v := visualizer{
		snapshots: snapshots,
	}

	if err := termbox.Init(); err != nil {
		panic(err)
	}

	defer termbox.Close()

MAINLOOP:
	for {
		width, height := termbox.Size()
		v.update(width, height)

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			case 'q':
				break MAINLOOP
			case 'j':
				v.stepSnapshot(-1)
			case 'l':
				v.stepSnapshot(+1)
			case rune(0):
				switch ev.Key {
				case termbox.KeyEsc:
					break MAINLOOP
				}
			}
		}
	}
}
