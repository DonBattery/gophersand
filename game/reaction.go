package game

// ============================================================================
// Generic reaction helpers
// ============================================================================

var (
	AlwaysSwap = SwapReaction(255)
)

// SwapReaction creates a new MaterialReaction function which swaps material A with material B with a given chance.
// chance 255 always swaps, 0 never swaps.
func SwapReaction(chance uint8) MaterialReaction {
	return func(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
		if chance == 255 || ca.rngChance256(chance) {
			ca.SwapCells(cidA, cidB)
			return true
		}
		return false
	}
}

// FireBurnReaction creates a fire reaction where Fire burns another material.
// - chanceToTransform: chance (0-255) for the affected material to turn into Fire
// - chanceFireToSmoke: chance (0-255) for the Fire to turn into Smoke
func FireBurnReaction(chanceToTransform, chanceFireToSmoke uint8) MaterialReaction {
	return func(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
		// Set the affected material to Burned
		ca.SetCell(cidB, matB.WithStatus(MaterialStatusBurned))

		// Check if Fire turns the target into Fire
		if ca.rngChance256(chanceToTransform) {
			ca.SetCellAsProcessed(cidB, MaterialFire.WithLife(3).WithFaceLeft(ca.rngBool()).WithStatus(uint8(ca.rngPick4(64, 128, 192))))
			// Also check if Fire turns into Smoke
			if ca.rngChance256(chanceFireToSmoke) {
				ca.CreateSmoke(cidB, 3)
			}
			return true
		}

		// Check if Fire turns into Smoke (without transforming target)
		if ca.rngChance256(chanceFireToSmoke) {
			ca.CreateSmoke(cidB, 3)
			return true
		}

		return false
	}
}

func RootGrowthReaction(growthChance uint8) MaterialReaction {
	return func(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
		life := matA.GetLife()
		if (matB.IsKind(MaterialKindAntHill) && !matB.GetIsPenetrable()) || life < 3 || !ca.rngChance256(growthChance) {
			return false
		}

		ca.SetCellAsProcessed(cidA, matA.WithLife(life-1))

		if ca.rngChance256(32) {
			ca.CreatePlant(cidB)
		} else {
			ca.CreateRoot(cidB)
		}

		return true
	}
}

func PlantGrowthReaction(growthChance uint8) MaterialReaction {
	return func(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
		lifeA := matA.GetLife()
		if lifeA == 0 || !ca.rngChance256(growthChance) {
			return false
		}

		ca.SetCellAsProcessed(cidA, matA.WithLife(lifeA-1))

		ca.CreatePlant(cidB)

		return true
	}
}

// ============================================================================
// Sand reactions (MaterialKind = 2)
// ============================================================================

func ReactionSandToAcid(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	switch rnd := ca.rngPick3(5, 150); rnd {
	case 0:
		ca.CreateSmoke(cidA, 1)
		if ca.rngChance256(128) {
			ca.CreateSmoke(cidB, 1)
		}
		return true
	case 1:
		ca.SwapCells(cidA, cidB)
		ca.SetCell(cidB, matA.WithStatus(MaterialStatusAcidic))
		return true
	default:
		ca.SetCell(cidA, matA.WithStatus(MaterialStatusAcidic))
		return true
	}
}

func ReactionSandToFire(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	switch rnd := ca.rngPick3(30, 220); rnd {
	case 0:
		ca.CreateSmoke(cidB, 1)
		return true
	case 1:
		ca.SwapCells(cidA, cidB)
		if ca.rngChance256(30) {
			ca.SetCell(cidB, matA.WithStatus(MaterialStatusBurned))
		}
		return true
	default:
		if ca.rngChance256(10) {
			ca.SetCell(cidA, matA.WithStatus(MaterialStatusBurned))
		}
		return false
	}
}

func ReactionSandToIce(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	if ca.rngChance256(80) {
		ca.SetCellAsProcessed(cidA, matA.WithStatus(MaterialStatusFrozen))
		return true
	}
	return false
}

// ============================================================================
// Water reactions (MaterialKind = 3)
// ============================================================================

func ReactionWaterToAcid(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	switch rnd := ca.rngPick3(15, 30); rnd {
	// Small chance to turn the water into steam
	case 0:
		// in this case there is also a chance to turn the Acid into Smoke
		if ca.rngBool() {
			ca.CreateSmoke(cidB, 1)
		}
		ca.SetCellAsProcessed(cidA, MaterialSteam.WithLife(ca.rngPick4(10, 20, 30)).WithFaceLeft(ca.rngBool()))
		return true

	// Small chance to swap
	case 1:
		ca.SwapCells(cidA, cidB)
		return true

	default:
		return false
	}
}

func ReactionWaterToFire(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Water has high chance to turn Fire into Smoke
	if ca.rngChance256(200) {
		ca.CreateSmoke(cidB, 1)
		// Also chance to turn Water into Steam
		if ca.rngChance256(150) {
			ca.SetCellAsProcessed(cidA, MaterialSteam.WithLife(uint8(ca.rngPick4(10, 20, 30))).WithFaceLeft(ca.rngBool()))
		}
		return true
	}
	return false
}

// ============================================================================
// Seed reactions (MaterialKind = 4)
// ============================================================================

func ReactionSeedToAcid(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Always turn the Seed's status into Acidic
	// ca.SetCell(cidA, matA.WithStatus(MaterialStatusAcidic))

	switch rnd := ca.rngPick3(1, 30); rnd {
	// The acid turns the Seed into Smoke or Steam
	case 0:
		// In this case there is also a slight chance to turn the Acid into Smoke
		if ca.rngChance256(60) {
			ca.CreateSmoke(cidB, 1)
		}
		// The Seed turns into Smoke or Steam
		if ca.rngBool() {
			ca.CreateSmoke(cidA, 1)
		} else {
			ca.SetCellAsProcessed(cidA, MaterialSteam.WithLife(uint8(ca.rngPick4(10, 20, 30))).WithFaceLeft(ca.rngBool()))
		}
		return true
	case 1:
		// The Seed turns into Acidic status and swaps with the Acid (sinks)
		ca.SetCell(cidA, matA.WithStatus(MaterialStatusAcidic))
		ca.SwapCells(cidA, cidB)
		return true
	default:
		// The Seed turns into Acidic status
		ca.SetCell(cidA, matA.WithStatus(MaterialStatusAcidic))
		return false
	}
}

func ReactionSeedToIce(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	if ca.rngChance256(100) {
		ca.SetCellAsProcessed(cidA, matA.WithStatus(MaterialStatusFrozen))
		return true
	}
	return false
}

// ============================================================================
// Acid reactions (MaterialKind = 6)
// ============================================================================

func ReactionAcidToSand(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	if ca.rngChance256(10) {
		ca.CreateSmoke(cidB, 1)
		if ca.rngChance256(128) {
			ca.CreateSmoke(cidB, 1)
		}
		return true
	}

	ca.SetCell(cidB, matB.WithStatus(MaterialStatusAcidic))
	return true
}

func ReactionAcidToStone(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	if matB.GetIsPenetrable() && ca.rngChance256(10) {
		ca.CreateSmoke(cidA, 1)
		ca.CreateSmoke(cidB, 1)
		return true
	}

	ca.SetCell(cidB, matB.WithStatus(MaterialStatusAcidic))
	return true
}

func ReactionAcidToWater(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	switch rnd := ca.rngPick3(15, 220); rnd {
	// Small chance to turn water into steam
	case 0:
		// in this case there is also a chance to turn the acid into smoke
		if ca.rngBool() {
			ca.CreateSmoke(cidB, 1)
		}
		ca.SetCellAsProcessed(cidB, MaterialSteam.WithLife(uint8(ca.rngPick4(10, 20, 30))).WithFaceLeft(ca.rngBool()))
		return true

	// Big chance to swap
	case 1:
		ca.SwapCells(cidA, cidB)
		return true

	default:
		return false
	}
}

func ReactionAcidToSeed(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Always turn the Acid's status into Normal
	// ca.SetCell(cidB, matB.WithStatus(MaterialStatusNormal))

	switch rnd := ca.rngPick3(1, 10); rnd {
	// The Acid turns the Seed into Smoke or Steam
	case 0:
		// In this case there is also a slight chance to turn the Acid into Smoke
		if ca.rngChance256(60) {
			ca.CreateSmoke(cidA, 1)
		}
		// Turn the Seed into Smoke or Steam
		if ca.rngBool() {
			ca.CreateSmoke(cidB, 1)
		} else {
			ca.SetCellAsProcessed(cidB, MaterialSteam.WithLife(uint8(ca.rngPick4(10, 20, 30))).WithFaceLeft(ca.rngBool()))
		}
		return true
	case 1:
		// There is a slight chance for the Acid to swap with the Seed (the Seed floats), in this case the Seed's stats is set to Acidic
		ca.SetCell(cidB, matB.WithStatus(MaterialStatusAcidic))
		ca.SwapCells(cidA, cidB)
		return true
	default:
		// The Seed's status is set to Acidic
		ca.SetCell(cidB, matB.WithStatus(MaterialStatusAcidic))
		return false
	}
}

func ReactionAcidToFire(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// 10% chance to turn Fire into Smoke
	if ca.rngChance256(26) {
		ca.CreateSmoke(cidB, 2)
		return true
	}
	// Otherwise they can swap
	ca.SwapCells(cidA, cidB)
	return true
}

func ReactionAcidToRoot(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Acid has a slight chance to deal damage to the Root
	// If no Life is left, the Root turns into Smoke, in this case there is also a chance for the Acid to turn into Smoke
	if ca.rngChance256(10) {
		life := matB.GetLife()
		if life > 0 {
			ca.SetCellAsProcessed(cidB, matB.WithLife(life-1).WithStatus(MaterialStatusAcidic))
		} else {
			ca.CreateSmoke(cidB, 1)
			if ca.rngBool() {
				ca.CreateSmoke(cidA, 1)
			}
		}
		return true
	}

	ca.SetCell(cidB, matB.WithStatus(MaterialStatusAcidic))

	return false
}

func ReactionAcidToPlant(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Acid has a slight chance to deal damage to the Plant.
	// If no Life is left, the Plant turns into Steam or Smoke (50/50), in this case there is also a chance for the Acid to turn into Smoke
	if ca.rngChance256(10) {
		life := matB.GetLife()
		if life > 0 {
			ca.SetCellAsProcessed(cidB, matB.WithLife(life-1).WithStatus(MaterialStatusAcidic))
		} else {
			if ca.rngBool() {
				ca.CreateSmoke(cidB, 1)
			} else {
				ca.SetCellAsProcessed(cidB, MaterialSteam.WithLife(uint8(ca.rngPick4(10, 20, 30))).WithFaceLeft(ca.rngBool()))
			}
			if ca.rngChance256(32) {
				ca.CreateSmoke(cidA, 1)
			}
		}
		return true
	}

	ca.SetCell(cidB, matB.WithStatus(MaterialStatusAcidic))

	return false
}

func ReactionAcidToFlower(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Acid has a slight chance to turn the Flower into Water
	if ca.rngChance256(10) {
		// In this case there is 50% chance to turn the Acid into Smoke
		if ca.rngBool() {
			ca.CreateSmoke(cidA, 1)
		}
		// The Flower turns into Water
		ca.SetCellAsProcessed(cidB, MaterialWater.WithFaceLeft(ca.rngBool()))
		return true
	}

	// Otherwise just turn the Flower into Acidic
	ca.SetCell(cidB, matB.WithStatus(MaterialStatusAcidic))

	return false
}

// ============================================================================
// Fire reactions (MaterialKind = 7)
// ============================================================================

func ReactionFireToWater(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// High chance to turn Water into Steam
	if ca.rngChance256(200) {
		ca.SetCellAsProcessed(cidB, MaterialSteam.WithLife(uint8(ca.rngPick4(10, 20, 30))).WithFaceLeft(ca.rngBool()))
		// Also high chance for Fire to turn into Smoke
		if ca.rngChance256(180) {
			ca.CreateSmoke(cidA, 1)
		}
		return true
	}
	return false
}

func ReactionFireToAcid(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// 30% chance to turn Acid into Fire
	if ca.rngChance256(77) {
		ca.SetCellAsProcessed(cidB, MaterialFire.WithLife(3).WithFaceLeft(ca.rngBool()).WithStatus(uint8(ca.rngPick4(64, 128, 192))))
		return true
	}
	// Otherwise they can swap
	ca.SwapCells(cidA, cidB)
	return true
}

func ReactionFireToPlant(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// 20% chance to turn Plant into Fire
	if ca.rngChance256(51) {
		ca.SetCellAsProcessed(cidB, MaterialFire.WithLife(3).WithFaceLeft(ca.rngBool()).WithStatus(uint8(ca.rngPick4(64, 128, 192))))
		return true
	}
	// 20% chance to turn Plant into Steam
	if ca.rngChance256(51) {
		ca.SetCellAsProcessed(cidB, MaterialSteam.WithLife(uint8(ca.rngPick4(10, 20, 30))).WithFaceLeft(ca.rngBool()))
		return true
	}
	// 20% chance to turn Plant into Smoke
	if ca.rngChance256(51) {
		ca.CreateSmoke(cidB, 1)
		return true
	}
	// Always set Plant to Burned status
	ca.SetCell(cidB, matB.WithStatus(MaterialStatusBurned))
	return false
}

// (no AntHill-specific reactions yet)
// ============================================================================
// Ice reactions (MaterialKind = 8)
// ============================================================================

func ReactionIceToWater(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	switch ca.rngPick4(10, 47, 60) {
	// Ice has a low chance to freeze Water, turning it into Ice with 0-1 life
	case 0:
		newLife := uint8(0)
		if ca.rngBool() {
			newLife = 1
		}
		ca.SetCellAsProcessed(cidB, MaterialIce.WithLife(newLife))
		return true
	// Ice has a slightly higher chance to swap with Water (helps it slowly sink).
	case 1:
		ca.SwapCells(cidA, cidB)
		return true
	// Water has a low chance to melt the Ice
	case 2:
		life := matA.GetLife()
		if life > 0 {
			ca.SetCellAsProcessed(cidA, MaterialIce.WithLife(life-1))
		} else {
			ca.SetCellAsProcessed(cidA, MaterialWater.WithFaceLeft(ca.rngBool()))
		}
		return true
	default:
		return false
	}
}

func ReactionIceToSeed(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	if matB.GetStatus() == MaterialStatusFrozen {
		return false
	}

	switch ca.rngPick3(20, 60) {
	case 0:
		ca.SetCellAsProcessed(cidB, matB.WithStatus(MaterialStatusFrozen))
		return true
	case 1:
		ca.SwapCells(cidA, cidB)
		return true
	default:
		return false

	}
}

func ReactionIceToRoot(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	if matB.GetStatus() == MaterialStatusFrozen {
		return false
	}
	if ca.rngChance256(12) {
		ca.SetCellAsProcessed(cidB, matB.WithStatus(MaterialStatusFrozen))
		return true
	}
	return false
}

func ReactionIceToPlant(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	if matB.GetStatus() == MaterialStatusFrozen {
		return false
	}
	if ca.rngChance256(14) {
		ca.SetCellAsProcessed(cidB, matB.WithStatus(MaterialStatusFrozen))
		return true
	}
	return false
}

func ReactionIceToFlower(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	if matB.GetStatus() == MaterialStatusFrozen {
		return false
	}
	if ca.rngChance256(18) {
		ca.SetCellAsProcessed(cidB, matB.WithStatus(MaterialStatusFrozen))
		return true
	}
	return false
}

func ReactionIceToWasp(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	if matB.GetStatus() == MaterialStatusFrozen {
		return false
	}

	switch ca.rngPick3(20, 60) {
	case 0:
		ca.SetCellAsProcessed(cidB, matB.WithStatus(MaterialStatusFrozen))
		return true
	case 1:
		ca.SwapCells(cidA, cidB)
		return true
	default:
		return false

	}
}

func ReactionRootToIce(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	if ca.rngChance256(12) {
		ca.SetCellAsProcessed(cidA, matA.WithStatus(MaterialStatusFrozen))
		return true
	}

	return false
}

func ReactionPlantToIce(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	if ca.rngChance256(14) {
		ca.SetCellAsProcessed(cidA, matA.WithStatus(MaterialStatusFrozen))
		return true
	}

	return false
}

func ReactionWaspToIce(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	switch ca.rngPick3(20, 60) {
	case 0:
		ca.SetCellAsProcessed(cidA, matA.WithStatus(MaterialStatusFrozen))
		return true
	case 1:
		ca.SwapCells(cidA, cidB)
		return true
	default:
		return false
	}
}

func ReactionWaterToIce(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	switch ca.rngPick4(11, 19, 35) {
	// Ice has a low chance to freeze Water, turning it into Ice with 0-1 life
	case 0:
		newLife := uint8(0)
		if ca.rngBool() {
			newLife = 1
		}
		ca.SetCellAsProcessed(cidA, MaterialIce.WithLife(newLife))
		return true
	// Ice has a slightly higher chance to swap with Water (helps it slowly sink).
	case 1:
		ca.SwapCells(cidA, cidB)
		return true
	// Water has a low chance to melt the Ice
	case 2:
		life := matB.GetLife()
		if life > 0 {
			ca.SetCellAsProcessed(cidB, MaterialIce.WithLife(life-1))
		} else {
			ca.SetCellAsProcessed(cidB, MaterialWater.WithFaceLeft(ca.rngBool()))
		}
		return true
	default:
		return false
	}
}

func ReactionIceToAcid(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Ice has a low chance to swap with Acid
	if ca.rngChance256(20) {
		ca.SwapCells(cidA, cidB)
		return true
	}

	return false
}

func ReactionAcidToIce(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	switch ca.rngPick3(20, 150) {
	// Acid has a low chance to swap with Ice
	case 0:
		ca.SwapCells(cidA, cidB)
		return true
	// Acid has a high chance to melt Ice by lowering its life until it turns into Water
	case 1:
		life := matB.GetLife()
		if life == 0 {
			ca.SetCellAsProcessed(cidB, MaterialWater.WithLife(uint8(ca.rngPick4(64, 128, 192))).WithFaceLeft(ca.rngBool()))
		} else {
			ca.SetCellAsProcessed(cidB, matB.WithLife(life-1))
		}
		return true
	default:
		return false

	}
}

func ReactionFireToIce(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Fire has a high chance to melt Ice
	if ca.rngChance256(200) {
		life := matB.GetLife()
		if life == 0 {
			// Ice melts into Water
			ca.SetCellAsProcessed(cidB, MaterialWater.WithLife(uint8(ca.rngPick4(64, 128, 192))).WithFaceLeft(ca.rngBool()))
			// Fire might turn into Smoke
			if ca.rngChance256(100) {
				ca.CreateSmoke(cidB, 1)
			}
		} else {
			ca.SetCellAsProcessed(cidB, matB.WithLife(life-1))
		}
		return true
	}

	return false
}

func ReactionIceToFire(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Ice has a slight chance to turn Fire into Smoke
	if ca.rngChance256(30) {
		ca.CreateSmoke(cidB, 1)
		return true
	}

	// Ice might melt when touching fire
	if ca.rngChance256(100) {
		life := matA.GetLife()
		if life == 0 {
			ca.SetCellAsProcessed(cidA, MaterialWater.WithLife(uint8(ca.rngPick4(64, 128, 192))).WithFaceLeft(ca.rngBool()))
		} else {
			ca.SetCellAsProcessed(cidA, matA.WithLife(life-1))
		}
		return true
	}

	return false
}

func ReactionIceToSteam(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Ice has a slight chance to turn Steam into Ice
	if ca.rngChance256(20) {
		newLife := uint8(0)
		if ca.rngBool() {
			newLife = 1
		}
		ca.SetCellAsProcessed(cidB, MaterialIce.WithLife(newLife))
		return true
	}

	// Ice can swap with Steam
	if ca.rngChance256(200) {
		ca.SwapCells(cidA, cidB)
		return true
	}

	return false
}

// ============================================================================
// Smoke reactions (MaterialKind = 9) - placeholder
// ============================================================================

// ============================================================================
// Steam reactions (MaterialKind = 10) - placeholder
// ============================================================================

// ============================================================================
// Root reactions (MaterialKind = 11)
// ============================================================================

func ReactionRootToSeed(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Root wakes up the Seed turning it into Root
	if ca.rngChance256(20) {
		if ca.rngBool() {
			ca.CreateRoot(cidB)
		} else {
			ca.CreateSand(cidB, false)
		}
		return true
	}
	return false
}

func ReactionRootToWater(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Root has a slight chance to gain Life from Water
	life := matA.GetLife()
	if life < 3 && ca.rngChance256(16) {
		ca.SetCellAsProcessed(cidA, matA.WithLife(life+1))
		if ca.rngChance256(64) {
			ca.SetCellAsProcessed(cidB, MaterialEmpty)
		}
		return true
	}
	return false
}

func ReactionRootToSand(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Root has a slight chance to gain Life from Sand
	life := matA.GetLife()
	if life < 3 && ca.rngChance256(8) {
		ca.SetCellAsProcessed(cidA, matA.WithLife(life+1))
		return true
	}

	// a Root on full life has a slight chance to grow into Sand (only if Sand is penetrable)
	if life == 3 && matB.GetIsPenetrable() {
		if ca.rngChance256(10) {
			ca.SetCellAsProcessed(cidA, matA.WithLife(2))
			ca.CreateRoot(cidB)
			return true
		}
	}

	return false
}

func ReactionRootToStone(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Root has a slight chance to gain Life from Stone
	life := matA.GetLife()
	if life < 3 && ca.rngChance256(60) {
		ca.SetCellAsProcessed(cidA, matA.WithLife(life+1))
		return true
	}

	// a Root on full life has a slight chance to grow into Stone (only if Stone is penetrable)
	if life == 3 && matB.GetIsPenetrable() {
		if ca.rngChance256(10) {
			ca.SetCellAsProcessed(cidA, matA.WithLife(2))
			ca.CreateRoot(cidB)
			return true
		}
	}

	return false
}

func ReactionRootToRoot(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Root has a slight chance to give life to the other Root
	lifeA := matA.GetLife()
	lifeB := matB.GetLife()
	if lifeA == lifeB {
		return false
	}
	if ca.rngChance256(50) {
		if lifeA < lifeB {
			ca.SetCellAsProcessed(cidA, matA.WithLife(lifeA+1))
			ca.SetCellAsProcessed(cidB, matB.WithLife(lifeB-1))
			return true
		}
		ca.SetCellAsProcessed(cidA, matA.WithLife(lifeA-1))
		ca.SetCellAsProcessed(cidB, matB.WithLife(lifeB+1))

		return true
	}

	return false
}

func ReactionRootToPlant(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	lifeA := matA.GetLife()
	lifeB := matB.GetLife()
	if lifeA == 0 || lifeB == 3 {
		return false
	}

	// Root has a small chance to Grow into the Plant (only if Plant is penetrable)
	if matB.GetIsPenetrable() && ca.rngChance256(60) {
		ca.SetCellAsProcessed(cidA, matA.WithLife(lifeA-1))
		ca.CreateRoot(cidB)
		return true
	}

	// Root has chance to give life to the Plant
	if ca.rngChance256(180) {
		ca.SetCellAsProcessed(cidA, matA.WithLife(lifeA-1))
		ca.SetCellAsProcessed(cidB, matB.WithLife(lifeB+1))
		return true
	}

	return false
}

// ============================================================================
// Plant reactions (MaterialKind = 12)
// ============================================================================

func ReactionPlantToWater(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// There is a slight chance for the Plant to "drink" the Water.
	if ca.rngChance256(10) {
		life := matA.GetLife()
		if life < 3 {
			ca.SetCellAsProcessed(cidA, matA.WithLife(life+1))
			ca.SetCellAsProcessed(cidB, MaterialEmpty)
			return true
		}
	}

	return false
}

func ReactionPlantToSeed(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	switch ca.rngPick3(20, 25) {
	case 0:
		ca.CreatePlant(cidB)
		return true
	case 1:
		ca.CreateRoot(cidB)
		return true
	default:
		return false
	}
}

func ReactionPlantToRoot(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Plant "sucks" Life from Root (Root loses 1 life, Plant gains 1 life).
	lifeA := matA.GetLife()
	lifeB := matB.GetLife()
	if lifeB == 0 || lifeA == 3 {
		return false
	}

	// Plant has a really small chance to Grow into the Root
	if ca.rngChance256(1) {
		ca.CreatePlant(cidB)
		return true
	}

	// Use same chance as Root->Plant life transfer, but in the opposite direction.
	if ca.rngChance256(180) {
		ca.SetCellAsProcessed(cidA, matA.WithLife(lifeA+1))
		ca.SetCellAsProcessed(cidB, matB.WithLife(lifeB-1))
		return true
	}
	return false
}

func ReactionPlantToPlant(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Plant pushes Life to other Plants (30% chance).
	lifeA := matA.GetLife()
	lifeB := matB.GetLife()
	if lifeA == 0 || lifeB == 3 {
		return false
	}
	if ca.rngChance256(77) {
		ca.SetCellAsProcessed(cidA, matA.WithLife(lifeA-1))
		ca.SetCellAsProcessed(cidB, matB.WithLife(lifeB+1))
		return true
	}
	return false
}

// ============================================================================
// Flower reactions (MaterialKind = 13) - placeholder
// ============================================================================

// ============================================================================
// Ant reactions (MaterialKind = 14) - placeholder
// ============================================================================

func ReactionAntToStone(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Very low chance to move into Stone, only if penetrable.
	if !matB.GetIsPenetrable() || !ca.rngChance256(2) {
		return false
	}
	// Swap, but the moved Stone becomes Sand.
	ca.SwapCells(cidA, cidB)

	// Create AntHill
	ca.CreateAntHill(cidA)
	return true
}

func ReactionAntToSand(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	if !matB.GetIsPenetrable() {
		return false
	}

	if ca.rngChance256(60) {
		ca.SetCellAsProcessed(cidB, matA)
		ca.CreateAntHill(cidA)
		return true
	}

	return false
}

func ReactionAntToWater(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Equal chance with Water to swap.
	if ca.rngChance256(16) {
		life := matA.GetLife()
		if life > 1 {
			ca.SetCellAsProcessed(cidA, matA.WithLife(life-1))
			return true
		}
		ca.CreateSand(cidA, true)
		return true
	}
	if ca.rngChance256(64) {
		ca.SwapCells(cidA, cidB)
		return true
	}
	return false
}

func ReactionWaterToAnt(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Equal chance with Ant to swap (symmetry).
	if ca.rngChance256(32) {
		ca.SwapCells(cidA, cidB)
		return true
	}
	return false
}

func ReactionWaterToAntHill(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	if ca.rngChance256(200) {
		return false
	}
	ca.SetCellAsProcessed(cidB, matA)
	ca.SetCellAsProcessed(cidA, MaterialEmpty)
	return true
}

func ReactionAcidToAntHill(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	switch ca.rngPick3(16, 64) {
	case 0:
		ca.SetCellAsProcessed(cidB, matA)
		ca.SetCellAsProcessed(cidA, MaterialEmpty)
		return true
	case 1:
		ca.CreateSmoke(cidB, 2)
		if ca.rngBool() {
			ca.CreateSmoke(cidA, 2)
		}
		return true
	default:
		return false
	}
}

func ReactionWaterToStone(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	if matB.GetIsPenetrable() && ca.rngChance256(2) {
		ca.CreateSand(cidB, ca.rngBool())
		ca.SwapCells(cidA, cidB)
		return true
	}

	return false
}

// Ant can only eat when its life is below 3.
// By eating the material it gains 1 life,
// replaces itself with AntHill,
// and moves into the material's place.
func AntEatReaction(chance uint8) MaterialReaction {
	return func(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
		life := matA.GetLife()
		if life == 3 || !ca.rngChance256(chance) {
			return false
		}

		ca.CreateAntHill(cidA)
		ca.SetCellAsProcessed(cidB, matA.WithLife(life+1))
		return true
	}
}

func ReactionAntToAcid(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Acid turns the Ant Acidic and can deal damage.
	ant := matA.WithStatus(MaterialStatusAcidic)

	// Slightly higher chance to deal damage / kill.
	if ca.rngChance256(20) {
		life := ant.GetLife()
		if life > 1 {
			ant = ant.WithLife(life - 1)
		} else {
			// Ant dies -> Smoke. Acid has 50% to turn into Smoke.
			ca.CreateSmoke(cidA, 1)
			if ca.rngBool() {
				ca.CreateSmoke(cidB, 2)
			}

			return true
		}
	}

	// Both can swap with each other.
	ca.SwapCells(cidA, cidB)
	ca.SetCell(cidB, ant)
	return true
}

func ReactionAcidToAnt(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Symmetric to ReactionAntToAcid: Acid turns Ant Acidic and can deal damage.
	ant := matB.WithStatus(MaterialStatusAcidic)

	// Slightly higher chance to deal damage / kill.
	if ca.rngChance256(20) {
		life := ant.GetLife()
		if life > 1 {
			ant = ant.WithLife(life - 1)
		} else {
			// Ant dies -> Smoke. Acid has 50% to turn into Smoke.
			ca.CreateSmoke(cidB, 1)
			if ca.rngBool() {
				ca.CreateSmoke(cidB, 2)
			}
			return true
		}
	}

	// Both can swap with each other.
	ca.SwapCells(cidA, cidB)
	ca.SetCell(cidA, ant)
	return true
}

func ReactionFireToAnt(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Fire kills Ant: Ant turns into Smoke.
	ca.CreateSmoke(cidB, 1)
	return true
}

func ReactionAntToFire(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Fire kills Ant: Ant turns into Smoke.
	ca.CreateSmoke(cidB, 1)
	return true
}

// ============================================================================
// Wasp reactions (MaterialKind = 15) - placeholder
// ============================================================================

func ReactionWaspToWater(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Wasp can "eat" one Water to enable egg laying later.
	// It consumes the Water (turns it into Empty) without moving.
	if matA.GetWaspHasWater() {
		// Already collected; avoid Water.
		return false
	}
	if !ca.rngChance256(200) {
		return false
	}
	ca.SetCellAsProcessed(cidB, MaterialEmpty)
	ca.SetCellAsProcessed(cidA, matA.WithWaspHasWater(true))
	return true
}

func ReactionWaspToAnt(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// If the Ant is still an egg (Ant with Life==0), Wasp can eat it.
	if matB.GetLife() == 0 {
		if !ca.rngChance256(160) {
			return false
		}
		wasp := matA.WithWaspHasAnt(true)
		life := wasp.GetLife()
		if life < 3 {
			life++
		}
		ca.SetCellAsProcessed(cidB, MaterialEmpty)
		ca.SetCellAsProcessed(cidA, wasp.WithLife(life))
		return true
	}

	// Instant fight resolution on contact:
	// - Wasp wins ~90%
	// - Ant wins otherwise
	// Loser disappears (becomes Empty).
	// Rewards:
	// - If Wasp wins: sets hasAnt=true and gains +1 Life (capped at 3)
	// - If Ant wins: Ant Life becomes 3
	waspWins := ca.rngChance256(230) // ~90%
	if waspWins {
		// Ant disappears.
		ca.SetCellAsProcessed(cidB, MaterialEmpty)

		// Wasp reward.
		wasp := matA.WithWaspHasAnt(true)
		life := wasp.GetLife()
		if life < 3 {
			life++
		}
		ca.SetCellAsProcessed(cidA, wasp.WithLife(life))
		return true
	}

	// Wasp disappears; Ant reward.
	ca.SetCellAsProcessed(cidA, MaterialEmpty)
	ca.SetCellAsProcessed(cidB, matB.WithLife(3))
	return true
}

func ReactionAntToWasp(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// If the Wasp is still an egg (Wasp with Life==0), Ant can eat it.
	if matB.GetLife() == 0 {
		if !ca.rngChance256(120) {
			return false
		}
		life := matA.GetLife()
		if life < 3 {
			life++
		}
		ca.SetCellAsProcessed(cidB, MaterialEmpty)
		ca.SetCellAsProcessed(cidA, matA.WithLife(life))
		return true
	}

	// Delegate to the same fight logic as ReactionWaspToAnt, but with swapped roles/cells.
	return ReactionWaspToAnt(ca, matB, matA, cidB, cidA)
}

func ReactionWaspToAcid(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	if ca.rngChance256(40) {
		life := matA.GetLife()
		if life > 1 {
			ca.SetCellAsProcessed(cidA, matA.WithLife(life-1).WithStatus(MaterialStatusAcidic))
		} else {
			ca.CreateSmoke(cidA, 1)
			if ca.rngBool() {
				ca.CreateSmoke(cidB, 2)
			}
			return true
		}
	}

	if ca.rngChance256(64) {
		ca.SwapCells(cidA, cidB)
	}

	return true
}

func ReactionAcidToWasp(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Symmetric to ReactionWaspToAcid.
	if ca.rngChance256(40) {
		life := matB.GetLife()
		if life > 1 {
			ca.SetCellAsProcessed(cidB, matB.WithLife(life-1).WithStatus(MaterialStatusAcidic))
		} else {
			ca.CreateSmoke(cidB, 2)
			if ca.rngBool() {
				ca.CreateSmoke(cidA, 1)
			}
			return true
		}
	}

	if ca.rngChance256(64) {
		ca.SwapCells(cidA, cidB)
	}

	return true
}

func ReactionFireToWasp(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// The Wasp always turns to Burned status when touched by Fire
	ca.SetCell(cidB, matB.WithStatus(MaterialStatusBurned))

	// Fire has a chance to lower the Wasp's Life
	if ca.rngChance256(80) {
		life := matB.GetLife()
		if life > 1 {
			ca.SetCellAsProcessed(cidB, matB.WithLife(life-1))
			return true
		}

		// Fire kills Wasp: Wasp turns into Burned Sand.
		ca.CreateSand(cidB, true)

		// Flip a coin to decide if the Fire turns into Smoke
		if ca.rngBool() {
			ca.CreateSmoke(cidB, 3)
		}
		return true
	}

	// Otherwise just swap (Fire moves through Wasp)
	ca.SwapCells(cidA, cidB)
	return true
}

func ReactionWaspToFire(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// The Wasp always turns to Burned status when touched by Fire
	ca.SetCell(cidA, matA.WithStatus(MaterialStatusBurned))

	// Fire has a chance to lower the Wasp's Life
	if ca.rngChance256(80) {
		life := matA.GetLife()
		if life > 1 {
			ca.SetCellAsProcessed(cidA, matA.WithLife(life-1))
			return true
		}

		// Fire kills Wasp: Wasp turns into Burned Sand.
		ca.CreateSand(cidA, true)

		// Flip a coin to decide if the Fire turns into Smoke
		if ca.rngBool() {
			ca.CreateSmoke(cidB, 3)
		}
		return true
	}

	// Otherwise just swap (Wasp moves through Fire)
	ca.SwapCells(cidA, cidB)
	return true
}

func ReactionWaspToSteam(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Wasp has a slight chance to lose life when touching Steam.
	if ca.rngChance256(25) {
		life := matA.GetLife()
		if life > 1 {
			ca.SetCellAsProcessed(cidA, matA.WithLife(life-1))
		} else {
			// Wasp dies -> turns into Sand
			ca.CreateSand(cidA, true)
		}
		return true
	}

	// Otherwise just swap (Wasp moves through Steam)
	if ca.rngChance256(220) {
		ca.SwapCells(cidA, cidB)
		return true
	}
	return false
}

func ReactionWaspToSmoke(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// Wasp has a slight chance to lose life when touching Smoke.
	if ca.rngChance256(50) {
		life := matA.GetLife()
		if life > 1 {
			ca.SetCellAsProcessed(cidA, matA.WithLife(life-1).WithStatus(MaterialStatusBurned))
		} else {
			// Wasp dies -> turns into Sand
			ca.CreateSand(cidA, true)
		}
		return true
	}

	// Otherwise just swap (Wasp moves through Smoke)
	if ca.rngChance256(220) {
		ca.SwapCells(cidA, cidB)
		return true
	}
	return false
}

func ReactionIceToSand(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	if ca.rngChance256(80) {
		ca.SetCellAsProcessed(cidB, matB.WithStatus(MaterialStatusFrozen))
		return true
	}
	return false
}

func ReactionFireToAntHill(ca *CellAutomata, matA, matB Material, cidA, cidB int) bool {
	// The AntHill always turns to Burned status when touched by Fire
	matB = matB.WithStatus(MaterialStatusBurned)

	switch ca.rngPick4(80, 120, 160) {
	// There is a slight chance for the AntHill to catch on Fire
	case 0:
		ca.SetCellAsProcessed(cidB, MaterialFire.WithLife(3).WithStatus(ca.rng0123()).WithFaceLeft(ca.rngBool()))
		return true
	// There is a slight chance for the AntHill to turn into Smoke
	case 1:
		if ca.rngBool() {
			ca.CreateSmoke(cidB, 3)
		} else {
			ca.CreateSand(cidB, true)
		}
		return true
	// There is a chance for the Fire to "move into the AntHill"
	case 2:
		ca.SetCellAsProcessed(cidB, matA)
		ca.SetCell(cidA, MaterialEmpty)
		return true
	// The default case is no reaction
	default:
		return false
	}
}
