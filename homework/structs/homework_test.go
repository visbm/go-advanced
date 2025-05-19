package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func getBits(data [10]byte, offset, bits int) uint64 {
	var val uint64
	byteOffset := offset / 8
	bitOffset := offset % 8
	totalBits := bitOffset + bits
	byteCount := (totalBits + 7) / 8

	for i := 0; i < byteCount; i++ {
		val |= uint64(data[byteOffset+i]) << (8 * i)
	}

	val >>= bitOffset
	mask := uint64((1 << bits) - 1)
	return val & mask
}

func setBits(data *[10]byte, offset, bits int, value uint64) {
	byteOffset := offset / 8
	bitOffset := offset % 8
	totalBits := bitOffset + bits
	byteCount := (totalBits + 7) / 8

	var val uint64 = 0
	for i := 0; i < byteCount; i++ {
		val |= uint64(data[byteOffset+i]) << (8 * i)
	}

	mask := ((uint64(1) << bits) - 1) << bitOffset
	val = (val &^ mask) | ((value << bitOffset) & mask)

	for i := 0; i < byteCount; i++ {
		data[byteOffset+i] = byte((val >> (8 * i)) & 0xFF)
	}
}

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		nameBytes := [42]byte{}
		for i := 0; i < len(name); i++ {
			nameBytes[i] = name[i]
		}
		person.name = nameBytes
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

const (
	goldOffset = 0
	goldBits   = 31

	manaOffset = 31
	manaBits   = 10

	healthOffset = 41
	healthBits   = 10

	respectOffset = 51
	respectBits   = 4

	strengthOffset = 55
	strengthBits   = 4

	experienceOffset = 59
	experienceBits   = 4

	levelOffset = 63
	levelBits   = 4

	withHouseOffset = 67
	withHouseBits   = 1

	withGunOffset = 68
	withGunBits   = 1

	withFamilyOffset = 69
	withFamilyBits   = 1

	personTypeOffset = 70
	personTypeBits   = 3
)

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(&person.attributes, goldOffset, goldBits, uint64(gold))
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(&person.attributes, manaOffset, manaBits, uint64(mana))
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(&person.attributes, healthOffset, healthBits, uint64(health))
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(&person.attributes, respectOffset, respectBits, uint64(respect))
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(&person.attributes, strengthOffset, strengthBits, uint64(strength))
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(&person.attributes, experienceOffset, experienceBits, uint64(experience))
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(&person.attributes, levelOffset, levelBits, uint64(level))
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(&person.attributes, withHouseOffset, withHouseBits, uint64(1))
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(&person.attributes, withGunOffset, withGunBits, uint64(1))
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(&person.attributes, withFamilyOffset, withFamilyBits, uint64(1))
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(&person.attributes, personTypeOffset, personTypeBits, uint64(personType))
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x, y, z int32 // coordinates
	// attributes — bitовое поле длиной 80 bit (10 байт), содержащее характеристики персонажа:
	// - gold          [0–30]   (31 bits)
	// - mana          [31–40]  (10 bits)
	// - health        [41–50]  (10 bits)
	// - respect       [51–54]  (4 bits)
	// - strength      [55–58]  (4 bits)
	// - experience    [59–62]  (4 bits)
	// - level         [63–66]  (4 bits)
	// - withHouse     [67]     (1 bit)
	// - withGun       [68]     (1 bit)
	// - withFamily    [69]     (1 bit)
	// - personType    [70–72]  (3 bits)
	attributes [10]byte
	name       [42]byte
}

type dtoGamePerson struct {
	Name       string `json:"name"`
	X          int    `json:"x"`
	Y          int    `json:"y"`
	Z          int    `json:"z"`
	Gold       int    `json:"gold"`
	Mana       int    `json:"mana"`
	Health     int    `json:"health"`
	Respect    int    `json:"respect"`
	Strength   int    `json:"strength"`
	Experience int    `json:"experience"`
	Level      int    `json:"level"`
	HasHouse   bool   `json:"has_house"`
	HasGun     bool   `json:"has_gun"`
	HasFamily  bool   `json:"has_family"`
	Type       int    `json:"type"`
}

func (p *GamePerson) MarshalJSON() ([]byte, error) {
	dto := dtoGamePerson{
		Name:       p.Name(),
		X:          p.X(),
		Y:          p.Y(),
		Z:          p.Z(),
		Gold:       p.Gold(),
		Mana:       p.Mana(),
		Health:     p.Health(),
		Respect:    p.Respect(),
		Strength:   p.Strength(),
		Experience: p.Experience(),
		Level:      p.Level(),
		HasHouse:   p.HasHouse(),
		HasGun:     p.HasGun(),
		HasFamily:  p.HasFamilty(),
		Type:       p.Type(),
	}
	return json.Marshal(dto)
}

func NewGamePerson(options ...Option) GamePerson {
	person := &GamePerson{}

	for _, opt := range options {
		opt(person)
	}
	return *person
}

func (p *GamePerson) Name() string {
	return string(p.name[:])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(getBits(p.attributes, goldOffset, goldBits))
}

func (p *GamePerson) Mana() int {
	return int(getBits(p.attributes, manaOffset, manaBits))
}

func (p *GamePerson) Health() int {
	return int(getBits(p.attributes, healthOffset, healthBits))
}

func (p *GamePerson) Respect() int {
	return int(getBits(p.attributes, respectOffset, respectBits))
}

func (p *GamePerson) Strength() int {
	return int(getBits(p.attributes, strengthOffset, strengthBits))
}

func (p *GamePerson) Experience() int {
	return int(getBits(p.attributes, experienceOffset, experienceBits))
}

func (p *GamePerson) Level() int {
	return int(getBits(p.attributes, levelOffset, levelBits))
}

func (p *GamePerson) HasHouse() bool {
	return int(getBits(p.attributes, withHouseOffset, withHouseBits)) == 1
}

func (p *GamePerson) HasGun() bool {
	return int(getBits(p.attributes, withGunOffset, withGunBits)) == 1
}

func (p *GamePerson) HasFamilty() bool {
	return int(getBits(p.attributes, withFamilyOffset, withFamilyBits)) == 1
}

func (p *GamePerson) Type() int {
	return int(getBits(p.attributes, personTypeOffset, personTypeBits))
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))
	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())

	b, _ := json.Marshal(&person)
	fmt.Println(string(b))

	b2, _ := MarshalJSON(getDTO(person))
	fmt.Println(string(b2))
}

func MarshalJSON(v interface{}) ([]byte, error) {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	b := bytes.Buffer{}
	b.Write([]byte("{"))

	if val.Kind() != reflect.Struct {
		return nil, errors.New("not a struct")
	}

	for i := range val.NumField() {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		tag := fieldType.Tag.Get("json")
		if tag == "" {
			continue
		}
		jsonValue := encodeJSONValue(fieldVal)
		if i == val.NumField()-1 {
			b.Write([]byte(`"` + tag + `":` + jsonValue))
			break
		}

		if i == val.NumField() {
			b.Write([]byte(`"` + tag + `":` + jsonValue + `,`))
			break
		}
		b.Write([]byte(`"` + tag + `":` + jsonValue + `,`))

	}

	b.Write([]byte("}"))
	return b.Bytes(), nil

}

func encodeJSONValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return `"` + v.String() + `"`
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case reflect.Slice:
		var parts []string
		for i := 0; i < v.Len(); i++ {
			parts = append(parts, encodeJSONValue(v.Index(i)))
		}
		return "[" + strings.Join(parts, ",") + "]"
	case reflect.Map:
		var parts []string
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			parts = append(parts, fmt.Sprintf(`"%v":%s`, key, encodeJSONValue(val)))
		}
		return "{" + strings.Join(parts, ",") + "}"
	case reflect.Struct:
		j, _ := MarshalJSON(v.Interface())
		return string(j)
	case reflect.Ptr:
		if v.IsNil() {
			return "null"
		}
		return encodeJSONValue(v.Elem())
	default:
		return "null"
	}
}

func getDTO(p GamePerson) dtoGamePerson {
	return dtoGamePerson{
		Name:       p.Name(),
		X:          p.X(),
		Y:          p.Y(),
		Z:          p.Z(),
		Gold:       p.Gold(),
		Mana:       p.Mana(),
		Health:     p.Health(),
		Respect:    p.Respect(),
		Strength:   p.Strength(),
		Experience: p.Experience(),
		Level:      p.Level(),
		HasHouse:   p.HasHouse(),
		HasGun:     p.HasGun(),
		HasFamily:  p.HasFamilty(),
		Type:       p.Type(),
	}
}
