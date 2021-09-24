package load

const Uint32Size = 32 << (^uint32(0) >> 32 & 1) // 32 or 64
const MaxInt32 = 1<<(Uint32Size-1) - 1
