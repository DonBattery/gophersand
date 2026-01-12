// material.go
package game

/*
Material is a 16-bit packed value containing all information about one cell (pixel).

Layout (LSB -> MSB):

	bits 0..3   : MaterialKind (0..15)
	bits 4..5   : MaterialLife (0..3)
	bits 6..7   : MaterialStatus (0..3)  0=Normal,1=Burned,2=Acidic,3=Frozen
	bits 8..15  : Material-specific state (8 bits)

State-byte canonical map (bits 8..15):

	bit 8  : FaceLeft
	bit 9  : FaceUp
	bit 10 : FlagA
	bit 11 : FlagB
	bit 12 : FlagC
	bit 13 : FlagD
	bit 14 : FlagE
	bit 15 : FlagF
*/
type Material uint16

// All the 16 Materials (base value == kind)
const (
	MaterialEmpty Material = iota
	MaterialStone
	MaterialSand
	MaterialWater
	MaterialSeed
	MaterialAnt
	MaterialWasp
	MaterialAcid
	MaterialFire
	MaterialIce
	MaterialSmoke
	MaterialSteam
	MaterialRoot
	MaterialPlant
	MaterialFlower
	MaterialAntHill
)

// MaterialKind is a 4-bit number representing the 16 Material Kinds (including EmptyKind==0).
type MaterialKind uint8

// All the 16 MaterialKinds
const (
	MaterialKindEmpty MaterialKind = iota
	MaterialKindStone
	MaterialKindSand
	MaterialKindWater
	MaterialKindSeed
	MaterialKindAnt
	MaterialKindWasp
	MaterialKindAcid
	MaterialKindFire
	MaterialKindIce
	MaterialKindSmoke
	MaterialKindSteam
	MaterialKindRoot
	MaterialKindPlant
	MaterialKindFlower
	MaterialKindAntHill
)

const (
	MaterialStatusNormal uint8 = iota
	MaterialStatusBurned
	MaterialStatusAcidic
	MaterialStatusFrozen
)

// Bit shifting constants and masks
const (
	// Core fields
	kindMask   Material = 0x000F // bits 0..3
	lifeMask   Material = 0x0030 // bits 4..5
	statusMask Material = 0x00C0 // bits 6..7

	lifeShift   = 4
	statusShift = 6
	dataShift   = 8 // start of state byte

	// Canonical state-byte bits (8..15)
	stateFaceLeft Material = 1 << 8
	stateFaceUp   Material = 1 << 9
	stateFlagA    Material = 1 << 10
	stateFlagB    Material = 1 << 11
	stateFlagC    Material = 1 << 12
	stateFlagD    Material = 1 << 13
	stateFlagE    Material = 1 << 14
	stateFlagF    Material = 1 << 15

	// Per-kind aliases (reused bits across kinds)
	// FlagA: “generic kind-flag” reused for Sand/Stone/IsTopPetal/etc.
	isPenetrableBit Material = stateFlagA // Sand/Stone/Plant (unified)

	// Plant must NOT share this with IsPenetrable, because Plant uses both.
	canBloomBit Material = stateFlagD // Plant (moved off FlagA)

	// Flower can keep FlagA (or move it too if you ever need penetrable on Flower)
	isTopPetalBit Material = stateFlagA // Flower

	// Wasp collection flags
	waspHasWaterBit Material = stateFlagB
	waspHasAntBit   Material = stateFlagC
)

// MaterialKindSet is a 16-bit bit-field. Each bit indicates if the corresponding MaterialKind is part of the set.
type MaterialKindSet uint16

var (
	// Materials which can be penetrated by Root growth.
	// Note: Sand/Stone/Plant are handled specially in reactions (checks IsPenetrable flag).
	RootGrowableKinds = NewMaterialKindSet(
		MaterialKindAntHill,
		MaterialKindEmpty,
		MaterialKindSteam,
		MaterialKindSmoke,
		MaterialKindPlant,
	)

	// Materials which can be overwritten by Plant growth
	PlantGrowableKinds = NewMaterialKindSet(
		MaterialKindEmpty,
		MaterialKindAntHill,
		MaterialKindSteam,
		MaterialKindSmoke,
		MaterialKindRoot,
	)

	PlantSupporterKinds = NewMaterialKindSet(
		MaterialKindStone,
		MaterialKindSand,
		MaterialKindSeed,
		MaterialKindRoot,
		MaterialKindPlant,
		MaterialKindFlower,
	)

	// Materials an Ant can "grip" to avoid falling when adjacent.
	AntSupporterKinds = NewMaterialKindSet(
		MaterialKindStone,
		MaterialKindSand,
		MaterialKindSeed,
		MaterialKindAntHill,
		MaterialKindRoot,
		MaterialKindPlant,
		MaterialKindFlower,
		MaterialKindAnt,
	)

	AntAliveKinds = NewMaterialKindSet(
		MaterialKindSeed,
		MaterialKindRoot,
		MaterialKindPlant,
		MaterialKindFlower,
		MaterialKindAntHill,
	)

	// Materials an Ant is allowed to lay eggs into (eggs are represented as Ant with Life==0).
	AntEggLayableKinds = NewMaterialKindSet(
		MaterialKindEmpty,
		MaterialKindAntHill,
		MaterialKindSteam,
		MaterialKindSmoke,
		MaterialKindRoot,
		MaterialKindPlant,
		MaterialKindFlower,
	)

	// Materials an Ant can fall through when unsupported (gravity simulation).
	AntFallableKinds = NewMaterialKindSet(
		MaterialKindEmpty,
		MaterialKindSteam,
		MaterialKindSmoke,
		MaterialKindFire,
		MaterialKindWater,
		MaterialKindAcid,
	)

	// Materials Wasp eggs (Wasp with Life==0) can stick to (sideways or when hanging under them).
	WaspEggStickyKinds = NewMaterialKindSet(
		MaterialKindStone,
		MaterialKindSand,
		MaterialKindSeed,
		MaterialKindRoot,
		MaterialKindPlant,
		MaterialKindFlower,
	)

	// Materials a Wasp is allowed to lay eggs into (eggs are represented as Wasp with Life==0).
	WaspEggLayableKinds = NewMaterialKindSet(
		MaterialKindEmpty,
		MaterialKindAntHill,
		MaterialKindSteam,
		MaterialKindSmoke,
		MaterialKindPlant,
	)

	// Materials Ice can Freeze.
	FreezableKinds = NewMaterialKindSet(
		MaterialKindEmpty,
		MaterialKindSteam,
		MaterialKindSmoke,
		MaterialKindPlant,
	)

	// Materials that can be "eaten" by plants (they are not destroyed)
	PlantFoodKinds = NewMaterialKindSet(
		MaterialKindWater,
		MaterialKindSand,
	)

	// Steam cannot condense into Water below these Materials
	NonCondensableKinds = NewMaterialKindSet(
		MaterialKindEmpty,
		MaterialKindSteam,
		MaterialKindSmoke,
		MaterialKindWater,
		MaterialKindAcid,
		MaterialKindFire,
	)
)

// NewMaterialKindSet creates a MaterialKindSet from a list of MaterialKind by setting the corresponding bits to 1.
func NewMaterialKindSet(kinds ...MaterialKind) MaterialKindSet {
	var set MaterialKindSet
	for _, k := range kinds {
		set |= (MaterialKindSet(1) << k)
	}
	return set
}

// IsIn returns true if the material kind is in the given MaterialKindSet.
func (mk MaterialKind) IsIn(set MaterialKindSet) bool {
	return (set & (MaterialKindSet(1) << mk)) != 0
}

// GetKind returns the Kind of the material.
func (m Material) GetKind() MaterialKind {
	return MaterialKind(m & kindMask)
}

// IsKind returns true if the material is of the given Kind.
func (m Material) IsKind(k MaterialKind) bool {
	return m.GetKind() == k
}

// IsIn returns true if the Material's Kind is in the given MaterialKindSet.
func (m Material) IsIn(set MaterialKindSet) bool {
	return m.GetKind().IsIn(set)
}

// GetLife returns the Life of the material.
func (m Material) GetLife() uint8 {
	return uint8((m & lifeMask) >> lifeShift)
}

// WithLife sets the Life of the material and returns the new material.
func (m Material) WithLife(life uint8) Material {
	life &= 3
	return (m &^ lifeMask) | (Material(life) << lifeShift)
}

// GetStatus returns the Status of the material.
func (m Material) GetStatus() uint8 {
	return uint8((m & statusMask) >> statusShift)
}

// WithStatus sets the Status of the material and returns the new material.
func (m Material) WithStatus(status uint8) Material {
	status &= 3
	return (m &^ statusMask) | (Material(status) << statusShift)
}

// GetColor returns the display color of this material based on its kind, status, and life.
// Index formula: kind*16 + status*4 + life.
func (m Material) GetColor() Color {
	return MaterialColors[int(m.GetKind())*16+int(m.GetStatus())*4+int(m.GetLife())]
}

// -----------------------------------------------------------------------------
// Shared direction bits (preferred movement / intent bits)
// -----------------------------------------------------------------------------

// FaceLeft at state bit 8.
func (m Material) GetFaceLeft() bool {
	return (m & stateFaceLeft) != 0
}
func (m Material) WithFaceLeft(on bool) Material {
	if on {
		return m | stateFaceLeft
	}
	return m &^ stateFaceLeft
}

// FaceUp at state bit 9.
// Intended for Ant/Wasp “vertical desire” (up vs down). Can be reused elsewhere if needed.
func (m Material) GetFaceUp() bool {
	return (m & stateFaceUp) != 0
}
func (m Material) WithFaceUp(on bool) Material {
	if on {
		return m | stateFaceUp
	}
	return m &^ stateFaceUp
}

// -----------------------------------------------------------------------------
// Plant / Flower / Sand+Stone state helpers (bit-reused per kind)
// -----------------------------------------------------------------------------

// CanBloom at state FlagA - indicates if a Plant can bloom into flowers.
func (m Material) GetCanBloom() bool {
	return (m & canBloomBit) != 0
}
func (m Material) WithCanBloom(on bool) Material {
	if on {
		return m | canBloomBit
	}
	return m &^ canBloomBit
}

// IsTopPetal at state FlagA - indicates if this Flower is the top petal that can drop seeds.
func (m Material) GetIsTopPetal() bool {
	return (m & isTopPetalBit) != 0
}
func (m Material) WithIsTopPetal(on bool) Material {
	if on {
		return m | isTopPetalBit
	}
	return m &^ isTopPetalBit
}

// IsPenetrable at state FlagA - indicates if this Sand/Stone/Plant can be penetrated by Root growth.
func (m Material) GetIsPenetrable() bool {
	return (m & isPenetrableBit) != 0
}
func (m Material) WithIsPenetrable(on bool) Material {
	if on {
		return m | isPenetrableBit
	}
	return m &^ isPenetrableBit
}

// -----------------------------------------------------------------------------
// Wasp-specific state helpers (collection flags)
// -----------------------------------------------------------------------------

func (m Material) GetWaspHasWater() bool {
	return (m & waspHasWaterBit) != 0
}
func (m Material) WithWaspHasWater(on bool) Material {
	if on {
		return m | waspHasWaterBit
	}
	return m &^ waspHasWaterBit
}

func (m Material) GetWaspHasAnt() bool {
	return (m & waspHasAntBit) != 0
}
func (m Material) WithWaspHasAnt(on bool) Material {
	if on {
		return m | waspHasAntBit
	}
	return m &^ waspHasAntBit
}
