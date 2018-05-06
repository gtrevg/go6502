package cpu_test

import (
	"strings"
	"testing"

	"github.com/beevik/go6502/asm"
	"github.com/beevik/go6502/cpu"
)

func runCPU(t *testing.T, asmString string, steps int) *cpu.CPU {
	b := strings.NewReader(asmString)
	r, sm, err := asm.Assemble(b, "test.asm", 0)
	if err != nil {
		t.Error(err)
		return nil
	}

	mem := cpu.NewFlatMemory()
	cpu := cpu.NewCPU(cpu.NMOS, mem)
	mem.StoreBytes(sm.Origin, r.Code)
	cpu.SetPC(sm.Origin)

	for i := 0; i < steps; i++ {
		cpu.Step()
	}

	return cpu
}

func expectPC(t *testing.T, cpu *cpu.CPU, pc uint16) {
	if cpu.Reg.PC != pc {
		t.Errorf("PC incorrect. exp: $%04X, got: $%04X", pc, cpu.Reg.PC)
	}
}

func expectCycles(t *testing.T, cpu *cpu.CPU, cycles uint64) {
	if cpu.Cycles != cycles {
		t.Errorf("Cycles incorrect. exp: %d, got: %d", cycles, cpu.Cycles)
	}
}

func expectACC(t *testing.T, cpu *cpu.CPU, acc byte) {
	if cpu.Reg.A != acc {
		t.Errorf("Accumulator incorrect. exp: $%02X, got: $%02X", acc, cpu.Reg.A)
	}
}

func expectMem(t *testing.T, cpu *cpu.CPU, addr uint16, v byte) {
	got := cpu.Mem.LoadByte(addr)
	if got != v {
		t.Errorf("Memory at $%04X incorrect. exp: $%02X, got: $%02X", addr, v, got)
	}
}

func TestAccumulator(t *testing.T) {
	asm := `
	.ORG $1000
	LDA #$5e
	STA $15
	STA $1500`

	cpu := runCPU(t, asm, 3)
	if cpu == nil {
		return
	}

	expectPC(t, cpu, 0x1007)
	expectCycles(t, cpu, 9)
	expectACC(t, cpu, 0x5e)
	expectMem(t, cpu, 0x15, 0x5e)
	expectMem(t, cpu, 0x1500, 0x5e)
}

func TestPageCross(t *testing.T) {
	asm := `
	.ORG $1000
	LDA #$55		; 2 cycles
	STA $1101		; 4 cycles
	LDA #$00		; 2 cycles
	LDX #$FF		; 2 cycles
	LDA $1002,X		; 5 cycles`

	cpu := runCPU(t, asm, 5)
	if cpu == nil {
		return
	}

	expectPC(t, cpu, 0x100c)
	expectCycles(t, cpu, 15)
	expectACC(t, cpu, 0x55)
	expectMem(t, cpu, 0x1101, 0x55)
}

func TestUnused65c02(t *testing.T) {
	asm := `
	.ORG $1000
	.ARCH 65c02
	.DH 0200
	.DH 03
	.DH 07
	.DH 0b
	.DH 0f
	.DH fc0102`

	cpu := runCPU(t, asm, 6)
	if cpu == nil {
		return
	}

	expectPC(t, cpu, 0x1009)
	expectCycles(t, cpu, 10)
}
