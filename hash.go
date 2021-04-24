package main

type HashType int

const (
  MD5 HashType = iota
)

func (ht HashType) String() string {
  t, ok := map[HashType]string{
    MD5: "md5",
  }[ht]
  if !ok {
    return "Invalid"
  }
  return t
}

type Hash struct {
  Value string
  Type HashType
}

func ParseHash(h string) (Hash, error) {
  return Hash{
    Value: h,
    Type:  MD5,
  }, nil
}