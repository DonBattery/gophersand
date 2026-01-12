// material_test.go
package game

import (
	"reflect"
	"testing"
)

func TestMaterialConstantsMatchKindValues(t *testing.T) {
	mats := []Material{
		MaterialEmpty,
		MaterialStone,
		MaterialSand,
		MaterialWater,
		MaterialSeed,
		MaterialAnt,
		MaterialWasp,
		MaterialAcid,
		MaterialFire,
		MaterialIce,
		MaterialSmoke,
		MaterialSteam,
		MaterialRoot,
		MaterialPlant,
		MaterialFlower,
		MaterialAntHill,
	}
	kinds := []MaterialKind{
		MaterialKindEmpty,
		MaterialKindStone,
		MaterialKindSand,
		MaterialKindWater,
		MaterialKindSeed,
		MaterialKindAnt,
		MaterialKindWasp,
		MaterialKindAcid,
		MaterialKindFire,
		MaterialKindIce,
		MaterialKindSmoke,
		MaterialKindSteam,
		MaterialKindRoot,
		MaterialKindPlant,
		MaterialKindFlower,
		MaterialKindAntHill,
	}

	if len(mats) != 16 || len(kinds) != 16 {
		t.Fatalf("expected 16 materials and 16 kinds")
	}

	for i := 0; i < 16; i++ {
		if mats[i] != Material(i) {
			t.Fatalf("material constant index %d expected %d got %d", i, i, mats[i])
		}
		if kinds[i] != MaterialKind(i) {
			t.Fatalf("kind constant index %d expected %d got %d", i, i, kinds[i])
		}

		// Base material's kind should equal its ordinal.
		if got := mats[i].GetKind(); got != kinds[i] {
			t.Fatalf("material %d GetKind expected %d got %d", i, kinds[i], got)
		}
		if !mats[i].IsKind(kinds[i]) {
			t.Fatalf("material %d expected IsKind(%d)=true", i, kinds[i])
		}
	}
}

func TestNewMaterialKindSetAndMaterialKindIsIn(t *testing.T) {
	set := NewMaterialKindSet(MaterialKindStone, MaterialKindWater, MaterialKindStone /* duplicate */)

	if !MaterialKindStone.IsIn(set) {
		t.Fatalf("expected Stone in set")
	}
	if !MaterialKindWater.IsIn(set) {
		t.Fatalf("expected Water in set")
	}
	if MaterialKindSand.IsIn(set) {
		t.Fatalf("did not expect Sand in set")
	}

	// Exhaustive: each kind should round-trip with a singleton set.
	for k := MaterialKind(0); k < 16; k++ {
		s := NewMaterialKindSet(k)
		for kk := MaterialKind(0); kk < 16; kk++ {
			want := kk == k
			got := kk.IsIn(s)
			if got != want {
				t.Fatalf("k=%d kk=%d expected IsIn=%v got %v", k, kk, want, got)
			}
		}
	}
}

func TestMaterialIsInUsesKind(t *testing.T) {
	set := NewMaterialKindSet(MaterialKindAnt, MaterialKindWasp)

	mAnt := MaterialAnt.WithLife(3).WithStatus(MaterialStatusFrozen).WithFaceLeft(true)
	if !mAnt.IsIn(set) {
		t.Fatalf("expected Ant material to be in set")
	}

	mOther := MaterialSand.WithLife(1).WithStatus(MaterialStatusNormal)
	if mOther.IsIn(set) {
		t.Fatalf("did not expect Sand material to be in set")
	}
}

func TestGetKindIgnoresHigherBits(t *testing.T) {
	// Force various higher bits on; kind must still be low 4 bits.
	for k := Material(0); k < 16; k++ {
		m := Material(0xFFF0) | k
		if got := m.GetKind(); got != MaterialKind(k) {
			t.Fatalf("k=%d: expected kind %d got %d (m=0x%04x)", k, k, got, uint16(m))
		}
		if !m.IsKind(MaterialKind(k)) {
			t.Fatalf("k=%d: expected IsKind true (m=0x%04x)", k, uint16(m))
		}
	}
}

func TestLifeGetSetMaskingAndPreservation(t *testing.T) {
	orig := MaterialKindWasp
	m0 := Material(orig).
		WithStatus(MaterialStatusAcidic).
		WithFaceLeft(true).
		WithFaceUp(true).
		WithWaspHasWater(true).
		WithWaspHasAnt(true)

	// Masking: life &= 3
	m1 := m0.WithLife(7)
	if got := m1.GetLife(); got != 3 {
		t.Fatalf("expected life=3 after masking, got %d", got)
	}

	// Preservation: kind/status/state bits must remain unchanged.
	if m1.GetKind() != MaterialKindWasp {
		t.Fatalf("kind changed after WithLife")
	}
	if m1.GetStatus() != MaterialStatusAcidic {
		t.Fatalf("status changed after WithLife")
	}
	if !m1.GetFaceLeft() || !m1.GetFaceUp() {
		t.Fatalf("face bits changed after WithLife")
	}
	if !m1.GetWaspHasWater() || !m1.GetWaspHasAnt() {
		t.Fatalf("wasp bits changed after WithLife")
	}

	// Round-trip over all life values (including out of range).
	for life := uint8(0); life < 8; life++ {
		m := MaterialSand.WithLife(life)
		want := life & 3
		if got := m.GetLife(); got != want {
			t.Fatalf("life=%d: expected %d got %d", life, want, got)
		}
	}
}

func TestStatusGetSetMaskingAndPreservation(t *testing.T) {
	m0 := MaterialKindAntHill
	m := Material(m0).
		WithLife(2).
		WithFaceLeft(true).
		WithFaceUp(false).
		WithIsPenetrable(true).
		WithCanBloom(true) // uses a different bit than IsPenetrable

	// Masking: status &= 3
	m1 := m.WithStatus(255)
	if got := m1.GetStatus(); got != 3 {
		t.Fatalf("expected status=3 after masking, got %d", got)
	}

	// Preservation checks.
	if m1.GetKind() != MaterialKindAntHill {
		t.Fatalf("kind changed after WithStatus")
	}
	if m1.GetLife() != 2 {
		t.Fatalf("life changed after WithStatus")
	}
	if !m1.GetFaceLeft() || m1.GetFaceUp() {
		t.Fatalf("face bits changed after WithStatus")
	}
	if !m1.GetIsPenetrable() {
		t.Fatalf("penetrable bit changed after WithStatus")
	}
	if !m1.GetCanBloom() {
		t.Fatalf("canBloom bit changed after WithStatus")
	}

	// Round-trip over all status values (including out of range).
	for st := uint8(0); st < 8; st++ {
		m := MaterialStone.WithStatus(st)
		want := st & 3
		if got := m.GetStatus(); got != want {
			t.Fatalf("status=%d: expected %d got %d", st, want, got)
		}
	}
}

func TestGetColorMatchesMaterialColorsIndexing(t *testing.T) {
	// This test is intentionally "indexing only": it verifies GetColor returns
	// exactly MaterialColors[k*16 + s*4 + l] for multiple combinations.
	//
	// If MaterialColors is not at least 256 entries, GetColor would panic; treat
	// that as a test failure with a clearer message.
	if len(MaterialColors) < 16*16 {
		t.Fatalf("MaterialColors too short: got %d need at least %d", len(MaterialColors), 16*16)
	}

	for k := MaterialKind(0); k < 16; k++ {
		for st := uint8(0); st < 4; st++ {
			for life := uint8(0); life < 4; life++ {
				m := Material(k).WithStatus(st).WithLife(life)
				idx := int(k)*16 + int(st)*4 + int(life)

				got := m.GetColor()
				want := MaterialColors[idx]

				if !reflect.DeepEqual(got, want) {
					t.Fatalf("k=%d st=%d life=%d idx=%d: GetColor mismatch", k, st, life, idx)
				}
			}
		}
	}
}

func TestFaceLeftAndFaceUpBits(t *testing.T) {
	m := MaterialSeed

	if m.GetFaceLeft() || m.GetFaceUp() {
		t.Fatalf("expected face bits off by default for base material")
	}

	m = m.WithFaceLeft(true)
	if !m.GetFaceLeft() || m.GetFaceUp() {
		t.Fatalf("expected FaceLeft on, FaceUp unchanged/off")
	}

	m = m.WithFaceUp(true)
	if !m.GetFaceLeft() || !m.GetFaceUp() {
		t.Fatalf("expected both FaceLeft and FaceUp on")
	}

	m = m.WithFaceLeft(false)
	if m.GetFaceLeft() || !m.GetFaceUp() {
		t.Fatalf("expected FaceLeft off, FaceUp still on")
	}

	m = m.WithFaceUp(false)
	if m.GetFaceLeft() || m.GetFaceUp() {
		t.Fatalf("expected both face bits off")
	}

	// Sanity on bit locations (internal constants).
	if stateFaceLeft != (Material(1) << 8) {
		t.Fatalf("stateFaceLeft expected bit 8, got 0x%04x", uint16(stateFaceLeft))
	}
	if stateFaceUp != (Material(1) << 9) {
		t.Fatalf("stateFaceUp expected bit 9, got 0x%04x", uint16(stateFaceUp))
	}
}

func TestAliasBitLayoutAndIndependence(t *testing.T) {
	// Verify intended aliasing.
	if isPenetrableBit != stateFlagA {
		t.Fatalf("isPenetrableBit expected to alias stateFlagA")
	}
	if isTopPetalBit != stateFlagA {
		t.Fatalf("isTopPetalBit expected to alias stateFlagA")
	}
	if canBloomBit != stateFlagD {
		t.Fatalf("canBloomBit expected to alias stateFlagD")
	}
	if canBloomBit == isPenetrableBit {
		t.Fatalf("canBloomBit must not share bit with isPenetrableBit")
	}

	// For Plant: penetrable uses FlagA; canBloom uses FlagD -> must be independent.
	m := MaterialPlant.WithIsPenetrable(true).WithCanBloom(true)
	if !m.GetIsPenetrable() || !m.GetCanBloom() {
		t.Fatalf("expected both IsPenetrable and CanBloom true")
	}

	m2 := m.WithIsPenetrable(false)
	if m2.GetIsPenetrable() {
		t.Fatalf("expected IsPenetrable false")
	}
	if !m2.GetCanBloom() {
		t.Fatalf("expected CanBloom unchanged true")
	}

	m3 := m.WithCanBloom(false)
	if !m3.GetIsPenetrable() {
		t.Fatalf("expected IsPenetrable unchanged true")
	}
	if m3.GetCanBloom() {
		t.Fatalf("expected CanBloom false")
	}
}

func TestWaspCollectionBits(t *testing.T) {
	// Verify aliasing to canonical flags.
	if waspHasWaterBit != stateFlagB {
		t.Fatalf("waspHasWaterBit expected to alias stateFlagB")
	}
	if waspHasAntBit != stateFlagC {
		t.Fatalf("waspHasAntBit expected to alias stateFlagC")
	}
	if waspHasWaterBit == waspHasAntBit {
		t.Fatalf("waspHasWaterBit and waspHasAntBit must be distinct")
	}

	m := MaterialWasp
	if m.GetWaspHasWater() || m.GetWaspHasAnt() {
		t.Fatalf("expected wasp flags off by default")
	}

	m = m.WithWaspHasWater(true)
	if !m.GetWaspHasWater() || m.GetWaspHasAnt() {
		t.Fatalf("expected WaspHasWater on and WaspHasAnt unchanged/off")
	}

	m = m.WithWaspHasAnt(true)
	if !m.GetWaspHasWater() || !m.GetWaspHasAnt() {
		t.Fatalf("expected both wasp flags on")
	}

	m = m.WithWaspHasWater(false)
	if m.GetWaspHasWater() || !m.GetWaspHasAnt() {
		t.Fatalf("expected WaspHasWater off, WaspHasAnt still on")
	}

	m = m.WithWaspHasAnt(false)
	if m.GetWaspHasWater() || m.GetWaspHasAnt() {
		t.Fatalf("expected both wasp flags off")
	}
}

func TestPredefinedKindSetsContainExpectedMembership(t *testing.T) {
	type setCase struct {
		name     string
		set      MaterialKindSet
		in       []MaterialKind
		notInAny []MaterialKind
	}

	cases := []setCase{
		{
			name: "RootGrowableKinds",
			set:  RootGrowableKinds,
			in: []MaterialKind{
				MaterialKindAntHill,
				MaterialKindEmpty,
				MaterialKindSteam,
				MaterialKindSmoke,
				MaterialKindPlant,
			},
			notInAny: []MaterialKind{MaterialKindWater, MaterialKindStone},
		},
		{
			name: "PlantGrowableKinds",
			set:  PlantGrowableKinds,
			in: []MaterialKind{
				MaterialKindEmpty,
				MaterialKindAntHill,
				MaterialKindSteam,
				MaterialKindSmoke,
				MaterialKindRoot,
			},
			notInAny: []MaterialKind{MaterialKindWater, MaterialKindFire},
		},
		{
			name: "PlantSupporterKinds",
			set:  PlantSupporterKinds,
			in: []MaterialKind{
				MaterialKindStone,
				MaterialKindSand,
				MaterialKindSeed,
				MaterialKindRoot,
				MaterialKindPlant,
				MaterialKindFlower,
			},
			notInAny: []MaterialKind{MaterialKindEmpty, MaterialKindWater},
		},
		{
			name: "AntSupporterKinds",
			set:  AntSupporterKinds,
			in: []MaterialKind{
				MaterialKindStone,
				MaterialKindSand,
				MaterialKindSeed,
				MaterialKindAntHill,
				MaterialKindRoot,
				MaterialKindPlant,
				MaterialKindFlower,
				MaterialKindAnt,
			},
			notInAny: []MaterialKind{MaterialKindEmpty, MaterialKindWater},
		},
		{
			name: "AntAliveKinds",
			set:  AntAliveKinds,
			in: []MaterialKind{
				MaterialKindSeed,
				MaterialKindRoot,
				MaterialKindPlant,
				MaterialKindFlower,
				MaterialKindAntHill,
			},
			notInAny: []MaterialKind{MaterialKindStone, MaterialKindEmpty},
		},
		{
			name: "AntEggLayableKinds",
			set:  AntEggLayableKinds,
			in: []MaterialKind{
				MaterialKindEmpty,
				MaterialKindAntHill,
				MaterialKindSteam,
				MaterialKindSmoke,
				MaterialKindRoot,
				MaterialKindPlant,
				MaterialKindFlower,
			},
			notInAny: []MaterialKind{MaterialKindStone, MaterialKindWater},
		},
		{
			name: "AntFallableKinds",
			set:  AntFallableKinds,
			in: []MaterialKind{
				MaterialKindEmpty,
				MaterialKindSteam,
				MaterialKindSmoke,
				MaterialKindFire,
				MaterialKindWater,
				MaterialKindAcid,
			},
			notInAny: []MaterialKind{MaterialKindStone, MaterialKindSand},
		},
		{
			name: "WaspEggStickyKinds",
			set:  WaspEggStickyKinds,
			in: []MaterialKind{
				MaterialKindStone,
				MaterialKindSand,
				MaterialKindSeed,
				MaterialKindRoot,
				MaterialKindPlant,
				MaterialKindFlower,
			},
			notInAny: []MaterialKind{MaterialKindEmpty, MaterialKindWater},
		},
		{
			name: "WaspEggLayableKinds",
			set:  WaspEggLayableKinds,
			in: []MaterialKind{
				MaterialKindEmpty,
				MaterialKindAntHill,
				MaterialKindSteam,
				MaterialKindSmoke,
				MaterialKindPlant,
			},
			notInAny: []MaterialKind{MaterialKindStone, MaterialKindWater},
		},
		{
			name: "FreezableKinds",
			set:  FreezableKinds,
			in: []MaterialKind{
				MaterialKindEmpty,
				MaterialKindSteam,
				MaterialKindSmoke,
				MaterialKindPlant,
			},
			notInAny: []MaterialKind{MaterialKindStone, MaterialKindWater},
		},
		{
			name: "PlantFoodKinds",
			set:  PlantFoodKinds,
			in: []MaterialKind{
				MaterialKindWater,
				MaterialKindSand,
			},
			notInAny: []MaterialKind{MaterialKindStone, MaterialKindEmpty},
		},
		{
			name: "NonCondensableKinds",
			set:  NonCondensableKinds,
			in: []MaterialKind{
				MaterialKindEmpty,
				MaterialKindSteam,
				MaterialKindSmoke,
				MaterialKindWater,
				MaterialKindAcid,
				MaterialKindFire,
			},
			notInAny: []MaterialKind{MaterialKindStone, MaterialKindSand},
		},
	}

	for _, tc := range cases {
		for _, k := range tc.in {
			if !k.IsIn(tc.set) {
				t.Fatalf("%s: expected %d to be in set", tc.name, k)
			}
		}
		for _, k := range tc.notInAny {
			if k.IsIn(tc.set) {
				t.Fatalf("%s: did not expect %d to be in set", tc.name, k)
			}
		}
	}
}
