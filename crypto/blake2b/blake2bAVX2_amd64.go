





package blake2b

import "golang.org/x/sys/cpu"

func init() {
	useAVX2 = cpu.X86.HasAVX2
	useAVX = cpu.X86.HasAVX
	useSSE4 = cpu.X86.HasSSE41
}


func fAVX2(h *[8]uint64, m *[16]uint64, c0, c1 uint64, flag uint64, rounds uint64)


func fAVX(h *[8]uint64, m *[16]uint64, c0, c1 uint64, flag uint64, rounds uint64)


func fSSE4(h *[8]uint64, m *[16]uint64, c0, c1 uint64, flag uint64, rounds uint64)

func f(h *[8]uint64, m *[16]uint64, c0, c1 uint64, flag uint64, rounds uint64) {
	switch {
	case useAVX2:
		fAVX2(h, m, c0, c1, flag, rounds)
	case useAVX:
		fAVX(h, m, c0, c1, flag, rounds)
	case useSSE4:
		fSSE4(h, m, c0, c1, flag, rounds)
	default:
		fGeneric(h, m, c0, c1, flag, rounds)
	}
}
