package game

import (
	"testing"
)

// Test basic material construction
func TestNewMaterial(t *testing.T) {
	tests := []struct {
		name     string
		matType  MaterialType
		flags    []MaterialKindFlag
		expected uint32
	}{
		{
			name:     "Empty material",
			matType:  TypeEmpty,
			flags:    nil,
			expected: 0,
		},
		{
			name:     "Stone with flags",
			matType:  TypeStone,
			flags:    []MaterialKindFlag{KindStatic, KindSolid},
			expected: uint32(TypeStone) | uint32(KindStatic) | uint32(KindSolid),
		},
		{
			name:     "Sand with solid flag",
			matType:  TypeSand,
			flags:    []MaterialKindFlag{KindSolid},
			expected: uint32(TypeSand) | uint32(KindSolid),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMaterial(tt.matType, tt.flags...)
			if uint32(m) != tt.expected {
				t.Errorf("NewMaterial() = %d, want %d", uint32(m), tt.expected)
			}
		})
	}
}

// Test GetType method
func TestMaterial_GetType(t *testing.T) {
	tests := []struct {
		name     string
		material Material
		want     byte
	}{
		{"Empty", MaterialEmpty, byte(TypeEmpty)},
		{"Stone", MaterialStone, byte(TypeStone)},
		{"Sand", MaterialSand, byte(TypeSand)},
		{"Water", MaterialWater, byte(TypeWater)},
		{"Smoke", MaterialSmoke, byte(TypeSmoke)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.material.GetType(); got != tt.want {
				t.Errorf("Material.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test IsType method
func TestMaterial_IsType(t *testing.T) {
	tests := []struct {
		name      string
		material  Material
		checkType MaterialType
		want      bool
	}{
		{"Stone is Stone", MaterialStone, TypeStone, true},
		{"Stone is not Sand", MaterialStone, TypeSand, false},
		{"Water is Water", MaterialWater, TypeWater, true},
		{"Empty is Empty", MaterialEmpty, TypeEmpty, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.material.IsType(tt.checkType); got != tt.want {
				t.Errorf("Material.IsType() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test FilterAny method
func TestMaterial_FilterAny(t *testing.T) {
	tests := []struct {
		name     string
		material Material
		filter   MaterialKindFilter
		want     bool
	}{
		{"Stone matches Static filter", MaterialStone, StaticMaterials, true},
		{"Stone matches Solid filter", MaterialStone, SolidMaterials, true},
		{"Stone doesn't match Liquid filter", MaterialStone, LiquidMaterials, false},
		{"Water matches Liquid filter", MaterialWater, LiquidMaterials, true},
		{"Water matches Penetrable filter", MaterialWater, PenetrableMaterials, true},
		{"Empty filter matches everything", MaterialStone, MaterialKindFilter(0), true},
		{"Sand doesn't match Static", MaterialSand, StaticMaterials, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.material.FilterAny(tt.filter); got != tt.want {
				t.Errorf("Material.FilterAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test FilterAll method
func TestMaterial_FilterAll(t *testing.T) {
	staticAndSolid := NewMaterialKindFilter(KindStatic, KindSolid)

	tests := []struct {
		name     string
		material Material
		filter   MaterialKindFilter
		want     bool
	}{
		{"Stone has both Static and Solid", MaterialStone, staticAndSolid, true},
		{"Sand doesn't have Static", MaterialSand, staticAndSolid, false},
		{"Water doesn't have Solid", MaterialWater, SolidMaterials, false},
		{"Stone has Static", MaterialStone, StaticMaterials, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.material.FilterAll(tt.filter); got != tt.want {
				t.Errorf("Material.FilterAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test FilterNone method
func TestMaterial_FilterNone(t *testing.T) {
	tests := []struct {
		name     string
		material Material
		filter   MaterialKindFilter
		want     bool
	}{
		{"Stone has Static, so none is false", MaterialStone, StaticMaterials, false},
		{"Sand has no Static flag", MaterialSand, StaticMaterials, true},
		{"Water has no Solid flag", MaterialWater, SolidMaterials, true},
		{"Empty has no flags", MaterialEmpty, SolidMaterials, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.material.FilterNone(tt.filter); got != tt.want {
				t.Errorf("Material.FilterNone() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test data bit manipulation
func TestMaterial_GetSetBits(t *testing.T) {
	m := MaterialSand

	// Initially, bit 0 should be false
	if m.Get(0) {
		t.Error("Bit 0 should initially be false")
	}

	// Set bit 0
	m = m.Set(0)
	if !m.Get(0) {
		t.Error("Bit 0 should be true after Set(0)")
	}

	// Set bit 5
	m = m.Set(5)
	if !m.Get(5) {
		t.Error("Bit 5 should be true after Set(5)")
	}

	// Bit 0 should still be set
	if !m.Get(0) {
		t.Error("Bit 0 should still be true")
	}

	// Unset bit 0
	m = m.Unset(0)
	if m.Get(0) {
		t.Error("Bit 0 should be false after Unset(0)")
	}

	// Bit 5 should still be set
	if !m.Get(5) {
		t.Error("Bit 5 should still be true")
	}

	// Type should be preserved
	if !m.IsType(TypeSand) {
		t.Error("Material type should still be Sand")
	}
}

// Test GetInt and SetInt
func TestMaterial_GetSetInt(t *testing.T) {
	m := MaterialWater

	// Set a 4-bit value at position 0
	m = m.SetInt(0, 4, 7) // Binary: 0111
	if got := m.GetInt(0, 4); got != 7 {
		t.Errorf("GetInt(0, 4) = %d, want 7", got)
	}

	// Set a 3-bit value at position 4
	m = m.SetInt(4, 3, 5) // Binary: 101
	if got := m.GetInt(4, 3); got != 5 {
		t.Errorf("GetInt(4, 3) = %d, want 5", got)
	}

	// First value should still be there
	if got := m.GetInt(0, 4); got != 7 {
		t.Errorf("GetInt(0, 4) = %d, want 7 (should be preserved)", got)
	}

	// Type should be preserved
	if !m.IsType(TypeWater) {
		t.Error("Material type should still be Water")
	}
}

// Test SetInt with overlapping writes
func TestMaterial_SetInt_Overlapping(t *testing.T) {
	m := MaterialEmpty

	// Write 8 bits: 11111111 (255)
	m = m.SetInt(0, 8, 255)
	if got := m.GetInt(0, 8); got != 255 {
		t.Errorf("GetInt(0, 8) = %d, want 255", got)
	}

	// Overwrite middle 4 bits with 0000
	m = m.SetInt(2, 4, 0)

	// Check the result: bits should be 11000011 (195)
	if got := m.GetInt(0, 8); got != 195 {
		t.Errorf("GetInt(0, 8) = %d, want 195", got)
	}
}

// Test GetField and SetField
func TestMaterial_GetSetField(t *testing.T) {
	flowField := MaterialDataField{Pos: 4, Len: 3}

	m := MaterialWater
	m = m.SetField(flowField, 6)

	if got := m.GetField(flowField); got != 6 {
		t.Errorf("GetField() = %d, want 6", got)
	}

	// Type should be preserved
	if !m.IsType(TypeWater) {
		t.Error("Material type should still be Water")
	}
}

// Test RandomField
func TestMaterial_RandomField(t *testing.T) {
	field := MaterialDataField{Pos: 0, Len: 4}
	m := MaterialSand

	// Run multiple times to check range
	for i := 0; i < 100; i++ {
		m = m.RandomField(field)
		val := m.GetField(field)

		if val < 0 || val >= 16 {
			t.Errorf("RandomField produced out-of-range value: %d (expected 0-15)", val)
		}

		// Type should be preserved
		if !m.IsType(TypeSand) {
			t.Error("Material type should still be Sand")
		}
	}
}

// Test DirectionData field
func TestMaterial_Direction(t *testing.T) {
	m := MaterialWater

	// Initially should face right (false)
	if m.FaceLeft() {
		t.Error("Should initially face right")
	}

	// Set direction to left
	m = m.Set(DirectionData.Pos)
	if !m.FaceLeft() {
		t.Error("Should face left after Set")
	}

	// Unset direction (face right)
	m = m.Unset(DirectionData.Pos)
	if m.FaceLeft() {
		t.Error("Should face right after Unset")
	}

	// RandomDirection should produce valid values
	m = m.RandomDirection()
	_ = m.FaceLeft() // Should not panic
}

// Test convenience methods
func TestMaterial_ConvenienceMethods(t *testing.T) {
	tests := []struct {
		name     string
		material Material
		method   string
		want     bool
	}{
		{"Empty.IsEmpty", MaterialEmpty, "IsEmpty", true},
		{"Stone.IsEmpty", MaterialStone, "IsEmpty", false},
		{"Stone.IsStatic", MaterialStone, "IsStatic", true},
		{"Sand.IsStatic", MaterialSand, "IsStatic", false},
		{"Stone.IsSolid", MaterialStone, "IsSolid", true},
		{"Water.IsSolid", MaterialWater, "IsSolid", false},
		{"Water.IsLiquid", MaterialWater, "IsLiquid", true},
		{"Stone.IsLiquid", MaterialStone, "IsLiquid", false},
		{"Smoke.IsGas", MaterialSmoke, "IsGas", true},
		{"Water.IsGas", MaterialWater, "IsGas", false},
		{"Empty.IsPenetrable", MaterialEmpty, "IsPenetrable", true},
		{"Water.IsPenetrable", MaterialWater, "IsPenetrable", true},
		{"Stone.IsPenetrable", MaterialStone, "IsPenetrable", false},
		{"Empty.IsFlowable", MaterialEmpty, "IsFlowable", true},
		{"Smoke.IsFlowable", MaterialSmoke, "IsFlowable", true},
		{"Water.IsFlowable", MaterialWater, "IsFlowable", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bool
			switch tt.method {
			case "IsEmpty":
				got = tt.material.IsEmpty()
			case "IsStatic":
				got = tt.material.IsStatic()
			case "IsSolid":
				got = tt.material.IsSolid()
			case "IsLiquid":
				got = tt.material.IsLiquid()
			case "IsGas":
				got = tt.material.IsGas()
			case "IsPenetrable":
				got = tt.material.IsPenetrable()
			case "IsFlowable":
				got = tt.material.IsFlowable()
			}

			if got != tt.want {
				t.Errorf("%s.%s() = %v, want %v", tt.name, tt.method, got, tt.want)
			}
		})
	}
}

// Test String method
func TestMaterial_String(t *testing.T) {
	tests := []struct {
		material Material
		want     string
	}{
		{MaterialEmpty, "Empty"},
		{MaterialStone, "Stone"},
		{MaterialSand, "Sand"},
		{MaterialWater, "Water"},
		{MaterialSmoke, "Smoke"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.material.String(); got != tt.want {
				t.Errorf("Material.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test that data modifications don't affect type or kind
func TestMaterial_DataIsolation(t *testing.T) {
	m := MaterialStone

	// Store original type and check flags
	originalType := m.GetType()
	wasStatic := m.IsStatic()
	wasSolid := m.IsSolid()

	// Modify all data bits
	for i := 0; i < 16; i++ {
		m = m.Set(i)
	}

	// Type and kind should be unchanged
	if m.GetType() != originalType {
		t.Error("Type changed after data modification")
	}
	if m.IsStatic() != wasStatic {
		t.Error("Static flag changed after data modification")
	}
	if m.IsSolid() != wasSolid {
		t.Error("Solid flag changed after data modification")
	}
}

// Benchmark material operations
func BenchmarkMaterial_GetType(b *testing.B) {
	m := MaterialStone
	for i := 0; i < b.N; i++ {
		_ = m.GetType()
	}
}

func BenchmarkMaterial_FilterAny(b *testing.B) {
	m := MaterialStone
	for i := 0; i < b.N; i++ {
		_ = m.FilterAny(SolidMaterials)
	}
}

func BenchmarkMaterial_GetInt(b *testing.B) {
	m := MaterialWater.SetInt(0, 8, 255)
	for i := 0; i < b.N; i++ {
		_ = m.GetInt(0, 8)
	}
}

func BenchmarkMaterial_SetInt(b *testing.B) {
	m := MaterialWater
	for i := 0; i < b.N; i++ {
		m = m.SetInt(0, 8, i&255)
	}
}
