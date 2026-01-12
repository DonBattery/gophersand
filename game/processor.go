package game

func getCold(ca *CellAutomata, x, y int) uint8 {
	mat := ca.GetMaterialAt(x, y)
	if mat.IsKind(MaterialKindIce) {
		return 63
	}
	if mat.GetStatus() == MaterialStatusFrozen {
		return 17
	}
	return 0
}

// Returns a value between 0 and 255 representing the number of cold neighbors
func getTemp(ca *CellAutomata, x, y int) uint8 {
	var count uint8 = 255
	count -= getCold(ca, x-1, y)
	count -= getCold(ca, x+1, y)
	count -= getCold(ca, x, y-1)
	count -= getCold(ca, x, y+1)
	return count
}

func ProcessSand(ca *CellAutomata, kind MaterialKind, mat Material, cid, x, y int) bool {
	if mat.GetStatus() == MaterialStatusFrozen {
		if ca.tp.Turn5shift1 && ca.rngChance256(getTemp(ca, x, y)) {
			ca.SetCellAsProcessed(cid, mat.WithStatus(MaterialStatusNormal))
		}
		return true
	}

	canReact := false

	// Choose a random horizontal direction to reduce bias
	dir := 1
	if ca.rngBool() {
		dir = -1
	}

	// Check downwards first, except in a few cases check diagonal first
	checkX := x
	off := 0
	if ca.rngChance256(30) {
		off = dir
		checkX += off
	}

	// Make the first movement check
	checkY := y + 1
	canReactAt, reacted := ca.TryReactionAt(cid, mat, kind, checkX, checkY)
	if canReactAt {
		canReact = true
		if reacted {
			return true
		}
	}

	// Second movement check
	if off != 0 {
		checkX -= dir
	} else {
		checkX += dir
	}
	canReactAt, reacted = ca.TryReactionAt(cid, mat, kind, checkX, checkY)
	if canReactAt {
		canReact = true
		if reacted {
			return true
		}
	}

	// Activity check, no movement this time
	if !canReact {
		checkX = x - dir
		if ca.CanReactAt(kind, checkX, checkY) {
			return true
		}
	}

	return canReact
}

func ProcessWater(ca *CellAutomata, kind MaterialKind, mat Material, cid, x, y int) bool {
	if mat.GetStatus() == MaterialStatusFrozen {
		if ca.tp.Turn5shift1 && ca.rngChance256(getTemp(ca, x, y)) {
			ca.SetCellAsProcessed(cid, mat.WithStatus(MaterialStatusNormal))
		}
		return true
	}

	canReact := false

	// check desired flow direction, with a bit of randomness
	left := mat.GetFaceLeft()
	oleft := left
	if ca.rngChance256(50) {
		left = !left
	}

	dir := 1
	if left {
		dir = -1
	}

	// base positions
	downX, downY := x, y+1
	diagX, diagY := x+dir, y+1
	horzX, horzY := x+dir, y

	// decide check order based on priority:
	// 0 (10/256):  horizontal -> diagonal -> down
	// 1 (50/256):  diagonal -> down -> horizontal
	// 2 (196/256): down -> diagonal -> horizontal
	var x1, y1, x2, y2, x3, y3 int
	switch ca.rngPick3(10, 60) {
	case 0:
		x1, y1, x2, y2, x3, y3 = horzX, horzY, diagX, diagY, downX, downY
	case 1:
		x1, y1, x2, y2, x3, y3 = diagX, diagY, downX, downY, horzX, horzY
	default:
		x1, y1, x2, y2, x3, y3 = downX, downY, diagX, diagY, horzX, horzY
	}
	// x1, y1, x2, y2, x3, y3 = downX, downY, diagX, diagY, horzX, horzY

	// check in order
	if canReactAt, reacted := ca.TryReactionAt(cid, mat, kind, x1, y1); canReactAt {
		canReact = true
		if reacted {
			return true
		}
	}
	if canReactAt, reacted := ca.TryReactionAt(cid, mat, kind, x2, y2); canReactAt {
		canReact = true
		if reacted {
			return true
		}
	}
	if canReactAt, reacted := ca.TryReactionAt(cid, mat, kind, x3, y3); canReactAt {
		canReact = true
		if reacted {
			return true
		}
	}

	// if cannot move in the flow direction, check the other direction
	if !canReact {
		if ca.CanReactAt(kind, x-dir, y) || ca.CanReactAt(kind, x-dir, y+1) {
			// turn around
			ca.materials[cid] = mat.WithFaceLeft(!oleft)
			canReact = true
		}
	}

	// If water cannot move and there is an empty cell above it, there is a slight chance it turn to Steam
	if !canReact && ca.tp.Turn3 {
		if y > 0 && ca.materials[cid-WorldWidth].IsKind(MaterialKindEmpty) {
			if ca.rngChance256(1) {
				ca.SetCellAsProcessed(cid, MaterialSteam.WithFaceLeft(mat.GetFaceLeft()))
				return true
			}
		}
	}

	return canReact
}

func ProcessSeed(ca *CellAutomata, kind MaterialKind, mat Material, cid, x, y int) bool {
	if mat.GetStatus() == MaterialStatusFrozen {
		if ca.tp.Turn5shift1 && ca.rngChance256(getTemp(ca, x, y)) {
			ca.SetCellAsProcessed(cid, mat.WithStatus(MaterialStatusNormal))
		}
		return true
	}

	canReact := false

	// Use desired flow direction (like Water): FaceLeft decides the preferred horizontal direction.
	left := mat.GetFaceLeft()
	oleft := left
	dir := 1
	if left {
		dir = -1
	}

	// Base positions
	downX, downY := x, y+1
	diagX, diagY := x+dir, y+1
	horzX, horzY := x+dir, y

	// Decide check order based on priority:
	// ~10%: horizontal -> diagonal down -> down
	// ~30%: diagonal down -> down
	// ~60%: down -> diagonal down
	//
	// Note: percentages are approximate due to 1/256 RNG granularity.
	type pos struct{ x, y int }
	var checks [3]pos
	n := 0
	switch ca.rngPick3(26, 103) {
	case 0:
		checks[0], checks[1], checks[2] = pos{horzX, horzY}, pos{diagX, diagY}, pos{downX, downY}
		n = 3
	case 1:
		checks[0], checks[1] = pos{diagX, diagY}, pos{downX, downY}
		n = 2
	default:
		checks[0], checks[1] = pos{downX, downY}, pos{diagX, diagY}
		n = 2
	}

	// Try movement checks in order
	for i := 0; i < n; i++ {
		canReactAt, reacted := ca.TryReactionAt(cid, mat, kind, checks[i].x, checks[i].y)
		if canReactAt {
			canReact = true
			if reacted {
				return true
			}
		}
	}

	// If cannot move in the flow direction, check the other direction and turn around (like Water).
	if !canReact {
		if ca.CanReactAt(kind, x-dir, y) || ca.CanReactAt(kind, x-dir, y+1) {
			ca.materials[cid] = mat.WithFaceLeft(!oleft)
			return true
		}
	}

	// Activity check, no movement this time
	if !canReact {
		checkX := x - dir
		checkY := y + 1
		if ca.CanReactAt(kind, checkX, checkY) {
			canReact = true
		}
	}

	// Every 3th tick check if the seed can be activated (it is touching sand and water simultaneously)
	if ca.tp.Turn3 {
		touchWater := false
		touchSand := false

		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				xx := x + dx
				yy := y + dy
				if ca.InBounds(xx, yy) {
					checkMat := ca.materials[yy*WorldWidth+xx]
					if checkMat.IsKind(MaterialKindWater) {
						touchWater = true
					} else if checkMat.IsKind(MaterialKindSand) {
						touchSand = true
					}
				}
				if touchWater && touchSand {
					goto EarlyExit
				}
			}
		}

	EarlyExit:

		if touchWater && touchSand {
			ca.SetCellAsProcessed(cid, MaterialRoot.WithLife(ca.rngPick4(120, 180, 200)))
			return true
		}
	}

	return canReact
}

func ProcessAntHill(ca *CellAutomata, kind MaterialKind, mat Material, cid, x, y int) bool {
	_, reacted := ca.TryReactionAt(cid, mat, kind, x, y+1)
	return reacted
}

func eggTrySwapDownLikeSand(ca *CellAutomata, cid, x, y int, swapWithWaterChance, swapWithSandChance uint8) bool {
	// Choose a random horizontal direction to reduce bias
	dir := 1
	if ca.rngBool() {
		dir = -1
	}

	// Check downwards first, except in a few cases check diagonal first
	checkX := x
	off := 0
	if ca.rngChance256(30) {
		off = dir
		checkX += off
	}

	checkY := y + 1
	if y >= WorldHeight-1 {
		return true
	}

	trySwap := func(tx int) bool {
		if !ca.InBounds(tx, checkY) {
			return false
		}
		targetCid := checkY*WorldWidth + tx
		if ca.processed[targetCid] == ca.tick {
			return false
		}
		tk := ca.materials[targetCid].GetKind()
		switch tk {
		case MaterialKindEmpty, MaterialKindSteam, MaterialKindSmoke:
			ca.SwapCells(cid, targetCid)
			return true
		case MaterialKindWater:
			if ca.rngChance256(swapWithWaterChance) {
				ca.SwapCells(cid, targetCid)
				return true
			}
		case MaterialKindSand:
			if ca.rngChance256(swapWithSandChance) {
				ca.SwapCells(cid, targetCid)
				return true
			}
		}
		return false
	}

	// First movement check
	if trySwap(checkX) {
		return true
	}

	// Second movement check
	if off != 0 {
		checkX -= dir
	} else {
		checkX += dir
	}
	_ = trySwap(checkX)
	return true
}

func ProcessAcid(ca *CellAutomata, kind MaterialKind, mat Material, cid, x, y int) bool {
	canReact := false

	// check desired flow direction, with a bit of randomness
	left := mat.GetFaceLeft()
	oleft := left
	if ca.rngChance256(20) {
		left = !left
	}

	dir := 1
	if left {
		dir = -1
	}

	// base positions
	downX, downY := x, y+1
	diagX, diagY := x+dir, y+1
	horzX, horzY := x+dir, y

	// decide check order based on priority:
	// 0 (10/256):  horizontal -> diagonal -> down
	// 1 (50/256):  diagonal -> down -> horizontal
	// 2 (196/256): down -> diagonal -> horizontal
	var x1, y1, x2, y2, x3, y3 int
	switch ca.rngPick3(5, 30) {
	case 0:
		x1, y1, x2, y2, x3, y3 = horzX, horzY, diagX, diagY, downX, downY
	case 1:
		x1, y1, x2, y2, x3, y3 = diagX, diagY, downX, downY, horzX, horzY
	default:
		x1, y1, x2, y2, x3, y3 = downX, downY, diagX, diagY, horzX, horzY
	}

	// check in order
	if canReactAt, reacted := ca.TryReactionAt(cid, mat, kind, x1, y1); canReactAt {
		canReact = true
		if reacted {
			return true
		}
	}
	if canReactAt, reacted := ca.TryReactionAt(cid, mat, kind, x2, y2); canReactAt {
		canReact = true
		if reacted {
			return true
		}
	}
	if canReactAt, reacted := ca.TryReactionAt(cid, mat, kind, x3, y3); canReactAt {
		canReact = true
		if reacted {
			return true
		}
	}

	// if cannot move in the flow direction, check the other direction
	if !canReact {
		if ca.CanReactAt(kind, x-dir, y) || ca.CanReactAt(kind, x-dir, y+1) {
			// turn around
			ca.materials[cid] = mat.WithFaceLeft(!oleft)
			return true
		}
	}

	// If acid cannot move and there is an empty cell above it, there is a chance it evaporates into Smoke
	// Acid evaporates faster than water (every 2nd tick with 2/256 chance vs water's every 3rd tick with 1/256 chance)
	if !canReact && ca.tp.Turn5 {
		if y > 0 && ca.materials[cid-WorldWidth].IsKind(MaterialKindEmpty) {
			if ca.rngChance256(5) {
				ca.CreateSmoke(cid, 2)
				return true
			}
		}
	}

	return canReact
}

func ProcessFire(ca *CellAutomata, kind MaterialKind, mat Material, cid, x, y int) bool {
	// Fire decays over time (faster decay)
	life := mat.GetLife()
	if ca.rngChance256(20) {
		if life == 0 {
			// Turn into Smoke when life is depleted
			ca.CreateSmoke(cid, 3)
			return true
		}
		ca.SetCellAsProcessed(cid, mat.WithLife(life-1))
		return true
	}

	canReact := false

	// check desired flow direction, with a bit of randomness
	left := mat.GetFaceLeft()
	oleft := left
	if ca.rngChance256(100) {
		left = !left
	}

	dir := 1
	if left {
		dir = -1
	}

	// base positions (Fire rises like Smoke/Steam)
	upX, upY := x, y-1
	diagX, diagY := x+dir, y-1
	horzX, horzY := x+dir, y

	// decide check order based on priority:
	// 0 (15/256):  horizontal -> diagonal -> up
	// 1 (70/256):  diagonal -> up -> horizontal
	// 2 (171/256): up -> diagonal -> horizontal
	var x1, y1, x2, y2, x3, y3 int
	switch ca.rngPick3(30, 128) {
	case 0:
		x1, y1, x2, y2, x3, y3 = horzX, horzY, diagX, diagY, upX, upY
	case 1:
		x1, y1, x2, y2, x3, y3 = diagX, diagY, upX, upY, horzX, horzY
	default:
		x1, y1, x2, y2, x3, y3 = upX, upY, diagX, diagY, horzX, horzY
	}

	// check in order
	if canReactAt, reacted := ca.TryReactionAt(cid, mat, kind, x1, y1); canReactAt {
		canReact = true
		if reacted {
			return true
		}
	}
	if canReactAt, reacted := ca.TryReactionAt(cid, mat, kind, x2, y2); canReactAt {
		canReact = true
		if reacted {
			return true
		}
	}
	if canReactAt, reacted := ca.TryReactionAt(cid, mat, kind, x3, y3); canReactAt {
		canReact = true
		if reacted {
			return true
		}
	}

	// if cannot move in the flow direction, check the other direction
	if !canReact {
		if ca.CanReactAt(kind, x-dir, y) || ca.CanReactAt(kind, x-dir, y-1) {
			// turn around
			ca.materials[cid] = mat.WithFaceLeft(!oleft)
			return true
		}
	}

	return canReact
}

func ProcessIce(ca *CellAutomata, kind MaterialKind, mat Material, cid, x, y int) (activity bool) {
	// Ice always reports activity
	activity = true

	// Melting/freezing is only processed every 5 ticks (Turn5), but movement is processed every tick.
	if ca.rngChance256(32) && ca.tp.Turn5 {
		life := mat.GetLife()
		if ca.rngChance256(getTemp(ca, x, y)) {
			if life == 0 {
				// Ice melts into Water
				ca.SetCellAsProcessed(cid, MaterialWater.WithFaceLeft(ca.rngBool()))
				return
			}
			ca.SetCellAsProcessed(cid, mat.WithLife(life-1))
			return
		}
	}

	// Choose a random horizontal direction to reduce bias
	dir := 1
	if ca.rngBool() {
		dir = -1
	}

	// Ice movement: attempt only ONE direction per tick; diagonal is rare.
	checkX := x
	if ca.rngChance256(10) { // low chance for diagonal first
		checkX += dir
	}

	// Make a single movement check
	checkY := y + 1
	ca.TryReactionAt(cid, mat, kind, checkX, checkY)

	return
}

func ProcessSmoke(ca *CellAutomata, kind MaterialKind, mat Material, cid, x, y int) bool {
	// smoke decays over time
	life := mat.GetLife()
	if ca.rngChance256(3) {
		if life == 0 {
			ca.SetCellAsProcessed(cid, MaterialEmpty)
			return true
		}
		ca.materials[cid] = mat.WithLife(life - 1)
	}

	canReact := false

	// check desired flow direction, with a bit of randomness
	left := mat.GetFaceLeft()
	oleft := left
	if ca.rngChance256(100) {
		left = !left
	}

	dir := 1
	if left {
		dir = -1
	}

	// base positions
	upX, upY := x, y-1
	diagX, diagY := x+dir, y-1
	horzX, horzY := x+dir, y

	// decide check order based on priority:
	// 0 (10/256):  horizontal -> diagonal -> down
	// 1 (50/256):  diagonal -> down -> horizontal
	// 2 (196/256): down -> diagonal -> horizontal
	var x1, y1, x2, y2, x3, y3 int
	switch ca.rngPick3(20, 80) {
	case 0:
		x1, y1, x2, y2, x3, y3 = horzX, horzY, diagX, diagY, upX, upY
	case 1:
		x1, y1, x2, y2, x3, y3 = diagX, diagY, upX, upY, horzX, horzY
	default:
		x1, y1, x2, y2, x3, y3 = upX, upY, diagX, diagY, horzX, horzY
	}

	// check in order
	if canReactAt, reacted := ca.TryReactionAt(cid, mat, kind, x1, y1); canReactAt {
		canReact = true
		if reacted {
			return true
		}
	}
	if canReactAt, reacted := ca.TryReactionAt(cid, mat, kind, x2, y2); canReactAt {
		canReact = true
		if reacted {
			return true
		}
	}
	if canReactAt, reacted := ca.TryReactionAt(cid, mat, kind, x3, y3); canReactAt {
		canReact = true
		if reacted {
			return true
		}
	}

	// if cannot move in the flow direction, check the other direction
	if !canReact {
		// if it is possible to move in the other direction, turn around
		if ca.CanReactAt(kind, x-dir, y) || ca.CanReactAt(kind, x-dir, y+1) {
			ca.materials[cid] = mat.WithFaceLeft(!oleft)
			return true
		}
	}

	return canReact
}

func ProcessSteam(ca *CellAutomata, kind MaterialKind, mat Material, cid, x, y int) bool {
	// check if we are on top of the world or the cell above is a condensable material
	if y == 0 || !ca.materials[cid-WorldWidth].IsIn(NonCondensableKinds) {
		if ca.rngChance256(5) {
			// condense into water or empty
			if ca.rngBool() {
				ca.SetCellAsProcessed(cid, MaterialWater.WithFaceLeft(ca.rngBool()))
			} else {
				ca.SetCellAsProcessed(cid, MaterialEmpty)
			}
			return true
		}
	}

	// check desired flow direction, with a bit of randomness
	left := mat.GetFaceLeft()
	oleft := left
	if ca.rngChance256(85) {
		left = !left
	}

	dir := 1
	if left {
		dir = -1
	}

	// base positions
	upX, upY := x, y-1
	diagX, diagY := x+dir, y-1
	horzX, horzY := x+dir, y

	// decide check order based on priority:
	// 0 (10/256):  horizontal -> diagonal -> down
	// 1 (50/256):  diagonal -> down -> horizontal
	// 2 (196/256): down -> diagonal -> horizontal
	var x1, y1, x2, y2, x3, y3 int
	switch ca.rngPick3(15, 65) {
	case 0:
		x1, y1, x2, y2, x3, y3 = horzX, horzY, diagX, diagY, upX, upY
	case 1:
		x1, y1, x2, y2, x3, y3 = diagX, diagY, upX, upY, horzX, horzY
	default:
		x1, y1, x2, y2, x3, y3 = upX, upY, diagX, diagY, horzX, horzY
	}

	canReact := false
	// check in order
	if canReactAt, reacted := ca.TryReactionAt(cid, mat, kind, x1, y1); canReactAt {
		if reacted {
			return true
		}
	}
	if canReactAt, reacted := ca.TryReactionAt(cid, mat, kind, x2, y2); canReactAt {
		if reacted {
			return true
		}
	}
	if canReactAt, reacted := ca.TryReactionAt(cid, mat, kind, x3, y3); canReactAt {
		if reacted {
			return true
		}
	}

	// if cannot move in the flow direction, check the other direction
	if !canReact {
		if ca.CanReactAt(kind, x-dir, y) || ca.CanReactAt(kind, x-dir, y+1) {
			// turn around
			ca.materials[cid] = mat.WithFaceLeft(!oleft)
			return true
		}
	}

	return true
}

func ProcessRoot(ca *CellAutomata, kind MaterialKind, mat Material, cid, x, y int) bool {
	if mat.GetStatus() == MaterialStatusFrozen {
		if ca.tp.Turn5shift1 && ca.rngChance256(getTemp(ca, x, y)) {
			ca.SetCellAsProcessed(cid, mat.WithStatus(MaterialStatusNormal))
		}
		return true
	}

	// IMPORTANT for the tile-based wake system:
	// Root only *attempts* growth every 3 ticks, but it must still report "potential activity"
	// on the other ticks; otherwise a quiet tile will go to sleep and Root will appear to stop
	// at 32x32 tile borders.
	if !ca.tp.Turn3shift1 {
		if ca.CanReactAt(kind, x, y-1) {
			return true
		}
		if ca.CanReactAt(kind, x, y+1) {
			return true
		}
		if ca.CanReactAt(kind, x-1, y) {
			return true
		}
		if ca.CanReactAt(kind, x+1, y) {
			return true
		}
		return false
	}

	// Growth logic using weighted random direction selection:
	// - Life 0-1: 30% down, 20% up, 25% left, 25% right (roots dig down)
	// - Life 2-3: 20% down, 30% up, 25% left, 25% right (roots grow up toward surface)
	// Direction: 0=up, 1=down, 2=left, 3=right
	life := mat.GetLife()

	var dir uint8
	if life <= 1 {
		// t0=51 (20% up), t1=128 (30% down), t2=192 (25% left), remaining=right (25%)
		dir = ca.rngPick4(51, 128, 192)
	} else {
		// t0=77 (30% up), t1=128 (20% down), t2=192 (25% left), remaining=right (25%)
		dir = ca.rngPick4(77, 128, 192)
	}

	// Calculate target position based on direction
	var checkX, checkY int
	switch dir {
	case 0: // up
		checkX, checkY = x, y-1
	case 1: // down
		checkX, checkY = x, y+1
	case 2: // left
		checkX, checkY = x-1, y
	default: // right
		checkX, checkY = x+1, y
	}

	// Try to react at the chosen direction; if out of bounds or non-interactable, just wait
	if canReact, reacted := ca.TryReactionAt(cid, mat, kind, checkX, checkY); reacted || canReact {
		return true
	}

	// If we couldn't react this tick, still report potential activity so the tile stays awake.
	if ca.CanReactAt(kind, x, y-1) {
		return true
	}
	if ca.CanReactAt(kind, x, y+1) {
		return true
	}
	if ca.CanReactAt(kind, x-1, y) {
		return true
	}
	if ca.CanReactAt(kind, x+1, y) {
		return true
	}

	return false
}

func ProcessPlant(ca *CellAutomata, kind MaterialKind, mat Material, cid, x, y int) bool {
	if mat.GetStatus() == MaterialStatusFrozen {
		if ca.tp.Turn5shift1 && ca.rngChance256(getTemp(ca, x, y)) {
			ca.SetCellAsProcessed(cid, mat.WithStatus(MaterialStatusNormal))
		}
		return true
	}

	// IMPORTANT for the tile-based wake system:
	// Plant only *attempts* actions every 3 ticks (shift2), but it must still report "potential activity"
	// on the other ticks; otherwise a quiet tile will go to sleep and Plant will appear to stop growing at tile borders.
	if !ca.tp.Turn3shift2 {
		if ca.CanReactAt(kind, x, y-1) {
			return true
		}
		if ca.CanReactAt(kind, x, y+1) {
			return true
		}
		if ca.CanReactAt(kind, x-1, y) {
			return true
		}
		if ca.CanReactAt(kind, x+1, y) {
			return true
		}
		return false
	}

	upCid := cid - WorldWidth
	downCid := cid + WorldWidth
	leftCid := cid - 1
	rightCid := cid + 1

	// Check if this Plant can bloom into a flower
	// Requirements: CanBloom=true, Life=3, pass random check, all 4 neighbors are Plants
	// Only ~5% of plants have CanBloom=true (set at creation, only if not on edge)
	if mat.GetCanBloom() && mat.GetLife() == 3 && x != 0 && x != WorldWidth-1 && y != 0 && y != WorldHeight-1 && ca.rngChance256(10) {

		upMat := ca.materials[upCid]
		downMat := ca.materials[downCid]
		leftMat := ca.materials[leftCid]
		rightMat := ca.materials[rightCid]

		// Check if all 4 neighbors are Plants
		if upMat.IsKind(MaterialKindPlant) &&
			downMat.IsKind(MaterialKindPlant) &&
			leftMat.IsKind(MaterialKindPlant) &&
			rightMat.IsKind(MaterialKindPlant) {

			// Turn this Plant into a Seed
			ca.SetCellAsProcessed(cid, MaterialSeed.WithLife(ca.rng0123()))

			// Choose a flower color (Life: 0=Blue, 1=Pink, 2=Magenta, 3=Yellow)
			color := ca.rngPick4(64, 128, 192)

			// Turn all 4 neighbors into Flowers (Status 0 = Normal, Life = color)
			// The top petal (up) is marked as IsTopPetal for seed dropping
			ca.SetCellAsProcessed(upCid, MaterialFlower.WithStatus(MaterialStatusNormal).WithLife(color).WithIsTopPetal(true))
			ca.SetCellAsProcessed(downCid, MaterialFlower.WithStatus(MaterialStatusNormal).WithLife(color))
			ca.SetCellAsProcessed(leftCid, MaterialFlower.WithStatus(MaterialStatusNormal).WithLife(color))
			ca.SetCellAsProcessed(rightCid, MaterialFlower.WithStatus(MaterialStatusNormal).WithLife(color))

			return true
		}
	}

	// If a Plant is not supported by any of its 4 neighbors, it turns into a Seed.
	// Plants support each other (MaterialKindPlant is included in PlantSupporterKinds), so this only triggers when
	// the Plant becomes isolated (e.g. surrounding plants were eaten).
	hasSupport := false
	if y > 0 && ca.materials[upCid].IsIn(PlantSupporterKinds) {
		hasSupport = true
	} else if y < WorldHeight-1 && ca.materials[downCid].IsIn(PlantSupporterKinds) {
		hasSupport = true
	} else if x > 0 && ca.materials[leftCid].IsIn(PlantSupporterKinds) {
		hasSupport = true
	} else if x < WorldWidth-1 && ca.materials[rightCid].IsIn(PlantSupporterKinds) {
		hasSupport = true
	}
	if !hasSupport {
		ca.SetCellAsProcessed(cid, MaterialSeed.WithLife(ca.rng0123()))
		return true
	}

	// Growth logic using weighted random direction selection:
	// 40% up, 20% down, 20% left, 20% right (plants grow upward)
	// Direction: 0=up, 1=down, 2=left, 3=right
	// t0=102 (40% up), t1=153 (20% down), t2=204 (20% left), remaining=right (20%)
	dir := ca.rngPick4(102, 153, 204)

	// Calculate target position based on direction
	var checkX, checkY int
	switch dir {
	case 0: // up
		checkX, checkY = x, y-1
	case 1: // down
		checkX, checkY = x, y+1
	case 2: // left
		checkX, checkY = x-1, y
	default: // right
		checkX, checkY = x+1, y
	}

	// Try to react at the chosen direction; if out of bounds or non-interactable, just wait
	if canReact, reacted := ca.TryReactionAt(cid, mat, kind, checkX, checkY); reacted || canReact {
		return true
	}

	// If we couldn't react this tick, still report potential activity so the tile stays awake.
	if ca.CanReactAt(kind, x, y-1) {
		return true
	}
	if ca.CanReactAt(kind, x, y+1) {
		return true
	}
	if ca.CanReactAt(kind, x-1, y) {
		return true
	}
	if ca.CanReactAt(kind, x+1, y) {
		return true
	}

	return false
}

func ProcessFlower(ca *CellAutomata, kind MaterialKind, mat Material, cid, x, y int) bool {
	if mat.GetStatus() == MaterialStatusFrozen {
		if ca.tp.Turn5shift1 && ca.rngChance256(getTemp(ca, x, y)) {
			ca.SetCellAsProcessed(cid, mat.WithStatus(MaterialStatusNormal))
		}
		return true
	}

	// Flowers are mostly static, but the top petal can drop seeds
	if ca.tp.Turn5 && mat.GetIsTopPetal() && y < WorldHeight-1 && ca.rngChance256(5) {

		belowCid := cid + WorldWidth
		if ca.materials[belowCid].IsKind(MaterialKindEmpty) {
			// Create a new seed below
			ca.SetCellAsProcessed(belowCid, MaterialSeed.WithLife(ca.rng0123()))
			return true
		}

	}
	return false
}

func ProcessAnt(ca *CellAutomata, kind MaterialKind, mat Material, cid, x, y int) (activity bool) {
	// Ant always reports activity as long as it lives
	activity = true

	if mat.GetStatus() == MaterialStatusFrozen {
		if ca.tp.Turn5shift1 && ca.rngChance256(getTemp(ca, x, y)) {
			ca.SetCellAsProcessed(cid, mat.WithStatus(MaterialStatusNormal))
		}
		return
	}

	// Egg state: Ant with Life==0 behaves as an AntEgg (falls like sand, can hatch).
	if mat.GetLife() == 0 {
		// Hatch chance: when hatching, it becomes an adult Ant with Life 1-2.
		if ca.tp.Turn5 && ca.rngChance256(2) {
			newLife := uint8(1)
			if ca.rngBool() {
				newLife = 2
			}
			ca.SetCellAsProcessed(cid, MaterialAnt.
				WithLife(newLife).
				WithFaceLeft(ca.rngBool()).
				WithFaceUp(ca.rngBool()))
			return
		}

		// Movement: like Sand (gravity + diagonal).
		return eggTrySwapDownLikeSand(ca, cid, x, y, 120, 4)
	}

	// Hunger: slight chance to lose 1 life naturally (lower than before).
	// If this would drop Life below 1, the Ant dies (turns into Sand) and never becomes an egg again.
	if ca.rngChance256(1) {
		life := mat.GetLife()
		if life <= 1 {
			ca.CreateSand(cid, true)
			return
		}
		if ca.HasNeighborKind(x, y, AntAliveKinds) {
			return
		}
		ca.SetCellAsProcessed(cid, mat.WithLife(life-1))
		return
	}

	faceLeft := mat.GetFaceLeft()
	faceUp := mat.GetFaceUp()

	// --- Gravity simulation (every tick) ---
	// If Ant has fallable material directly below it, is not on left/right world boundary,
	// and has no AntSupporterKinds on (W, SW, E, SE), it tries to react with the cell below.
	if x != 0 && x != WorldWidth-1 && y < WorldHeight-1 {
		belowKind := ca.GetMaterialAt(x, y+1).GetKind()
		if belowKind.IsIn(AntFallableKinds) {
			w := ca.GetMaterialAt(x-1, y)
			sw := ca.GetMaterialAt(x-1, y+1)
			se := ca.GetMaterialAt(x+1, y+1)
			e := ca.GetMaterialAt(x+1, y)
			hasSupport := w.IsIn(AntSupporterKinds) || sw.IsIn(AntSupporterKinds) || se.IsIn(AntSupporterKinds) || e.IsIn(AntSupporterKinds)
			if !hasSupport {
				ca.TryReactionAt(cid, mat, kind, x, y+1)
				return
			}
		}
	}

	// Every 3rd tick (but not the same tick as Root/Plant; they use shift1/shift2),
	// an Ant with Life 2 or 3 may try to lay an egg (spawns an Ant egg: Ant with Life==0),
	// then loses 1 Life.
	if ca.tp.Turn3 {
		life := mat.GetLife()
		// Lower chance to lay Eggs.
		if life >= 2 && ca.rngChance256(1) {
			// 0=up, 1=down, 2=left, 3=right
			d := ca.rngPick4(64, 128, 192)
			tx, ty := x, y
			switch d {
			case 0:
				ty = y - 1
			case 1:
				ty = y + 1
			case 2:
				tx = x - 1
			default:
				tx = x + 1
			}
			// check if the target is inside the world
			if ca.InBounds(tx, ty) {
				targetCid := ty*WorldWidth + tx
				// and it is an AntEggLayableKind
				if ca.materials[targetCid].GetKind().IsIn(AntEggLayableKinds) {
					// lay the egg
					ca.SetCellAsProcessed(targetCid, MaterialAnt.WithLife(0))
					// and lower the Ant's life :)
					ca.materials[cid] = mat.WithLife(life - 1)
				}
			}
		}

		// --- Free will movement (every 3rd tick) ---
		vx := 1
		if faceLeft {
			vx = -1
		}
		vy := 1
		if faceUp {
			vy = -1
		}

		// Try to move: pick one of (horizontal, diagonal, vertical) with ~33% chance each.
		// Example: down+left => W / SW / S.
		tx, ty := x, y
		switch ca.rng012() {
		case 0:
			tx = x + vx
		case 1:
			tx = x + vx
			ty = y + vy
		default:
			ty = y + vy
		}

		if canReact, _ := ca.TryReactionAt(cid, mat, kind, tx, ty); canReact {
			return
		}

		// Impossible terrain: flip either horizontal or vertical desire (50/50).
		if ca.rngBool() {
			faceLeft = !faceLeft
		} else {
			faceUp = !faceUp
		}
		ca.materials[cid] = mat.WithFaceLeft(faceLeft).WithFaceUp(faceUp)
	}

	return
}

// Wasp
func ProcessWasp(ca *CellAutomata, kind MaterialKind, mat Material, cid, x, y int) (activity bool) {
	// Wasp always reports activity
	activity = true

	if mat.GetStatus() == MaterialStatusFrozen {
		if ca.tp.Turn5shift1 && ca.rngChance256(getTemp(ca, x, y)) {
			ca.SetCellAsProcessed(cid, mat.WithStatus(MaterialStatusNormal))
		}
		return
	}

	// Egg state: Wasp with Life==0 behaves as a WaspEgg (sticky + falls like sand, can hatch).
	if mat.GetLife() == 0 {

		// Hatch chance: when hatching, it becomes an adult Wasp with Life 1-2.
		if ca.tp.Turn5 && ca.rngChance256(2) {
			newLife := uint8(1)
			if ca.rngBool() {
				newLife = 2
			}
			ca.SetCellAsProcessed(cid, MaterialWasp.
				WithLife(newLife).
				WithFaceLeft(ca.rngBool()).
				WithFaceUp(ca.rngBool()))
			return
		}

		// Sticky: don't fall if touching sticky materials horizontally (left/right) or hanging under them (cell above).
		// World edges are also sticky: left, right, top.
		if x == 0 || x == WorldWidth-1 || y == 0 {
			return
		}

		up := ca.GetMaterialAt(x, y-1)
		left := ca.GetMaterialAt(x-1, y)
		right := ca.GetMaterialAt(x+1, y)
		if up.IsIn(WaspEggStickyKinds) || left.IsIn(WaspEggStickyKinds) || right.IsIn(WaspEggStickyKinds) {
			return
		}

		// Movement: like Sand (gravity + diagonal).
		return eggTrySwapDownLikeSand(ca, cid, x, y, 120, 4)
	}

	// Wasps are computed every 2 ticks.
	if !ca.tp.Turn2 {
		return
	}

	faceLeft := mat.GetFaceLeft()
	faceUp := mat.GetFaceUp()

	// Egg laying: only after the Wasp has drunk 1 Water and defeated 1 Ant (or ate 1 Ant egg: Ant with Life==0).
	// Can only place into Empty/Steam/Smoke/Plant, and only if the target has a sticky neighbor (U/D/L/R).
	if mat.GetWaspHasWater() && mat.GetWaspHasAnt() && ca.rngChance256(20) {
		// pick a random cardinal direction and try once
		d := ca.rngPick4(64, 128, 192) // 0=up,1=down,2=left,3=right
		tx, ty := x, y
		switch d {
		case 0:
			ty = y - 1
		case 1:
			ty = y + 1
		case 2:
			tx = x - 1
		default:
			tx = x + 1
		}
		if ca.InBounds(tx, ty) {
			targetCid := ty*WorldWidth + tx
			if ca.materials[targetCid].GetKind().IsIn(WaspEggLayableKinds) {
				// Sticky neighbor check (edges: left/right/top count as sticky too; bottom does not).
				hasStickyNeighbor := false
				// left
				if tx == 0 {
					hasStickyNeighbor = true
				} else if ca.GetMaterialAt(tx-1, ty).IsIn(WaspEggStickyKinds) {
					hasStickyNeighbor = true
				}
				// right
				if !hasStickyNeighbor {
					if tx == WorldWidth-1 {
						hasStickyNeighbor = true
					} else if ca.GetMaterialAt(tx+1, ty).IsIn(WaspEggStickyKinds) {
						hasStickyNeighbor = true
					}
				}
				// up
				if !hasStickyNeighbor {
					if ty == 0 {
						hasStickyNeighbor = true
					} else if ca.GetMaterialAt(tx, ty-1).IsIn(WaspEggStickyKinds) {
						hasStickyNeighbor = true
					}
				}
				// down (bottom edge is NOT sticky)
				if !hasStickyNeighbor {
					if ty < WorldHeight-1 && ca.GetMaterialAt(tx, ty+1).IsIn(WaspEggStickyKinds) {
						hasStickyNeighbor = true
					}
				}

				if hasStickyNeighbor {
					// lay the egg
					ca.SetCellAsProcessed(targetCid, MaterialWasp)
					// reset hunger flags
					ca.materials[cid] = mat.WithWaspHasWater(false).WithWaspHasAnt(false)
					return
				}
			}
		}
	}

	// Wasps do not fall down. They only try to move in their directed 3-way fan:
	// (vertical, diagonal, horizontal) with ~33% each.
	vx := 1
	if faceLeft {
		vx = -1
	}
	vy := -1
	if !faceUp {
		vy = 1
	}

	tx, ty := x, y
	switch ca.rng012() {
	case 0: // vertical (N/S)
		ty = y + vy
	case 1: // diagonal (NE/NW/SE/SW)
		tx = x + vx
		ty = y + vy
	default: // horizontal (E/W)
		tx = x + vx
	}

	// Edge of world => flip a random switch and stop.
	if !ca.InBounds(tx, ty) {
		if ca.rngBool() {
			faceLeft = !faceLeft
		} else {
			faceUp = !faceUp
		}

		ca.materials[cid] = mat.WithFaceLeft(faceLeft).WithFaceUp(faceUp)
		return
	}

	// For Water/Flower, we want "failed eating" (already collected or RNG miss) to behave like a wall:
	// bounce/flip direction instead of stopping in place due to "canReact but not reacted".
	targetKind := ca.GetMaterialAt(tx, ty).GetKind()
	canReact, reacted := ca.TryReactionAt(cid, mat, kind, tx, ty)
	if reacted {
		return
	}
	if canReact {
		// Bounce off Water when not drunk (already has water or RNG miss).
		if targetKind == MaterialKindWater {
			if ca.rngBool() {
				faceLeft = !faceLeft
			} else {
				faceUp = !faceUp
			}
			ca.materials[cid] = mat.WithFaceLeft(faceLeft).WithFaceUp(faceUp)
			return
		}
		// Otherwise (e.g. Steam/Smoke), keep the old behavior: interacting but failing still counts as activity.
		return
	}

	// Impossible terrain => flip a random switch.
	if ca.rngBool() {
		faceLeft = !faceLeft
	} else {
		faceUp = !faceUp
	}

	ca.materials[cid] = mat.WithFaceLeft(faceLeft).WithFaceUp(faceUp)
	return true
}
