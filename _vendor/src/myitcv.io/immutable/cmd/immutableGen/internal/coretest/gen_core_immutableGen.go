// Code generated by immutableGen. DO NOT EDIT.

// My favourite license

package coretest

//go:generate echo "hello world"
//immutableVet:skipFile

import (
	"myitcv.io/immutable"

	"myitcv.io/immutable/cmd/immutableGen/internal/coretest/pkga"
	"myitcv.io/immutable/cmd/immutableGen/internal/coretest/pkgb"
	"time"
)

// a comment about MyMap
//
// MyMap is an immutable type and has the following template:
//
// 	map[string]int
//
type MyMap struct {
	theMap  map[string]int
	mutable bool
	__tmpl  *_Imm_MyMap
}

var _ immutable.Immutable = new(MyMap)
var _ = new(MyMap).__tmpl

func NewMyMap(inits ...func(m *MyMap)) *MyMap {
	res := NewMyMapCap(0)
	if len(inits) == 0 {
		return res
	}

	return res.WithMutable(func(m *MyMap) {
		for _, i := range inits {
			i(m)
		}
	})
}

func NewMyMapCap(l int) *MyMap {
	return &MyMap{
		theMap: make(map[string]int, l),
	}
}

func (m *MyMap) Mutable() bool {
	return m.mutable
}

func (m *MyMap) Len() int {
	if m == nil {
		return 0
	}

	return len(m.theMap)
}

func (m *MyMap) Get(k string) (int, bool) {
	v, ok := m.theMap[k]
	return v, ok
}

func (m *MyMap) AsMutable() *MyMap {
	if m == nil {
		return nil
	}

	if m.Mutable() {
		return m
	}

	res := m.dup()
	res.mutable = true

	return res
}

func (m *MyMap) dup() *MyMap {
	resMap := make(map[string]int, len(m.theMap))

	for k := range m.theMap {
		resMap[k] = m.theMap[k]
	}

	res := &MyMap{
		theMap: resMap,
	}

	return res
}

func (m *MyMap) AsImmutable(v *MyMap) *MyMap {
	if m == nil {
		return nil
	}

	if v == m {
		return m
	}

	m.mutable = false
	return m
}

func (m *MyMap) Range() map[string]int {
	if m == nil {
		return nil
	}

	return m.theMap
}

func (mr *MyMap) WithMutable(f func(m *MyMap)) *MyMap {
	res := mr.AsMutable()
	f(res)
	res = res.AsImmutable(mr)

	return res
}

func (mr *MyMap) WithImmutable(f func(m *MyMap)) *MyMap {
	prev := mr.mutable
	mr.mutable = false
	f(mr)
	mr.mutable = prev

	return mr
}

func (m *MyMap) Set(k string, v int) *MyMap {
	if m.mutable {
		m.theMap[k] = v
		return m
	}

	res := m.dup()
	res.theMap[k] = v

	return res
}

func (m *MyMap) Del(k string) *MyMap {
	if _, ok := m.theMap[k]; !ok {
		return m
	}

	if m.mutable {
		delete(m.theMap, k)
		return m
	}

	res := m.dup()
	delete(res.theMap, k)

	return res
}
func (s *MyMap) IsDeeplyNonMutable(seen map[interface{}]bool) bool {
	if s == nil {
		return true
	}

	if s.Mutable() {
		return false
	}
	return true
}

//
// AM is an immutable type and has the following template:
//
// 	map[*A]*A
//
type AM struct {
	theMap  map[*A]*A
	mutable bool
	__tmpl  *_Imm_AM
}

var _ immutable.Immutable = new(AM)
var _ = new(AM).__tmpl

func NewAM(inits ...func(m *AM)) *AM {
	res := NewAMCap(0)
	if len(inits) == 0 {
		return res
	}

	return res.WithMutable(func(m *AM) {
		for _, i := range inits {
			i(m)
		}
	})
}

func NewAMCap(l int) *AM {
	return &AM{
		theMap: make(map[*A]*A, l),
	}
}

func (m *AM) Mutable() bool {
	return m.mutable
}

func (m *AM) Len() int {
	if m == nil {
		return 0
	}

	return len(m.theMap)
}

func (m *AM) Get(k *A) (*A, bool) {
	v, ok := m.theMap[k]
	return v, ok
}

func (m *AM) AsMutable() *AM {
	if m == nil {
		return nil
	}

	if m.Mutable() {
		return m
	}

	res := m.dup()
	res.mutable = true

	return res
}

func (m *AM) dup() *AM {
	resMap := make(map[*A]*A, len(m.theMap))

	for k := range m.theMap {
		resMap[k] = m.theMap[k]
	}

	res := &AM{
		theMap: resMap,
	}

	return res
}

func (m *AM) AsImmutable(v *AM) *AM {
	if m == nil {
		return nil
	}

	if v == m {
		return m
	}

	m.mutable = false
	return m
}

func (m *AM) Range() map[*A]*A {
	if m == nil {
		return nil
	}

	return m.theMap
}

func (mr *AM) WithMutable(f func(a *AM)) *AM {
	res := mr.AsMutable()
	f(res)
	res = res.AsImmutable(mr)

	return res
}

func (mr *AM) WithImmutable(f func(a *AM)) *AM {
	prev := mr.mutable
	mr.mutable = false
	f(mr)
	mr.mutable = prev

	return mr
}

func (m *AM) Set(k *A, v *A) *AM {
	if m.mutable {
		m.theMap[k] = v
		return m
	}

	res := m.dup()
	res.theMap[k] = v

	return res
}

func (m *AM) Del(k *A) *AM {
	if _, ok := m.theMap[k]; !ok {
		return m
	}

	if m.mutable {
		delete(m.theMap, k)
		return m
	}

	res := m.dup()
	delete(res.theMap, k)

	return res
}
func (s *AM) IsDeeplyNonMutable(seen map[interface{}]bool) bool {
	if s == nil {
		return true
	}

	if s.Mutable() {
		return false
	}
	if s.Len() == 0 {
		return true
	}

	if seen == nil {
		return s.IsDeeplyNonMutable(make(map[interface{}]bool))
	}

	if seen[s] {
		return true
	}

	seen[s] = true

	for k, v := range s.theMap {
		if k != nil && !k.IsDeeplyNonMutable(seen) {
			return false
		}
		if v != nil && !v.IsDeeplyNonMutable(seen) {
			return false
		}
	}
	return true
}

// a comment about Slice
//
// MySlice is an immutable type and has the following template:
//
// 	[]string
//
type MySlice struct {
	theSlice []string
	mutable  bool
	__tmpl   *_Imm_MySlice
}

var _ immutable.Immutable = new(MySlice)
var _ = new(MySlice).__tmpl

func NewMySlice(s ...string) *MySlice {
	c := make([]string, len(s))
	copy(c, s)

	return &MySlice{
		theSlice: c,
	}
}

func NewMySliceLen(l int) *MySlice {
	c := make([]string, l)

	return &MySlice{
		theSlice: c,
	}
}

func (m *MySlice) Mutable() bool {
	return m.mutable
}

func (m *MySlice) Len() int {
	if m == nil {
		return 0
	}

	return len(m.theSlice)
}

func (m *MySlice) Get(i int) string {
	return m.theSlice[i]
}

func (m *MySlice) AsMutable() *MySlice {
	if m == nil {
		return nil
	}

	if m.Mutable() {
		return m
	}

	res := m.dup()
	res.mutable = true

	return res
}

func (m *MySlice) dup() *MySlice {
	resSlice := make([]string, len(m.theSlice))

	for i := range m.theSlice {
		resSlice[i] = m.theSlice[i]
	}

	res := &MySlice{
		theSlice: resSlice,
	}

	return res
}

func (m *MySlice) AsImmutable(v *MySlice) *MySlice {
	if m == nil {
		return nil
	}

	if v == m {
		return m
	}

	m.mutable = false
	return m
}

func (m *MySlice) Range() []string {
	if m == nil {
		return nil
	}

	return m.theSlice
}

func (m *MySlice) WithMutable(f func(mi *MySlice)) *MySlice {
	res := m.AsMutable()
	f(res)
	res = res.AsImmutable(m)

	return res
}

func (m *MySlice) WithImmutable(f func(mi *MySlice)) *MySlice {
	prev := m.mutable
	m.mutable = false
	f(m)
	m.mutable = prev

	return m
}

func (m *MySlice) Set(i int, v string) *MySlice {
	if m.mutable {
		m.theSlice[i] = v
		return m
	}

	res := m.dup()
	res.theSlice[i] = v

	return res
}

func (m *MySlice) Append(v ...string) *MySlice {
	if m.mutable {
		m.theSlice = append(m.theSlice, v...)
		return m
	}

	res := m.dup()
	res.theSlice = append(res.theSlice, v...)

	return res
}
func (s *MySlice) IsDeeplyNonMutable(seen map[interface{}]bool) bool {
	if s == nil {
		return true
	}

	if s.Mutable() {
		return false
	}
	return true
}

//
// AS is an immutable type and has the following template:
//
// 	[]*A
//
type AS struct {
	theSlice []*A
	mutable  bool
	__tmpl   *_Imm_AS
}

var _ immutable.Immutable = new(AS)
var _ = new(AS).__tmpl

func NewAS(s ...*A) *AS {
	c := make([]*A, len(s))
	copy(c, s)

	return &AS{
		theSlice: c,
	}
}

func NewASLen(l int) *AS {
	c := make([]*A, l)

	return &AS{
		theSlice: c,
	}
}

func (m *AS) Mutable() bool {
	return m.mutable
}

func (m *AS) Len() int {
	if m == nil {
		return 0
	}

	return len(m.theSlice)
}

func (m *AS) Get(i int) *A {
	return m.theSlice[i]
}

func (m *AS) AsMutable() *AS {
	if m == nil {
		return nil
	}

	if m.Mutable() {
		return m
	}

	res := m.dup()
	res.mutable = true

	return res
}

func (m *AS) dup() *AS {
	resSlice := make([]*A, len(m.theSlice))

	for i := range m.theSlice {
		resSlice[i] = m.theSlice[i]
	}

	res := &AS{
		theSlice: resSlice,
	}

	return res
}

func (m *AS) AsImmutable(v *AS) *AS {
	if m == nil {
		return nil
	}

	if v == m {
		return m
	}

	m.mutable = false
	return m
}

func (m *AS) Range() []*A {
	if m == nil {
		return nil
	}

	return m.theSlice
}

func (m *AS) WithMutable(f func(mi *AS)) *AS {
	res := m.AsMutable()
	f(res)
	res = res.AsImmutable(m)

	return res
}

func (m *AS) WithImmutable(f func(mi *AS)) *AS {
	prev := m.mutable
	m.mutable = false
	f(m)
	m.mutable = prev

	return m
}

func (m *AS) Set(i int, v *A) *AS {
	if m.mutable {
		m.theSlice[i] = v
		return m
	}

	res := m.dup()
	res.theSlice[i] = v

	return res
}

func (m *AS) Append(v ...*A) *AS {
	if m.mutable {
		m.theSlice = append(m.theSlice, v...)
		return m
	}

	res := m.dup()
	res.theSlice = append(res.theSlice, v...)

	return res
}
func (s *AS) IsDeeplyNonMutable(seen map[interface{}]bool) bool {
	if s == nil {
		return true
	}

	if s.Mutable() {
		return false
	}
	if s.Len() == 0 {
		return true
	}

	if seen == nil {
		return s.IsDeeplyNonMutable(make(map[interface{}]bool))
	}

	if seen[s] {
		return true
	}

	seen[s] = true

	for _, v := range s.theSlice {
		if v != nil && !v.IsDeeplyNonMutable(seen) {
			return false
		}
	}
	return true
}

// a comment about myStruct
//
// MyStruct is an immutable type and has the following template:
//
// 	struct {
// 		Key	MyStructKey
//
// 		Name, surname	string
// 		age		int
//
// 		string
//
// 		fieldWithoutTag	bool
// 	}
//
type MyStruct struct {
	field_Key             MyStructKey
	field_Name            string `tag:"value"`
	field_surname         string `tag:"value"`
	field_age             int    `tag:"age"`
	anonfield_string      string
	field_fieldWithoutTag bool

	mutable bool
	__tmpl  *_Imm_MyStruct
}

var _ immutable.Immutable = new(MyStruct)
var _ = new(MyStruct).__tmpl

func (s *MyStruct) AsMutable() *MyStruct {
	if s.Mutable() {
		return s
	}

	res := *s
	res.field_Key.Version++
	res.mutable = true
	return &res
}

func (s *MyStruct) AsImmutable(v *MyStruct) *MyStruct {
	if s == nil {
		return nil
	}

	if s == v {
		return s
	}

	s.mutable = false
	return s
}

func (s *MyStruct) Mutable() bool {
	return s.mutable
}

func (s *MyStruct) WithMutable(f func(si *MyStruct)) *MyStruct {
	res := s.AsMutable()
	f(res)
	res = res.AsImmutable(s)

	return res
}

func (s *MyStruct) WithImmutable(f func(si *MyStruct)) *MyStruct {
	prev := s.mutable
	s.mutable = false
	f(s)
	s.mutable = prev

	return s
}

func (s *MyStruct) IsDeeplyNonMutable(seen map[interface{}]bool) bool {
	if s == nil {
		return true
	}

	if s.Mutable() {
		return false
	}

	if seen == nil {
		return s.IsDeeplyNonMutable(make(map[interface{}]bool))
	}

	if seen[s] {
		return true
	}

	seen[s] = true
	return true
}
func (s *MyStruct) Key() MyStructKey {
	return s.field_Key
}

// SetKey is the setter for Key()
func (s *MyStruct) SetKey(n MyStructKey) *MyStruct {
	if s.mutable {
		s.field_Key = n
		return s
	}

	res := *s
	res.field_Key.Version++
	res.field_Key = n
	return &res
}

// my field comment
//somethingspecial
/*

	Heelo

*/
func (s *MyStruct) Name() string {
	return s.field_Name
}

// SetName is the setter for Name()
func (s *MyStruct) SetName(n string) *MyStruct {
	if s.mutable {
		s.field_Name = n
		return s
	}

	res := *s
	res.field_Key.Version++
	res.field_Name = n
	return &res
}
func (s *MyStruct) age() int {
	return s.field_age
}

// setAge is the setter for Age()
func (s *MyStruct) setAge(n int) *MyStruct {
	if s.mutable {
		s.field_age = n
		return s
	}

	res := *s
	res.field_Key.Version++
	res.field_age = n
	return &res
}
func (s *MyStruct) fieldWithoutTag() bool {
	return s.field_fieldWithoutTag
}

// setFieldWithoutTag is the setter for FieldWithoutTag()
func (s *MyStruct) setFieldWithoutTag(n bool) *MyStruct {
	if s.mutable {
		s.field_fieldWithoutTag = n
		return s
	}

	res := *s
	res.field_Key.Version++
	res.field_fieldWithoutTag = n
	return &res
}
func (s *MyStruct) string() string {
	return s.anonfield_string
}

// setString is the setter for String()
func (s *MyStruct) setString(n string) *MyStruct {
	if s.mutable {
		s.anonfield_string = n
		return s
	}

	res := *s
	res.field_Key.Version++
	res.anonfield_string = n
	return &res
}

// my field comment
//somethingspecial
/*

	Heelo

*/
func (s *MyStruct) surname() string {
	return s.field_surname
}

// setSurname is the setter for Surname()
func (s *MyStruct) setSurname(n string) *MyStruct {
	if s.mutable {
		s.field_surname = n
		return s
	}

	res := *s
	res.field_Key.Version++
	res.field_surname = n
	return &res
}

//
// A is an immutable type and has the following template:
//
// 	struct {
// 		Name	string
// 		A	*A
//
// 		Blah
// 	}
//
type A struct {
	field_Name     string
	field_A        *A
	anonfield_Blah Blah

	mutable bool
	__tmpl  *_Imm_A
}

var _ immutable.Immutable = new(A)
var _ = new(A).__tmpl

func (s *A) AsMutable() *A {
	if s.Mutable() {
		return s
	}

	res := *s
	res.mutable = true
	return &res
}

func (s *A) AsImmutable(v *A) *A {
	if s == nil {
		return nil
	}

	if s == v {
		return s
	}

	s.mutable = false
	return s
}

func (s *A) Mutable() bool {
	return s.mutable
}

func (s *A) WithMutable(f func(si *A)) *A {
	res := s.AsMutable()
	f(res)
	res = res.AsImmutable(s)

	return res
}

func (s *A) WithImmutable(f func(si *A)) *A {
	prev := s.mutable
	s.mutable = false
	f(s)
	s.mutable = prev

	return s
}

func (s *A) IsDeeplyNonMutable(seen map[interface{}]bool) bool {
	if s == nil {
		return true
	}

	if s.Mutable() {
		return false
	}

	if seen == nil {
		return s.IsDeeplyNonMutable(make(map[interface{}]bool))
	}

	if seen[s] {
		return true
	}

	seen[s] = true
	{
		v := s.field_A

		if v != nil && !v.IsDeeplyNonMutable(seen) {
			return false
		}
	}
	{
		v := s.anonfield_Blah

		if v != nil && !v.IsDeeplyNonMutable(seen) {
			return false
		}
	}
	return true
}
func (s *A) A() *A {
	return s.field_A
}

// SetA is the setter for A()
func (s *A) SetA(n *A) *A {
	if s.mutable {
		s.field_A = n
		return s
	}

	res := *s
	res.field_A = n
	return &res
}
func (s *A) Blah() Blah {
	return s.anonfield_Blah
}

// SetBlah is the setter for Blah()
func (s *A) SetBlah(n Blah) *A {
	if s.mutable {
		s.anonfield_Blah = n
		return s
	}

	res := *s
	res.anonfield_Blah = n
	return &res
}
func (s *A) Name() string {
	return s.field_Name
}

// SetName is the setter for Name()
func (s *A) SetName(n string) *A {
	if s.mutable {
		s.field_Name = n
		return s
	}

	res := *s
	res.field_Name = n
	return &res
}

//
// BlahUse is an immutable type and has the following template:
//
// 	struct {
// 		Blah
// 	}
//
type BlahUse struct {
	anonfield_Blah Blah

	mutable bool
	__tmpl  *_Imm_BlahUse
}

var _ immutable.Immutable = new(BlahUse)
var _ = new(BlahUse).__tmpl

func (s *BlahUse) AsMutable() *BlahUse {
	if s.Mutable() {
		return s
	}

	res := *s
	res.mutable = true
	return &res
}

func (s *BlahUse) AsImmutable(v *BlahUse) *BlahUse {
	if s == nil {
		return nil
	}

	if s == v {
		return s
	}

	s.mutable = false
	return s
}

func (s *BlahUse) Mutable() bool {
	return s.mutable
}

func (s *BlahUse) WithMutable(f func(si *BlahUse)) *BlahUse {
	res := s.AsMutable()
	f(res)
	res = res.AsImmutable(s)

	return res
}

func (s *BlahUse) WithImmutable(f func(si *BlahUse)) *BlahUse {
	prev := s.mutable
	s.mutable = false
	f(s)
	s.mutable = prev

	return s
}

func (s *BlahUse) IsDeeplyNonMutable(seen map[interface{}]bool) bool {
	if s == nil {
		return true
	}

	if s.Mutable() {
		return false
	}

	if seen == nil {
		return s.IsDeeplyNonMutable(make(map[interface{}]bool))
	}

	if seen[s] {
		return true
	}

	seen[s] = true
	{
		v := s.anonfield_Blah

		if v != nil && !v.IsDeeplyNonMutable(seen) {
			return false
		}
	}
	return true
}
func (s *BlahUse) Blah() Blah {
	return s.anonfield_Blah
}

// SetBlah is the setter for Blah()
func (s *BlahUse) SetBlah(n Blah) *BlahUse {
	if s.mutable {
		s.anonfield_Blah = n
		return s
	}

	res := *s
	res.anonfield_Blah = n
	return &res
}

//
// Clash1 is an immutable type and has the following template:
//
// 	struct {
// 		Clash		string
// 		NoClash1	string
// 	}
//
type Clash1 struct {
	field_Clash    string
	field_NoClash1 string

	mutable bool
	__tmpl  *_Imm_Clash1
}

var _ immutable.Immutable = new(Clash1)
var _ = new(Clash1).__tmpl

func (s *Clash1) AsMutable() *Clash1 {
	if s.Mutable() {
		return s
	}

	res := *s
	res.mutable = true
	return &res
}

func (s *Clash1) AsImmutable(v *Clash1) *Clash1 {
	if s == nil {
		return nil
	}

	if s == v {
		return s
	}

	s.mutable = false
	return s
}

func (s *Clash1) Mutable() bool {
	return s.mutable
}

func (s *Clash1) WithMutable(f func(si *Clash1)) *Clash1 {
	res := s.AsMutable()
	f(res)
	res = res.AsImmutable(s)

	return res
}

func (s *Clash1) WithImmutable(f func(si *Clash1)) *Clash1 {
	prev := s.mutable
	s.mutable = false
	f(s)
	s.mutable = prev

	return s
}

func (s *Clash1) IsDeeplyNonMutable(seen map[interface{}]bool) bool {
	if s == nil {
		return true
	}

	if s.Mutable() {
		return false
	}

	if seen == nil {
		return s.IsDeeplyNonMutable(make(map[interface{}]bool))
	}

	if seen[s] {
		return true
	}

	seen[s] = true
	return true
}
func (s *Clash1) Clash() string {
	return s.field_Clash
}

// SetClash is the setter for Clash()
func (s *Clash1) SetClash(n string) *Clash1 {
	if s.mutable {
		s.field_Clash = n
		return s
	}

	res := *s
	res.field_Clash = n
	return &res
}
func (s *Clash1) NoClash1() string {
	return s.field_NoClash1
}

// SetNoClash1 is the setter for NoClash1()
func (s *Clash1) SetNoClash1(n string) *Clash1 {
	if s.mutable {
		s.field_NoClash1 = n
		return s
	}

	res := *s
	res.field_NoClash1 = n
	return &res
}

// types for testing embedding
//
// Embed1 is an immutable type and has the following template:
//
// 	struct {
// 		Name	string
// 		*Embed2
// 		*pkga.PkgA
// 		*Clash1
// 		*pkga.Clash2
// 		NonImmStruct
// 		pkga.NonImmStructA
// 	}
//
type Embed1 struct {
	field_Name              string
	anonfield_Embed2        *Embed2
	anonfield_PkgA          *pkga.PkgA
	anonfield_Clash1        *Clash1
	anonfield_Clash2        *pkga.Clash2
	anonfield_NonImmStruct  NonImmStruct
	anonfield_NonImmStructA pkga.NonImmStructA

	mutable bool
	__tmpl  *_Imm_Embed1
}

var _ immutable.Immutable = new(Embed1)
var _ = new(Embed1).__tmpl

func (s *Embed1) AsMutable() *Embed1 {
	if s.Mutable() {
		return s
	}

	res := *s
	res.mutable = true
	return &res
}

func (s *Embed1) AsImmutable(v *Embed1) *Embed1 {
	if s == nil {
		return nil
	}

	if s == v {
		return s
	}

	s.mutable = false
	return s
}

func (s *Embed1) Mutable() bool {
	return s.mutable
}

func (s *Embed1) WithMutable(f func(si *Embed1)) *Embed1 {
	res := s.AsMutable()
	f(res)
	res = res.AsImmutable(s)

	return res
}

func (s *Embed1) WithImmutable(f func(si *Embed1)) *Embed1 {
	prev := s.mutable
	s.mutable = false
	f(s)
	s.mutable = prev

	return s
}

func (s *Embed1) IsDeeplyNonMutable(seen map[interface{}]bool) bool {
	if s == nil {
		return true
	}

	if s.Mutable() {
		return false
	}

	if seen == nil {
		return s.IsDeeplyNonMutable(make(map[interface{}]bool))
	}

	if seen[s] {
		return true
	}

	seen[s] = true
	{
		v := s.anonfield_Embed2

		if v != nil && !v.IsDeeplyNonMutable(seen) {
			return false
		}
	}
	{
		v := s.anonfield_PkgA

		if v != nil && !v.IsDeeplyNonMutable(seen) {
			return false
		}
	}
	{
		v := s.anonfield_Clash1

		if v != nil && !v.IsDeeplyNonMutable(seen) {
			return false
		}
	}
	{
		v := s.anonfield_Clash2

		if v != nil && !v.IsDeeplyNonMutable(seen) {
			return false
		}
	}
	return true
}
func (s *Embed1) Address() string {
	return s.PkgA().Address()
}
func (s *Embed1) SetAddress(n string) *Embed1 {
	v1 := s.PkgA().SetAddress(n)
	v0 := s.SetPkgA(v1)
	return v0
}
func (s *Embed1) Age() int {
	return s.Embed2().Age()
}
func (s *Embed1) SetAge(n int) *Embed1 {
	v1 := s.Embed2().SetAge(n)
	v0 := s.SetEmbed2(v1)
	return v0
}
func (s *Embed1) Clash1() *Clash1 {
	return s.anonfield_Clash1
}

// SetClash1 is the setter for Clash1()
func (s *Embed1) SetClash1(n *Clash1) *Embed1 {
	if s.mutable {
		s.anonfield_Clash1 = n
		return s
	}

	res := *s
	res.anonfield_Clash1 = n
	return &res
}
func (s *Embed1) Clash2() *pkga.Clash2 {
	return s.anonfield_Clash2
}

// SetClash2 is the setter for Clash2()
func (s *Embed1) SetClash2(n *pkga.Clash2) *Embed1 {
	if s.mutable {
		s.anonfield_Clash2 = n
		return s
	}

	res := *s
	res.anonfield_Clash2 = n
	return &res
}
func (s *Embed1) Embed2() *Embed2 {
	return s.anonfield_Embed2
}

// SetEmbed2 is the setter for Embed2()
func (s *Embed1) SetEmbed2(n *Embed2) *Embed1 {
	if s.mutable {
		s.anonfield_Embed2 = n
		return s
	}

	res := *s
	res.anonfield_Embed2 = n
	return &res
}
func (s *Embed1) Name() string {
	return s.field_Name
}

// SetName is the setter for Name()
func (s *Embed1) SetName(n string) *Embed1 {
	if s.mutable {
		s.field_Name = n
		return s
	}

	res := *s
	res.field_Name = n
	return &res
}
func (s *Embed1) NoClash1() string {
	return s.Clash1().NoClash1()
}
func (s *Embed1) SetNoClash1(n string) *Embed1 {
	v1 := s.Clash1().SetNoClash1(n)
	v0 := s.SetClash1(v1)
	return v0
}
func (s *Embed1) NoClash2() string {
	return s.Clash2().NoClash2()
}
func (s *Embed1) SetNoClash2(n string) *Embed1 {
	v1 := s.Clash2().SetNoClash2(n)
	v0 := s.SetClash2(v1)
	return v0
}
func (s *Embed1) NonImmStruct() NonImmStruct {
	return s.anonfield_NonImmStruct
}

// SetNonImmStruct is the setter for NonImmStruct()
func (s *Embed1) SetNonImmStruct(n NonImmStruct) *Embed1 {
	if s.mutable {
		s.anonfield_NonImmStruct = n
		return s
	}

	res := *s
	res.anonfield_NonImmStruct = n
	return &res
}
func (s *Embed1) NonImmStructA() pkga.NonImmStructA {
	return s.anonfield_NonImmStructA
}

// SetNonImmStructA is the setter for NonImmStructA()
func (s *Embed1) SetNonImmStructA(n pkga.NonImmStructA) *Embed1 {
	if s.mutable {
		s.anonfield_NonImmStructA = n
		return s
	}

	res := *s
	res.anonfield_NonImmStructA = n
	return &res
}
func (s *Embed1) Now() time.Time {
	return s.NonImmStruct().Now
}
func (s *Embed1) SetNow(n time.Time) *Embed1 {
	v1 := s.NonImmStruct()
	v1.Now = n
	v0 := s.SetNonImmStruct(v1)
	return v0
}
func (s *Embed1) NowA() time.Time {
	return s.NonImmStructA().NowA
}
func (s *Embed1) SetNowA(n time.Time) *Embed1 {
	v1 := s.NonImmStructA()
	v1.NowA = n
	v0 := s.SetNonImmStructA(v1)
	return v0
}
func (s *Embed1) Other() *Other {
	return s.NonImmStruct().Other
}
func (s *Embed1) SetOther(n *Other) *Embed1 {
	v1 := s.NonImmStruct()
	v1.Other = n
	v0 := s.SetNonImmStruct(v1)
	return v0
}
func (s *Embed1) OtherA() *pkga.OtherA {
	return s.NonImmStructA().OtherA
}
func (s *Embed1) SetOtherA(n *pkga.OtherA) *Embed1 {
	v1 := s.NonImmStructA()
	v1.OtherA = n
	v0 := s.SetNonImmStructA(v1)
	return v0
}
func (s *Embed1) OtherName() string {
	return s.NonImmStruct().Other.OtherName()
}
func (s *Embed1) SetOtherName(n string) *Embed1 {
	v2 := s.NonImmStruct().Other.SetOtherName(n)
	v1 := s.NonImmStruct()
	v1.Other = v2
	v0 := s.SetNonImmStruct(v1)
	return v0
}
func (s *Embed1) OtherNameA() string {
	return s.NonImmStructA().OtherA.OtherNameA()
}
func (s *Embed1) SetOtherNameA(n string) *Embed1 {
	v2 := s.NonImmStructA().OtherA.SetOtherNameA(n)
	v1 := s.NonImmStructA()
	v1.OtherA = v2
	v0 := s.SetNonImmStructA(v1)
	return v0
}
func (s *Embed1) PkgA() *pkga.PkgA {
	return s.anonfield_PkgA
}

// SetPkgA is the setter for PkgA()
func (s *Embed1) SetPkgA(n *pkga.PkgA) *Embed1 {
	if s.mutable {
		s.anonfield_PkgA = n
		return s
	}

	res := *s
	res.anonfield_PkgA = n
	return &res
}
func (s *Embed1) PkgB() *pkgb.PkgB {
	return s.PkgA().PkgB()
}
func (s *Embed1) SetPkgB(n *pkgb.PkgB) *Embed1 {
	v1 := s.PkgA().SetPkgB(n)
	v0 := s.SetPkgA(v1)
	return v0
}
func (s *Embed1) Postcode() string {
	return s.PkgA().PkgB().Postcode()
}
func (s *Embed1) SetPostcode(n string) *Embed1 {
	v2 := s.PkgA().PkgB().SetPostcode(n)
	v1 := s.PkgA().SetPkgB(v2)
	v0 := s.SetPkgA(v1)
	return v0
}

//
// Embed2 is an immutable type and has the following template:
//
// 	struct {
// 		Age int
// 	}
//
type Embed2 struct {
	field_Age int

	mutable bool
	__tmpl  *_Imm_Embed2
}

var _ immutable.Immutable = new(Embed2)
var _ = new(Embed2).__tmpl

func (s *Embed2) AsMutable() *Embed2 {
	if s.Mutable() {
		return s
	}

	res := *s
	res.mutable = true
	return &res
}

func (s *Embed2) AsImmutable(v *Embed2) *Embed2 {
	if s == nil {
		return nil
	}

	if s == v {
		return s
	}

	s.mutable = false
	return s
}

func (s *Embed2) Mutable() bool {
	return s.mutable
}

func (s *Embed2) WithMutable(f func(si *Embed2)) *Embed2 {
	res := s.AsMutable()
	f(res)
	res = res.AsImmutable(s)

	return res
}

func (s *Embed2) WithImmutable(f func(si *Embed2)) *Embed2 {
	prev := s.mutable
	s.mutable = false
	f(s)
	s.mutable = prev

	return s
}

func (s *Embed2) IsDeeplyNonMutable(seen map[interface{}]bool) bool {
	if s == nil {
		return true
	}

	if s.Mutable() {
		return false
	}

	if seen == nil {
		return s.IsDeeplyNonMutable(make(map[interface{}]bool))
	}

	if seen[s] {
		return true
	}

	seen[s] = true
	return true
}
func (s *Embed2) Age() int {
	return s.field_Age
}

// SetAge is the setter for Age()
func (s *Embed2) SetAge(n int) *Embed2 {
	if s.mutable {
		s.field_Age = n
		return s
	}

	res := *s
	res.field_Age = n
	return &res
}

//
// Other is an immutable type and has the following template:
//
// 	struct {
// 		OtherName string
// 	}
//
type Other struct {
	field_OtherName string

	mutable bool
	__tmpl  *_Imm_Other
}

var _ immutable.Immutable = new(Other)
var _ = new(Other).__tmpl

func (s *Other) AsMutable() *Other {
	if s.Mutable() {
		return s
	}

	res := *s
	res.mutable = true
	return &res
}

func (s *Other) AsImmutable(v *Other) *Other {
	if s == nil {
		return nil
	}

	if s == v {
		return s
	}

	s.mutable = false
	return s
}

func (s *Other) Mutable() bool {
	return s.mutable
}

func (s *Other) WithMutable(f func(si *Other)) *Other {
	res := s.AsMutable()
	f(res)
	res = res.AsImmutable(s)

	return res
}

func (s *Other) WithImmutable(f func(si *Other)) *Other {
	prev := s.mutable
	s.mutable = false
	f(s)
	s.mutable = prev

	return s
}

func (s *Other) IsDeeplyNonMutable(seen map[interface{}]bool) bool {
	if s == nil {
		return true
	}

	if s.Mutable() {
		return false
	}

	if seen == nil {
		return s.IsDeeplyNonMutable(make(map[interface{}]bool))
	}

	if seen[s] {
		return true
	}

	seen[s] = true
	return true
}
func (s *Other) OtherName() string {
	return s.field_OtherName
}

// SetOtherName is the setter for OtherName()
func (s *Other) SetOtherName(n string) *Other {
	if s.mutable {
		s.field_OtherName = n
		return s
	}

	res := *s
	res.field_OtherName = n
	return &res
}
