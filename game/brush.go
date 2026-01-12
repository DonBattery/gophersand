package game

import "math/rand"

// This file contains the default brush actions for each material kind.
// Brushes are managed by Game and applied onto a CellAutomata when the user paints.

// BrushAction is invoked for each painted cell coordinate.
// The CellAutomata pointer is provided for context (e.g. neighbor queries)
type BrushAction func(ca *CellAutomata, x, y int) Material

// BrushActions are the two brush passes associated with a material kind.
// FirstAction is required; SecondAction is optional (nil means no second pass).
type BrushActions struct {
	FirstAction  BrushAction
	SecondAction BrushAction
}

// --- Brush Actions (First/Second pass) ---

func brushEmpty(_ *CellAutomata, _, _ int) Material {
	return MaterialEmpty
}

func brushStonePass1(_ *CellAutomata, _, _ int) Material {
	return MaterialStone.WithIsPenetrable(rand.Intn(100) < 51)
}

// brushStonePass2 colors stones based on how many non-stone neighbors they have above and below
func brushStonePass2(ca *CellAutomata, x, y int) Material {
	current := ca.GetMaterialAt(x, y)
	if !current.IsKind(MaterialKindStone) {
		return current
	}

	nonStoneAbove := 0
	for i := 1; i <= 3; i++ {
		if !ca.GetMaterialAt(x, y-i).IsKind(MaterialKindStone) {
			nonStoneAbove++
		}
	}

	nonStoneBelow := 0
	for i := 1; i <= 3; i++ {
		if !ca.GetMaterialAt(x, y+i).IsKind(MaterialKindStone) {
			nonStoneBelow++
		}
	}

	var life uint8 = 1

	if rand.Intn(2) == 1 {
		life = 2
	}
	if rand.Intn(256) < nonStoneAbove*75+20 {
		life = 0
	}
	if rand.Intn(256) < nonStoneBelow*75+20 {
		life = 3
	}

	return current.WithLife(life)
}

func brushSand(_ *CellAutomata, _, _ int) Material {
	return MaterialSand.
		WithLife(uint8(rand.Intn(4))).
		WithIsPenetrable(rand.Intn(100) < 66)
}

func brushWater(_ *CellAutomata, _, _ int) Material {
	life := uint8(rand.Intn(4))
	return MaterialWater.
		WithLife(life).
		WithFaceLeft(rand.Intn(2) == 1)
}

func brushSeed(_ *CellAutomata, _, _ int) Material {
	return MaterialSeed.WithLife(uint8(rand.Intn(4)))
}

// Just a placeholder, the user cannot paint AntHill
func brushAntHill(_ *CellAutomata, _, _ int) Material {
	return MaterialAntHill
}

func brushAcid(_ *CellAutomata, _, _ int) Material {
	return MaterialAcid.
		WithFaceLeft(rand.Intn(2) == 1)
}

func brushFire(_ *CellAutomata, _, _ int) Material {
	return MaterialFire.
		WithLife(3).
		WithFaceLeft(rand.Intn(2) == 1).
		WithStatus(uint8(rand.Intn(4)))
}

func brushIce(_ *CellAutomata, _, _ int) Material {
	// 90% chance for life 3, 10% chance for life 2
	life := uint8(3)
	if rand.Intn(100) < 10 {
		life = 2
	}
	return MaterialIce.
		WithLife(life)
}

func brushWasp(_ *CellAutomata, _, _ int) Material {
	// Spawn adult Wasps by default (Life 1-3).
	return MaterialWasp.
		WithLife(uint8(rand.Intn(3) + 1)).
		WithFaceLeft(rand.Intn(2) == 0).
		WithFaceUp(rand.Intn(2) == 0)
}

func brushAnt(_ *CellAutomata, _, _ int) Material {
	return MaterialAnt.
		WithLife(0).
		WithFaceLeft(rand.Intn(2) == 0).
		WithFaceUp(rand.Intn(2) == 0)
}

// Just a placeholder, the user cannot paint Steam
func brushSteam(_ *CellAutomata, _, _ int) Material {
	return MaterialSteam
}

// Just a placeholder, the user cannot paint Smoke
func brushSmoke(_ *CellAutomata, _, _ int) Material {
	return MaterialSmoke
}

// Just a placeholder, the user cannot paint Flower
func brushFlower(_ *CellAutomata, _, _ int) Material {
	return MaterialFlower
}

// Just a placeholder, the user cannot paint Root
func brushRoot(ca *CellAutomata, x, y int) Material {
	return MaterialRoot
}

// Just a placeholder, the user cannot paint Plant
func brushPlant(ca *CellAutomata, x, y int) Material {
	return MaterialPlant
}
