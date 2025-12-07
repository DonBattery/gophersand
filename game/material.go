package game

import "math/rand"

// Material is a 32-bit value encoding type, kind flags, and data
// Layout:
//   Byte 0 (bits 0-7):   Material Type (256 possible materials)
//   Byte 1 (bits 8-15):  Material Kind Flags (8 category flags)
//   Bytes 2-3 (bits 16-31): Material Data (for state machines)
type Material uint32

type MaterialType byte
type MaterialKindFlag uint32
type MaterialKindFilter uint32

// Bit layout constants
const (
	MaterialTypeMask  uint32 = 0xFF       // Bits 0-7
	MaterialKindMask  uint32 = 0xFF00     // Bits 8-15
	MaterialDataMask  uint32 = 0xFFFF0000 // Bits 16-31
	MaterialKindShift        = 8
	MaterialDataShift        = 16
)

// Material types (byte values)
const (
	TypeEmpty MaterialType = 0
	TypeStone MaterialType = 1
	TypeSand  MaterialType = 2
	TypeWater MaterialType = 3
	TypeSmoke MaterialType = 4
)

// Material kind flags (bit positions in the kind byte)
const (
	KindStatic MaterialKindFlag = 1 << (MaterialKindShift + 0)
	KindSolid  MaterialKindFlag = 1 << (MaterialKindShift + 1)
	KindLiquid MaterialKindFlag = 1 << (MaterialKindShift + 2)
	KindGas    MaterialKindFlag = 1 << (MaterialKindShift + 3)
	// Can add up to 4 more flags (bits 12-15)
)

// Predefined materials
var (
	MaterialEmpty = NewMaterial(TypeEmpty)
	MaterialStone = NewMaterial(TypeStone, KindStatic, KindSolid)
	MaterialSand  = NewMaterial(TypeSand, KindSolid)
	MaterialWater = NewMaterial(TypeWater, KindLiquid)
	MaterialSmoke = NewMaterial(TypeSmoke, KindGas)
)

// Predefined filters
var (
	StaticMaterials     = NewMaterialKindFilter(KindStatic)
	SolidMaterials      = NewMaterialKindFilter(KindSolid)
	LiquidMaterials     = NewMaterialKindFilter(KindLiquid)
	GasMaterials        = NewMaterialKindFilter(KindGas)
	PenetrableMaterials = NewMaterialKindFilter(KindLiquid, KindGas)
)

// Material names and colors
var (
	MaterialNames = map[MaterialType]string{
		TypeEmpty: "Empty",
		TypeStone: "Stone",
		TypeSand:  "Sand",
		TypeWater: "Water",
		TypeSmoke: "Smoke",
	}

	MaterialColors = map[MaterialType]Color{
		TypeEmpty: ColorNull,
		TypeStone: ColorGray,
		TypeSand:  ColorLighterYellow,
		TypeWater: ColorLighterBlue,
		TypeSmoke: ColorWhite,
	}
)

// MaterialDataField defines a range of bits in the data portion
type MaterialDataField struct {
	Pos int // Position relative to bit 16 (start of data section)
	Len int // Length in bits
}

// Common data fields
var (
	DirectionData = MaterialDataField{Pos: 0, Len: 1}
)

// Constructor functions
func NewMaterial(t MaterialType, flags ...MaterialKindFlag) Material {
	m := Material(t)
	for _, f := range flags {
		m |= Material(f)
	}
	return m
}

func NewMaterialKindFilter(flags ...MaterialKindFlag) MaterialKindFilter {
	var f MaterialKindFilter
	for _, flag := range flags {
		f |= MaterialKindFilter(flag)
	}
	return f
}

// Type and Kind component methods

// GetType returns the material type (first byte)
func (m Material) GetType() byte {
	return byte(uint32(m) & MaterialTypeMask)
}

// IsType checks if the material matches the given type
func (m Material) IsType(t MaterialType) bool {
	return m.GetType() == byte(t)
}

// FilterAny returns true if any flags from the filter are set
func (m Material) FilterAny(filter MaterialKindFilter) bool {
	if filter == 0 {
		return true // Empty filter matches everything
	}
	return (uint32(m) & uint32(filter)) != 0
}

// FilterAll returns true if all flags from the filter are set
func (m Material) FilterAll(filter MaterialKindFilter) bool {
	return (uint32(m) & uint32(filter)) == uint32(filter)
}

// FilterNone returns true if none of the flags from the filter are set
func (m Material) FilterNone(filter MaterialKindFilter) bool {
	return (uint32(m) & uint32(filter)) == 0
}

// Data component methods (operate on bits 16-31)

// Get returns true if the n-th bit in the data section is set (0-15)
func (m Material) Get(n int) bool {
	return (m & (1 << (MaterialDataShift + n))) != 0
}

// Set returns a new Material with the n-th data bit set to 1
func (m Material) Set(n int) Material {
	return m | (1 << (MaterialDataShift + n))
}

// Unset returns a new Material with the n-th data bit set to 0
func (m Material) Unset(n int) Material {
	return m &^ (1 << (MaterialDataShift + n))
}

// SetBool returns a new Material with the n-th data bit set to val
func (m Material) SetBool(n int, val bool) Material {
	if val {
		return m.Set(n)
	}
	return m.Unset(n)
}

// GetInt returns the integer value of i bits starting from bit n in the data section
func (m Material) GetInt(n, i int) int {
	mask := (uint32(1) << i) - 1
	return int((uint32(m) >> (MaterialDataShift + n)) & mask)
}

// SetInt returns a new Material with val written using i bits starting at bit n in the data section
func (m Material) SetInt(n, i, val int) Material {
	shift := MaterialDataShift + n
	mask := ((uint32(1) << i) - 1) << shift
	return Material((uint32(m) &^ mask) | ((uint32(val) << shift) & mask))
}

// GetField returns the integer value of the MaterialDataField
func (m Material) GetField(field MaterialDataField) int {
	return m.GetInt(field.Pos, field.Len)
}

// SetField returns a new Material with val written to the MaterialDataField
func (m Material) SetField(field MaterialDataField, val int) Material {
	return m.SetInt(field.Pos, field.Len, val)
}

// RandomField returns a new Material with a random value in the MaterialDataField
func (m Material) RandomField(field MaterialDataField) Material {
	return m.SetInt(field.Pos, field.Len, rand.Intn(1<<field.Len))
}

// Convenience methods

func (m Material) IsEmpty() bool {
	return m.IsType(TypeEmpty)
}

func (m Material) IsStatic() bool {
	return m.FilterAny(StaticMaterials)
}

func (m Material) IsSolid() bool {
	return m.FilterAny(SolidMaterials)
}

func (m Material) IsLiquid() bool {
	return m.FilterAny(LiquidMaterials)
}

func (m Material) IsGas() bool {
	return m.FilterAny(GasMaterials)
}

func (m Material) IsPenetrable() bool {
	return m.IsEmpty() || m.FilterAny(PenetrableMaterials)
}

func (m Material) IsFlowable() bool {
	return m.IsEmpty() || m.IsGas()
}

func (m Material) String() string {
	return MaterialNames[MaterialType(m.GetType())]
}

func (m Material) Color() Color {
	return MaterialColors[MaterialType(m.GetType())]
}

// Domain-specific data helpers

func (m Material) RandomDirection() Material {
	return m.RandomField(DirectionData)
}

func (m Material) FaceLeft() bool {
	return m.Get(DirectionData.Pos)
}
